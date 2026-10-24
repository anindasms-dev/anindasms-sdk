package anindasms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://anindasms.com/api/v1"
	Version        = "1.0.0"
)

type APIError struct {
	Message    string
	StatusCode int
	Data       json.RawMessage
}

func (e *APIError) Error() string {
	return fmt.Sprintf("anindasms api error (status %d): %s", e.StatusCode, e.Message)
}

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string, opts ...func(*Client)) *Client {
	c := &Client{
		apiKey:     strings.TrimSpace(apiKey),
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithBaseURL(u string) func(*Client) {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(u, "/")
	}
}

func WithHTTPClient(hc *http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = hc
	}
}

type apiEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
}

func (c *Client) request(method, path string, body any, out any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "anindasms-go/"+Version)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("invalid json response (status %d): %w", res.StatusCode, err)
	}

	if !env.Success {
		return &APIError{
			Message:    env.Error,
			StatusCode: res.StatusCode,
			Data:       env.Data,
		}
	}

	if out != nil && len(env.Data) > 0 {
		return json.Unmarshal(env.Data, out)
	}
	return nil
}

type MeResponse struct {
	ID              int     `json:"id"`
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	Balance         float64 `json:"balance"`
	Currency        string  `json:"currency"`
	APIKeyPrefix    string  `json:"api_key_prefix"`
	APIKeyCreatedAt string  `json:"api_key_created_at"`
}

func (c *Client) Me() (*MeResponse, error) {
	var out MeResponse
	if err := c.request("GET", "/me", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type BalanceResponse struct {
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	Formatted string  `json:"formatted"`
}

func (c *Client) Balance() (*BalanceResponse, error) {
	var out BalanceResponse
	if err := c.request("GET", "/balance", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type ServiceItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (c *Client) Services() ([]ServiceItem, error) {
	var out struct {
		Services []ServiceItem `json:"services"`
	}
	if err := c.request("GET", "/services", nil, &out); err != nil {
		return nil, err
	}
	return out.Services, nil
}

type OrderResponse struct {
	OrderID   int     `json:"order_id"`
	Phone     string  `json:"phone"`
	Service   string  `json:"service"`
	Country   string  `json:"country"`
	Operator  *string `json:"operator"`
	Cost      float64 `json:"cost"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"status"`
	Code      *string `json:"code"`
	SMSText   *string `json:"sms_text"`
	Remaining int     `json:"remaining"`
	ExpiresAt string  `json:"expires_at"`
}

func (c *Client) CreateOrder(service, country, operator string) (*OrderResponse, error) {
	if operator == "" {
		operator = "any"
	}
	body := map[string]string{
		"service":  service,
		"country":  country,
		"operator": operator,
	}
	var out OrderResponse
	if err := c.request("POST", "/order", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetOrder(orderID int) (*OrderResponse, error) {
	var out OrderResponse
	if err := c.request("GET", fmt.Sprintf("/order/%d", orderID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CancelOrder(orderID int) (*OrderResponse, error) {
	var out OrderResponse
	if err := c.request("POST", fmt.Sprintf("/order/%d/cancel", orderID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CompleteOrder(orderID int) (*OrderResponse, error) {
	var out OrderResponse
	if err := c.request("POST", fmt.Sprintf("/order/%d/done", orderID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) WaitForCode(orderID int, timeout, interval time.Duration) (*OrderResponse, error) {
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}
	if interval < 2*time.Second {
		interval = 5 * time.Second
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		order, err := c.GetOrder(orderID)
		if err != nil {
			return nil, err
		}

		if order.Code != nil && *order.Code != "" && (order.Status == "received" || order.Status == "done") {
			return order, nil
		}

		if order.Status == "cancelled" || order.Status == "expired" || order.Status == "done" {
			return nil, fmt.Errorf("order ended without code (status: %s)", order.Status)
		}

		time.Sleep(interval)
	}

	return nil, errors.New("timed out waiting for SMS code")
}
