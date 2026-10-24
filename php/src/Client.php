<?php
declare(strict_types=1);

namespace AnindaSms;

/**
 * AnındaSMS REST API istemcisi (bağımlılıksız, yalnızca ext-curl + ext-json).
 *
 *   $client = new \AnindaSms\Client('YOUR_API_KEY');
 *   $order  = $client->createOrder('wa', '62');
 *   $sms    = $client->waitForCode($order['order_id']);
 */
class Client
{
    public const DEFAULT_BASE_URL = 'https://anindasms.com/api/v1';
    public const VERSION = '1.0.0';

    private string $apiKey;
    private string $baseUrl;
    private int $timeout;

    public function __construct(string $apiKey, string $baseUrl = self::DEFAULT_BASE_URL, int $timeout = 30)
    {
        if (trim($apiKey) === '') {
            throw new \InvalidArgumentException('API key is required.');
        }
        $this->apiKey  = trim($apiKey);
        $this->baseUrl = rtrim($baseUrl, '/');
        $this->timeout = $timeout;
    }

    /** Hesap özeti */
    public function me(): array
    {
        return $this->request('GET', '/me');
    }

    /** Güncel bakiye: ['balance' => float, 'currency' => 'TRY', 'formatted' => string] */
    public function balance(): array
    {
        return $this->request('GET', '/balance');
    }

    /** Aktif servis listesi: [['code' => 'wa', 'name' => 'WhatsApp'], ...] */
    public function services(): array
    {
        return $this->request('GET', '/services')['services'] ?? [];
    }

    /**
     * Fiyat ve stok. $country verilirse tek ülkenin fiyatı döner,
     * verilmezse ['service' => ..., 'prices' => [...]] döner.
     */
    public function prices(string $service, ?string $country = null): array
    {
        $query = ['service' => $service];
        if ($country !== null && $country !== '') {
            $query['country'] = $country;
        }
        return $this->request('GET', '/prices?' . http_build_query($query));
    }

    /** Numara siparişi oluşturur. */
    public function createOrder(string $service, string $country, string $operator = 'any'): array
    {
        return $this->request('POST', '/order', [
            'service'  => $service,
            'country'  => $country,
            'operator' => $operator,
        ]);
    }

    /** Sipariş durumu + SMS kodu */
    public function getOrder(int $orderId): array
    {
        return $this->request('GET', '/order/' . $orderId);
    }

    /** SMS gelmemiş siparişi iptal eder ve ücreti iade eder. */
    public function cancelOrder(int $orderId): array
    {
        return $this->request('POST', '/order/' . $orderId . '/cancel');
    }

    /** Siparişi tamamlandı olarak işaretler. */
    public function completeOrder(int $orderId): array
    {
        return $this->request('POST', '/order/' . $orderId . '/done');
    }

    /**
     * SMS kodu gelene kadar siparişi sorgular.
     * Kod gelirse sipariş verisini döner; süre dolarsa, iptal/expire olursa istisna fırlatır.
     *
     * @param int $timeout  Saniye cinsinden en uzun bekleme
     * @param int $interval Sorgu aralığı (saniye)
     */
    public function waitForCode(int $orderId, int $timeout = 600, int $interval = 5): array
    {
        $deadline = time() + $timeout;
        $interval = max(2, $interval);

        while (true) {
            $order = $this->getOrder($orderId);
            $status = (string)($order['status'] ?? '');

            if (!empty($order['code']) && in_array($status, ['received', 'done'], true)) {
                return $order;
            }
            if (in_array($status, ['cancelled', 'expired', 'done'], true)) {
                throw new AnindaSmsException('Order ended without a code (status: ' . $status . ').', 0, $order);
            }
            if (time() + $interval > $deadline) {
                throw new AnindaSmsException('Timed out while waiting for the SMS code.', 0, $order);
            }
            sleep($interval);
        }
    }

    /**
     * @throws AnindaSmsException
     */
    private function request(string $method, string $path, ?array $body = null): array
    {
        $ch = curl_init($this->baseUrl . $path);
        $headers = [
            'Authorization: Bearer ' . $this->apiKey,
            'Accept: application/json',
            'User-Agent: anindasms-php/' . self::VERSION,
        ];
        $options = [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_TIMEOUT        => $this->timeout,
            CURLOPT_CONNECTTIMEOUT => 10,
            CURLOPT_CUSTOMREQUEST  => $method,
        ];
        if ($body !== null) {
            $headers[] = 'Content-Type: application/json';
            $options[CURLOPT_POSTFIELDS] = json_encode($body);
        } elseif ($method === 'POST') {
            $options[CURLOPT_POSTFIELDS] = '';
        }
        $options[CURLOPT_HTTPHEADER] = $headers;
        curl_setopt_array($ch, $options);

        $raw    = curl_exec($ch);
        $status = (int)curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $error  = curl_error($ch);
        curl_close($ch);

        if ($raw === false) {
            throw new AnindaSmsException('Connection error: ' . $error);
        }

        $json = json_decode((string)$raw, true);
        if (!is_array($json)) {
            throw new AnindaSmsException('Invalid response from API (HTTP ' . $status . ').', $status);
        }
        if (empty($json['success'])) {
            throw new AnindaSmsException((string)($json['error'] ?? 'Unknown error'), $status, $json['data'] ?? null);
        }
        return is_array($json['data'] ?? null) ? $json['data'] : [];
    }
}
