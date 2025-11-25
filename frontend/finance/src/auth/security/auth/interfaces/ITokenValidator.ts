/**
 * Token Validator Interface
 * Defines contract for token validation operations
 * Follows Interface Segregation Principle (ISP)
 */
export interface ITokenValidator {
  /**
   * Check if token is expired (frontend validation)
   * Decodes JWT and checks 'exp' claim
   * @returns boolean - true if token expired
   */
  isTokenExpired(): boolean;

  /**
   * Check if backend validation is needed
   * Returns true if last validation was more than threshold ago
   * @returns boolean - true if backend validation needed
   */
  needsBackendValidation(): boolean;

  /**
   * Validate token with backend API
   * Calls /validate endpoint
   * @returns Promise<{ success: boolean; user: any }>
   * @throws Error if validation fails
   */
  validateWithBackend(): Promise<{ success: boolean; user: any }>;

  /**
   * Mark that token was successfully validated with backend
   * Stores timestamp for smart caching
   */
  markValidated(): void;

  /**
   * Get the validation interval in milliseconds
   * @returns number - interval in ms (e.g., 5 minutes = 300000)
   */
  getValidationInterval(): number;
}
