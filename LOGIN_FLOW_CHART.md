# Login Flow Chart

## Complete Authentication Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          USER INITIATES LOGIN                                │
│                         (LoginForm.vue Component)                            │
└──────────────────────────────────┬──────────────────────────────────────────┘
                                   │
                                   │ 1. User enters email + password
                                   ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         LoginService.login()                                 │
│                        (Singleton Service Layer)                             │
│  • Orchestrates the entire login process                                     │
│  • Delegates to LoginStrategyFactory                                         │
│  • Manages authentication state via AuthService                              │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                               │ 2. Pass credentials
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    LoginStrategyFactory.login()                              │
│                      (Strategy Auto-Detection)                               │
│  • Analyzes credential fields (email + password)                             │
│  • Determines strategy: 'email'                                              │
│  • Builds endpoint: /api/login/email                                         │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                               │ 3. POST /api/login/email
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           ApiService.post()                                  │
│                         (HTTP Client Layer)                                  │
│  • Adds /api prefix: /api/login/email                                        │
│  • Sends request to backend via Vite proxy                                   │
│  • Proxy forwards to nginx → PHP-FPM                                         │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                               │ 4. Request reaches backend
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Backend: LoginController.login()                          │
│                    Route: /api/login/{strategy}                              │
│  • Receives strategy = 'email'                                               │
│  • Routes to loginWithEmail()                                                │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                               │ 5. Email login handler
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│               LoginController.loginWithEmail()                               │
│                    (Main Authentication Logic)                               │
│                                                                               │
│  Step 1: Rate Limit Check                                                    │
│  ├─ RateLimiterService.isRateLimited(ip, userId)                            │
│  ├─ Check Redis: IP attempts (max 10) & User attempts (max 5)               │
│  └─ If locked → Return 429 Too Many Requests                                 │
│                                                                               │
│  Step 2: User Lookup                                                         │
│  ├─ UserRepository.findByEmail(email)                                        │
│  └─ If not found → Dispatch UserLoginFailedEvent → Return 401                │
│                                                                               │
│  Step 3: Password Verification                                               │
│  ├─ UserPasswordHasher.isPasswordValid(user, password)                      │
│  └─ If invalid → Dispatch UserLoginFailedEvent → Return 401                  │
│                                                                               │
│  Step 4: Success Flow                                                        │
│  ├─ Dispatch UserLoginSuccessEvent                                           │
│  ├─ JWTAuthenticator.generateToken(userId, email, role)                     │
│  └─ Return JSON: { token, user data, success: true }                         │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                               │ 6. Events dispatched
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SecurityEventListener (Background)                        │
│                                                                               │
│  On UserLoginSuccessEvent:                                                   │
│  ├─ RateLimiterService.clearAttempts(ip, userId)                            │
│  ├─ Update User.lastLoginAt                                                  │
│  └─ Create SecurityEvent record (LOGIN_SUCCESS)                              │
│                                                                               │
│  On UserLoginFailedEvent:                                                    │
│  ├─ RateLimiterService.recordFailedAttempt(ip, userId)                      │
│  ├─ Increment Redis counters                                                 │
│  └─ Create SecurityEvent record (LOGIN_FAILED with details)                  │
└───────────────────────────────────────────────────────────────────────────────┘
                               │
                               │ 7. Response returns to frontend
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    LoginService receives response                            │
│  • Validates response.success && response.token                              │
│  • Calls AuthService.setToken(token)                                         │
│  • Token stored in sessionStorage                                            │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                               │ 8. Login complete
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         LoginForm.vue                                        │
│  • Receives success response                                                 │
│  • Stores user data in sessionStorage (temp_user)                            │
│  • Redirects to /dashboard                                                   │
└─────────────────────────────────────────────────────────────────────────────┘
                               │
                               │ 9. Route guard check
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Router Guard (router/index.ts)                          │
│  • Calls LoginService.validateSession()                                      │
│  • Checks if token exists                                                    │
│  • Validates token structure (JWT decode)                                    │
│  • Checks if backend validation is needed (5-min cache)                      │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                               │ 10. If 5+ minutes passed
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    TokenValidator.validateWithBackend()                      │
│  • GET /api/validate with Bearer token                                       │
│  • Backend: LoginController.validate()                                       │
│  • JWTAuthenticator.requireAuth() - validates JWT                            │
│  • Returns { success: true, user: {...} }                                    │
│  • Frontend caches validation timestamp                                      │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                               │ 11. Access granted
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     User accesses dashboard                                  │
│  • AuthenticatedLayout.vue loads                                             │
│  • DashboardLayout.vue renders                                               │
│  • Protected content displayed                                               │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Key Components

