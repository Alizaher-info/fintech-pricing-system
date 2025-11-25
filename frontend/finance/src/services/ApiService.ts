/**
 * API Service
 * HTTP client for communicating with backend API
 * 
 * This service provides methods for making HTTP requests to the backend.
 * Handles GET, POST, PUT, DELETE operations with optional authentication.
 */

const API_BASE = '/api';

/**
 * API Service for making HTTP requests to backend
 * Provides simple methods for common HTTP operations
 */
class ApiService {
  /**
   * Make a GET request
   * @param endpoint - API endpoint (e.g., '/users', '/transactions')
   * @param token - Optional auth token
   */
  async get(endpoint: string, token?: string): Promise<any> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      method: 'GET',
      headers,
    });

    const result = await response.json();

    if (!response.ok) {
      throw result;
    }

    return result;
  }

  /**
   * Make a POST request
   * @param endpoint - API endpoint
   * @param data - Request body
   * @param token - Optional auth token
   */
  async post(endpoint: string, data: any, token?: string): Promise<any> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      method: 'POST',
      headers,
      body: JSON.stringify(data),
    });

    const result = await response.json();

    if (!response.ok) {
      throw result;
    }

    return result;
  }

  /**
   * Make a PUT request
   * @param endpoint - API endpoint
   * @param data - Request body
   * @param token - Optional auth token
   */
  async put(endpoint: string, data: any, token?: string): Promise<any> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      method: 'PUT',
      headers,
      body: JSON.stringify(data),
    });

    const result = await response.json();

    if (!response.ok) {
      throw result;
    }

    return result;
  }

  /**
   * Make a DELETE request
   * @param endpoint - API endpoint
   * @param token - Optional auth token
   */
  async delete(endpoint: string, token?: string): Promise<any> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      method: 'DELETE',
      headers,
    });

    const result = await response.json();

    if (!response.ok) {
      throw result;
    }

    return result;
  }
}

export const apiService = new ApiService();
