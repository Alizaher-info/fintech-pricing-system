import { LoginStrategyFactory } from '../auth/login/factory/LoginStrategyFactory';
import { AuthService } from '../auth/security/auth/AuthService';
import { TokenValidator } from '../auth/security/auth/TokenValidator';
import type { LoginRequest, LoginResponse } from '../auth/login/interface/ILogin';

/**
 * Login Service
 * Orchestrates login flow: uses factory to authenticate, manages auth state
 * Follows Single Responsibility Principle
 */
export class LoginService {
  private static instance: LoginService;
  private loginFactory: LoginStrategyFactory;
  private authService: AuthService;
  private tokenValidator: TokenValidator;

  /**
   * Private constructor to enforce Singleton pattern
   */
  private constructor() {
    this.loginFactory = new LoginStrategyFactory();
    this.authService = AuthService.getInstance();
    this.tokenValidator = TokenValidator.getInstance();
    console.log('🔐 LoginService instance created');
  }

  /**
   * Get singleton instance of LoginService
   * @returns LoginService - The singleton instance
   */
  public static getInstance(): LoginService {
    if (!LoginService.instance) {
      LoginService.instance = new LoginService();
    }
    return LoginService.instance;
  }

  /**
   * Authenticate user with provided credentials
   * Automatically detects login strategy (email, google, otp, etc.)
   * 
   * @param credentials - Login credentials
   * @returns Promise<LoginResponse> - Response with token and user data
   * @throws Error if authentication fails
   */
  public async login(credentials: LoginRequest): Promise<LoginResponse> {
    const startTime = performance.now();
    
    try {
      console.log('🚀 Starting login process...');

      // Use factory to determine strategy and authenticate
      const apiStartTime = performance.now();
      const response = await this.loginFactory.login(credentials);
      const apiEndTime = performance.now();
      console.log(`⏱️ API call took: ${(apiEndTime - apiStartTime).toFixed(2)}ms`);

      if (!response.success || !response.token) {
        throw new Error('Login failed: Invalid response from server');
      }

      // Store token using AuthService
      const tokenStartTime = performance.now();
      this.authService.setToken(response.token);
      const tokenEndTime = performance.now();
      console.log(`⏱️ Token storage took: ${(tokenEndTime - tokenStartTime).toFixed(2)}ms`);

      const totalTime = performance.now() - startTime;
      console.log(`✅ Login successful, user authenticated`);
      console.log(`⏱️ TOTAL LOGIN TIME: ${totalTime.toFixed(2)}ms (${(totalTime / 1000).toFixed(2)}s)`);
      
      return response;

    } catch (error) {
      const totalTime = performance.now() - startTime;
      console.error(`❌ Login failed after ${totalTime.toFixed(2)}ms:`, error);
      // Clear any existing token on failed login
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

  /**
   * Validate current session with backend if needed
   * Checks if 5 minutes passed since last validation
   * Automatically clears token if validation fails
   * 
   * @returns Promise<boolean> - true if session is valid
   */
  public async validateSession(): Promise<boolean> {
    try {
      // First check if token exists and not expired (frontend check)
      if (!this.isAuthenticated()) {
        console.log('❌ No valid token found');
        return false;
      }

      // Check if backend validation is needed (5 minutes interval)
      if (this.tokenValidator.needsBackendValidation()) {
        console.log('🔄 Validating session with backend...');
        await this.tokenValidator.validateWithBackend();
        console.log('✅ Session validated with backend');
      }

      return true;
    } catch (error) {
      console.error('❌ Session validation failed:', error);
      // Clear invalid token
      this.authService.clearToken();
      return false;
    }
  }
}
