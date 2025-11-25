import { apiService } from '../../../services/ApiService';
import { RegisterRequest, RegisterResponse } from '../interface/IRegister';
import { IRegisterStrategy } from '../interface/IRegisterStrategy';

/**
 * Generic Registration Strategy Implementation
 * Dynamically determines registration method from request data
 */
export class RegisterStrategyFactory implements IRegisterStrategy {
  private strategyName: string = 'email';
  private endpoint: string = '/register/email';

  async register(data: RegisterRequest): Promise<RegisterResponse> {
    this.strategyName = this.getStrategyName(data);
    this.endpoint = `/register/${this.strategyName}`;
    
    console.log(`🔵 Register with ${this.strategyName}:`, data.email || data.phone || 'token-based');
    const response = await apiService.post(this.endpoint, data);
    console.log(`✅ ${this.strategyName} registration successful`);
    
    return response;
  }

  getStrategyName(...args: any[]): string {
    // If no args, return the stored strategy name
    if (args.length === 0) {
      return this.strategyName;
    }
    
    // Get registration data from arguments
    const data = args[0] as RegisterRequest;
    
    // Get array of fields that have values
    const presentFields = Object.keys(data).filter(key => 
      data[key] !== undefined && data[key] !== null
    );
    
    // Determine strategy based on present fields (priority order)
    if (presentFields.includes('token')) {
      const validProviders = ['google', 'microsoft', 'facebook'];
      if (!data.oauthProvider || !validProviders.includes(data.oauthProvider)) {
        throw new Error(`Invalid OAuth provider: ${data.oauthProvider}`);
      }
      return data.oauthProvider;
    }
    
    if (presentFields.includes('phone') && presentFields.includes('otp')) {
      return 'phone';
    }
    
    if (presentFields.includes('email') && presentFields.includes('password')) {
      return 'email';
    }
    
    if (presentFields.includes('inviteCode')) {
      return 'invite';
    }
    
    throw new Error(`Cannot determine registration strategy from fields: ${presentFields.join(', ')}`);
  }
}
