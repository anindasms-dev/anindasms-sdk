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
