import { AuthService } from '../auth/security/auth/AuthService';
import { apiService } from './ApiService';

/**
 * Logout Service
 * Handles complete logout flow: backend invalidation + frontend cleanup
 * Follows Single Responsibility Principle
 */
export class LogoutService {
  private static instance: LogoutService;
  private authService: AuthService;

  /**
   * Private constructor to enforce Singleton pattern
   */
  private constructor() {
    this.authService = AuthService.getInstance();
    console.log('🔐 LogoutService instance created');
  }

  /**
   * Get singleton instance of LogoutService
   * @returns LogoutService - The singleton instance
   */
  public static getInstance(): LogoutService {
    if (!LogoutService.instance) {
      LogoutService.instance = new LogoutService();
    }
    return LogoutService.instance;
  }

  /**
   * Logout from current device only
   * Invalidates token on backend and clears frontend state
   * 
   * @returns Promise<void>
   * @throws Error if backend logout fails
   */
  public async logoutCurrentDevice(): Promise<void> {
    try {
      const token = this.authService.getToken();

      if (!token) {
        console.log('⚠️ No token found, clearing local state only');
        this.authService.clearToken();
        return;
      }

      console.log('🔵 Logging out from current device...');

      // Call backend to invalidate token
      await apiService.post('/auth/logout', {}, token);

      console.log('✅ Backend logout successful');

      // Clear frontend token
      this.authService.clearToken();

      console.log('✅ Logout complete');

    } catch (error) {
      console.error('❌ Logout failed:', error);
      // Clear token even if backend call fails
      this.authService.clearToken();
      throw error;
    }
  }

  /**
   * Logout from all devices
   * Invalidates all user tokens on backend
   * 
   * @returns Promise<void>
   * @throws Error if backend logout fails
   */
  public async logoutAllDevices(): Promise<void> {
    try {
      const token = this.authService.getToken();

      if (!token) {
        console.log('⚠️ No token found, clearing local state only');
        this.authService.clearToken();
        return;
      }

      console.log('🔵 Logging out from all devices...');

      // Call backend to invalidate all tokens
      await apiService.post('/auth/logout/all', {}, token);

      console.log('✅ Backend logout from all devices successful');

      // Clear frontend token
      this.authService.clearToken();

      console.log('✅ Logout from all devices complete');

    } catch (error) {
      console.error('❌ Logout from all devices failed:', error);
      // Clear token even if backend call fails
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
}
