# Contributing to Portfolio Finance

## Development Workflow

### 1. Feature Development
```bash
git checkout -b feature/your-feature-name
# Make your changes
git commit -m "feat: add new feature"
git push origin feature/your-feature-name
```

### 2. Pull Request Process
- Create PR against `develop` branch
- Ensure all CI checks pass
- Add tests for new features
- Update documentation if needed

### 3. CI/CD Pipeline

#### On Push to Feature Branch:
- Quick build tests
- Dockerfile linting
- Performance tests (if `[perf-test]` in commit message)

#### On Pull Request:
- Full code quality checks
- PHP and Go syntax validation
- Docker build tests
- Documentation checks

#### On Push to Main:
- Full test suite
- Security scanning
- Docker image building and publishing
- Production deployment

### 4. Local Development

#### Start Services:
```bash
cd deploy/compose
docker-compose up -d --build
```

#### Test Endpoints:
```bash
curl http://localhost:8080/quote-test
```

#### Stop Services:
```bash
docker-compose down -v
```

### 5. Code Quality Standards

#### PHP (Symfony):
- PSR-12 coding standards
- No direct service injection for gRPC clients
- Use Client/ namespace for external service clients

#### Go (gRPC):
- Standard Go formatting (`go fmt`)
- Proper error handling
- Unit tests for business logic

#### Docker:
- Multi-stage builds
- Minimal base images
- Health checks

### 6. Commit Message Format
```
type(scope): description

Types: feat, fix, docs, style, refactor, test, chore
```

Examples:
- `feat(pricing): add new rate calculation`
- `fix(grpc): resolve connection timeout`
- `docs(readme): update setup instructions`

### 7. Environment Variables

#### Required:
- `PRICING_GRPC_TARGET` - gRPC service endpoint

#### Optional:
- `APP_ENV` - Application environment (dev/prod)
- `DATABASE_URL` - Database connection string

### 8. Testing

#### Manual Testing:
```bash
# Health check
curl http://localhost:8080/quote-test

# Expected response:
{"interestRate":0.049,"apr":0.054,"monthlyPayment":457.5}
```

#### Automated Testing:
Tests run automatically in CI pipeline on every push/PR.

### 9. Troubleshooting

#### Common Issues:
1. **gRPC Connection Failed**: Check if pricing-api container is running
2. **Container Build Failed**: Ensure Docker has enough memory
3. **Port Conflicts**: Check if ports 8080/50051 are available

#### Debug Commands:
```bash
# Check container logs
docker-compose logs pricing-api
docker-compose logs php

# Check container status
docker-compose ps
```

### 10. Security

- All dependencies are scanned with Trivy
- No secrets in code (use environment variables)
- Regular security updates via Dependabot