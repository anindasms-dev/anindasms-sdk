# AnındaSMS — Resmi REST API v1 & Çoklu Dil SDK Kütüphaneleri

[![API Sürümü](https://img.shields.io/badge/API-v1-blue.svg)](https://anindasms.com/api)
[![Ülke](https://img.shields.io/badge/%C3%9Clke-200%2B-orange.svg)](https://anindasms.com/api)
[![Servis](https://img.shields.io/badge/Servis-5000%2B-red.svg)](https://anindasms.com/api)
[![Lisans: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PHP](https://img.shields.io/badge/PHP-7.4%2B%20%7C%208.x-8892BF.svg)](php/)
[![Python](https://img.shields.io/badge/Python-3.8%2B-3776AB.svg)](python/)
[![Node.js](https://img.shields.io/badge/Node.js-18%2B-339933.svg)](node/)
[![.NET](https://img.shields.io/badge/.NET-6.0%2B-512BD4.svg)](csharp/)
[![Go](https://img.shields.io/badge/Go-1.18%2B-00ADD8.svg)](go/)

> 🌍 **Dil / Language:** [🇬🇧 Switch to English Documentation (README.md)](README.md) | **Türkçe**

Bu depo, [AnındaSMS](https://anindasms.com) platformunun **REST API v1** uç noktalarının detaylı dokümantasyonunu ve 5 farklı dilde hazırlanmış resmi açık kaynak istemci kütüphanelerini (**SDK**) içerir. AnındaSMS, **200'den fazla ülke** ve **5000'den fazla online servis** için tek kullanımlık sanal SMS doğrulama numaraları sağlar.

---

## 📑 İçindekiler
1. [API Genel Bilgileri](#-api-genel-bilgileri)
   - [Temel URL (Base URL)](#temel-url-base-url)
   - [Kimlik Doğrulama (Authentication)](#kimlik-doğrulama-authentication)
   - [Standart Yanıt Formatı](#standart-yanıt-formatı)
   - [Rate Limit (Hız Sınırları)](#rate-limit-hız-sınırları)
2. [REST API Uç Noktaları (Endpoints)](#-rest-api-uç-noktaları)
   - [1. Hesap Bilgileri (`GET /me`)](#1-hesap-bilgileri-get-me)
   - [2. Bakiye Sorgulama (`GET /balance`)](#2-bakiye-sorgulama-get-balance)
   - [3. Aktif Servis Listesi (`GET /services`)](#3-aktif-servis-listesi-get-services)
   - [4. Fiyat ve Stok Sorgulama (`GET /prices`)](#4-fiyat-ve-stok-sorgulama-get-prices)
   - [5. Numara Siparişi Oluşturma (`POST /order`)](#5-numara-siparişi-oluşturma-post-order)
   - [6. Sipariş ve SMS Kodu Sorgulama (`GET /order/{id}`)](#6-sipariş-ve-sms-kodu-sorgulama-get-orderid)
   - [7. Sipariş İptali (`POST /order/{id}/cancel`)](#7-sipariş-iptali-post-orderidcancel)
   - [8. Siparişi Tamamlama (`POST /order/{id}/done`)](#8-siparişi-tamamlama-post-orderiddone)
3. [Sipariş Yaşam Döngüsü & Durumlar](#-sipariş-yaşam-döngüsü--durumlar)
4. [Hata Kodları](#-hata-kodları)
5. [Resmi SDK Kütüphaneleri](#-resmi-sdk-kütüphaneleri)
   - [PHP](#-php-sdk)
   - [Python](#-python-sdk)
   - [Node.js](#-nodejs-sdk)
   - [C# / .NET](#-c--net-sdk)
   - [Go](#-go-sdk)
6. [Lisans & Destek](#-lisans--destek)

---

## 🌐 API Genel Bilgileri

### Temel URL (Base URL)
Tüm istekler HTTPS üzerinden aşağıdaki temel adrese gönderilmelidir:
```
https://anindasms.com/api/v1
```

### Kimlik Doğrulama (Authentication)
API anahtarınızı [anindasms.com/profile](https://anindasms.com/profile) sayfasından oluşturabilirsiniz.

Anahtarınızı her istekte HTTP başlığı (header) olarak iletmeniz gerekir. İki yöntem desteklenir:
```http
Authorization: Bearer SIZIN_API_KEYINIZ
```
veya
```http
X-API-Key: SIZIN_API_KEYINIZ
```
> 🔒 **Güvenlik Notu:** Güvenlik gerekçesiyle `api_key` değerini URL parametresinde (`?api_key=...`) veya istek gövdesinde göndermeyin. Bu tür istekler `401 Unauthorized` ile reddedilir.

### Standart Yanıt Formatı
Tüm yanıtlar UTF-8 kodlamasında JSON olarak döner:

**Başarılı Yanıt:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Hatalı Yanıt:**
```json
{
  "success": false,
  "error": "Açıklayıcı hata mesajı.",
  "data": null
}
```

### Rate Limit (Hız Sınırları)
Sistem kararlılığı için istekler dakika bazında sınırlandırılmıştır:
* **Genel İstekler:** Dakikada maksimum **120 istek**
* **Sipariş Oluşturma (`POST /order`):** Dakikada maksimum **30 istek**

İstek yanıtlarında hız limitinizi takip edebileceğiniz başlıklar yer alır:
* `X-RateLimit-Limit`: İzin verilen dakikalık toplam istek sayısı.
* `X-RateLimit-Remaining`: Mevcut dakikada kalan istek hakkınız.
* `X-RateLimit-Reset`: Sayacın sıfırlanacağı Unix zaman damgası.

---

## 📡 REST API Uç Noktaları

### 1. Hesap Bilgileri (`GET /me`)
Hesap ID'si, kullanıcı adı, kayıtlı e-posta, anlık bakiye ve API anahtarı meta verilerini döndürür.

```bash
curl -X GET "https://anindasms.com/api/v1/me" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

**Örnek Yanıt:**
```json
{
  "success": true,
  "data": {
    "id": 2,
    "username": "kullaniciadi",
    "email": "kullanici@example.com",
    "balance": 150.00,
    "currency": "TRY",
    "api_key_prefix": "abab247a",
    "api_key_created_at": "2026-10-24 12:00:00",
    "api_key_last_used_at": "2026-10-24 15:30:10"
  }
}
```

---

### 2. Bakiye Sorgulama (`GET /balance`)
Hesaptaki mevcut kullanılabilir bakiyeyi döner.

```bash
curl -X GET "https://anindasms.com/api/v1/balance" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

**Örnek Yanıt:**
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

### 3. Aktif Servis Listesi (`GET /services`)
Sistemde aktif olarak SMS onayı alınabilen tüm servislerin kodlarını ve adlarını döndürür.

```bash
curl -X GET "https://anindasms.com/api/v1/services" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

**Örnek Yanıt:**
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

### 4. Fiyat ve Stok Sorgulama (`GET /prices`)
Belirli bir servis için veya tüm servisler için ülke bazlı anlık fiyat ve stok durumunu listeler.

**Parametreler (Query String):**
* `service` (string, zorunlu): Servis kodu (örn: `wa`, `tg`)
* `country` (string, opsiyonel): Ülke kodu (örn: `62`, `tr`)

```bash
curl -X GET "https://anindasms.com/api/v1/prices?service=wa&country=62" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

**Örnek Yanıt:**
```json
{
  "success": true,
  "data": {
    "service": "wa",
    "prices": [
      {
        "country": "62",
        "country_name": "Endonezya",
        "price": 12.50,
        "count": 482
      }
    ]
  }
}
```

---

### 5. Numara Siparişi Oluşturma (`POST /order`)
Belirtilen servis ve ülke için yeni bir sanal numara tahsis eder. Bakiye hesabınızdan rezerve edilir.

**Gövde (JSON):**
* `service` (string, zorunlu): Servis kodu (örn: `"wa"`)
* `country` (string, zorunlu): Ülke kodu (örn: `"62"`)
* `operator` (string, opsiyonel): Belirli bir operatör tercihi

```bash
curl -X POST "https://anindasms.com/api/v1/order" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ" \
  -H "Content-Type: application/json" \
  -d '{"service":"wa", "country":"62"}'
```

**Örnek Yanıt:**
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

### 6. Sipariş ve SMS Kodu Sorgulama (`GET /order/{id}`)
Tahsis edilen numaranın durumunu ve gelen SMS onay kodunu sorgulamak için kullanılır. SMS kodu gelene kadar belirli aralıklarla (tavsiye: 3-5 saniyede bir) bu uç nokta sorgulanır.

```bash
curl -X GET "https://anindasms.com/api/v1/order/58231" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

**Beklemede Yanıtı (`pending`):**
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

**Kod Geldiğindeki Yanıt (`received`):**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "phone": "+628123456789",
    "status": "received",
    "sms_code": "489201",
    "full_sms": "WhatsApp kodunuz: 489-201. Kimseyle paylasmayin.",
    "remaining_seconds": 490
  }
}
```

---

### 7. Sipariş İptali (`POST /order/{id}/cancel`)
Henüz SMS gelmemiş olan bir siparişi iptal eder. Rezerve edilen ücret **anında ve kesintisiz** olarak hesap bakiyenize iade edilir.

```bash
curl -X POST "https://anindasms.com/api/v1/order/58231/cancel" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

**Örnek Yanıt:**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "status": "cancelled",
    "message": "Sipariş iptal edildi ve bakiye iade edildi."
  }
}
```

---

### 8. Siparişi Tamamlama (`POST /order/{id}/done`)
SMS kodu başarıyla alındıktan sonra işlemi sonlandırmak için kullanılır.

```bash
curl -X POST "https://anindasms.com/api/v1/order/58231/done" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

**Örnek Yanıt:**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "status": "completed",
    "message": "Sipariş tamamlandı."
  }
}
```

---

## 🔄 Sipariş Yaşam Döngüsü & Durumlar

| Durum (`status`) | Açıklama |
| :--- | :--- |
| `pending` | Numara ayrıldı, gelen SMS bekleniyor. |
| `received` | SMS başarıyla ulaştı. `sms_code` alanında doğrulama kodu mevcuttur. |
| `cancelled` | Müşteri veya sistem tarafından iptal edildi. Ücret tam iade edildi. |
| `timeout` | Zaman aşımı (varsayılan 10 dk) doldu. Ücret otomatik olarak hesaba iade edildi. |
| `completed` | İşlem müşteri tarafından başarıyla onaylandı ve kapatıldı. |

---

## ⚠️ Hata Kodları

| HTTP Kodu | Sebep | Açıklama |
| :--- | :--- | :--- |
| `400 Bad Request` | Eksik / Geçersiz Parametre | Servis veya ülke parametresi hatalı ya da eksik gönderildi. |
| `401 Unauthorized` | Kimlik Doğrulama Hatası | API anahtarı eksik, geçersiz veya formatı hatalı. |
| `402 Payment Required` | Yetersiz Bakiye | Hesabınızdaki bakiye sipariş tutarını karşılamıyor. |
| `403 Forbidden` | Yetkisiz Erişim | Hesap askıya alınmış veya IP izin listesi dışında bir IP'den istek yapıldı. |
| `404 Not Found` | Bulunamadı | Belirtilen sipariş ID'si hesaba ait değil veya bulunamadı. |
| `409 Conflict` | Stok Yok / Durum Çakışması | İlgili ülkede stok tükenmiş veya sipariş zaten kapatılmış. |
| `429 Too Many Requests`| Rate Limit Aşıldı | Dakikalık istek limiti aşıldı. `Retry-After` süresi kadar bekleyin. |
| `503 Service Unavailable`| Bakım Modu | Sistem genel bakım çalışmasında. |

---

## 🛠 Resmi SDK Kütüphaneleri

Tüm SDK'lar hazır hata yönetimi ve SMS kodunu sizin yerinize periyodik sorgulayan `waitForCode` metoduna sahiptir.

### 🐘 PHP SDK
Kurulum:
```bash
composer config repositories.anindasms vcs https://github.com/anindasms-dev/anindasms-sdk
composer require anindasms/sdk:dev-main
```
Örnek Kullanım:
```php
<?php
require 'vendor/autoload.php';

use AnindaSms\Client;

$client = new Client('SIZIN_API_KEYINIZ');

// Bakiye sorgula
$balance = $client->balance();
echo "Bakiye: " . $balance['formatted'] . PHP_EOL;

// Numara al (WhatsApp, Endonezya)
$order = $client->createOrder('wa', '62');
echo "Numara: " . $order['phone'] . " (Sipariş: #" . $order['order_id'] . ")" . PHP_EOL;

// SMS kodunu otomatik bekle
$sms = $client->waitForCode($order['order_id']);
echo "Gelen SMS Kodu: " . $sms['code'] . PHP_EOL;

// Tamamla
$client->completeOrder($order['order_id']);
```

---

### 🐍 Python SDK
Kurulum:
```bash
pip install "git+https://github.com/anindasms-dev/anindasms-sdk.git#subdirectory=python"
```
Örnek Kullanım:
```python
from anindasms import Client

client = Client("SIZIN_API_KEYINIZ")

# Bakiye sorgula
print("Bakiye:", client.balance()["formatted"])

# Numara al
order = client.create_order("wa", "62")
print("Numara:", order["phone"])

# SMS kodunu bekle
sms = client.wait_for_code(order["order_id"])
print("SMS Kodu:", sms["code"])

# Tamamla
client.complete_order(order["order_id"])
```

---

### 🟩 Node.js SDK
Kurulum:
```bash
npm install github:anindasms-dev/anindasms-sdk
```
Örnek Kullanım:
```javascript
const { AnindaSMS } = require("anindasms");

const client = new AnindaSMS("SIZIN_API_KEYINIZ");

async function main() {
  const balance = await client.balance();
  console.log("Bakiye:", balance.formatted);

  const order = await client.createOrder("wa", "62");
  console.log("Numara:", order.phone);

  const sms = await client.waitForCode(order.order_id);
  console.log("SMS Kodu:", sms.code);

  await client.completeOrder(order.order_id);
}

main().catch(console.error);
```

---

### 🔷 C# / .NET SDK
Projenize `csharp/AnindaSms.cs` dosyasını dahil edin (.NET 6+):
```csharp
using System;
using AnindaSms;

var client = new AnindaSmsClient("SIZIN_API_KEYINIZ");

var balance = await client.GetBalanceAsync();
Console.WriteLine($"Bakiye: {balance.GetProperty("formatted").GetString()}");

var order = await client.CreateOrderAsync("wa", "62");
Console.WriteLine($"Numara: {order.GetProperty("phone").GetString()}");

var sms = await client.WaitForCodeAsync(order.GetProperty("order_id").GetInt32());
Console.WriteLine($"SMS Kodu: {sms.GetProperty("code").GetString()}");

await client.DoneOrderAsync(order.GetProperty("order_id").GetInt32());
```

---

### 🐹 Go SDK
Kurulum:
```bash
go get github.com/anindasms-dev/anindasms-sdk/go
```
Örnek Kullanım:
```go
package main

import (
	"fmt"
	anindasms "github.com/anindasms-dev/anindasms-sdk/go"
)

func main() {
	client := anindasms.NewClient("SIZIN_API_KEYINIZ")

	bal, err := client.Balance()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bakiye: %s\n", bal.Formatted)

	order, err := client.CreateOrder("wa", "62", "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Numara: %s (ID: %d)\n", order.Phone, order.OrderID)

	sms, err := client.WaitForCode(order.OrderID, 0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("SMS Kodu: %s\n", *sms.Code)

	_ = client.DoneOrder(order.OrderID)
}
```

---

## 📄 Lisans & Destek

Bu kütüphane ve dokümantasyon [MIT Lisansı](LICENSE) altında açık kaynak olarak sunulmuştur.

* **Web Sitesi:** [anindasms.com](https://anindasms.com)
* **API Belgeleri Web Görünümü:** [anindasms.com/api](https://anindasms.com/api)
* **Destek & İletişim:** destek@anindasms.com / support@anindasms.com
