import type { ITokenValidator } from './interfaces/ITokenValidator';
import { apiService } from '../../../services/ApiService';


export class TokenValidator implements ITokenValidator {
  private static instance: TokenValidator;
  private readonly VALIDATION_INTERVAL = 5 * 60 * 1000; // 5 minutes

  /**
   * Private constructor to enforce Singleton pattern
   */
  private constructor() {
    console.log('🔐 TokenValidator instance created');
  }

  /**
   * Get singleton instance of TokenValidator
   * @returns TokenValidator - The singleton instance
   */
  public static getInstance(): TokenValidator {
    if (!TokenValidator.instance) {
      TokenValidator.instance = new TokenValidator();
    }
    return TokenValidator.instance;
  }

  /**
   * Check if token is expired (frontend validation)
   * Decodes JWT and checks expiration timestamp
   * 
   * @returns boolean - true if token is expired
   */
  public isTokenExpired(): boolean {
    const token = this.getTokenFromStorage();
    if (!token) return true;

    try {
      // JWT format: header.payload.signature
      const payload = JSON.parse(atob(token.split('.')[1]));
      const expirationTime = payload.exp * 1000; // Convert to milliseconds
      const now = Date.now();
      
      const isExpired = expirationTime < now;
      
      if (isExpired) {
        console.log('⏰ Token expired at:', new Date(expirationTime).toISOString());
      }
      
      return isExpired;
    } catch (error) {
      console.error('❌ Failed to decode token:', error);
      return true; // Assume expired if can't decode
    }
  }

  /**
   * Check if we need to validate with backend
   * Returns true if last validation was more than 5 minutes ago
   * 
   * @returns boolean - true if backend validation needed
   */
  public needsBackendValidation(): boolean {
    const lastValidation = sessionStorage.getItem('last_token_validation');
    if (!lastValidation) {
      console.log('🔄 No previous validation found, backend check needed');
      return true;
    }

    const lastValidationTime = parseInt(lastValidation);
    const now = Date.now();
    const timeSinceValidation = now - lastValidationTime;
    const needsValidation = timeSinceValidation > this.VALIDATION_INTERVAL;
    
    if (needsValidation) {
      const minutesAgo = Math.floor(timeSinceValidation / 60000);
      console.log(`🔄 Last validated ${minutesAgo} minutes ago, backend check needed`);
    }
    
    return needsValidation;
  }

  /**
   * Validate token with backend API
   * Uses ApiService for backend communication
   * Calls /validate endpoint
   * 
   * @returns Promise<{ success: boolean; user: any }>
   * @throws Error if validation fails
   */
  public async validateWithBackend(): Promise<{ success: boolean; user: any }> {
    const token = this.getTokenFromStorage();
    
    if (!token) {
      throw new Error('No token found for validation');
    }

    console.log('🔄 Validating token with backend...');

    try {
      const result = await apiService.get('/validate', token);
      
      console.log('✅ Backend validation successful');
      this.markValidated();

      return result;
    } catch (error) {
      console.error('❌ Backend validation failed:', error);
      // Token will be cleared by AuthService when needed
      throw error;
    }
  }

  /**
   * Mark that token was successfully validated with backend
   * Stores timestamp for smart caching
   */
  public markValidated(): void {
    const timestamp = Date.now();
    sessionStorage.setItem('last_token_validation', timestamp.toString());
    console.log('✅ Validation timestamp saved:', new Date(timestamp).toISOString());
  }

  /**
   * Get the validation interval in milliseconds
   * @returns number - interval in ms (5 minutes = 300000)
   */
  public getValidationInterval(): number {
    return this.VALIDATION_INTERVAL;
  }

  /**
   * Get token from sessionStorage
   * @private
   */
  private getTokenFromStorage(): string | null {
    return sessionStorage.getItem('auth_token');
  }
}
