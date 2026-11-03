class AnindaSMSError extends Error {
  constructor(message, status = 0, data = null) {
    super(message);
    this.name = 'AnindaSMSError';
    this.status = status;
    this.data = data;
  }
}

class AnindaSMS {
  /**
   * @param {string} apiKey
   * @param {string} [baseUrl='https://anindasms.com/api/v1']
   */
  constructor(apiKey, baseUrl = 'https://anindasms.com/api/v1') {
    if (!apiKey || typeof apiKey !== 'string' || !apiKey.trim()) {
      throw new Error('API key is required.');
    }
    this.apiKey = apiKey.trim();
    this.baseUrl = baseUrl.replace(/\/+$/, '');
  }

  async _request(method, path, body = null) {
    const url = `${this.baseUrl}${path}`;
    const headers = {
      'Authorization': `Bearer ${this.apiKey}`,
      'Accept': 'application/json',
      'User-Agent': 'anindasms-node/1.0.0'
    };

    const options = { method, headers };
    if (body) {
      headers['Content-Type'] = 'application/json';
      options.body = JSON.stringify(body);
    }

    let res;
    try {
      res = await fetch(url, options);
    } catch (err) {
      throw new AnindaSMSError(`Connection error: ${err.message}`);
    }

    let json;
    try {
      json = await res.json();
    } catch (err) {
      throw new AnindaSMSError(`Invalid JSON from API (HTTP ${res.status})`, res.status);
    }

    if (!json || !json.success) {
      throw new AnindaSMSError(json?.error || 'Unknown error', res.status, json?.data || null);
    }

    return json.data;
  }

  /** Hesap özeti */
  async me() {
    return this._request('GET', '/me');
  }

  /** Güncel bakiye */
  async balance() {
    return this._request('GET', '/balance');
  }

  /** Aktif servisler */
  async services() {
    const res = await this._request('GET', '/services');
    return res.services || [];
  }

  /** Fiyat ve stok bilgisi */
  async prices(service, country = null) {
    const qs = new URLSearchParams({ service });
    if (country) qs.set('country', country);
    return this._request('GET', `/prices?${qs.toString()}`);
  }

  /** Numara siparişi */
  async createOrder(service, country, operator = 'any') {
    return this._request('POST', '/order', { service, country, operator });
  }

  /** Sipariş durumu */
  async getOrder(orderId) {
    return this._request('GET', `/order/${orderId}`);
  }

  /** Siparişi iptal et */
  async cancelOrder(orderId) {
    return this._request('POST', `/order/${orderId}/cancel`);
  }

  /** Siparişi tamamlandı olarak işaretle */
  async completeOrder(orderId) {
    return this._request('POST', `/order/${orderId}/done`);
  }

  /**
   * SMS kodu gelene kadar bekler.
   * @param {number} orderId
   * @param {number} [timeout=600] Toplam bekleme süresi (saniye)
   * @param {number} [interval=5] Kontrol sıklığı (saniye)
   */
  async waitForCode(orderId, timeout = 600, interval = 5) {
    const deadline = Date.now() + timeout * 1000;
    const delay = Math.max(2, interval) * 1000;

    while (Date.now() < deadline) {
      const order = await this.getOrder(orderId);
      const status = order.status;

      if (order.code && (status === 'received' || status === 'done')) {
        return order;
      }
      if (['cancelled', 'expired', 'done'].includes(status)) {
        throw new AnindaSMSError(`Order ended without code (status: ${status})`, 0, order);
      }

      await new Promise(r => setTimeout(r, delay));
    }

    throw new AnindaSMSError('Timed out waiting for SMS code');
  }

  // ---- Rent API (Numara Kiralama) -----------------------------------
  /** Kiralama sureleri */
  async rentDurations() {
    const res = await this._request('GET', '/rent/durations');
    return res.durations || [];
  }

  /** Kiralik ulkeler */
  async rentCountries(duration = 24) {
    const res = await this._request('GET', `/rent/countries?duration=${duration}`);
    return res.countries || [];
  }

  /** Kiralanabilir servisler */
  async rentServices(country, duration = 24) {
    const res = await this._request('GET', `/rent/services?country=${encodeURIComponent(country)}&duration=${duration}`);
    return res.services || [];
  }

  /** Kiralama fiyati */
  async rentPrice(service, country, duration = 24) {
    return this._request('GET', `/rent/prices?service=${encodeURIComponent(service)}&country=${encodeURIComponent(country)}&duration=${duration}`);
  }

  /** Numara kirala */
  async rentBuy(service, country, duration = 24) {
    return this._request('POST', '/rent/buy', { service, country, duration });
  }

