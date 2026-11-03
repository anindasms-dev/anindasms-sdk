export class AnindaSMSError extends Error {
  status: number;
  data: any;
  constructor(message: string, status?: number, data?: any);
}

export interface BalanceResponse {
  balance: number;
  currency: string;
  formatted: string;
}

export interface UserResponse {
  id: number;
  username: string;
  email: string;
  balance: number;
  currency: string;
  api_key_prefix: string;
  api_key_created_at: string;
  api_key_last_used_at: string | null;
}

export interface ServiceItem {
  code: string;
  name: string;
}

export interface OrderResponse {
  order_id: number;
  phone: string;
  service: string;
  country: string;
  status: 'waiting' | 'received' | 'done' | 'cancelled' | 'expired';
  code?: string | null;
  sms_text?: string | null;
  cost?: number;
  balance?: number;
  remaining?: number;
  expires_at?: string;
}

export interface ProxyBuyOptions {
  type: string;
  country: string;
  period?: number;
  quantity?: number;
  rotation_interval?: number;
}

export interface EsimBuyOptions {
  package_id: string;
  quantity?: number;
}

export interface MailCreateOptions {
  duration?: number;
  domain?: string;
  prefix?: string;
}

export class AnindaSMS {
  apiKey: string;
  baseUrl: string;

  constructor(apiKey: string, baseUrl?: string);

  // Account
  me(): Promise<UserResponse>;
  balance(): Promise<BalanceResponse>;

  // SMS Verification
  services(): Promise<ServiceItem[]>;
  prices(service: string, country?: string | null): Promise<any>;
  createOrder(service: string, country: string, operator?: string): Promise<OrderResponse>;
  getOrder(orderId: number): Promise<OrderResponse>;
  cancelOrder(orderId: number): Promise<any>;
  completeOrder(orderId: number): Promise<any>;
  waitForCode(orderId: number, timeout?: number, interval?: number): Promise<OrderResponse>;

  // Rent API (Numara Kiralama)
  rentDurations(): Promise<Array<{ hours: number; label: string }>>;
  rentCountries(duration?: number): Promise<Array<{ id: string; code: string; name: string }>>;
  rentServices(country: string, duration?: number): Promise<any[]>;
  rentPrice(service: string, country: string, duration?: number): Promise<any>;
  rentBuy(service: string, country: string, duration?: number): Promise<any>;
  rentOrders(): Promise<any[]>;
  rentOrder(orderId: number): Promise<any>;
  rentCancel(orderId: number): Promise<any>;
  rentReactivate(orderId: number, duration?: number): Promise<any>;

  // Proxy API
  proxyTypes(): Promise<any[]>;
  proxyCountries(type?: string | null): Promise<any[]>;
  proxyPrice(type: string, country?: string | null, period?: number, quantity?: number): Promise<any>;
  proxyBuy(type: string, country: string, period?: number, quantity?: number, rotationInterval?: number): Promise<any>;
  proxyOrders(page?: number, limit?: number): Promise<any[]>;
  proxyOrder(orderId: number): Promise<any>;

  // eSIM API
  esimCountries(): Promise<any[]>;
  esimPackages(country?: string | null): Promise<any[]>;
  esimBuy(packageId: string, quantity?: number): Promise<any>;
  esimOrders(page?: number, limit?: number): Promise<any[]>;
  esimOrder(orderId: number): Promise<any>;
  esimTopup(orderId: number, packageId: string): Promise<any>;

  // Temporary Mail API
  mailDomains(): Promise<any[]>;
  mailCreate(duration?: number, domain?: string | null, prefix?: string | null): Promise<any>;
  mailOrders(page?: number, limit?: number): Promise<any[]>;
  mailOrder(orderId: number): Promise<any>;
  mailMessages(orderId: number): Promise<any[]>;
  mailExtend(orderId: number, duration?: number): Promise<any>;
  mailCancel(orderId: number): Promise<any>;
}

