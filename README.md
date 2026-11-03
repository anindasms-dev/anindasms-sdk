# AnindaSMS — Official REST API v1 & Multi-Language SDKs

[![API Version](https://img.shields.io/badge/API-v1-blue.svg)](https://anindasms.com/api)
[![Countries](https://img.shields.io/badge/Countries-200%2B-orange.svg)](https://anindasms.com/api)
[![Services](https://img.shields.io/badge/Services-5000%2B-red.svg)](https://anindasms.com/api)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PHP](https://img.shields.io/badge/PHP-7.4%2B%20%7C%208.x-8892BF.svg)](php/)
[![Python](https://img.shields.io/badge/Python-3.8%2B-3776AB.svg)](python/)
[![Node.js](https://img.shields.io/badge/Node.js-18%2B-339933.svg)](node/)
[![.NET](https://img.shields.io/badge/.NET-6.0%2B-512BD4.svg)](csharp/)
[![Go](https://img.shields.io/badge/Go-1.18%2B-00ADD8.svg)](go/)

> 🌍 **Language / Dil:** **English** | [🇹🇷 Türkçe Dokümantasyon için tıklayın (README.tr.md)](README.tr.md)

Welcome to the official repository for the [AnindaSMS](https://anindasms.com) **REST API v1** and open-source SDK libraries. AnindaSMS provides temporary, disposable virtual numbers for SMS verification across **5000+ online services** and **200+ countries** worldwide.

---

## 📑 Table of Contents
1. [Overview](#-overview)
   - [Base URL](#base-url)
   - [Authentication](#authentication)
   - [Standard Response Envelope](#standard-response-envelope)
   - [Rate Limiting](#rate-limiting)
2. [REST API Endpoints](#-rest-api-endpoints)
   - [1. User Profile (`GET /me`)](#1-user-profile-get-me)
   - [2. Balance Query (`GET /balance`)](#2-balance-query-get-balance)
   - [3. Service Catalog (`GET /services`)](#3-service-catalog-get-services)
   - [4. Prices & Live Stock (`GET /prices`)](#4-prices--live-stock-get-prices)
   - [5. Create Order / Buy Number (`POST /order`)](#5-create-order--buy-number-post-order)
   - [6. Check Order & SMS Code (`GET /order/{id}`)](#6-check-order--sms-code-get-orderid)
   - [7. Cancel Order (`POST /order/{id}/cancel`)](#7-cancel-order-post-orderidcancel)
   - [8. Complete Order (`POST /order/{id}/done`)](#8-complete-order-post-orderiddone)
3. [Order Lifecycle & Statuses](#-order-lifecycle--statuses)
4. [Error Codes Reference](#-error-codes-reference)
5. [Official SDK Clients](#-official-sdk-clients)
   - [PHP](#-php-sdk)
   - [Python](#-python-sdk)
   - [Node.js](#-nodejs-sdk)
   - [C# / .NET](#-c--net-sdk)
   - [Go](#-go-sdk)
6. [License & Support](#-license--support)

---

## 🌐 Overview

### Base URL
All API requests must be made over HTTPS to the following endpoint:
```
https://anindasms.com/api/v1
```

### Authentication
You can generate your personal API Key in your dashboard at [anindasms.com/profile](https://anindasms.com/profile).

Authenticate your requests by passing your key in an HTTP header using either method:
```http
Authorization: Bearer YOUR_API_KEY
```
or
```http
X-API-Key: YOUR_API_KEY
```
> 🔒 **Security Note:** Never pass the API key in the URL query string (`?api_key=...`) or within request bodies. Such requests are rejected with `401 Unauthorized` for your account security.

### Standard Response Envelope
All responses are formatted in UTF-8 JSON:

**Success Response:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Descriptive error message.",
  "data": null
}
```

### Rate Limiting
To ensure fair resource distribution and platform stability, requests are capped per minute:
* **Global Rate Limit:** Maximum **120 requests / minute**
* **Order Creation (`POST /order`):** Maximum **30 requests / minute**

Response headers indicate your current quota status:
* `X-RateLimit-Limit`: Maximum requests allowed per minute.
* `X-RateLimit-Remaining`: Remaining request allowance in the current window.
* `X-RateLimit-Reset`: Unix timestamp when the limit resets.

---

## 📡 REST API Endpoints

### 1. User Profile (`GET /me`)
Returns account identifier, username, registered email, current balance, and API key metadata.

```bash
curl -X GET "https://anindasms.com/api/v1/me" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Sample Response:**
```json
{
  "success": true,
  "data": {
    "id": 2,
    "username": "developer",
    "email": "developer@example.com",
    "balance": 150.00,
    "currency": "TRY",
    "api_key_prefix": "abab247a",
    "api_key_created_at": "2026-10-24 12:00:00",
    "api_key_last_used_at": "2026-10-24 15:30:10"
  }
}
```

---

### 2. Balance Query (`GET /balance`)
Retrieves the real-time available balance for the authenticated account.

```bash
curl -X GET "https://anindasms.com/api/v1/balance" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Sample Response:**
```json
{
  "success": true,
  "data": {
    "balance": 150.00,
    "currency": "TRY",
    "formatted": "₺150,00"
  }
}
```

---

### 3. Service Catalog (`GET /services`)
Returns the list of all currently supported services (WhatsApp, Telegram, Google, etc.).

```bash
curl -X GET "https://anindasms.com/api/v1/services" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Sample Response:**
```json
{
  "success": true,
  "data": {
    "services": [
      { "code": "wa", "name": "WhatsApp" },
      { "code": "tg", "name": "Telegram" },
      { "code": "go", "name": "Google / Gmail" },
      { "code": "ig", "name": "Instagram" },
      { "code": "ds", "name": "Discord" },
      { "code": "tt", "name": "TikTok" },
      { "code": "oi", "name": "OpenAI / ChatGPT" }
    ]
  }
}
```

---

### 4. Prices & Live Stock (`GET /prices`)
Queries country-based live inventory counts and unit pricing for any service.

**Query Parameters:**
* `service` (string, required): Service identifier code (e.g. `wa`, `tg`)
* `country` (string, optional): Specific country code (e.g. `62`, `tr`)

```bash
curl -X GET "https://anindasms.com/api/v1/prices?service=wa&country=62" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Sample Response:**
```json
{
  "success": true,
  "data": {
    "service": "wa",
    "prices": [
      {
        "country": "62",
        "country_name": "Indonesia",
        "price": 12.50,
        "count": 482
      }
    ]
  }
}
```

---

### 5. Create Order / Buy Number (`POST /order`)
Allocates a new virtual number for the specified service and country. The price is reserved from your balance.

**Request Body (JSON):**
* `service` (string, required): Service code (e.g. `"wa"`)
* `country` (string, required): Country code (e.g. `"62"`)
* `operator` (string, optional): Specific operator code

```bash
curl -X POST "https://anindasms.com/api/v1/order" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"service":"wa", "country":"62"}'
```

**Sample Response:**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "phone": "+628123456789",
    "country": "62",
    "service": "wa",
    "price": 12.50,
    "status": "pending",
    "expires_at": 1792845600,
    "remaining_seconds": 600
  }
}
```

---

### 6. Check Order & SMS Code (`GET /order/{id}`)
Polls the order status and receives the inbound verification code. You can poll this endpoint periodically (every 3–5 seconds recommended).

```bash
curl -X GET "https://anindasms.com/api/v1/order/58231" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Pending State (`status: pending`):**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "phone": "+628123456789",
    "status": "pending",
    "sms_code": null,
    "full_sms": null,
    "remaining_seconds": 540
  }
}
```

**Code Received State (`status: received`):**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "phone": "+628123456789",
    "status": "received",
    "sms_code": "489201",
    "full_sms": "Your WhatsApp verification code is: 489-201",
    "remaining_seconds": 490
  }
}
```

---

### 7. Cancel Order (`POST /order/{id}/cancel`)
Cancels an active order before any SMS code has been received. The reserved fee is **instantly and fully refunded** to your balance.

```bash
curl -X POST "https://anindasms.com/api/v1/order/58231/cancel" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Sample Response:**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "status": "cancelled",
    "message": "Order cancelled and balance refunded."
  }
}
```

---

### 8. Complete Order (`POST /order/{id}/done`)
Marks an order as successfully completed after receiving the verification code.

```bash
curl -X POST "https://anindasms.com/api/v1/order/58231/done" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Sample Response:**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "status": "completed",
    "message": "Order completed."
  }
}
```

---

### 9. Number Rental (Rent) API Endpoints
Rent dedicated numbers for extended durations (4 hours up to 30 days) and receive multiple inbound SMS messages during the rental window.

#### A. Rental Durations (GET /rent/durations)
`ash
curl -X GET "https://anindasms.com/api/v1/rent/durations" \
  -H "Authorization: Bearer YOUR_API_KEY"
`

#### B. Available Rental Countries (GET /rent/countries)
`ash
curl -X GET "https://anindasms.com/api/v1/rent/countries?duration=24" \
  -H "Authorization: Bearer YOUR_API_KEY"
`

#### C. Rental Services & Quotes (GET /rent/services, GET /rent/prices)
`ash
curl -X GET "https://anindasms.com/api/v1/rent/services?country=62&duration=24" \
  -H "Authorization: Bearer YOUR_API_KEY"

curl -X GET "https://anindasms.com/api/v1/rent/prices?service=wa&country=62&duration=24" \
  -H "Authorization: Bearer YOUR_API_KEY"
`

#### D. Rent Number (POST /rent/buy)
`ash
curl -X POST "https://anindasms.com/api/v1/rent/buy" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"service":"wa","country":"62","duration":24}'
`
**Sample Response:**
`json
{
  "success": true,
  "data": {
    "order_id": 74102,
    "phone": "+905559876543",
    "service": "wa",
    "country": "62",
    "duration": 24,
    "cost": 125.0,
    "balance": 250.0,
    "status": "waiting",
    "expires_at": "2026-10-25 14:00:00"
  }
}
`

#### E. Rental Orders & Polling SMS List (GET /rent/orders, GET /rent/orders/{id})
`ash
curl -X GET "https://anindasms.com/api/v1/rent/orders" \
  -H "Authorization: Bearer YOUR_API_KEY"

curl -X GET "https://anindasms.com/api/v1/rent/orders/74102" \
  -H "Authorization: Bearer YOUR_API_KEY"
`

#### F. Cancel or Reactivate Rental (POST /rent/orders/{id}/cancel, POST /rent/orders/{id}/reactivate)
`ash
# Cancel within the first 20 minutes if no SMS arrived
curl -X POST "https://anindasms.com/api/v1/rent/orders/74102/cancel" \
  -H "Authorization: Bearer YOUR_API_KEY"

# Reactivate / extend same rented number
curl -X POST "https://anindasms.com/api/v1/rent/orders/74102/reactivate" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"duration":24}'
`

---

### 10. Proxy API Endpoints
Datacenter, Residential, IPv6, and Mobile rotating proxies.

#### A. Proxy Types & Countries (GET /proxy/types, GET /proxy/countries)
`ash
curl -X GET "https://anindasms.com/api/v1/proxy/types" -H "Authorization: Bearer YOUR_API_KEY"
curl -X GET "https://anindasms.com/api/v1/proxy/countries?type=datacenter" -H "Authorization: Bearer YOUR_API_KEY"
`

#### B. Price Calculation (GET /proxy/price)
`ash
curl -X GET "https://anindasms.com/api/v1/proxy/price?type=datacenter&country=TR&period=30&quantity=1" \
  -H "Authorization: Bearer YOUR_API_KEY"
`

#### C. Buy Proxy (POST /proxy/buy)
`ash
curl -X POST "https://anindasms.com/api/v1/proxy/buy" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"type":"datacenter","country":"TR","period":30,"quantity":1}'
`
**Sample Response:**
`json
{
  "success": true,
  "data": {
    "order_id": 1042,
    "proxy_type": "datacenter",
    "ip": "185.193.124.50",
    "port": 8080,
    "username": "user_1042",
    "password": "sec_984f2b",
    "expires_at": "2026-11-24 14:00:00"
  }
}
`

#### D. Proxy Orders & Details (GET /proxy/orders, GET /proxy/orders/{id})
`ash
curl -X GET "https://anindasms.com/api/v1/proxy/orders" -H "Authorization: Bearer YOUR_API_KEY"
curl -X GET "https://anindasms.com/api/v1/proxy/orders/1042" -H "Authorization: Bearer YOUR_API_KEY"
`

---

### 11. eSIM API Endpoints
High-speed travel and global eSIM data packages for 150+ countries.

#### A. eSIM Countries & Packages (GET /esim/countries, GET /esim/packages)
`ash
curl -X GET "https://anindasms.com/api/v1/esim/countries" -H "Authorization: Bearer YOUR_API_KEY"
curl -X GET "https://anindasms.com/api/v1/esim/packages?country=TR" -H "Authorization: Bearer YOUR_API_KEY"
`

#### B. Buy eSIM Package (POST /esim/buy)
`ash
curl -X POST "https://anindasms.com/api/v1/esim/buy" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"package_id":"pkg_tr_1gb_7d","quantity":1}'
`
**Sample Response:**
`json
{
  "success": true,
  "data": {
    "order_id": 329,
    "iccid": "8985231024829102481",
    "qr_code_url": "https://anindasms.com/uploads/esim/qr_329.png",
    "lpa_string": "LPA:1$rsp.example.com$MATCHING_ID_329",
    "status": "active"
  }
}
`

#### C. eSIM Orders & Topup (GET /esim/orders/{id}, POST /esim/orders/{id}/topup)
`ash
curl -X GET "https://anindasms.com/api/v1/esim/orders/329" -H "Authorization: Bearer YOUR_API_KEY"
curl -X POST "https://anindasms.com/api/v1/esim/orders/329/topup" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"package_id":"pkg_tr_1gb_7d"}'
`

---

### 12. Temporary Mail API Endpoints
Disposable inboxes with instant message retrieval.

#### A. Mail Domains & Create Mailbox (GET /mail/domains, POST /mail/create)
`ash
curl -X GET "https://anindasms.com/api/v1/mail/domains" -H "Authorization: Bearer YOUR_API_KEY"
curl -X POST "https://anindasms.com/api/v1/mail/create" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"duration":60}'
`
**Sample Response:**
`json
{
  "success": true,
  "data": {
    "order_id": 892,
    "email": "box_892@fastinbox.net",
    "cost": 5.0,
    "expires_at": "2026-10-24 15:00:00"
  }
}
`

#### B. Read Inbox Messages (GET /mail/orders/{id}/messages)
`ash
curl -X GET "https://anindasms.com/api/v1/mail/orders/892/messages" -H "Authorization: Bearer YOUR_API_KEY"
`

#### C. Extend or Cancel Mailbox (POST /mail/orders/{id}/extend, POST /mail/orders/{id}/cancel)
`ash
curl -X POST "https://anindasms.com/api/v1/mail/orders/892/extend" -H "Authorization: Bearer YOUR_API_KEY" -d '{"duration":60}'
curl -X POST "https://anindasms.com/api/v1/mail/orders/892/cancel" -H "Authorization: Bearer YOUR_API_KEY"
`


---

## 🔄 Order Lifecycle & Statuses

| Status | Description |
| :--- | :--- |
| `pending` | Number allocated. Waiting for the incoming SMS verification code. |
| `received` | SMS delivered. The code is available in `sms_code` and full text in `full_sms`. |
| `cancelled` | Order cancelled by the client. Full amount is refunded immediately. |
| `timeout` | 10-minute expiration window elapsed without SMS arrival. Full amount automatically refunded. |
| `completed` | Order successfully concluded by the client. |

---

## ⚠️ Error Codes Reference

| HTTP Code | Error Reason | Description |
| :--- | :--- | :--- |
| `400 Bad Request` | Missing / Invalid Parameters | Required service or country parameters are missing or invalid. |
| `401 Unauthorized` | Invalid Authentication | API Key is missing, invalid, or sent through unsupported transport. |
| `402 Payment Required`| Insufficient Funds | Account balance is lower than the order unit price. |
| `403 Forbidden` | Access Prohibited | Account suspended or request IP is not in your configured allowlist. |
| `404 Not Found` | Order Not Found | The requested order ID does not exist or does not belong to your account. |
| `409 Conflict` | Out of Stock / Invalid State | No numbers available for this country/service or order already finished. |
| `429 Too Many Requests`| Rate Limit Exceeded | Exceeded requests/min quota. Wait for the `Retry-After` duration. |
| `503 Service Unavailable`| Maintenance Mode | Platform is temporarily in scheduled maintenance. |

---

## 🛠 Official SDK Clients

All official SDKs include built-in automatic retry handling and a `waitForCode` helper method that polls for the SMS code for you.

### 🐘 PHP SDK
Installation:
```bash
composer config repositories.anindasms vcs https://github.com/anindasms-dev/anindasms-sdk
composer require anindasms/sdk:dev-main
```
Usage Example:
```php
<?php
require 'vendor/autoload.php';

use AnindaSms\Client;

$client = new Client('YOUR_API_KEY');

// 1. Check account balance
$balance = $client->balance();
echo "Balance: " . $balance['formatted'] . PHP_EOL;

// 2. Buy a virtual number (WhatsApp, Indonesia)
$order = $client->createOrder('wa', '62');
echo "Allocated Phone: " . $order['phone'] . " (#" . $order['order_id'] . ")" . PHP_EOL;

// 3. Wait for the SMS code automatically
$sms = $client->waitForCode($order['order_id']);
echo "Received Verification Code: " . $sms['code'] . PHP_EOL;

// 4. Mark order as done
$client->completeOrder($order['order_id']);
```

---

### 🐍 Python SDK
Installation:
```bash
pip install "git+https://github.com/anindasms-dev/anindasms-sdk.git#subdirectory=python"
```
Usage Example:
```python
from anindasms import Client

client = Client("YOUR_API_KEY")

# Check balance
print("Balance:", client.balance()["formatted"])

# Buy number
order = client.create_order("wa", "62")
print("Allocated Phone:", order["phone"])

# Wait for SMS code
sms = client.wait_for_code(order["order_id"])
print("Received Code:", sms["code"])

# Finish
client.complete_order(order["order_id"])
```

---

### 🟩 Node.js SDK
Installation:
```bash
npm install github:anindasms-dev/anindasms-sdk
```
Usage Example:
```javascript
const { AnindaSMS } = require("anindasms");

const client = new AnindaSMS("YOUR_API_KEY");

async function main() {
  const balance = await client.balance();
  console.log("Balance:", balance.formatted);

  const order = await client.createOrder("wa", "62");
  console.log("Allocated Phone:", order.phone);

  const sms = await client.waitForCode(order.order_id);
  console.log("Received Code:", sms.code);

  await client.completeOrder(order.order_id);
}

main().catch(console.error);
```

---

### 🔷 C# / .NET SDK
Add `csharp/AnindaSms.cs` into your .NET 6+ project:
```csharp
using System;
using AnindaSms;

var client = new AnindaSmsClient("YOUR_API_KEY");

var balance = await client.GetBalanceAsync();
Console.WriteLine($"Balance: {balance.GetProperty("formatted").GetString()}");

var order = await client.CreateOrderAsync("wa", "62");
Console.WriteLine($"Allocated Phone: {order.GetProperty("phone").GetString()}");

var sms = await client.WaitForCodeAsync(order.GetProperty("order_id").GetInt32());
Console.WriteLine($"Received Code: {sms.GetProperty("code").GetString()}");

await client.DoneOrderAsync(order.GetProperty("order_id").GetInt32());
```

---

### 🐹 Go SDK
Installation:
```bash
go get github.com/anindasms-dev/anindasms-sdk/go
```
Usage Example:
```go
package main

import (
	"fmt"
	anindasms "github.com/anindasms-dev/anindasms-sdk/go"
)

func main() {
	client := anindasms.NewClient("YOUR_API_KEY")

	bal, err := client.Balance()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Balance: %s\n", bal.Formatted)

	order, err := client.CreateOrder("wa", "62", "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Allocated Phone: %s (#%d)\n", order.Phone, order.OrderID)

	sms, err := client.WaitForCode(order.OrderID, 0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received Code: %s\n", *sms.Code)

	_ = client.DoneOrder(order.OrderID)
}
```

---

## 📄 License & Support

This SDK and documentation are open-source and released under the [MIT License](LICENSE).

* **Official Website:** [anindasms.com](https://anindasms.com)
* **Web Documentation:** [anindasms.com/api](https://anindasms.com/api)
* **Customer Support:** support@anindasms.com / destek@anindasms.com
