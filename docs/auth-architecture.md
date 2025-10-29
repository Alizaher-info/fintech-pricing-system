services/
├── auth-api/           # ← NEW: Dedicated Auth Service
│   ├── main.go
│   ├── Dockerfile
│   ├── handlers/
│   │   ├── login.go
│   │   ├── register.go
│   │   └── validate.go
│   └── jwt/
│       └── manager.go
├── pricing-api/        # ← Existing
└── user-api/          # ← Future: User management
