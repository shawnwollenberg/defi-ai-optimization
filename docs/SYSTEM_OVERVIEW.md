# System Overview

## High-Level Architecture

```mermaid
graph TB
    subgraph "User Interface"
        WEB[Web Browser<br/>React Dashboard]
        MOBILE[Future: Mobile App]
    end
    
    subgraph "API Layer"
        GATEWAY[API Gateway<br/>Port 8080<br/>• Authentication<br/>• Routing<br/>• Rate Limiting<br/>• WebSocket Hub]
    end
    
    subgraph "Business Logic Services"
        DEFI[DeFi Data Service<br/>Port 8081<br/>• Protocol Integration<br/>• APY Monitoring<br/>• Position Tracking]
        ML[ML Service<br/>Port 8001<br/>• Risk Forecasting<br/>• APY Prediction<br/>• Model Inference]
        WALLET[Wallet Service<br/>Port 8082<br/>• Transaction Building<br/>• Wallet Integration<br/>• Multi-chain Support]
        AUTO[Automation Engine<br/>Port 8083<br/>• Rule Monitoring<br/>• Trigger Evaluation<br/>• Action Execution]
    end
    
    subgraph "Data Layer"
        PG[(PostgreSQL<br/>Port 5432<br/>• Users<br/>• Portfolios<br/>• Rules<br/>• Transactions)]
        REDIS[(Redis<br/>Port 6379<br/>• APY Cache<br/>• Price Cache<br/>• Session Data)]
    end
    
    subgraph "External Integrations"
        ETH[Ethereum Mainnet<br/>RPC Provider]
        BASE[Base L2<br/>RPC Provider]
        STRIPE[Stripe<br/>Payment Processing]
        WC[WalletConnect<br/>Wallet Connections]
        PROTOCOLS[DeFi Protocols<br/>Aave, Compound, EigenLayer]
    end
    
    WEB -->|HTTPS| GATEWAY
    MOBILE -.->|Future| GATEWAY
    
    GATEWAY -->|HTTP| DEFI
    GATEWAY -->|HTTP| ML
    GATEWAY -->|HTTP| WALLET
    GATEWAY -->|HTTP| AUTO
    GATEWAY -->|WebSocket| WEB
    GATEWAY -->|Read/Write| PG
    GATEWAY -->|Cache| REDIS
    
    AUTO -->|HTTP| DEFI
    AUTO -->|HTTP| ML
    AUTO -->|HTTP| WALLET
    AUTO -->|Read/Write| PG
    AUTO -->|Cache| REDIS
    
    DEFI -->|RPC| ETH
    DEFI -->|RPC| BASE
    DEFI -->|Contract Calls| PROTOCOLS
    DEFI -->|Cache| REDIS
    
    WALLET -->|RPC| ETH
    WALLET -->|RPC| BASE
    WALLET -->|WebSocket| WC
    
    STRIPE -->|Webhooks| GATEWAY
    WC -->|Wallet Events| WEB
```

## Request Flow Example: User Views Portfolio

```mermaid
sequenceDiagram
    autonumber
    participant U as User Browser
    participant FE as Frontend
    participant API as API Gateway
    participant DEFI as DeFi Service
    participant ML as ML Service
    participant ETH as Ethereum RPC
    participant REDIS as Redis
    participant DB as PostgreSQL
    
    U->>FE: Navigate to Dashboard
    FE->>API: GET /portfolios<br/>(JWT: Bearer token)
    
    API->>API: Validate JWT
    API->>API: Extract user_id
    
    API->>DB: SELECT portfolios<br/>WHERE user_id = ?
    DB-->>API: Portfolio records
    
    par Fetch Live Data
        API->>DEFI: GET /protocols/aave/positions<br/>?user_address=0x...
        DEFI->>REDIS: Check cache
        alt Cache Miss
            DEFI->>ETH: eth_call (Aave Pool)
            ETH-->>DEFI: Position data
            DEFI->>REDIS: Cache (TTL: 30s)
        end
        DEFI-->>API: Positions
    and Risk Assessment
        API->>ML: POST /risk/forecast<br/>(positions, health_factor)
        ML->>ML: Run ML model
        ML-->>API: Risk: 0.15 (medium)
    end
    
    API->>DB: INSERT portfolio_snapshot
    API-->>FE: Complete response
    FE->>FE: Render dashboard
    FE-->>U: Display portfolio
```

