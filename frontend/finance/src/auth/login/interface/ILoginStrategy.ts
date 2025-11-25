import { LoginRequest, LoginResponse } from './ILogin';
export interface ILoginStrategy {
  /**
   * Execute login with provided credentials
   * @param credentials - Login credentials (varies by strategy)
   * @returns Promise<LoginResponse>
   */
  login(credentials: LoginRequest): Promise<LoginResponse>;
  
  /**
   * Get the name of this login strategy
   * @returns string - Strategy name (e.g., 'email', 'google', 'microsoft')
   */
  getStrategyName(...args: any[]): string;
}