  /** Kiralik siparisler */
  async rentOrders() {
    const res = await this._request('GET', '/rent/orders');
    return res.orders || [];
  }

  /** Kiralik siparis detayi & gelen SMS'ler */
  async rentOrder(orderId) {
    return this._request('GET', `/rent/orders/${orderId}`);
  }

  /** Kiralik siparisi iptal et */
  async rentCancel(orderId) {
    return this._request('POST', `/rent/orders/${orderId}/cancel`);
  }

  /** Kiralik numarayi uzat / yeniden kirala */
  async rentReactivate(orderId, duration = 24) {
    return this._request('POST', `/rent/orders/${orderId}/reactivate`, { duration });
  }
  // ---- Proxy API ----------------------------------------------------
  /** Mevcut proxy turleri */
  async proxyTypes() {
    const res = await this._request('GET', '/proxy/types');
    return res.types || [];
  }

  /** Proxy ulkeleri */
  async proxyCountries(type = null) {
    const qs = type ? `?type=${encodeURIComponent(type)}` : '';
    const res = await this._request('GET', `/proxy/countries${qs}`);
    return res.countries || [];
  }

  /** Proxy fiyat hesaplama */
  async proxyPrice(type, country = null, period = 30, quantity = 1) {
    const qs = new URLSearchParams({ type, period: String(period), quantity: String(quantity) });
    if (country) qs.set('country', country);
    return this._request('GET', `/proxy/price?${qs.toString()}`);
  }

  /** Proxy satin al */
  async proxyBuy(type, country, period = 30, quantity = 1, rotationInterval = 0) {
    const body = { type, country, period, quantity };
    if (rotationInterval > 0) body.rotation_interval = rotationInterval;
    return this._request('POST', '/proxy/buy', body);
  }

  /** Proxy siparisleri */
  async proxyOrders(page = 1, limit = 20) {
    const res = await this._request('GET', `/proxy/orders?page=${page}&limit=${limit}`);
    return res.orders || [];
  }

  /** Proxy siparis detayi */
  async proxyOrder(orderId) {
    return this._request('GET', `/proxy/orders/${orderId}`);
  }

  // ---- eSIM API -----------------------------------------------------
  /** eSIM ulkeleri */
  async esimCountries() {
    const res = await this._request('GET', '/esim/countries');
    return res.countries || [];
  }

  /** eSIM veri paketleri */
  async esimPackages(country = null) {
    const qs = country ? `?country=${encodeURIComponent(country)}` : '';
    const res = await this._request('GET', `/esim/packages${qs}`);
    return res.packages || [];
  }

  /** eSIM paketi satin al */
  async esimBuy(packageId, quantity = 1) {
    return this._request('POST', '/esim/buy', { package_id: packageId, quantity });
  }

  /** eSIM siparisleri */
  async esimOrders(page = 1, limit = 20) {
    const res = await this._request('GET', `/esim/orders?page=${page}&limit=${limit}`);
    return res.orders || [];
  }

  /** eSIM siparis detayi */
  async esimOrder(orderId) {
    return this._request('GET', `/esim/orders/${orderId}`);
  }

  /** eSIM paket yukle (topup) */
  async esimTopup(orderId, packageId) {
    return this._request('POST', `/esim/orders/${orderId}/topup`, { package_id: packageId });
  }

  // ---- Gecici Mail API ----------------------------------------------
  /** Gecici mail alan adlari */
  async mailDomains() {
    const res = await this._request('GET', '/mail/domains');
    return res.domains || [];
  }

  /** Yeni gecici mail olustur */
  async mailCreate(duration = 60, domain = null, prefix = null) {
    const body = { duration };
    if (domain) body.domain = domain;
    if (prefix) body.prefix = prefix;
    return this._request('POST', '/mail/create', body);
  }

  /** Gecici mail kutulari */
  async mailOrders(page = 1, limit = 20) {
    const res = await this._request('GET', `/mail/orders?page=${page}&limit=${limit}`);
    return res.orders || [];
  }

  /** Gecici mail kutusu detayi */
  async mailOrder(orderId) {
    return this._request('GET', `/mail/orders/${orderId}`);
  }

  /** Gelen kutusu mesajlari */
  async mailMessages(orderId) {
    const res = await this._request('GET', `/mail/orders/${orderId}/messages`);
    return res.messages || [];
  }

  /** Gecici mail suresini uzat */
  async mailExtend(orderId, duration = 60) {
    return this._request('POST', `/mail/orders/${orderId}/extend`, { duration });
  }

  /** Gecici mail kutusunu kapat / sil */
  async mailCancel(orderId) {
    return this._request('POST', `/mail/orders/${orderId}/cancel`);
  }
}

module.exports = { AnindaSMS, AnindaSMSError };
