export interface IAuthService {
  /**
   * Get current authentication token
   * @returns string | null - Token if authenticated, null otherwise
   */
  getToken(): string | null;

  /**
   * Check if user is currently authenticated
   * @returns boolean - true if user has valid token
   */
  isAuthenticated(): boolean;

  /**
   * Set authentication token
   * @param token - JWT token to store
   */
  setToken(token: string): void;

  /**
   * Clear authentication token
   */
  clearToken(): void;
}