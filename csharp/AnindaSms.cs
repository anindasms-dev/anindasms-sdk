using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;

namespace AnindaSms
{
    public class AnindaSmsException : Exception
    {
        public int StatusCode { get; }
        public JsonElement? Data { get; }

        public AnindaSmsException(string message, int statusCode = 0, JsonElement? data = null)
            : base(message)
        {
            StatusCode = statusCode;
            Data = data;
        }
    }

    public class AnindaSmsClient
    {
        public const string DefaultBaseUrl = "https://anindasms.com/api/v1";
        private readonly HttpClient _http;
        private readonly string _baseUrl;

        public AnindaSmsClient(string apiKey, string baseUrl = DefaultBaseUrl, HttpClient httpClient = null)
        {
            if (string.IsNullOrWhiteSpace(apiKey))
                throw new ArgumentException("API key is required.", nameof(apiKey));

            _baseUrl = baseUrl.TrimEnd('/');
            _http = httpClient ?? new HttpClient();
            _http.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", apiKey.Trim());
            _http.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
        }

        private async Task<JsonElement> RequestAsync(HttpMethod method, string path, object body = null)
        {
            using var req = new HttpRequestMessage(method, $"{_baseUrl}{path}");
            if (body != null)
            {
                var json = JsonSerializer.Serialize(body);
                req.Content = new StringContent(json, Encoding.UTF8, "application/json");
            }

            HttpResponseMessage res;
            try
            {
                res = await _http.SendAsync(req);
            }
            catch (Exception ex)
            {
                throw new AnindaSmsException($"Connection error: {ex.Message}");
            }

            var text = await res.Content.ReadAsStringAsync();
            JsonDocument doc;
            try
            {
                doc = JsonDocument.Parse(text);
            }
            catch
            {
                throw new AnindaSmsException($"Invalid JSON from API (HTTP {(int)res.StatusCode})", (int)res.StatusCode);
            }

            var root = doc.RootElement;
            bool success = root.TryGetProperty("success", out var s) && s.GetBoolean();
            if (!success)
            {
                string err = root.TryGetProperty("error", out var e) ? e.GetString() : "Unknown error";
                JsonElement? extra = root.TryGetProperty("data", out var d) ? d : null;
                throw new AnindaSmsException(err, (int)res.StatusCode, extra);
            }

            return root.GetProperty("data");
        }

        public Task<JsonElement> MeAsync() => RequestAsync(HttpMethod.Get, "/me");
        public Task<JsonElement> BalanceAsync() => RequestAsync(HttpMethod.Get, "/balance");
        public Task<JsonElement> ServicesAsync() => RequestAsync(HttpMethod.Get, "/services");

        public Task<JsonElement> PricesAsync(string service, string country = null)
        {
            var url = $"/prices?service={Uri.EscapeDataString(service)}";
            if (!string.IsNullOrEmpty(country))
                url += $"&country={Uri.EscapeDataString(country)}";
            return RequestAsync(HttpMethod.Get, url);
        }

        public Task<JsonElement> CreateOrderAsync(string service, string country, string operatorName = "any")
        {
            var body = new Dictionary<string, string>
            {
                { "service", service },
                { "country", country },
                { "operator", operatorName ?? "any" }
            };
            return RequestAsync(HttpMethod.Post, "/order", body);
        }

        public Task<JsonElement> GetOrderAsync(int orderId) => RequestAsync(HttpMethod.Get, $"/order/{orderId}");
        public Task<JsonElement> CancelOrderAsync(int orderId) => RequestAsync(HttpMethod.Post, $"/order/{orderId}/cancel");
        public Task<JsonElement> CompleteOrderAsync(int orderId) => RequestAsync(HttpMethod.Post, $"/order/{orderId}/done");

        public async Task<JsonElement> WaitForCodeAsync(int orderId, int timeoutSeconds = 600, int intervalSeconds = 5, CancellationToken ct = default)
        {
            var deadline = DateTime.UtcNow.AddSeconds(timeoutSeconds);
            var delay = TimeSpan.FromSeconds(Math.Max(2, intervalSeconds));

            while (DateTime.UtcNow < deadline)
            {
                ct.ThrowIfCancellationRequested();
                var order = await GetOrderAsync(orderId);
                var status = order.TryGetProperty("status", out var st) ? st.GetString() : "";

                if (order.TryGetProperty("code", out var code) && code.ValueKind == JsonValueKind.String && !string.IsNullOrEmpty(code.GetString()))
                {
                    if (status == "received" || status == "done")
                        return order;
                }

                if (status == "cancelled" || status == "expired" || status == "done")
                    throw new AnindaSmsException($"Order ended without a code (status: {status}).", 0, order);

                await Task.Delay(delay, ct);
            }

            throw new AnindaSmsException("Timed out while waiting for SMS code.");
        }
    }
}
