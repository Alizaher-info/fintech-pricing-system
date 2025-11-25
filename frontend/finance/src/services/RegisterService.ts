import { RegisterStrategyFactory } from '../auth/register/factory/RegisterStrategyFactory';
import { AuthService } from '../auth/security/auth/AuthService';
import type { RegisterRequest, RegisterResponse } from '../auth/register/interface/IRegister';

/**
 * Register Service
 * Orchestrates registration flow: uses factory to register, manages auth state
 * Follows Single Responsibility Principle
 */
export class RegisterService {
  private static instance: RegisterService;
  private registerFactory: RegisterStrategyFactory;
  private authService: AuthService;

  /**
   * Private constructor to enforce Singleton pattern
   */
  private constructor() {
    this.registerFactory = new RegisterStrategyFactory();
    this.authService = AuthService.getInstance();
    console.log('🔐 RegisterService instance created');
  }

  /**
   * Get singleton instance of RegisterService
   * @returns RegisterService - The singleton instance
   */
  public static getInstance(): RegisterService {
    if (!RegisterService.instance) {
      RegisterService.instance = new RegisterService();
    }
    return RegisterService.instance;
  }

  /**
   * Register new user with provided data
   * Automatically detects registration strategy (email, google, phone, etc.)
   * 
   * @param data - Registration data
   * @returns Promise<RegisterResponse> - Response with token and user data
   * @throws Error if registration fails
   */
  public async register(data: RegisterRequest): Promise<RegisterResponse> {
    try {
      console.log('🔵 Starting registration process...');

      // Use factory to determine strategy and register
      const response = await this.registerFactory.register(data);

      if (!response.success || !response.token) {
        throw new Error('Registration failed: Invalid response from server');
      }

      // Store token using AuthService (auto-login after registration)
      this.authService.setToken(response.token);

      console.log('✅ Registration successful, user authenticated');
      return response;

    } catch (error) {
      console.error('❌ Registration failed:', error);
      // Clear any existing token on failed registration
      this.authService.clearToken();
      throw error;
    }
  }

  /**
   * Check if user is currently authenticated
   * @returns boolean - true if user has valid token
   */
  public isAuthenticated(): boolean {
    return this.authService.isAuthenticated();
  }

  /**
   * Get current authentication token
   * @returns string | null - Token if authenticated
   */
  public getToken(): string | null {
    return this.authService.getToken();
  }
}
