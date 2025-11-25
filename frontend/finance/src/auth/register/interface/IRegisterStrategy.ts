import { RegisterRequest, RegisterResponse } from './IRegister';

export interface IRegisterStrategy {
  /**
   * Execute registration with provided data
   * @param data - Registration data (varies by strategy)
   * @returns Promise<RegisterResponse>
   */
  register(data: RegisterRequest): Promise<RegisterResponse>;
  
  /**
   * Get the name of this registration strategy
   * @returns string - Strategy name (e.g., 'email', 'google', 'phone')
   */
  getStrategyName(...args: any[]): string;
}