### Frontend Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                          FRONTEND LAYERS                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  📱 Presentation Layer                                               │
│  ├─ LoginForm.vue          → UI Component                           │
│  ├─ RegisterForm.vue       → UI Component                           │
│  └─ AuthenticatedLayout.vue → Protected Route Wrapper               │
│                                                                      │
│  🔧 Service Layer (Singleton)                                        │
│  ├─ LoginService           → Orchestrates login                     │
│  ├─ RegisterService        → Orchestrates registration              │
│  └─ LogoutService          → Handles logout (single/all devices)    │
│                                                                      │
│  🏭 Factory Layer                                                    │
│  ├─ LoginStrategyFactory   → Auto-detects login strategy            │
│  └─ RegisterStrategyFactory → Auto-detects register strategy        │
│                                                                      │
│  🔐 Security Layer (Singleton)                                       │
│  ├─ AuthService            → Token management (get/set/clear)       │
│  └─ TokenValidator         → JWT decode + backend validation        │
│                                                                      │
│  🌐 HTTP Layer                                                       │
│  └─ ApiService             → HTTP client (GET/POST/PUT/DELETE)      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### Backend Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                          BACKEND LAYERS                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  🎯 Controller Layer                                                 │
│  ├─ LoginController        → /api/login/{strategy}, /api/validate   │
│  ├─ RegisterController     → /api/register/{strategy}               │
│  └─ LogoutController       → /api/auth/logout, /api/auth/logout/all │
│                                                                      │
│  🔒 Security Layer                                                   │
│  └─ JWTAuthenticator       → Token generation & validation          │
│                                                                      │
│  📊 Service Layer                                                    │
│  └─ RateLimiterService     → Redis-based rate limiting              │
│                                                                      │
│  🗄️ Repository Layer                                                 │
│  ├─ UserRepository         → User CRUD operations                   │
│  └─ SecurityEventRepository → Security audit logging                │
│                                                                      │
│  📡 Event Layer                                                      │
│  ├─ UserLoginSuccessEvent  → Dispatched on successful login         │
│  ├─ UserLoginFailedEvent   → Dispatched on failed login             │
│  └─ SecurityEventListener  → Handles security events                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Data Flow

### 1. Login Request Flow
```
User Input → LoginForm.vue → LoginService → LoginStrategyFactory 
→ ApiService → Vite Proxy → Nginx → PHP-FPM → LoginController 
→ RateLimiterService → UserRepository → Password Check 
→ JWTAuthenticator → SecurityEventListener → Response
```

### 2. Token Storage Flow
```
Backend JWT → LoginService → AuthService.setToken() 
→ sessionStorage → Available for all requests
```

### 3. Validation Flow (Every 5 Minutes)
```
Router Guard → LoginService.validateSession() 
→ TokenValidator.needsBackendValidation() → If true:
  → ApiService.get('/validate') → LoginController.validate() 
  → JWTAuthenticator.requireAuth() → Response
→ TokenValidator.markValidated() → Continue
```

## Security Features

### Rate Limiting
```
┌─────────────────────────────────────────────────────────┐
│             Rate Limiting (Redis-based)                  │
├─────────────────────────────────────────────────────────┤
│  IP-based:    10 attempts per 15 minutes                │
│  User-based:  5 attempts per 15 minutes                 │
│  Storage:     Redis with TTL                            │
│  On lock:     Return 429 with remaining time            │
└─────────────────────────────────────────────────────────┘
```

### JWT Token
```
┌─────────────────────────────────────────────────────────┐
│              JWT Token Structure                         │
├─────────────────────────────────────────────────────────┤
│  Issuer:      fintech-api                               │
│  Algorithm:   HS256                                      │
│  Expiration:  24 hours                                   │
│  Payload:     user_id, email, role                      │
│  Storage:     sessionStorage (frontend)                 │
└─────────────────────────────────────────────────────────┘
```

### Session Validation
```
┌─────────────────────────────────────────────────────────┐
│           5-Minute Validation Cycle                      │
├─────────────────────────────────────────────────────────┤
│  1. User logs in → Token cached                         │
│  2. < 5 minutes: Local validation only (JWT decode)     │
│  3. ≥ 5 minutes: Backend validation required            │
│  4. Backend validates → Cache timestamp updated         │
│  5. Repeat from step 2                                  │
└─────────────────────────────────────────────────────────┘
```

## Endpoints Summary

| Method | Endpoint                | Description                  | Handler                        |
|--------|-------------------------|------------------------------|--------------------------------|
| POST   | /api/login/email        | Email login                  | LoginController.loginWithEmail |
| POST   | /api/register/email     | Email registration           | RegisterController.register    |
| GET    | /api/validate           | Token validation             | LoginController.validate       |
| POST   | /api/auth/logout        | Logout current device        | LogoutController.logout        |
| POST   | /api/auth/logout/all    | Logout all devices           | LogoutController.logoutAll     |

## Error Handling

### Frontend
```javascript
try {
  await loginService.login(credentials)
  // Success → redirect to dashboard
} catch (error) {
  // Error → display to user
  // Token automatically cleared
}
```

### Backend
```php
try {
  // Authentication logic
  return new JsonResponse(['success' => true, 'token' => $token]);
} catch (Exception $e) {
  // Log event, return error response
  return new JsonResponse(['error' => $e->getMessage()], 401);
}
```

## Current Implementation Status

✅ **Implemented (Email Only)**
- Email login with rate limiting
- Email registration with auto-login
- JWT token generation & validation
- 5-minute backend validation caching
- Security event logging
- Session management
- Route guards

⏳ **Planned (Future)**
- OAuth login (Google, Microsoft, Facebook)
- OTP/SMS login
- Biometric login
- Phone registration
- Invite code registration
- Token blacklisting for logout
