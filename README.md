# AnindaSMS — Official REST API v1 & Multi-Language SDKs

[![API Version](https://img.shields.io/badge/API-v1-blue.svg)](https://anindasms.com/api)
[![Countries](https://img.shields.io/badge/Countries-200%2B-orange.svg)](https://anindasms.com/api)
[![Services](https://img.shields.io/badge/Services-5000%2B-red.svg)](https://anindasms.com/api)
[![Features](https://img.shields.io/badge/Features-SMS%20%7C%20Rent%20%7C%20Proxy%20%7C%20eSIM%20%7C%20Mail-success.svg)](https://anindasms.com/api)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PHP](https://img.shields.io/badge/PHP-7.4%2B%20%7C%208.x-8892BF.svg)](php/)
[![Python](https://img.shields.io/badge/Python-3.8%2B-3776AB.svg)](python/)
[![Node.js](https://img.shields.io/badge/Node.js-18%2B-339933.svg)](node/)
[![.NET](https://img.shields.io/badge/.NET-6.0%2B-512BD4.svg)](csharp/)
[![Go](https://img.shields.io/badge/Go-1.18%2B-00ADD8.svg)](go/)

> 🌍 **Language / Dil:** **English** | [🇹🇷 Türkçe Dokümantasyon için tıklayın (README.tr.md)](README.tr.md)

Welcome to the official repository for the [AnindaSMS](https://anindasms.com) **REST API v1** and open-source multi-language SDK libraries. AnindaSMS provides an all-in-one developer platform for:

* 📱 **SMS Verification:** Temporary, disposable virtual numbers for SMS verification across **5000+ services** and **200+ countries**.
* ⏳ **Number Rental (Rent API):** Dedicated numbers rented from **4 hours up to 30 days** with support for receiving multiple incoming SMS messages.
* 🛡️ **High-Speed Proxies:** Static Datacenter, Residential (ISP), IPv6, and Rotating 4G/5G Mobile proxies worldwide.
* 📶 **Global eSIMs:** Instant digital travel data packages across **150+ countries** with instant QR & LPA provisioning.
* 📬 **Temporary Mail:** Disposable, programmable email inboxes with real-time webhook & polling message retrieval.

---

## 📑 Table of Contents

1. [Overview](#-overview)
   - [Base URL](#base-url)
   - [Authentication](#authentication)
   - [Standard Response Envelope](#standard-response-envelope)
   - [Rate Limiting](#rate-limiting)
2. [Account Endpoints](#-account-endpoints)
   - [Profile Overview (`GET /me`)](#1-user-profile-get-me)
   - [Balance Query (`GET /balance`)](#2-balance-query-get-balance)
3. [SMS Verification API](#-sms-verification-api)
   - [Service Catalog (`GET /services`)](#3-service-catalog-get-services)
   - [Prices & Live Stock (`GET /prices`)](#4-prices--live-stock-get-prices)
   - [Buy Number (`POST /order`)](#5-create-order--buy-number-post-order)
   - [Check Status & OTP (`GET /order/{id}`)](#6-check-order--sms-code-get-orderid)
   - [Cancel Order (`POST /order/{id}/cancel`)](#7-cancel-order-post-orderidcancel)
   - [Complete Order (`POST /order/{id}/done`)](#8-complete-order-post-orderiddone)
4. [Number Rental (Rent) API](#-number-rental-rent-api)
   - [Rental Durations (`GET /rent/durations`)](#rental-durations-get-rentdurations)
   - [Available Countries (`GET /rent/countries`)](#available-countries-get-rentcountries)
   - [Services & Quotes (`GET /rent/services`, `GET /rent/prices`)](#services--quotes-get-rentservices-get-rentprices)
   - [Rent Number (`POST /rent/buy`)](#rent-number-post-rentbuy)
   - [Rental Orders & SMS List (`GET /rent/orders`, `GET /rent/orders/{id}`)](#rental-orders--sms-list-get-rentorders-get-rentordersid)
   - [Cancel or Reactivate (`POST /rent/orders/{id}/cancel`, `POST /rent/orders/{id}/reactivate`)](#cancel-or-reactivate-post-rentordersidcancel-post-rentordersidreactivate)
5. [Proxy API](#-proxy-api)
   - [Proxy Types & Countries (`GET /proxy/types`, `GET /proxy/countries`)](#proxy-types--countries-get-proxytypes-get-proxycountries)
   - [Price Quote (`GET /proxy/price`)](#price-quote-get-proxyprice)
   - [Buy Proxy (`POST /proxy/buy`)](#buy-proxy-post-proxybuy)
   - [Proxy Orders (`GET /proxy/orders`, `GET /proxy/orders/{id}`)](#proxy-orders-get-proxyorders-get-proxyordersid)
6. [eSIM API](#-esim-api)
   - [eSIM Countries & Packages (`GET /esim/countries`, `GET /esim/packages`)](#esim-countries--packages-get-esimcountries-get-esimpackages)
   - [Buy eSIM Package (`POST /esim/buy`)](#buy-esim-package-post-esimbuy)
   - [eSIM Orders & Topup (`GET /esim/orders/{id}`, `POST /esim/orders/{id}/topup`)](#esim-orders--topup-get-esimordersid-post-esimordersidtopup)
7. [Temporary Mail API](#-temporary-mail-api)
   - [Domains & Create Mailbox (`GET /mail/domains`, `POST /mail/create`)](#domains--create-mailbox-get-maildomains-post-mailcreate)
   - [Read Inbox Messages (`GET /mail/orders/{id}/messages`)](#read-inbox-messages-get-mailordersidmessages)
   - [Extend or Cancel (`POST /mail/orders/{id}/extend`, `POST /mail/orders/{id}/cancel`)](#extend-or-cancel-post-mailordersidextend-post-mailordersidcancel)
8. [Order Lifecycle & Statuses](#-order-lifecycle--statuses)
9. [Error Codes Reference](#-error-codes-reference)
10. [Official SDK Clients](#-official-sdk-clients)
    - [PHP](#-php-sdk)
    - [Python](#-python-sdk)
    - [Node.js](#-nodejs-sdk)
    - [C# / .NET](#-c--net-sdk)
    - [Go](#-go-sdk)
11. [License & Support](#-license--support)

---

## 🌐 Overview

### Base URL
All API requests must be made over HTTPS to the following endpoint:
```
https://anindasms.com/api/v1
```

### Authentication
Obtain your API Key from your dashboard at [anindasms.com/profile](https://anindasms.com/profile).

Authenticate every request by passing your API Key in the HTTP headers:
```http
Authorization: Bearer YOUR_API_KEY
```
or
```http
X-API-Key: YOUR_API_KEY
```

> 🔒 **Security Notice:** Do not expose your API key in URL query strings or public repositories.

### Standard Response Envelope
All responses return UTF-8 encoded JSON:

**Success (`HTTP 200 / 201`):**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error (`HTTP 4xx / 5xx`):**
```json
{
  "success": false,
  "error": "Detailed error message.",
  "data": null
}
```

### Rate Limiting
* **General Limit:** 120 requests / minute
* **Order Creation:** 30 requests / minute
* Headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`

---

## 👤 Account Endpoints

### 1. User Profile (`GET /me`)
```bash
curl -X GET "https://anindasms.com/api/v1/me" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### 2. Balance Query (`GET /balance`)
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

## 📱 SMS Verification API

### 3. Service Catalog (`GET /services`)
```bash
curl -X GET "https://anindasms.com/api/v1/services" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### 4. Prices & Live Stock (`GET /prices`)
```bash
curl -X GET "https://anindasms.com/api/v1/prices?service=wa&country=62" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### 5. Create Order / Buy Number (`POST /order`)
```bash
curl -X POST "https://anindasms.com/api/v1/order" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"service":"wa","country":"62"}'
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
    "cost": 12.50,
    "status": "waiting",
    "expires_at": "2026-10-24 14:12:10"
  }
}
```

### 6. Check Order & SMS Code (`GET /order/{id}`)
Poll every 3–5 seconds until the code arrives:
```bash
curl -X GET "https://anindasms.com/api/v1/order/58231" \
  -H "Authorization: Bearer YOUR_API_KEY"
```
**Sample Response:**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "status": "received",
    "code": "489201",
    "sms_text": "Your WhatsApp verification code is: 489-201"
  }
}
```

### 7. Cancel Order (`POST /order/{id}/cancel`)
Cancels order before SMS arrives and refunds balance immediately:
```bash
curl -X POST "https://anindasms.com/api/v1/order/58231/cancel" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### 8. Complete Order (`POST /order/{id}/done`)
```bash
curl -X POST "https://anindasms.com/api/v1/order/58231/done" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

---

## ⏳ Number Rental (Rent) API

Dedicated numbers for extended periods with unlimited incoming SMS reception.

### Rental Durations (`GET /rent/durations`)
```bash
curl -X GET "https://anindasms.com/api/v1/rent/durations" -H "Authorization: Bearer YOUR_API_KEY"
```

### Available Countries (`GET /rent/countries`)
```bash
curl -X GET "https://anindasms.com/api/v1/rent/countries?duration=24" -H "Authorization: Bearer YOUR_API_KEY"
```

### Services & Quotes (`GET /rent/services`, `GET /rent/prices`)
```bash
curl -X GET "https://anindasms.com/api/v1/rent/services?country=62&duration=24" -H "Authorization: Bearer YOUR_API_KEY"
curl -X GET "https://anindasms.com/api/v1/rent/prices?service=wa&country=62&duration=24" -H "Authorization: Bearer YOUR_API_KEY"
```

### Rent Number (`POST /rent/buy`)
```bash
curl -X POST "https://anindasms.com/api/v1/rent/buy" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"service":"wa","country":"62","duration":24}'
```
**Sample Response:**
```json
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
```

### Rental Orders & SMS List (`GET /rent/orders`, `GET /rent/orders/{id}`)
```bash
curl -X GET "https://anindasms.com/api/v1/rent/orders/74102" -H "Authorization: Bearer YOUR_API_KEY"
```
**Sample Response:**
```json
{
  "success": true,
  "data": {
    "order_id": 74102,
    "phone": "+905559876543",
    "status": "received",
    "remaining": 82400,
    "sms_list": [
      {
        "code": "948210",
        "text": "Your WhatsApp verification code: 948-210",
        "date": "2026-10-24 14:12:00"
      }
    ]
  }
}
```

### Cancel or Reactivate (`POST /rent/orders/{id}/cancel`, `POST /rent/orders/{id}/reactivate`)
```bash
# Cancel within first 20 minutes if no SMS received
curl -X POST "https://anindasms.com/api/v1/rent/orders/74102/cancel" -H "Authorization: Bearer YOUR_API_KEY"

# Extend same number for another 24 hours
curl -X POST "https://anindasms.com/api/v1/rent/orders/74102/reactivate" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"duration":24}'
```

---

## 🛡️ Proxy API

Static Datacenter, Residential, IPv6, and Rotating 4G Mobile proxies.

### Proxy Types & Countries (`GET /proxy/types`, `GET /proxy/countries`)
```bash
curl -X GET "https://anindasms.com/api/v1/proxy/types" -H "Authorization: Bearer YOUR_API_KEY"
curl -X GET "https://anindasms.com/api/v1/proxy/countries?type=datacenter" -H "Authorization: Bearer YOUR_API_KEY"
```

### Price Quote (`GET /proxy/price`)
```bash
curl -X GET "https://anindasms.com/api/v1/proxy/price?type=datacenter&country=TR&period=30&quantity=1" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### Buy Proxy (`POST /proxy/buy`)
```bash
curl -X POST "https://anindasms.com/api/v1/proxy/buy" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"type":"datacenter","country":"TR","period":30,"quantity":1}'
```
**Sample Response:**
```json
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
```

### Proxy Orders (`GET /proxy/orders`, `GET /proxy/orders/{id}`)
```bash
curl -X GET "https://anindasms.com/api/v1/proxy/orders" -H "Authorization: Bearer YOUR_API_KEY"
curl -X GET "https://anindasms.com/api/v1/proxy/orders/1042" -H "Authorization: Bearer YOUR_API_KEY"
```

---

## 📶 eSIM API

High-speed digital travel data for 150+ countries.

### eSIM Countries & Packages (`GET /esim/countries`, `GET /esim/packages`)
```bash
curl -X GET "https://anindasms.com/api/v1/esim/countries" -H "Authorization: Bearer YOUR_API_KEY"
curl -X GET "https://anindasms.com/api/v1/esim/packages?country=TR" -H "Authorization: Bearer YOUR_API_KEY"
```

### Buy eSIM Package (`POST /esim/buy`)
```bash
curl -X POST "https://anindasms.com/api/v1/esim/buy" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"package_id":"pkg_tr_1gb_7d","quantity":1}'
```
**Sample Response:**
```json
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
```

### eSIM Orders & Topup (`GET /esim/orders/{id}`, `POST /esim/orders/{id}/topup`)
```bash
curl -X GET "https://anindasms.com/api/v1/esim/orders/329" -H "Authorization: Bearer YOUR_API_KEY"
curl -X POST "https://anindasms.com/api/v1/esim/orders/329/topup" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"package_id":"pkg_tr_1gb_7d"}'
```

---

## 📬 Temporary Mail API

Disposable inboxes with instant message reading.

### Domains & Create Mailbox (`GET /mail/domains`, `POST /mail/create`)
```bash
curl -X GET "https://anindasms.com/api/v1/mail/domains" -H "Authorization: Bearer YOUR_API_KEY"
curl -X POST "https://anindasms.com/api/v1/mail/create" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"duration":60}'
```
**Sample Response:**
```json
{
  "success": true,
  "data": {
    "order_id": 892,
    "email": "box_892@fastinbox.net",
    "cost": 5.0,
    "expires_at": "2026-10-24 15:00:00"
  }
}
```

### Read Inbox Messages (`GET /mail/orders/{id}/messages`)
```bash
curl -X GET "https://anindasms.com/api/v1/mail/orders/892/messages" -H "Authorization: Bearer YOUR_API_KEY"
```

### Extend or Cancel (`POST /mail/orders/{id}/extend`, `POST /mail/orders/{id}/cancel`)
```bash
curl -X POST "https://anindasms.com/api/v1/mail/orders/892/extend" -H "Authorization: Bearer YOUR_API_KEY" -d '{"duration":60}'
curl -X POST "https://anindasms.com/api/v1/mail/orders/892/cancel" -H "Authorization: Bearer YOUR_API_KEY"
```

---

## 🔄 Order Lifecycle & Statuses

| Status | Description |
| :--- | :--- |
| `waiting` / `pending` | Number allocated. Waiting for incoming SMS. |
| `received` | SMS delivered. Verification code available in `code` or `sms_list`. |
| `done` / `completed` | Order successfully completed. |
| `cancelled` | Cancelled by user. Cost refunded immediately. |
| `expired` | Time expired without SMS. Cost refunded automatically. |

---

## ⚠️ Error Codes Reference

| HTTP Code | Error Reason | Description |
| :--- | :--- | :--- |
| `400 Bad Request` | Missing / Invalid Parameters | Required parameters missing or invalid format. |
| `401 Unauthorized` | Invalid Authentication | API Key missing, invalid, or unauthorized. |
| `402 Payment Required`| Insufficient Funds | Account balance lower than the item price. |
| `403 Forbidden` | Access Prohibited | Service disabled, country closed, or IP blocked. |
| `404 Not Found` | Resource Not Found | Order ID or catalog item not found. |
| `409 Conflict` | Out of Stock | No inventory available for requested service/country. |
| `429 Too Many Requests`| Rate Limit Exceeded | Exceeded quota. Check `Retry-After` header. |
| `500 Server Error` | Temporary Provider Error | Upstream provider error. Retry request. |

---

## 🛠 Official SDK Clients

Official SDK clients provide typed interfaces, built-in retry handling, and automatic polling helpers for all 5 services.

### 🐘 PHP SDK

**Install:**
```bash
composer config repositories.anindasms vcs https://github.com/anindasms-dev/anindasms-sdk
composer require anindasms/sdk:dev-main
```

**Example:**
```php
<?php
require 'vendor/autoload.php';

use AnindaSms\Client;

$client = new Client('YOUR_API_KEY');

// 1. SMS Verification
$order = $client->createOrder('wa', '62');
echo "Phone: " . $order['phone'] . PHP_EOL;
$sms = $client->waitForCode($order['order_id']);
echo "SMS Code: " . $sms['code'] . PHP_EOL;

// 2. Number Rental
$rent = $client->rentBuy('wa', '62', 24);
echo "Rented Phone: " . $rent['phone'] . PHP_EOL;

// 3. Buy Proxy
$proxy = $client->proxyBuy('datacenter', 'TR', 30, 1);
echo "Proxy IP: " . $proxy['ip'] . ":" . $proxy['port'] . PHP_EOL;

// 4. Buy eSIM
$esim = $client->esimBuy('pkg_tr_1gb_7d');
echo "eSIM QR: " . $esim['qr_code_url'] . PHP_EOL;

// 5. Temporary Mail
$mail = $client->mailCreate(60);
echo "Temp Email: " . $mail['email'] . PHP_EOL;
```

---

### 🐍 Python SDK

**Install:**
```bash
pip install "git+https://github.com/anindasms-dev/anindasms-sdk.git#subdirectory=python"
```

**Example:**
```python
from anindasms import Client

client = Client("YOUR_API_KEY")

# 1. SMS Verification
order = client.create_order("wa", "62")
print("Phone:", order["phone"])
sms = client.wait_for_code(order["order_id"])
print("SMS Code:", sms["code"])

# 2. Number Rental
rent = client.rent_buy("wa", "62", 24)
print("Rented Phone:", rent["phone"])

# 3. Buy Proxy
proxy = client.proxy_buy("datacenter", "TR", 30, 1)
print(f"Proxy: {proxy['ip']}:{proxy['port']}")

# 4. Buy eSIM
esim = client.esim_buy("pkg_tr_1gb_7d")
print("eSIM QR:", esim["qr_code_url"])

# 5. Temporary Mail
mail = client.mail_create(60)
print("Temp Mail:", mail["email"])
```

---

### 🟩 Node.js SDK

**Install:**
```bash
npm install github:anindasms-dev/anindasms-sdk
```

**Example:**
```javascript
const { AnindaSMS } = require("anindasms");

const client = new AnindaSMS("YOUR_API_KEY");

async function main() {
  // 1. SMS Verification
  const order = await client.createOrder("wa", "62");
  console.log("Phone:", order.phone);
  const sms = await client.waitForCode(order.order_id);
  console.log("SMS Code:", sms.code);

  // 2. Number Rental
  const rent = await client.rentBuy("wa", "62", 24);
  console.log("Rented Phone:", rent.phone);

  // 3. Buy Proxy
  const proxy = await client.proxyBuy("datacenter", "TR", 30, 1);
  console.log(`Proxy: ${proxy.ip}:${proxy.port}`);

  // 4. Buy eSIM
  const esim = await client.esimBuy("pkg_tr_1gb_7d");
  console.log("eSIM QR:", esim.qr_code_url);

  // 5. Temporary Mail
  const mail = await client.mailCreate(60);
  console.log("Temp Mail:", mail.email);
}

main().catch(console.error);
```

---

### 🔷 C# / .NET SDK

Add `csharp/AnindaSms.csproj` or reference `AnindaSms.cs` in your .NET project:

```csharp
using System;
using AnindaSms;

var client = new AnindaSmsClient("YOUR_API_KEY");

// 1. SMS Verification
var order = await client.CreateOrderAsync("wa", "62");
Console.WriteLine($"Phone: {order.GetProperty("phone").GetString()}");
var sms = await client.WaitForCodeAsync(order.GetProperty("order_id").GetInt32());
Console.WriteLine($"Code: {sms.GetProperty("code").GetString()}");

// 2. Number Rental
var rent = await client.RentBuyAsync("wa", "62", 24);
Console.WriteLine($"Rented Phone: {rent.GetProperty("phone").GetString()}");

// 3. Buy Proxy
var proxy = await client.ProxyBuyAsync("datacenter", "TR", 30, 1);
Console.WriteLine($"Proxy: {proxy.GetProperty("ip").GetString()}:{proxy.GetProperty("port").GetInt32()}");

// 4. Buy eSIM
var esim = await client.EsimBuyAsync("pkg_tr_1gb_7d");
Console.WriteLine($"eSIM QR: {esim.GetProperty("qr_code_url").GetString()}");

// 5. Temporary Mail
var mail = await client.MailCreateAsync(60);
Console.WriteLine($"Temp Mail: {mail.GetProperty("email").GetString()}");
```

---

### 🐹 Go SDK

**Install:**
```bash
go get github.com/anindasms-dev/anindasms-sdk/go
```

**Example:**
```go
package main

import (
	"fmt"
	anindasms "github.com/anindasms-dev/anindasms-sdk/go"
)

func main() {
	client := anindasms.NewClient("YOUR_API_KEY")

	// 1. SMS Verification
	order, err := client.CreateOrder("wa", "62", "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Phone: %s (#%d)\n", order.Phone, order.OrderID)

	sms, err := client.WaitForCode(order.OrderID, 0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Code: %s\n", *sms.Code)

	// 2. Number Rental
	rent, _ := client.RentBuy("wa", "62", 24)
	fmt.Printf("Rent Response: %s\n", string(rent))

	// 3. Proxy
	proxy, _ := client.ProxyBuy("datacenter", "TR", 30, 1, 0)
	fmt.Printf("Proxy Response: %s\n", string(proxy))

	// 4. eSIM
	esim, _ := client.EsimBuy("pkg_tr_1gb_7d", 1)
	fmt.Printf("eSIM Response: %s\n", string(esim))

	// 5. Temp Mail
	mail, _ := client.MailCreate(60, "", "")
	fmt.Printf("Mail Response: %s\n", string(mail))
}
```

---

## 📄 License & Support

This SDK and documentation are open-source and released under the [MIT License](LICENSE).

* **Official Website:** [anindasms.com](https://anindasms.com)
* **Web Documentation:** [anindasms.com/api](https://anindasms.com/api)
* **Customer Support:** support@anindasms.com / destek@anindasms.com
