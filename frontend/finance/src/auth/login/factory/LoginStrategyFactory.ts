import { apiService } from '../../../services/ApiService';
import { LoginRequest, LoginResponse } from '../interface/ILogin';
import { ILoginStrategy } from '../interface/ILoginStrategy';

/**
 * Generic Login Strategy Implementation
 * Uses endpoint URL to determine login method
 */
export class LoginStrategyFactory implements ILoginStrategy {
  private strategyName: string = 'email';
  private endpoint: string = '/login/email';

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    this.strategyName = this.getStrategyName(credentials);
    this.endpoint = `/login/${this.strategyName}`;
    
    console.log(`🔵 Login with ${this.strategyName}:`, credentials.email || 'token-based');
    const response = await apiService.post(this.endpoint, credentials);
    console.log(`✅ ${this.strategyName} login successful`);
    
    return response;
  }


  getStrategyName(...args: any[]): string {
    // If no args, return the stored strategy name
    if (args.length === 0) {
      return this.strategyName;
    }
    
    // Get credentials from arguments
    const credentials = args[0] as LoginRequest;
    
    // Get array of fields that have values
    const presentFields = Object.keys(credentials).filter(key => 
      credentials[key] !== undefined && credentials[key] !== null
    );
    
    // Determine strategy based on present fields (priority order)
    if (presentFields.includes('biometricData')) {
      return 'biometric';
    }
    
    // Check for Google OAuth credential (ID token from Google Sign-In)
    if (presentFields.includes('credential') && credentials.oauthProvider === 'google') {
      return 'google';
    }
    
    if (presentFields.includes('token')) {
      if (credentials.oauthProvider === 'google') {
        return 'google';
      }
      if (credentials.oauthProvider === 'microsoft') {
        return 'microsoft';
      }
      if (credentials.oauthProvider === 'facebook') {
        return 'facebook';
      }
      throw new Error(`Invalid OAuth provider: ${credentials.oauthProvider}`);
    }
    
    if (presentFields.includes('otp') && presentFields.includes('email')) {
      return 'otp';
    }
    
    if (presentFields.includes('email') && presentFields.includes('password')) {
      return 'email';
    }
    
    throw new Error(`Cannot determine login strategy from fields: ${presentFields.join(', ')}`);
  }
}
