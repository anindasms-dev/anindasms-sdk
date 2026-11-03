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

// ---- Rent API (Numara Kiralama) -----------------------------------

func (c *Client) RentDurations() (json.RawMessage, error) {
	var out struct {
		Durations json.RawMessage `json:"durations"`
	}
	if err := c.request("GET", "/rent/durations", nil, &out); err != nil {
		return nil, err
	}
	return out.Durations, nil
}

func (c *Client) RentCountries(duration int) (json.RawMessage, error) {
	if duration <= 0 {
		duration = 24
	}
	var out struct {
		Countries json.RawMessage `json:"countries"`
	}
	if err := c.request("GET", fmt.Sprintf("/rent/countries?duration=%d", duration), nil, &out); err != nil {
		return nil, err
	}
	return out.Countries, nil
}

func (c *Client) RentServices(country string, duration int) (json.RawMessage, error) {
	if duration <= 0 {
		duration = 24
	}
	var out struct {
		Services json.RawMessage `json:"services"`
	}
	u := fmt.Sprintf("/rent/services?country=%s&duration=%d", url.QueryEscape(country), duration)
	if err := c.request("GET", u, nil, &out); err != nil {
		return nil, err
	}
	return out.Services, nil
}

func (c *Client) RentPrice(service, country string, duration int) (json.RawMessage, error) {
	if duration <= 0 {
		duration = 24
	}
	u := fmt.Sprintf("/rent/prices?service=%s&country=%s&duration=%d", url.QueryEscape(service), url.QueryEscape(country), duration)
	var out json.RawMessage
	if err := c.request("GET", u, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) RentBuy(service, country string, duration int) (json.RawMessage, error) {
	if duration <= 0 {
		duration = 24
	}
	body := map[string]any{
		"service":  service,
		"country":  country,
		"duration": duration,
	}
	var out json.RawMessage
	if err := c.request("POST", "/rent/buy", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) RentOrders() (json.RawMessage, error) {
	var out struct {
		Orders json.RawMessage `json:"orders"`
	}
	if err := c.request("GET", "/rent/orders", nil, &out); err != nil {
		return nil, err
	}
	return out.Orders, nil
}

func (c *Client) RentOrder(orderID int) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.request("GET", fmt.Sprintf("/rent/orders/%d", orderID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) RentCancel(orderID int) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.request("POST", fmt.Sprintf("/rent/orders/%d/cancel", orderID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) RentReactivate(orderID int, duration int) (json.RawMessage, error) {
	if duration <= 0 {
		duration = 24
	}
	body := map[string]int{
		"duration": duration,
	}
	var out json.RawMessage
	if err := c.request("POST", fmt.Sprintf("/rent/orders/%d/reactivate", orderID), body, &out); err != nil {
		return nil, err
	}
	return out, nil
}
// ---- Proxy API ----------------------------------------------------

func (c *Client) ProxyTypes() (json.RawMessage, error) {
	var out struct {
		Types json.RawMessage `json:"types"`
	}
	if err := c.request("GET", "/proxy/types", nil, &out); err != nil {
		return nil, err
	}
	return out.Types, nil
}

func (c *Client) ProxyCountries(proxyType string) (json.RawMessage, error) {
	path := "/proxy/countries"
	if proxyType != "" {
		path += "?type=" + url.QueryEscape(proxyType)
	}
	var out struct {
		Countries json.RawMessage `json:"countries"`
	}
	if err := c.request("GET", path, nil, &out); err != nil {
		return nil, err
	}
	return out.Countries, nil
}

func (c *Client) ProxyPrice(proxyType, country string, period, quantity int) (json.RawMessage, error) {
	if period <= 0 {
		period = 30
	}
	if quantity <= 0 {
		quantity = 1
	}
	u := fmt.Sprintf("/proxy/price?type=%s&period=%d&quantity=%d", url.QueryEscape(proxyType), period, quantity)
	if country != "" {
		u += "&country=" + url.QueryEscape(country)
	}
	var out json.RawMessage
	if err := c.request("GET", u, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) ProxyBuy(proxyType, country string, period, quantity, rotationInterval int) (json.RawMessage, error) {
	body := map[string]any{
		"type":     proxyType,
		"country":  country,
		"period":   period,
		"quantity": quantity,
	}
	if rotationInterval > 0 {
		body["rotation_interval"] = rotationInterval
	}
	var out json.RawMessage
	if err := c.request("POST", "/proxy/buy", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) ProxyOrders(page, limit int) (json.RawMessage, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	var out struct {
		Orders json.RawMessage `json:"orders"`
	}
	if err := c.request("GET", fmt.Sprintf("/proxy/orders?page=%d&limit=%d", page, limit), nil, &out); err != nil {
		return nil, err
	}
	return out.Orders, nil
}

func (c *Client) ProxyOrder(orderID int) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.request("GET", fmt.Sprintf("/proxy/orders/%d", orderID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---- eSIM API -----------------------------------------------------

func (c *Client) EsimCountries() (json.RawMessage, error) {
	var out struct {
		Countries json.RawMessage `json:"countries"`
	}
	if err := c.request("GET", "/esim/countries", nil, &out); err != nil {
		return nil, err
	}
	return out.Countries, nil
}

func (c *Client) EsimPackages(country string) (json.RawMessage, error) {
	path := "/esim/packages"
	if country != "" {
		path += "?country=" + url.QueryEscape(country)
	}
	var out struct {
		Packages json.RawMessage `json:"packages"`
	}
	if err := c.request("GET", path, nil, &out); err != nil {
		return nil, err
	}
	return out.Packages, nil
}

func (c *Client) EsimBuy(packageID string, quantity int) (json.RawMessage, error) {
	if quantity <= 0 {
		quantity = 1
	}
	body := map[string]any{
		"package_id": packageID,
		"quantity":   quantity,
	}
	var out json.RawMessage
	if err := c.request("POST", "/esim/buy", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) EsimOrders(page, limit int) (json.RawMessage, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	var out struct {
		Orders json.RawMessage `json:"orders"`
	}
	if err := c.request("GET", fmt.Sprintf("/esim/orders?page=%d&limit=%d", page, limit), nil, &out); err != nil {
		return nil, err
	}
	return out.Orders, nil
}

func (c *Client) EsimOrder(orderID int) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.request("GET", fmt.Sprintf("/esim/orders/%d", orderID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) EsimTopup(orderID int, packageID string) (json.RawMessage, error) {
	body := map[string]string{
		"package_id": packageID,
	}
	var out json.RawMessage
	if err := c.request("POST", fmt.Sprintf("/esim/orders/%d/topup", orderID), body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---- Gecici Mail API ----------------------------------------------

func (c *Client) MailDomains() (json.RawMessage, error) {
	var out struct {
		Domains json.RawMessage `json:"domains"`
	}
	if err := c.request("GET", "/mail/domains", nil, &out); err != nil {
		return nil, err
	}
	return out.Domains, nil
}

func (c *Client) MailCreate(duration int, domain, prefix string) (json.RawMessage, error) {
	if duration <= 0 {
		duration = 60
	}
	body := map[string]any{
		"duration": duration,
	}
	if domain != "" {
		body["domain"] = domain
	}
	if prefix != "" {
		body["prefix"] = prefix
	}
	var out json.RawMessage
	if err := c.request("POST", "/mail/create", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) MailOrders(page, limit int) (json.RawMessage, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	var out struct {
		Orders json.RawMessage `json:"orders"`
	}
	if err := c.request("GET", fmt.Sprintf("/mail/orders?page=%d&limit=%d", page, limit), nil, &out); err != nil {
		return nil, err
	}
	return out.Orders, nil
}

func (c *Client) MailOrder(orderID int) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.request("GET", fmt.Sprintf("/mail/orders/%d", orderID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) MailMessages(orderID int) (json.RawMessage, error) {
	var out struct {
		Messages json.RawMessage `json:"messages"`
	}
	if err := c.request("GET", fmt.Sprintf("/mail/orders/%d/messages", orderID), nil, &out); err != nil {
		return nil, err
	}
	return out.Messages, nil
}

func (c *Client) MailExtend(orderID int, duration int) (json.RawMessage, error) {
	if duration <= 0 {
		duration = 60
	}
	body := map[string]int{
		"duration": duration,
	}
	var out json.RawMessage
	if err := c.request("POST", fmt.Sprintf("/mail/orders/%d/extend", orderID), body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) MailCancel(orderID int) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.request("POST", fmt.Sprintf("/mail/orders/%d/cancel", orderID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
