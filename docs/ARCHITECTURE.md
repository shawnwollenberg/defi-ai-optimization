# Architecture Documentation

## System Architecture Overview

The DeFi AI Optimization Platform is built as a microservices architecture with clear separation of concerns.

```mermaid
graph TB
    subgraph "Frontend Layer"
        FE[React/TypeScript Frontend<br/>Port 3000]
    end
    
    subgraph "API Gateway Layer"
        API[API Gateway<br/>Go - Port 8080<br/>Authentication, Routing, Rate Limiting]
    end
    
    subgraph "Core Services"
        DEFI[DeFi Data Service<br/>Go - Port 8081<br/>Protocol Integrations]
        ML[ML Service<br/>Python - Port 8001<br/>Risk Forecasting, APY Prediction]
        WALLET[Wallet Service<br/>Go - Port 8082<br/>Transaction Building]
        AUTO[Automation Engine<br/>Go - Port 8083<br/>Rule Execution]
    end
    
    subgraph "Data Layer"
        PG[(PostgreSQL<br/>Port 5432<br/>User Data, Portfolios, Rules)]
        REDIS[(Redis<br/>Port 6379<br/>Caching, Real-time Data)]
    end
    
    subgraph "External Services"
        ETH[Ethereum RPC<br/>Alchemy/Infura]
        BASE[Base RPC<br/>Alchemy/Infura]
        STRIPE[Stripe<br/>Payments]
        WC[WalletConnect<br/>Wallet Connections]
    end
    
    FE -->|HTTP/WebSocket| API
    API -->|HTTP| DEFI
    API -->|HTTP| ML
    API -->|HTTP| WALLET
    API -->|HTTP| AUTO
    API -->|Read/Write| PG
    API -->|Cache| REDIS
    
    AUTO -->|HTTP| DEFI
    AUTO -->|HTTP| ML
    AUTO -->|HTTP| WALLET
    AUTO -->|Read/Write| PG
    AUTO -->|Cache| REDIS
    
    DEFI -->|RPC Calls| ETH
    DEFI -->|RPC Calls| BASE
    DEFI -->|Cache| REDIS
    
    WALLET -->|RPC Calls| ETH
    WALLET -->|RPC Calls| BASE
    WALLET -->|WebSocket| WC
    
    API -->|Webhooks| STRIPE
    STRIPE -->|Webhooks| API
    
    FE -->|WalletConnect| WC
```

## Data Flow Diagrams

### User Authentication Flow

```mermaid
sequenceDiagram
    participant U as User
    participant FE as Frontend
    participant API as API Gateway
    participant DB as PostgreSQL
    
    U->>FE: Connect Wallet (MetaMask)
    FE->>FE: Sign Message
    FE->>API: POST /api/v1/auth/wallet<br/>(wallet_address, signature)
    API->>API: Verify Signature
    API->>DB: Create/Get User
    DB-->>API: User Data
    API->>API: Generate JWT Token
    API-->>FE: JWT Token + User Info
    FE->>FE: Store Token
    FE-->>U: Authenticated
```

### Portfolio Monitoring Flow

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant API as API Gateway
    participant DEFI as DeFi Service
    participant ML as ML Service
    participant ETH as Ethereum RPC
    participant REDIS as Redis
    
    FE->>API: GET /api/v1/portfolios<br/>(with JWT)
    API->>API: Validate JWT
    API->>DEFI: GET /api/v1/protocols/aave/positions<br/>?user_address=0x...
    DEFI->>REDIS: Check Cache
    alt Cache Hit
        REDIS-->>DEFI: Cached Data
    else Cache Miss
        DEFI->>ETH: Query Aave Contracts
        ETH-->>DEFI: Position Data
        DEFI->>REDIS: Store Cache
    end
    DEFI-->>API: Positions
    API->>ML: POST /api/v1/risk/forecast<br/>(positions, health_factor)
    ML->>ML: Run ML Model
    ML-->>API: Risk Assessment
    API-->>FE: Portfolio + Risk Data
    FE->>FE: Display Dashboard
```

### Automation Execution Flow

```mermaid
sequenceDiagram
    participant AUTO as Automation Engine
    participant DEFI as DeFi Service
    participant ML as ML Service
    participant WALLET as Wallet Service
    participant DB as PostgreSQL
    participant ETH as Ethereum
    
    loop Every 30 seconds
        AUTO->>DB: Fetch Enabled Rules
        DB-->>AUTO: Automation Rules
        
        AUTO->>DEFI: Check APY for Protocol
        DEFI-->>AUTO: Current APY: 3.5%
        
        alt APY Below Threshold
            AUTO->>ML: Forecast Risk
            ML-->>AUTO: Risk Assessment
            
            AUTO->>WALLET: Build Transaction<br/>(Rebalance: Aave → EigenLayer)
            WALLET-->>AUTO: Transaction Data
            
            AUTO->>DB: Create Transaction Record
            AUTO->>DB: Update Rule (execution_count++)
            
            Note over AUTO,ETH: Transaction sent to user<br/>for approval via WalletConnect
        end
    end
