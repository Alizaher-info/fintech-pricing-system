const API_BASE = '/api';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface AuthResponse {
  success: boolean;
  token: string;
  user: {
    id: number;
    email: string;
    firstName: string;
    lastName: string;
    roles: string[];
  };
}

export interface ErrorResponse {
  error: string;
  remaining_attempts?: {
    ip: number;
    user: number;
  };
  lockout_remaining?: number;
  retry_after?: number;
}

class ApiService {
  private token: string | null = null;

  constructor() {
    this.token = localStorage.getItem('auth_token');
  }

  setToken(token: string) {
    this.token = token;
    localStorage.setItem('auth_token', token);
  }

  clearToken() {
    this.token = null;
    localStorage.removeItem('auth_token');
  }

  getToken(): string | null {
    return this.token;
  }

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    const result = await response.json();

    if (!response.ok) {
      throw result as ErrorResponse;
    }

    if (result.token) {
      this.setToken(result.token);
    }

    return result;
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    console.log('🔵 Login Request:', {
      url: `${API_BASE}/auth/login`,
      email: data.email,
    });

    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    console.log('🔵 Login Response Status:', response.status, response.statusText);

    const result = await response.json();
    console.log('🔵 Login Response Data:', result);

    if (!response.ok) {
      console.error('❌ Login Failed:', result);
      throw result as ErrorResponse;
    }

    if (result.token) {
      this.setToken(result.token);
      console.log('✅ Token saved successfully');
    }

    return result;
  }

  async validateToken(): Promise<{ success: boolean; user: any }> {
    if (!this.token) {
      throw { error: 'No token found' };
    }

    const response = await fetch(`${API_BASE}/auth/validate`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.token}`,
      },
    });

    const result = await response.json();

    if (!response.ok) {
      this.clearToken();
      throw result;
    }

    return result;
  }

  logout() {
    this.clearToken();
  }
}

export const api = new ApiService();
