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
}

module.exports = { AnindaSMS, AnindaSMSError };
