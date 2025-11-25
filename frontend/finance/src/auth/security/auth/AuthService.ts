import { TokenValidator } from './TokenValidator';
import { IAuthService } from './interfaces/IAuthService';


export class AuthService implements IAuthService {
  private static instance: AuthService;
  private token: string | null = null;
  private tokenValidator: TokenValidator;

  /**
   * Private constructor to enforce Singleton pattern
   */
  private constructor() {
    this.token = sessionStorage.getItem('auth_token');
    this.tokenValidator = TokenValidator.getInstance();
    console.log('🔐 AuthService instance created');
  }

  /**
   * Get singleton instance of AuthService
   * @returns AuthService - The singleton instance
   */
  public static getInstance(): AuthService {
    if (!AuthService.instance) {
      AuthService.instance = new AuthService();
    }
    return AuthService.instance;
  }

  /**
   * Get current authentication token
   * @returns string | null - Token if authenticated, null otherwise
   */
  public getToken(): string | null {
    if (!this.token) {
      this.token = sessionStorage.getItem('auth_token');
    }
    return this.token;
  }

  /**
   * Check if user is currently authenticated
   * Checks both token existence and expiration
   * 
   * @returns boolean - true if user has valid token
   */
  public isAuthenticated(): boolean {
    const token = this.getToken();
    if (!token) return false;
    
    return !this.tokenValidator.isTokenExpired();
  }

  /**
   * Set authentication token
   * Stores in memory and sessionStorage
   * Called by LoginService or RegisterService after successful authentication
   */
  public setToken(token: string): void {
    this.token = token;
    sessionStorage.setItem('auth_token', token);
    console.log('✅ Token saved');
  }

  /**
   * Clear authentication token
   * Removes from memory and sessionStorage
   * Called by LogoutService
   */
  public clearToken(): void {
    this.token = null;
    sessionStorage.removeItem('auth_token');
    sessionStorage.removeItem('last_token_validation');
    console.log('✅ Token cleared');
  }
}