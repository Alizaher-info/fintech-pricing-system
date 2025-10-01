# Portfolio Finance - Microservices System 🚀

[![System Status](https://img.shields.io/badge/Status-Production%20Ready-brightgreen)](https://github.com/yourusername/portfolio-finance)
[![Architecture](https://img.shields.io/badge/Architecture-Microservices-blue)](https://github.com/yourusername/portfolio-finance)
[![Tech Stack](https://img.shields.io/badge/Stack-Go%20%2B%20PHP%20%2B%20gRPC-orange)](https://github.com/yourusername/portfolio-finance)

A production-ready microservices system for financial calculations using **Go gRPC** backend with **PHP Symfony** HTTP gateway.

## 🏗️ **Architecture**

```
HTTP Request → nginx → PHP Symfony → gRPC Client → Go Service → Response
```

### **Services:**
- **pricing-api** (Go 1.22) - gRPC server on port 50051
- **backend** (PHP 8.2 + Symfony 7.3) - HTTP gateway on port 80  
- **nginx** - Reverse proxy on port 8080
- **database** - MySQL 8.0 on port 3307

## 🎯 **Key Features**

- ✅ **gRPC Communication** - High-performance inter-service communication
- ✅ **Factory Pattern** - Elegant solution for gRPC + Symfony DI challenges
- ✅ **Docker Compose** - Full containerization with multi-service orchestration
- ✅ **Protocol Buffers** - Type-safe service contracts
- ✅ **Clean Architecture** - Proper separation of concerns

## 🚀 **Quick Start**

### **Prerequisites:**
- Docker & Docker Compose
- Git

### **Run the System:**

```bash
# Clone repository
git clone https://github.com/yourusername/portfolio-finance.git
cd portfolio-finance

# Start all services
cd deploy/compose
docker-compose up -d

# Test the system
curl http://localhost:8080/quote-test
```

**Expected Response:**
```json
{"interestRate":0.049,"apr":0.054,"monthlyPayment":457.5}
```

## 📁 **Project Structure**

```
portfolio-finance/
├── backend/                    # PHP Symfony HTTP Gateway
│   ├── src/
│   │   ├── Client/            # gRPC Client (NOT Symfony Services!)
│   │   │   ├── PricingGrpcClient.php
│   │   │   └── PricingGrpcClientFactory.php
│   │   └── Controller/
│   │       └── QuoteTestController.php
│   └── config/
│       └── services.yaml
├── services/
│   ├── pricing-api/           # Go gRPC Service
│   │   ├── main.go
│   │   └── Dockerfile
│   ├── event-gateway/         # Future: Event handling
│   └── offers-worker/         # Future: Offer processing
├── protos/
│   └── pricing/v1/
│       └── pricing.proto      # Service contracts
└── deploy/
    └── compose/
        └── docker-compose.yml
```

## 🔧 **Technical Highlights**

### **gRPC + Symfony Integration Challenge:**

**Problem:** Symfony's Dependency Injection compiles at build-time, but gRPC protobuf classes are autoloaded at runtime.

**Solution:** Factory Pattern with static methods:

```php
// ❌ This fails - DI can't resolve at compile-time
public function __construct(PricingGrpcClient $client) {}

// ✅ This works - Factory creates at runtime  
$client = PricingGrpcClientFactory::createFromEnv();
```

### **Clean Client Architecture:**

```php
// Clear naming - not in Service/ folder (no DI confusion)
use App\Client\PricingGrpcClientFactory;

// Environment-based configuration
$client = PricingGrpcClientFactory::createFromEnv();
$result = $client->quote(10000.0, 24, null);
```

## 🛠️ **Development**

### **Environment Variables:**
```bash
# .env
PRICING_GRPC_TARGET=pricing-api:50051
```

### **Code Standards & Quality**

This project enforces consistent code style across all languages:

#### **PHP Standards (PSR-12 + Symfony)**
```bash
# Check code style (dry-run)
cd backend
vendor/bin/php-cs-fixer fix --dry-run --diff

# Fix code style automatically
vendor/bin/php-cs-fixer fix

# Configuration files:
# .php-cs-fixer.php       → Active project configuration
# .php-cs-fixer.dist.php  → Symfony standard template
```

**Key PHP Rules:**
- ✅ PSR-12 coding standards
- ✅ Symfony coding conventions  
- ✅ Short array syntax `[]` instead of `array()`
- ✅ Alphabetical import ordering
- ✅ Automatic unused import removal
- ✅ Excludes `src/Grpc/` (generated protobuf files)

#### **Go Standards (fmt + vet)**
```bash
cd services/pricing-api

# Format Go code
go fmt ./...

# Static analysis & error detection
go vet ./...

# Run tests
go test ./...
```

**Go Standards:**
- ✅ Standard Go formatting (`gofmt`)
- ✅ Static analysis (`go vet`)
- ✅ Proper error handling patterns
- ✅ Unit tests for business logic

### **Protobuf Generation:**
```bash
# Generate Go stubs
protoc --go_out=. --go-grpc_out=. protos/pricing/v1/pricing.proto

# Generate PHP stubs  
protoc --php_out=backend/src/Grpc --grpc_out=backend/src/Grpc protos/pricing/v1/pricing.proto
```

## 📊 **Performance**

- **Response Time:** ~50ms HTTP to gRPC roundtrip
- **Throughput:** Handles concurrent requests efficiently
- **Reliability:** Stateless microservice design

## 🎓 **Lessons Learned**

1. **gRPC + Symfony DI incompatibility** - Solved with Factory Pattern
2. **Compile-time vs Runtime** class loading - Use static factories
3. **Clear folder structure** - `/Client` vs `/Service` for team clarity
4. **Environment-based configuration** - Flexible deployment

## 🔄 **CI/CD Workflows**

This project uses a **3-tier workflow system** for optimal development efficiency and production safety.

### **📋 Workflow Overview**

| Workflow | Trigger | Duration | Purpose |
|----------|---------|----------|---------|
| `development.yml` | Push to `feature/*` | 2-3 min | Fast developer feedback |
| `pr-checks.yml` | Pull Request created | 5-7 min | Quality gatekeeper |
| `ci-cd.yml` | Push to `main/develop` | 10-15 min | Production deployment |

### **🎯 Why 3 Separate Workflows?**

#### **Efficiency Comparison:**
```
❌ Single Workflow (naive approach):
20 feature pushes × 15 min = 5 hours CI time per day

✅ Multi-tier Workflows (our approach):  
20 feature pushes × 2 min + 1 PR × 7 min + 1 merge × 15 min = 62 minutes
63% time savings! 🚀
```

### **1️⃣ development.yml - Developer Feedback**

**Triggers:** Push to `feature/*` branches  
**Purpose:** Lightning-fast syntax validation for active development

```yaml
Jobs:
├── quick-test (30s)
│   ├── Docker compose syntax check
│   └── Basic configuration validation
├── lint-dockerfile (15s)  
│   └── Dockerfile best practices
└── performance-test (optional)
    └── Triggered by commit message [perf-test]
```

**Developer Experience:**
```bash
git push origin feature/new-pricing
# ⚡️ 2 minutes later: "All checks passed!" 
# Continue coding with confidence
```

### **2️⃣ pr-checks.yml - Quality Gatekeeper**

**Triggers:** Pull Request to `main` or `develop`  
**Purpose:** Comprehensive quality validation before merge

```yaml
Jobs:
├── code-quality (3-4 min)
│   ├── PHP syntax check + Composer validation
│   ├── Go formatting (go fmt) + static analysis (go vet)
│   └── Go unit tests (go test)
├── docker-build-test (2-3 min)
│   ├── Test Docker builds (without push)
│   └── Container runtime validation  
└── documentation (30s)
    ├── README.md presence check
    └── API documentation validation
```

**Code Review Process:**
```
1. Developer creates PR → Automatic pr-checks.yml
2. All checks must pass ✅ before review
3. Reviewer sees "All PR checks passed" 
4. Focus on business logic, not syntax errors
5. Merge blocked if any check fails ❌
```

### **3️⃣ ci-cd.yml - Production Pipeline**

**Triggers:** Push to `main` (production) or `develop` (staging)  
**Purpose:** Full integration testing, security scanning, and deployment

```yaml
Jobs:
├── test (5 min)
│   ├── Full docker-compose integration test
│   ├── HTTP endpoint validation (curl tests)
│   └── gRPC service health checks
├── build-and-push (4 min, parallel)
│   ├── Multi-service Docker builds (pricing-api + backend)
│   ├── Container registry push (ghcr.io)
│   └── GitHub Actions cache optimization
├── security-scan (3 min, parallel)
│   ├── Trivy vulnerability scanner
│   ├── Dependency security analysis
│   └── SARIF upload to GitHub Security tab
├── deploy-staging (develop branch)
│   └── Automatic staging deployment
└── deploy-production (main branch)
    └── Production deployment (after all checks pass)
```

### **🏗️ Complete Development Lifecycle**

#### **Day 1: Feature Development**
```bash
# Developer workflow
git checkout -b feature/enterprise-pricing
# Code changes...
git commit -m "Add enterprise pricing tiers"
git push origin feature/enterprise-pricing

🚀 development.yml runs automatically
├── ✅ Syntax check (30s)
├── ✅ Docker lint (15s)  
└── 📧 "Checks passed" notification
```

#### **Day 2: Pull Request**
```bash
# Create PR via GitHub UI
feature/enterprise-pricing → main

🚀 pr-checks.yml runs automatically  
├── ✅ PHP + Go syntax & tests (4 min)
├── ✅ Docker build test (3 min)
├── ✅ Documentation check (30s)
└── 🔒 "Ready for review" status
```

#### **Day 3: Production Deployment**
```bash
# After code review + approval
git merge pull-request

🚀 ci-cd.yml runs automatically
├── ✅ Integration tests (5 min)
├── ✅ Security scan (3 min, parallel)  
├── ✅ Build + push images (4 min, parallel)
├── ✅ Deploy to production (2 min)
└── 🌍 Feature is LIVE!
```

### **🛡️ Security & Quality Gates**

#### **Branch Protection:**
- ✅ All PR checks must pass before merge
- ✅ Security scan blocks production deployment
- ✅ No direct pushes to main branch
- ✅ Required code review approval

#### **Container Registry:**
```
Only stable, tested images reach production:
├── ghcr.io/hp/portfolio-finance-pricing-api:latest
└── ghcr.io/hp/portfolio-finance-backend:latest
```

#### **Environment Promotion:**
```
feature → develop → staging → main → production
   ↓        ↓         ↓        ↓         ↓
 dev.yml  pr+ci.yml  🚀    ci-cd.yml   🌍
```

### **⚡️ Performance Benefits**

| Metric | Single Workflow | Multi-tier Workflows |
|--------|----------------|---------------------|
| Developer feedback | 15 min | 2 min (87% faster) |
| Daily CI cost | 300 min | 62 min (79% savings) |
| Failed merge prevention | After merge | Before merge |
| Developer productivity | Low | High |

### **📈 GitHub Actions Usage**

```bash
# View workflow runs
gh run list

# Check specific workflow
gh run view --job=code-quality

# Monitor costs
# Settings → Billing → Actions usage
```

## 🚀 **Production Ready**

- ✅ Containerized with Docker
- ✅ Environment configuration
- ✅ Error handling & graceful degradation
- ✅ Clean, maintainable architecture
- ✅ Type-safe service contracts
- ✅ **3-tier CI/CD pipeline with security gates**
- ✅ **Automated quality assurance**
- ✅ **Production deployment automation**



## 🤝 **Contributing**

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request



---

**Built with Ali Zaher using Go, PHP, gRPC, and Docker**