```

## Component Interaction Diagram

```mermaid
graph LR
    subgraph "Request Flow"
        A[User Request] --> B[API Gateway]
        B --> C{Authentication}
        C -->|Valid| D[Route to Service]
        C -->|Invalid| E[401 Unauthorized]
        D --> F[Service Processing]
        F --> G[Response]
    end
    
    subgraph "Service Communication"
        H[API Gateway] -->|HTTP| I[DeFi Service]
        H -->|HTTP| J[ML Service]
        H -->|HTTP| K[Wallet Service]
        L[Automation Engine] -->|HTTP| I
        L -->|HTTP| J
        L -->|HTTP| K
    end
    
    subgraph "Data Persistence"
        M[Services] -->|GORM| N[(PostgreSQL)]
        M -->|Cache| O[(Redis)]
    end
```

## Technology Stack by Layer

### Frontend
- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite
- **State Management**: Zustand + React Query
- **Blockchain**: ethers.js
- **Charts**: Recharts
- **Routing**: React Router

### Backend Services (Go)
- **HTTP Framework**: Gin
- **Database ORM**: GORM
- **Authentication**: JWT
- **WebSocket**: Gorilla WebSocket
- **Blockchain**: go-ethereum
- **Payments**: Stripe Go SDK

### ML Service (Python)
- **Framework**: FastAPI
- **ML Libraries**: scikit-learn
- **Data Processing**: pandas, numpy
- **Blockchain**: web3.py

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Message Queue**: (Future: RabbitMQ/NATS)

## Service Responsibilities

### API Gateway
- **Purpose**: Single entry point for all client requests
- **Responsibilities**:
  - Request routing to appropriate services
  - JWT authentication and authorization
  - Rate limiting
  - CORS handling
  - WebSocket hub for real-time updates
  - Stripe webhook handling

### DeFi Data Service
- **Purpose**: Aggregate data from DeFi protocols
- **Responsibilities**:
  - Monitor APYs across protocols (Aave, Compound, EigenLayer)
  - Track user positions and health factors
  - Fetch asset prices
  - Cache frequently accessed data
  - Support multiple chains (Ethereum, Base)

### ML Service
- **Purpose**: AI-powered risk and trend analysis
- **Responsibilities**:
  - Liquidation risk prediction
  - APY trend forecasting
  - Portfolio optimization recommendations
  - Model inference and serving

### Wallet Service
- **Purpose**: Handle wallet interactions and transactions
- **Responsibilities**:
  - WalletConnect/MetaMask integration
  - Transaction building and signing
  - Multi-chain transaction support
  - Gas estimation and optimization

### Automation Engine
- **Purpose**: Execute automated DeFi strategies
- **Responsibilities**:
  - Monitor automation rules
  - Evaluate trigger conditions
  - Execute rebalancing actions
  - Coordinate with other services
  - Track execution history

## Security Architecture

```mermaid
graph TB
    subgraph "Security Layers"
        A[Frontend] -->|HTTPS| B[API Gateway]
        B -->|JWT Validation| C[Authentication]
        C -->|Role-Based| D[Authorization]
        D -->|Rate Limiting| E[Service Access]
        E -->|Internal Network| F[Microservices]
    end
    
    subgraph "Data Protection"
        G[Encrypted at Rest<br/>PostgreSQL] 
        H[Encrypted in Transit<br/>TLS/HTTPS]
        I[Secrets Management<br/>Environment Variables]
    end
```

## Scalability Considerations

- **Horizontal Scaling**: Each service can be scaled independently
- **Caching Strategy**: Redis for frequently accessed data
- **Database**: PostgreSQL with connection pooling
- **Load Balancing**: (Future: Add nginx/HAProxy in front of API Gateway)
- **Message Queue**: (Future: For async job processing)

## Deployment Architecture

```mermaid
graph TB
    subgraph "Production Environment"
        LB[Load Balancer]
        LB --> API1[API Gateway Instance 1]
        LB --> API2[API Gateway Instance 2]
        
        API1 --> DEFI1[DeFi Service]
        API2 --> DEFI2[DeFi Service]
        
        DEFI1 --> PG[(PostgreSQL<br/>Primary)]
        DEFI2 --> PG
        
        DEFI1 --> REDIS[(Redis Cluster)]
        DEFI2 --> REDIS
        
        ML1[ML Service] --> REDIS
        ML2[ML Service] --> REDIS
    end
```

