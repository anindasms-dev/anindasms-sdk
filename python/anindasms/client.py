import json
import time
import urllib.error
import urllib.parse
import urllib.request
from typing import Any, Dict, List, Optional

DEFAULT_BASE_URL = "https://anindasms.com/api/v1"
VERSION = "1.0.0"


class AnindaSMSError(Exception):
    """API hatası. `status` HTTP kodunu, `data` varsa ek bilgiyi içerir."""

    def __init__(self, message: str, status: int = 0, data: Optional[Dict[str, Any]] = None):
        super().__init__(message)
        self.status = status
        self.data = data


class Client:
    """AnındaSMS REST API istemcisi.

    >>> client = Client("YOUR_API_KEY")
    >>> order = client.create_order("wa", "62")
    >>> sms = client.wait_for_code(order["order_id"])
    """

    def __init__(self, api_key: str, base_url: str = DEFAULT_BASE_URL, timeout: int = 30):
        if not api_key or not api_key.strip():
            raise ValueError("API key is required.")
        self.api_key = api_key.strip()
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout

    # ---- Hesap -------------------------------------------------------
    def me(self) -> Dict[str, Any]:
        return self._request("GET", "/me")

    def balance(self) -> Dict[str, Any]:
        return self._request("GET", "/balance")

    # ---- Servis & fiyat ----------------------------------------------
    def services(self) -> List[Dict[str, Any]]:
        return self._request("GET", "/services").get("services", [])

    def prices(self, service: str, country: Optional[str] = None) -> Dict[str, Any]:
        query = {"service": service}
        if country:
            query["country"] = country
        return self._request("GET", "/prices?" + urllib.parse.urlencode(query))

    # ---- Siparişler --------------------------------------------------
    def create_order(self, service: str, country: str, operator: str = "any") -> Dict[str, Any]:
        return self._request("POST", "/order", {"service": service, "country": country, "operator": operator})

    def get_order(self, order_id: int) -> Dict[str, Any]:
        return self._request("GET", f"/order/{int(order_id)}")

    def cancel_order(self, order_id: int) -> Dict[str, Any]:
        return self._request("POST", f"/order/{int(order_id)}/cancel")

    def complete_order(self, order_id: int) -> Dict[str, Any]:
        return self._request("POST", f"/order/{int(order_id)}/done")

    def wait_for_code(self, order_id: int, timeout: int = 600, interval: int = 5) -> Dict[str, Any]:
        """SMS kodu gelene kadar bekler; kod gelince sipariş verisini döner."""
        deadline = time.time() + timeout
        interval = max(2, interval)
        while True:
            order = self.get_order(order_id)
            status = order.get("status", "")
            if order.get("code") and status in ("received", "done"):
                return order
            if status in ("cancelled", "expired", "done"):
                raise AnindaSMSError(f"Order ended without a code (status: {status}).", 0, order)
            if time.time() + interval > deadline:
                raise AnindaSMSError("Timed out while waiting for the SMS code.", 0, order)
            time.sleep(interval)

    # ---- Rent API (Numara Kiralama) -----------------------------------
    def rent_durations(self) -> List[Dict[str, Any]]:
        """Kiralama surelerini listeler."""
        return self._request("GET", "/rent/durations").get("durations", [])

    def rent_countries(self, duration: int = 24) -> List[Dict[str, Any]]:
        """Sureye gore kiralik ulkeleri listeler."""
        return self._request("GET", f"/rent/countries?duration={duration}").get("countries", [])

    def rent_services(self, country: str, duration: int = 24) -> List[Dict[str, Any]]:
        """Kiralanabilir servisleri, stok ve fiyatlari doner."""
        path = f"/rent/services?country={urllib.parse.quote(country)}&duration={duration}"
        return self._request("GET", path).get("services", [])

    def rent_price(self, service: str, country: str, duration: int = 24) -> Dict[str, Any]:
        """Kiralama anlik fiyatini doner."""
        path = f"/rent/prices?service={urllib.parse.quote(service)}&country={urllib.parse.quote(country)}&duration={duration}"
        return self._request("GET", path)

    def rent_buy(self, service: str, country: str, duration: int = 24) -> Dict[str, Any]:
        """Numara kiralar ve telefon numarasini doner."""
        body = {"service": service, "country": country, "duration": duration}
        return self._request("POST", "/rent/buy", body)

    def rent_orders(self) -> List[Dict[str, Any]]:
        """Kiralik numara siparislerini listeler."""
        return self._request("GET", "/rent/orders").get("orders", [])

    def rent_order(self, order_id: int) -> Dict[str, Any]:
        """Kiralik numara durumunu ve gelen tum SMS'leri doner."""
        return self._request("GET", f"/rent/orders/{int(order_id)}")

    def rent_cancel(self, order_id: int) -> Dict[str, Any]:
        """Kiralik numarayi iptal eder (ilk 20 dk SMS gelmediyse)."""
        return self._request("POST", f"/rent/orders/{int(order_id)}/cancel")

    def rent_reactivate(self, order_id: int, duration: int = 24) -> Dict[str, Any]:
        """Ayni numarayi yeniden kiralar / uzatir."""
        return self._request("POST", f"/rent/orders/{int(order_id)}/reactivate", {"duration": duration})
    # ---- Proxy API ---------------------------------------------------
    def proxy_types(self) -> List[Dict[str, Any]]:
        """Mevcut proxy turlerini listeler."""
        return self._request("GET", "/proxy/types").get("types", [])

    def proxy_countries(self, proxy_type: Optional[str] = None) -> List[Dict[str, Any]]:
        """Secilen proxy turu icin mevcut ulkeleri listeler."""
        path = "/proxy/countries"
        if proxy_type:
            path += "?" + urllib.parse.urlencode({"type": proxy_type})
        return self._request("GET", path).get("countries", [])

    def proxy_price(self, proxy_type: str, country: Optional[str] = None, period: int = 30, quantity: int = 1) -> Dict[str, Any]:
        """Proxy maliyetini hesaplar."""
        params = {"type": proxy_type, "period": period, "quantity": quantity}
        if country:
            params["country"] = country
        return self._request("GET", "/proxy/price?" + urllib.parse.urlencode(params))

    def proxy_buy(self, proxy_type: str, country: str, period: int = 30, quantity: int = 1, rotation_interval: int = 0) -> Dict[str, Any]:
        """Proxy satin alir ve baglanti bilgilerini doner."""
        body = {"type": proxy_type, "country": country, "period": period, "quantity": quantity}
        if rotation_interval > 0:
            body["rotation_interval"] = rotation_interval
        return self._request("POST", "/proxy/buy", body)

    def proxy_orders(self, page: int = 1, limit: int = 20) -> List[Dict[str, Any]]:
        """Proxy siparislerini listeler."""
        return self._request("GET", f"/proxy/orders?page={page}&limit={limit}").get("orders", [])

    def proxy_order(self, order_id: int) -> Dict[str, Any]:
        """Proxy siparis detayini ve sifresini doner."""
        return self._request("GET", f"/proxy/orders/{int(order_id)}")

    # ---- eSIM API ----------------------------------------------------
    def esim_countries(self) -> List[Dict[str, Any]]:
        """eSIM desteklenen ulkeleri listeler."""
        return self._request("GET", "/esim/countries").get("countries", [])

    def esim_packages(self, country: Optional[str] = None) -> List[Dict[str, Any]]:
        """eSIM veri paketlerini listeler."""
        path = "/esim/packages"
        if country:
            path += "?" + urllib.parse.urlencode({"country": country})
        return self._request("GET", path).get("packages", [])

    def esim_buy(self, package_id: str, quantity: int = 1) -> Dict[str, Any]:
        """eSIM paketi satin alir ve QR / LPA doner."""
        return self._request("POST", "/esim/buy", {"package_id": package_id, "quantity": quantity})

    def esim_orders(self, page: int = 1, limit: int = 20) -> List[Dict[str, Any]]:
        """eSIM siparislerini listeler."""
        return self._request("GET", f"/esim/orders?page={page}&limit={limit}").get("orders", [])

    def esim_order(self, order_id: int) -> Dict[str, Any]:
        """eSIM siparis detayini ve QR linkini doner."""
        return self._request("GET", f"/esim/orders/{int(order_id)}")

    def esim_topup(self, order_id: int, package_id: str) -> Dict[str, Any]:
        """eSIM siparisini yukler (topup)."""
        return self._request("POST", f"/esim/orders/{int(order_id)}/topup", {"package_id": package_id})

    # ---- Gecici Mail API ---------------------------------------------
    def mail_domains(self) -> List[Dict[str, Any]]:
        """Gecici mail alan adlarini listeler."""
        return self._request("GET", "/mail/domains").get("domains", [])

    def mail_create(self, duration: int = 60, domain: Optional[str] = None, prefix: Optional[str] = None) -> Dict[str, Any]:
        """Yeni gecici e-posta kutusu olusturur."""
        body = {"duration": duration}
        if domain:
            body["domain"] = domain
        if prefix:
            body["prefix"] = prefix
        return self._request("POST", "/mail/create", body)

    def mail_orders(self, page: int = 1, limit: int = 20) -> List[Dict[str, Any]]:
        """Gecici mail kutularini listeler."""
        return self._request("GET", f"/mail/orders?page={page}&limit={limit}").get("orders", [])

    def mail_order(self, order_id: int) -> Dict[str, Any]:
        """Gecici mail kutusu detayini doner."""
        return self._request("GET", f"/mail/orders/{int(order_id)}")

    def mail_messages(self, order_id: int) -> List[Dict[str, Any]]:
        """Gelen kutusundaki e-postalari doner."""
        return self._request("GET", f"/mail/orders/{int(order_id)}/messages").get("messages", [])

    def mail_extend(self, order_id: int, duration: int = 60) -> Dict[str, Any]:
        """Gecici e-posta suresini uzatir."""
        return self._request("POST", f"/mail/orders/{int(order_id)}/extend", {"duration": duration})

    def mail_cancel(self, order_id: int) -> Dict[str, Any]:
        """Gecici e-posta kutusunu kapatir/siler."""
        return self._request("POST", f"/mail/orders/{int(order_id)}/cancel")

    # ---- İç ----------------------------------------------------------
    def _request(self, method: str, path: str, body: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Accept": "application/json",
            "User-Agent": f"anindasms-python/{VERSION}",
        }
        data = None
        if body is not None:
            headers["Content-Type"] = "application/json"
            data = json.dumps(body).encode("utf-8")
        elif method == "POST":
            data = b""

        req = urllib.request.Request(self.base_url + path, data=data, headers=headers, method=method)
        try:
            with urllib.request.urlopen(req, timeout=self.timeout) as res:
                status, raw = res.status, res.read()
        except urllib.error.HTTPError as err:
            status, raw = err.code, err.read()
        except urllib.error.URLError as err:
            raise AnindaSMSError(f"Connection error: {err.reason}") from err

        try:
            payload = json.loads(raw.decode("utf-8"))
        except ValueError:
            raise AnindaSMSError(f"Invalid response from API (HTTP {status}).", status)

        if not payload.get("success"):
            raise AnindaSMSError(payload.get("error", "Unknown error"), status, payload.get("data"))
        return payload.get("data") or {}