## Component Interaction Matrix

| Component | Communicates With | Protocol | Purpose |
|-----------|------------------|----------|---------|
| Frontend | API Gateway | HTTP/WebSocket | User requests, real-time updates |
| API Gateway | All Services | HTTP | Request routing |
| API Gateway | PostgreSQL | SQL (via GORM) | Data persistence |
| API Gateway | Redis | Redis Protocol | Session caching |
| DeFi Service | Blockchain RPCs | JSON-RPC | Fetch protocol data |
| DeFi Service | Redis | Redis Protocol | Cache APY/price data |
| ML Service | Redis | Redis Protocol | Cache model inputs |
| Automation Engine | DeFi Service | HTTP | Check trigger conditions |
| Automation Engine | ML Service | HTTP | Get risk assessments |
| Automation Engine | Wallet Service | HTTP | Build transactions |
| Automation Engine | PostgreSQL | SQL (via GORM) | Store execution history |
| Wallet Service | Blockchain RPCs | JSON-RPC | Build/send transactions |
| Wallet Service | WalletConnect | WebSocket | Wallet connections |

## Data Flow Patterns

### Read Pattern (Cached)
```
User Request → API Gateway → Check Redis → Return Cached Data
                                    ↓ (miss)
                              Service → External API → Store in Redis → Return Data
```

### Write Pattern
```
User Action → API Gateway → Validate → Service → Database → WebSocket Broadcast → Frontend Update
```

### Automation Pattern
```
Background Loop → Check Rules → Evaluate Triggers → Execute Actions → Update Database → Notify User
```

## Technology Decisions

### Why Go for Backend?
- **Performance**: Fast compilation and execution
- **Concurrency**: Excellent goroutine support for handling multiple requests
- **Blockchain**: Strong ecosystem (go-ethereum)
- **Type Safety**: Compile-time error checking

### Why Python for ML?
- **ML Ecosystem**: Rich libraries (scikit-learn, pandas, numpy)
- **Rapid Development**: Easy model prototyping
- **FastAPI**: Modern async framework
- **Integration**: Easy to call from Go services

### Why React/TypeScript?
- **Type Safety**: Catch errors at compile time
- **Ecosystem**: Rich library ecosystem
- **Performance**: Virtual DOM and optimizations
- **Developer Experience**: Great tooling

### Why Microservices?
- **Scalability**: Scale services independently
- **Technology Flexibility**: Use best tool for each service
- **Fault Isolation**: One service failure doesn't bring down entire system
- **Team Structure**: Different teams can own different services

## Performance Characteristics

### Expected Latencies

| Operation | Typical Latency | Notes |
|-----------|----------------|-------|
| API Gateway (cached) | < 10ms | Redis cache hit |
| API Gateway (uncached) | 50-200ms | External service call |
| DeFi Data Fetch | 100-500ms | Blockchain RPC call |
| ML Inference | 50-150ms | Model prediction |
| Database Query | 5-20ms | Simple queries |
| WebSocket Message | < 5ms | Real-time push |

### Scalability Limits

- **API Gateway**: ~1000 req/s per instance
- **DeFi Service**: Limited by RPC provider rate limits
- **ML Service**: ~100 predictions/s per instance
- **Database**: PostgreSQL can handle 1000s of concurrent connections
- **Redis**: Very high throughput (100k+ ops/s)

## Security Architecture

### Authentication Flow
```
Wallet Connection → Message Signing → Signature Verification → JWT Generation → Token Storage
```

### Authorization
- JWT tokens contain user_id
- Middleware validates tokens on protected routes
- WebSocket connections require valid JWT

### Data Protection
- Database: Encrypted at rest (PostgreSQL)
- Network: HTTPS/TLS for all external communication
- Secrets: Environment variables (never in code)
- API Keys: Stored securely, rotated regularly

## Monitoring & Observability

### Health Checks
- All services expose `/health` endpoints
- Docker health checks configured
- Database connection monitoring

### Logging
- Structured logging in all services
- Docker logs: `docker compose logs [service]`
- Future: Centralized logging (ELK stack)

### Metrics (Future)
- Request rates
- Error rates
- Response times
- Cache hit rates
- Database query performance

