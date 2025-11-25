/**
 * Registration Strategy Interface
 * Defines contract for all registration methods
 * Each registration strategy must implement this interface
 */

export interface RegisterRequest {
  email?: string;
  password?: string;
  firstName?: string;
  lastName?: string;
  phone?: string;
  otp?: string;
  token?: string;
  oauthProvider?: string;
  inviteCode?: string;
  [key: string]: any; // Allow additional fields
}

export interface RegisterResponse {
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