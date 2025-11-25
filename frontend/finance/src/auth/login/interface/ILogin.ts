export interface LoginRequest {
  email?: string;
  password?: string;
  token?: string;
  credential?: string; // Google ID token from Sign-In with Google
  otp?: string;
  biometricData?: any;
  oauthProvider?: string; // 'google' | 'microsoft' | 'facebook'
  [key: string]: any; // Allow additional fields
}

export interface LoginResponse {
  success: boolean;
  token: string;
  user: {
    id: number;
    email: string;
    firstName: string;
    lastName: string;
    roles: string[];
  };
}