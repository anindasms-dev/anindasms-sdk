# AnındaSMS — Resmi REST API v1 & Çoklu Dil SDK Kütüphaneleri

[![API Sürümü](https://img.shields.io/badge/API-v1-blue.svg)](https://anindasms.com/api)
[![Ülke](https://img.shields.io/badge/%C3%9Clke-200%2B-orange.svg)](https://anindasms.com/api)
[![Servis](https://img.shields.io/badge/Servis-5000%2B-red.svg)](https://anindasms.com/api)
[![Özellikler](https://img.shields.io/badge/%C3%96zellikler-SMS%20%7C%20Kiralama%20%7C%20Proxy%20%7C%20eSIM%20%7C%20Mail-success.svg)](https://anindasms.com/api)
[![Lisans: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PHP](https://img.shields.io/badge/PHP-7.4%2B%20%7C%208.x-8892BF.svg)](php/)
[![Python](https://img.shields.io/badge/Python-3.8%2B-3776AB.svg)](python/)
[![Node.js](https://img.shields.io/badge/Node.js-18%2B-339933.svg)](node/)
[![.NET](https://img.shields.io/badge/.NET-6.0%2B-512BD4.svg)](csharp/)
[![Go](https://img.shields.io/badge/Go-1.18%2B-00ADD8.svg)](go/)

> 🌍 **Dil / Language:** **Türkçe** | [🇬🇧 Switch to English Documentation (README.md)](README.md)

Bu depo, [AnındaSMS](https://anindasms.com) platformunun **REST API v1** uç noktalarının detaylı dokümantasyonunu ve 5 farklı dilde hazırlanmış resmi açık kaynak istemci kütüphanelerini (**SDK**) içerir. AnındaSMS, geliştiriciler için tek bir çatı altında 5 temel servis sunar:

* 📱 **SMS Onay (Sanal Numara):** **200'den fazla ülke** ve **5000'den fazla online servis** için tek kullanımlık SMS onay numaraları.
* ⏳ **Numara Kiralama (Rent API):** **4 saatten 30 güne kadar**, aynı numaraya sınırsız sayıda SMS alabileceğiniz tahsisli numaralar.
* 🛡️ **Yüksek Hızlı Proxyler:** Datacenter, Residential (Ev/ISP), IPv6 ve Dönen 4G/5G Mobil proxy seçenekleri.
* 📶 **Global eSIM Paketleri:** **150+ ülke** için anında QR kod ve LPA formatında teslim edilen dijital seyahat internet paketleri.
* 📬 **Geçici Mail:** Webhook bildirimleri ve anlık API okuması destekli, tek kullanımlık ve programlanabilir e-posta kutuları.

---

## 📑 İçindekiler

1. [API Genel Bilgileri](#-api-genel-bilgileri)
   - [Temel URL (Base URL)](#temel-url-base-url)
   - [Kimlik Doğrulama (Authentication)](#kimlik-doğrulama-authentication)
   - [Standart Yanıt Formatı](#standart-yanıt-formatı)
   - [Rate Limit (İstek Sınırları)](#rate-limit-istek-sınırları)
2. [Hesap API Uç Noktaları](#-hesap-api-uç-noktaları)
   - [Kullanıcı Profili (`GET /me`)](#1-kullanıcı-profili-get-me)
   - [Bakiye Sorgulama (`GET /balance`)](#2-bakiye-sorgulama-get-balance)
3. [SMS Onay API](#-sms-onay-api)
   - [Aktif Servis Listesi (`GET /services`)](#3-aktif-servis-listesi-get-services)
   - [Ülke Fiyatları ve Stok (`GET /prices`)](#4-ülke-fiyatları-ve-stok-get-prices)
   - [Numara Siparişi Oluştur (`POST /order`)](#5-numara-siparişi-oluştur-post-order)
   - [Sipariş Durumu ve SMS Kodu (`GET /order/{id}`)](#6-sipariş-durumu-ve-sms-kodu-get-orderid)
   - [Sipariş İptali ve İade (`POST /order/{id}/cancel`)](#7-sipariş-iptali-ve-iade-post-orderidcancel)
   - [Siparişi Tamamla (`POST /order/{id}/done`)](#8-siparişi-tamamla-post-orderiddone)
4. [Numara Kiralama (Rent) API](#-numara-kiralama-rent-api)
   - [Kiralama Süreleri (`GET /rent/durations`)](#kiralama-süreleri-get-rentdurations)
   - [Kullanılabilir Ülkeler (`GET /rent/countries`)](#kullanılabilir-ülkeler-get-rentcountries)
   - [Servisler ve Fiyat Teklifi (`GET /rent/services`, `GET /rent/prices`)](#servisler-ve-fiyat-teklifi-get-rentservices-get-rentprices)
   - [Numara Kirala (`POST /rent/buy`)](#numara-kirala-post-rentbuy)
   - [Kiralık Siparişler ve Gelen SMS Listesi (`GET /rent/orders`, `GET /rent/orders/{id}`)](#kiralık-siparişler-ve-gelen-sms-listesi-get-rentorders-get-rentordersid)
   - [İptal veya Süre Uzatma (`POST /rent/orders/{id}/cancel`, `POST /rent/orders/{id}/reactivate`)](#iptal-veya-süre-uzatma-post-rentordersidcancel-post-rentordersidreactivate)
5. [Proxy API](#-proxy-api)
   - [Proxy Tipleri ve Ülkeler (`GET /proxy/types`, `GET /proxy/countries`)](#proxy-tipleri-ve-ülkeler-get-proxytypes-get-proxycountries)
   - [Fiyat Teklifi Hesaplama (`GET /proxy/price`)](#fiyat-teklifi-hesaplama-get-proxyprice)
   - [Proxy Satın Al (`POST /proxy/buy`)](#proxy-satın-al-post-proxybuy)
   - [Proxy Siparişleri (`GET /proxy/orders`, `GET /proxy/orders/{id}`)](#proxy-siparişleri-get-proxyorders-get-proxyordersid)
6. [eSIM API](#-esim-api)
   - [eSIM Ülkeleri ve Paketleri (`GET /esim/countries`, `GET /esim/packages`)](#esim-ülkeleri-ve-paketleri-get-esimcountries-get-esimpackages)
   - [eSIM Paketi Satın Al (`POST /esim/buy`)](#esim-paketi-satın-al-post-esimbuy)
   - [eSIM Detayı ve Bakiye Yükleme (`GET /esim/orders/{id}`, `POST /esim/orders/{id}/topup`)](#esim-detayı-ve-bakiye-yükleme-get-esimordersid-post-esimordersidtopup)
7. [Geçici Mail API](#-geçici-mail-api)
   - [Domain Listesi ve Mail Kutusu Oluşturma (`GET /mail/domains`, `POST /mail/create`)](#domain-listesi-ve-mail-kutusu-oluşturma-get-maildomains-post-mailcreate)
   - [Gelen Kutusu Mesajlarını Okuma (`GET /mail/orders/{id}/messages`)](#gelen-kutusu-mesajlarını-okuma-get-mailordersidmessages)
   - [Süre Uzatma veya İptal (`POST /mail/orders/{id}/extend`, `POST /mail/orders/{id}/cancel`)](#süre-uzatma-veya-iptal-post-mailordersidextend-post-mailordersidcancel)
8. [Sipariş Yaşam Döngüsü & Durumlar](#-sipariş-yaşam-döngüsü--durumlar)
9. [Hata Kodları Referansı](#-hata-kodları-referansı)
10. [Resmi SDK İstemcileri](#-resmi-sdk-istemcileri)
    - [PHP](#-php-sdk)
    - [Python](#-python-sdk)
    - [Node.js](#-nodejs-sdk)
    - [C# / .NET](#-c--net-sdk)
    - [Go](#-go-sdk)
11. [Lisans & Destek](#-lisans--destek)

---

## 🌐 API Genel Bilgileri

### Temel URL (Base URL)
Tüm istekler HTTPS üzerinden aşağıdaki temel adrese gönderilmelidir:
```
https://anindasms.com/api/v1
```

### Kimlik Doğrulama (Authentication)
API anahtarınızı [anindasms.com/profile](https://anindasms.com/profile) adresindeki panelinizden temin edebilirsiniz.

Her istekte API anahtarınızı HTTP başlığı (`Header`) ile gönderiniz:
```http
Authorization: Bearer SIZIN_API_KEYINIZ
```
veya
```http
X-API-Key: SIZIN_API_KEYINIZ
```

> 🔒 **Güvenlik Uyarısı:** API anahtarınızı URL sorgu parametrelerinde taşımayınız ve herkese açık kod depolarına yüklemeyiniz.

### Standart Yanıt Formatı
Tüm yanıtlar UTF-8 formatında JSON olarak döndürülür:

**Başarılı Yanıt (`HTTP 200 / 201`):**
```json
{
  "success": true,
  "data": { ... }
}
```

**Hatalı Yanıt (`HTTP 4xx / 5xx`):**
```json
{
  "success": false,
  "error": "Açıklayıcı hata mesajı.",
  "data": null
}
```

### Rate Limit (İstek Sınırları)
* **Genel Limit:** 120 istek / dakika
* **Sipariş Oluşturma:** 30 istek / dakika
* Başlıklar: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`

---

## 👤 Hesap API Uç Noktaları

### 1. Kullanıcı Profili (`GET /me`)
Kullanıcı ID, e-posta, mevcut bakiye ve kayıt tarihi bilgilerini döner.
```bash
curl -X GET "https://anindasms.com/api/v1/me" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### 2. Bakiye Sorgulama (`GET /balance`)
Güncel hesap bakiyesini döner.
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

## 📱 SMS Onay API

### 3. Aktif Servis Listesi (`GET /services`)
SMS onayı desteklenen tüm servislerin listesini döner.
```bash
curl -X GET "https://anindasms.com/api/v1/services" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### 4. Ülke Fiyatları ve Stok (`GET /prices`)
Belirli bir servis veya ülke için güncel fiyat ve stok durumunu döner.
```bash
curl -X GET "https://anindasms.com/api/v1/prices?service=wa&country=62" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### 5. Numara Siparişi Oluştur (`POST /order`)
Belirtilen servis ve ülke için geçici sanal numara tahsis eder.
```bash
curl -X POST "https://anindasms.com/api/v1/order" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ" \
  -H "Content-Type: application/json" \
  -d '{"service":"wa","country":"62"}'
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
    "cost": 12.50,
    "status": "waiting",
    "expires_at": "2026-10-24 14:12:10"
  }
}
```

### 6. Sipariş Durumu ve SMS Kodu (`GET /order/{id}`)
Kod gelene kadar her 3–5 saniyede bir sorgulayın:
```bash
curl -X GET "https://anindasms.com/api/v1/order/58231" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```
**Örnek Yanıt:**
```json
{
  "success": true,
  "data": {
    "order_id": 58231,
    "status": "received",
    "code": "489201",
    "sms_text": "WhatsApp doğrulama kodunuz: 489-201"
  }
}
```

### 7. Sipariş İptali ve İade (`POST /order/{id}/cancel`)
SMS gelmeden önce siparişi iptal eder ve ücreti anında bakiye olarak iade eder:
```bash
curl -X POST "https://anindasms.com/api/v1/order/58231/cancel" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### 8. Siparişi Tamamla (`POST /order/{id}/done`)
```bash
curl -X POST "https://anindasms.com/api/v1/order/58231/done" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

---

## ⏳ Numara Kiralama (Rent) API

Aynı numarayı uzun süreliğine kiralayarak sınırsız gelen SMS mesajı alabilirsiniz.

### Kiralama Süreleri (`GET /rent/durations`)
Kullanılabilir kiralama periyotlarını (saat veya gün) listeler.
```bash
curl -X GET "https://anindasms.com/api/v1/rent/durations" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### Kullanılabilir Ülkeler (`GET /rent/countries`)
```bash
curl -X GET "https://anindasms.com/api/v1/rent/countries?duration=24" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### Servisler ve Fiyat Teklifi (`GET /rent/services`, `GET /rent/prices`)
```bash
curl -X GET "https://anindasms.com/api/v1/rent/services?country=62&duration=24" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
curl -X GET "https://anindasms.com/api/v1/rent/prices?service=wa&country=62&duration=24" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### Numara Kirala (`POST /rent/buy`)
```bash
curl -X POST "https://anindasms.com/api/v1/rent/buy" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ" \
  -H "Content-Type: application/json" \
  -d '{"service":"wa","country":"62","duration":24}'
```
**Örnek Yanıt:**
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

### Kiralık Siparişler ve Gelen SMS Listesi (`GET /rent/orders`, `GET /rent/orders/{id}`)
```bash
curl -X GET "https://anindasms.com/api/v1/rent/orders/74102" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```
**Örnek Yanıt:**
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
        "text": "WhatsApp doğrulama kodunuz: 948-210",
        "date": "2026-10-24 14:12:00"
      }
    ]
  }
}
```

### İptal veya Süre Uzatma (`POST /rent/orders/{id}/cancel`, `POST /rent/orders/{id}/reactivate`)
```bash
# SMS gelmemişse ilk 20 dakika içinde ücretsiz iptal
curl -X POST "https://anindasms.com/api/v1/rent/orders/74102/cancel" -H "Authorization: Bearer SIZIN_API_KEYINIZ"

# Aynı numarayı 24 saat daha uzat
curl -X POST "https://anindasms.com/api/v1/rent/orders/74102/reactivate" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ" \
  -H "Content-Type: application/json" \
  -d '{"duration":24}'
```

---

## 🛡️ Proxy API

Statik Datacenter, Residential (Ev/ISP), IPv6 ve Dönen 4G Mobil proxy seçenekleri.

### Proxy Tipleri ve Ülkeler (`GET /proxy/types`, `GET /proxy/countries`)
```bash
curl -X GET "https://anindasms.com/api/v1/proxy/types" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
curl -X GET "https://anindasms.com/api/v1/proxy/countries?type=datacenter" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### Fiyat Teklifi Hesaplama (`GET /proxy/price`)
```bash
curl -X GET "https://anindasms.com/api/v1/proxy/price?type=datacenter&country=TR&period=30&quantity=1" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### Proxy Satın Al (`POST /proxy/buy`)
```bash
curl -X POST "https://anindasms.com/api/v1/proxy/buy" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ" \
  -H "Content-Type: application/json" \
  -d '{"type":"datacenter","country":"TR","period":30,"quantity":1}'
```
**Örnek Yanıt:**
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

### Proxy Siparişleri (`GET /proxy/orders`, `GET /proxy/orders/{id}`)
```bash
curl -X GET "https://anindasms.com/api/v1/proxy/orders" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
curl -X GET "https://anindasms.com/api/v1/proxy/orders/1042" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

---

## 📶 eSIM API

150'den fazla ülke için anında teslim edilen dijital seyahat internet paketleri.

### eSIM Ülkeleri ve Paketleri (`GET /esim/countries`, `GET /esim/packages`)
```bash
curl -X GET "https://anindasms.com/api/v1/esim/countries" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
curl -X GET "https://anindasms.com/api/v1/esim/packages?country=TR" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### eSIM Paketi Satın Al (`POST /esim/buy`)
```bash
curl -X POST "https://anindasms.com/api/v1/esim/buy" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ" \
  -H "Content-Type: application/json" \
  -d '{"package_id":"pkg_tr_1gb_7d","quantity":1}'
```
**Örnek Yanıt:**
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

### eSIM Detayı ve Bakiye Yükleme (`GET /esim/orders/{id}`, `POST /esim/orders/{id}/topup`)
```bash
curl -X GET "https://anindasms.com/api/v1/esim/orders/329" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
curl -X POST "https://anindasms.com/api/v1/esim/orders/329/topup" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ" \
  -H "Content-Type: application/json" \
  -d '{"package_id":"pkg_tr_1gb_7d"}'
```

---

## 📬 Geçici Mail API

Tek kullanımlık gelen kutusu oluşturma ve anlık gelen e-posta okuma.

### Domain Listesi ve Mail Kutusu Oluşturma (`GET /mail/domains`, `POST /mail/create`)
```bash
curl -X GET "https://anindasms.com/api/v1/mail/domains" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
curl -X POST "https://anindasms.com/api/v1/mail/create" \
  -H "Authorization: Bearer SIZIN_API_KEYINIZ" \
  -H "Content-Type: application/json" \
  -d '{"duration":60}'
```
**Örnek Yanıt:**
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

### Gelen Kutusu Mesajlarını Okuma (`GET /mail/orders/{id}/messages`)
```bash
curl -X GET "https://anindasms.com/api/v1/mail/orders/892/messages" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

### Süre Uzatma veya İptal (`POST /mail/orders/{id}/extend`, `POST /mail/orders/{id}/cancel`)
```bash
curl -X POST "https://anindasms.com/api/v1/mail/orders/892/extend" -H "Authorization: Bearer SIZIN_API_KEYINIZ" -d '{"duration":60}'
curl -X POST "https://anindasms.com/api/v1/mail/orders/892/cancel" -H "Authorization: Bearer SIZIN_API_KEYINIZ"
```

---

## 🔄 Sipariş Yaşam Döngüsü & Durumlar

| Durum | Açıklama |
| :--- | :--- |
| `waiting` / `pending` | Numara tahsis edildi. Gelen SMS bekleniyor. |
| `received` | SMS kodu başarıyla teslim alındı. Doğrulama kodu `code` veya `sms_list` alanında mevcuttur. |
| `done` / `completed` | Sipariş başarıyla tamamlandı. |
| `cancelled` | Kullanıcı tarafından iptal edildi. Bakiye anında iade edildi. |
| `expired` | Süre doldu (SMS gelmedi). Bakiye otomatik olarak hesaba iade edildi. |

---

## ⚠️ Hata Kodları Referansı

| HTTP Kodu | Hata Nedeni | Açıklama |
| :--- | :--- | :--- |
| `400 Bad Request` | Eksik / Geçersiz Parametre | Zorunlu alanlar eksik veya format hatalı. |
| `401 Unauthorized` | Yetkisiz İstek | API anahtarı eksik, geçersiz veya yetkisiz. |
| `402 Payment Required`| Yetersiz Bakiye | Hesap bakiyesi işlem tutarından düşük. |
| `403 Forbidden` | Erişim Reddedildi | Servis kapalı, ülke kısıtlı veya IP engellenmiş. |
| `404 Not Found` | Kaynak Bulunamadı | Belirtilen sipariş ID veya katalog öğesi bulunamadı. |
| `409 Conflict` | Stok Tükendi | Talep edilen ülke ve serviste anlık stok bulunmuyor. |
| `429 Too Many Requests`| İstek Sınırı Aşıldı | Dakikalık istek kotası doldu. `Retry-After` başlığını kontrol edin. |
| `500 Server Error` | Sunucu / Sağlayıcı Hatası| Üst sağlayıcı veya sunucu hatası. Birkaç saniye sonra tekrar deneyin. |

---

## 🛠 Resmi SDK İstemcileri

Resmi açık kaynak SDK istemcileri tüm 5 servis (SMS Onay, Kiralama, Proxy, eSIM, Geçici Mail) için otomatik tip denetimi, hata yakalama ve otomatik kod bekleme (`wait`) döngüleri sunar.

### 🐘 PHP SDK

**Kurulum:**
```bash
composer config repositories.anindasms vcs https://github.com/anindasms-dev/anindasms-sdk
composer require anindasms/sdk:dev-main
```

**Kullanım Örneği:**
```php
<?php
require 'vendor/autoload.php';

use AnindaSms\Client;

$client = new Client('SIZIN_API_KEYINIZ');

// 1. SMS Onay
$order = $client->createOrder('wa', '62');
echo "Telefon: " . $order['phone'] . PHP_EOL;
$sms = $client->waitForCode($order['order_id']);
echo "SMS Kodu: " . $sms['code'] . PHP_EOL;

// 2. Numara Kiralama
$rent = $client->rentBuy('wa', '62', 24);
echo "Kiralık Numara: " . $rent['phone'] . PHP_EOL;

// 3. Proxy Satın Al
$proxy = $client->proxyBuy('datacenter', 'TR', 30, 1);
echo "Proxy IP: " . $proxy['ip'] . ":" . $proxy['port'] . PHP_EOL;

// 4. eSIM Satın Al
$esim = $client->esimBuy('pkg_tr_1gb_7d');
echo "eSIM QR: " . $esim['qr_code_url'] . PHP_EOL;

// 5. Geçici Mail
$mail = $client->mailCreate(60);
echo "Geçici E-Posta: " . $mail['email'] . PHP_EOL;
```

---

### 🐍 Python SDK

**Kurulum:**
```bash
pip install "git+https://github.com/anindasms-dev/anindasms-sdk.git#subdirectory=python"
```

**Kullanım Örneği:**
```python
from anindasms import Client

client = Client("SIZIN_API_KEYINIZ")

# 1. SMS Onay
order = client.create_order("wa", "62")
print("Telefon:", order["phone"])
sms = client.wait_for_code(order["order_id"])
print("SMS Kodu:", sms["code"])

# 2. Numara Kiralama
rent = client.rent_buy("wa", "62", 24)
print("Kiralık Numara:", rent["phone"])

# 3. Proxy Satın Al
proxy = client.proxy_buy("datacenter", "TR", 30, 1)
print(f"Proxy: {proxy['ip']}:{proxy['port']}")

# 4. eSIM Satın Al
esim = client.esim_buy("pkg_tr_1gb_7d")
print("eSIM QR:", esim["qr_code_url"])

# 5. Geçici Mail
mail = client.mail_create(60)
print("Geçici Mail:", mail["email"])
```

---

### 🟩 Node.js SDK

**Kurulum:**
```bash
npm install github:anindasms-dev/anindasms-sdk
```

**Kullanım Örneği:**
```javascript
const { AnindaSMS } = require("anindasms");

const client = new AnindaSMS("SIZIN_API_KEYINIZ");

async function main() {
  // 1. SMS Onay
  const order = await client.createOrder("wa", "62");
  console.log("Telefon:", order.phone);
  const sms = await client.waitForCode(order.order_id);
  console.log("SMS Kodu:", sms.code);

  // 2. Numara Kiralama
  const rent = await client.rentBuy("wa", "62", 24);
  console.log("Kiralık Numara:", rent.phone);

  // 3. Proxy Satın Al
  const proxy = await client.proxyBuy("datacenter", "TR", 30, 1);
  console.log(`Proxy: ${proxy.ip}:${proxy.port}`);

  // 4. eSIM Satın Al
  const esim = await client.esimBuy("pkg_tr_1gb_7d");
  console.log("eSIM QR:", esim.qr_code_url);

  // 5. Geçici Mail
  const mail = await client.mailCreate(60);
  console.log("Geçici Mail:", mail.email);
}

main().catch(console.error);
```

---

### 🔷 C# / .NET SDK

Projenize `csharp/AnindaSms.csproj` dosyasını dahil edin veya `AnindaSms.cs` dosyasını referans verin:

```csharp
using System;
using AnindaSms;

var client = new AnindaSmsClient("SIZIN_API_KEYINIZ");

// 1. SMS Onay
var order = await client.CreateOrderAsync("wa", "62");
Console.WriteLine($"Telefon: {order.GetProperty("phone").GetString()}");
var sms = await client.WaitForCodeAsync(order.GetProperty("order_id").GetInt32());
Console.WriteLine($"Kod: {sms.GetProperty("code").GetString()}");

// 2. Numara Kiralama
var rent = await client.RentBuyAsync("wa", "62", 24);
Console.WriteLine($"Kiralık Numara: {rent.GetProperty("phone").GetString()}");

// 3. Proxy Satın Al
var proxy = await client.ProxyBuyAsync("datacenter", "TR", 30, 1);
Console.WriteLine($"Proxy: {proxy.GetProperty("ip").GetString()}:{proxy.GetProperty("port").GetInt32()}");

// 4. eSIM Satın Al
var esim = await client.EsimBuyAsync("pkg_tr_1gb_7d");
Console.WriteLine($"eSIM QR: {esim.GetProperty("qr_code_url").GetString()}");

// 5. Geçici Mail
var mail = await client.MailCreateAsync(60);
Console.WriteLine($"Geçici Mail: {mail.GetProperty("email").GetString()}");
```

---

### 🐹 Go SDK

**Kurulum:**
```bash
go get github.com/anindasms-dev/anindasms-sdk/go
```

**Kullanım Örneği:**
```go
package main

import (
	"fmt"
	anindasms "github.com/anindasms-dev/anindasms-sdk/go"
)

func main() {
	client := anindasms.NewClient("SIZIN_API_KEYINIZ")

	// 1. SMS Onay
	order, err := client.CreateOrder("wa", "62", "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Telefon: %s (#%d)\n", order.Phone, order.OrderID)

	sms, err := client.WaitForCode(order.OrderID, 0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("SMS Kodu: %s\n", *sms.Code)

	// 2. Numara Kiralama
	rent, _ := client.RentBuy("wa", "62", 24)
	fmt.Printf("Kiralama Yanıtı: %s\n", string(rent))

	// 3. Proxy
	proxy, _ := client.ProxyBuy("datacenter", "TR", 30, 1, 0)
	fmt.Printf("Proxy Yanıtı: %s\n", string(proxy))

	// 4. eSIM
	esim, _ := client.EsimBuy("pkg_tr_1gb_7d", 1)
	fmt.Printf("eSIM Yanıtı: %s\n", string(esim))

	// 5. Geçici Mail
	mail, _ := client.MailCreate(60, "", "")
	fmt.Printf("Mail Yanıtı: %s\n", string(mail))
}
```

---

## 📄 Lisans & Destek

Bu SDK ve dokümantasyon [MIT Lisansı](LICENSE) kapsamında açık kaynak olarak yayınlanmıştır.

* **Resmi Web Sitesi:** [anindasms.com](https://anindasms.com)
* **Web Dokümantasyonu:** [anindasms.com/api](https://anindasms.com/api)
* **Müşteri & API Destek:** support@anindasms.com / destek@anindasms.com
