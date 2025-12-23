# Data Flow Documentation

## End-to-End User Journey

### 1. User Registration & Authentication

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant API
    participant DB
    participant Wallet
    
    User->>Frontend: Opens app
    Frontend->>User: Shows "Connect Wallet"
    User->>Wallet: Approves connection
    Wallet->>Frontend: Returns wallet address
    Frontend->>User: Requests message signature
    User->>Wallet: Signs message
    Wallet->>Frontend: Returns signature
    Frontend->>API: POST /auth/wallet<br/>(address, signature, message)
    API->>API: Verify signature
    API->>DB: Find or create user
    DB-->>API: User record
    API->>API: Generate JWT (7-day expiry)
    API-->>Frontend: JWT token + user data
    Frontend->>Frontend: Store token in localStorage
    Frontend->>API: Establish WebSocket connection
    API-->>Frontend: WebSocket connected
```

### 2. Portfolio Discovery & Risk Assessment

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant API
    participant DEFI
    participant ML
    participant Blockchain
    participant Redis
    
    User->>Frontend: Views dashboard
    Frontend->>API: GET /portfolios<br/>(with JWT)
    API->>API: Validate JWT, extract user_id
    
    par Fetch Positions
        API->>DEFI: GET /protocols/aave/positions<br/>?user_address=0x...
        DEFI->>Redis: Check cache
        alt Cache Hit (< 30s old)
            Redis-->>DEFI: Cached positions
        else Cache Miss
            DEFI->>Blockchain: Query Aave Pool contract
            Blockchain-->>DEFI: User positions
            DEFI->>Redis: Cache positions (TTL: 30s)
        end
        DEFI-->>API: Positions data
    and Fetch APY Data
        API->>DEFI: GET /protocols/aave/apy?asset=USDC
        DEFI->>Blockchain: Query Aave interest rates
        Blockchain-->>DEFI: Current APY
        DEFI-->>API: APY: 4.5%
    end
    
    API->>ML: POST /risk/forecast<br/>(positions, health_factor)
    ML->>ML: Load ML model
    ML->>ML: Calculate risk features
    ML->>ML: Predict liquidation risk
    ML-->>API: Risk: 0.15 (medium)
    
    API->>DB: Save portfolio snapshot
    API-->>Frontend: Complete portfolio data
    Frontend->>Frontend: Render dashboard
    Frontend-->>User: Display portfolio + risk
```

### 3. Automation Rule Execution

```mermaid
sequenceDiagram
    participant AutoEngine
    participant DB
    participant DEFI
    participant ML
    participant WALLET
    participant User
    participant Blockchain
    
    loop Every 30 seconds
        AutoEngine->>DB: SELECT * FROM automation_rules<br/>WHERE enabled = true
        DB-->>AutoEngine: List of active rules
        
        loop For each rule
            AutoEngine->>DEFI: Check trigger condition<br/>(e.g., GET /apy)
            DEFI->>Blockchain: Query current APY
            Blockchain-->>DEFI: APY: 3.8%
            DEFI-->>AutoEngine: APY data
            
            alt Trigger Condition Met (APY < 4.0%)
                AutoEngine->>ML: Assess risk of action
                ML-->>AutoEngine: Risk assessment
                
                AutoEngine->>WALLET: Build transaction<br/>(rebalance: Aave → EigenLayer)
                WALLET->>WALLET: Estimate gas
                WALLET->>WALLET: Build transaction data
                WALLET-->>AutoEngine: Transaction ready
                
                AutoEngine->>DB: Create transaction record<br/>(status: pending)
                AutoEngine->>DB: Update rule<br/>(execution_count++, last_executed_at)
                
                AutoEngine->>User: Send transaction for approval<br/>(via WalletConnect)
                
                alt User Approves
                    User->>Blockchain: Sign & broadcast transaction
                    Blockchain-->>WALLET: Transaction hash
                    WALLET->>DB: Update transaction<br/>(status: confirmed, tx_hash)
                    AutoEngine->>User: Notification: "Rebalancing complete"
                else User Rejects
                    WALLET->>DB: Update transaction<br/>(status: rejected)
                    AutoEngine->>User: Notification: "Transaction cancelled"
                end
            end
        end
    end
```

### 4. Real-time Updates via WebSocket

```mermaid
sequenceDiagram
    participant Frontend
    participant API
    participant WebSocketHub
    participant DEFI
    participant Blockchain
    
    Frontend->>API: WebSocket connection<br/>(with JWT)
    API->>WebSocketHub: Register client
    WebSocketHub-->>Frontend: Connection established
    
    loop Background Monitoring
        DEFI->>Blockchain: Poll for changes<br/>(every 10s)
        Blockchain-->>DEFI: Updated APY/positions
        
        alt Significant Change Detected
            DEFI->>WebSocketHub: Broadcast update
            WebSocketHub->>Frontend: Message: portfolio_update
            Frontend->>Frontend: Update UI in real-time
        end
    end
    
    alt Risk Alert
        DEFI->>WebSocketHub: Risk threshold exceeded
        WebSocketHub->>Frontend: Message: risk_alert
        Frontend->>Frontend: Show warning notification
    end
```

## Service Communication Patterns

### Synchronous HTTP Calls

```mermaid
graph LR
    A[API Gateway] -->|HTTP REST| B[DeFi Service]
    A -->|HTTP REST| C[ML Service]
    A -->|HTTP REST| D[Wallet Service]
    E[Automation] -->|HTTP REST| B
    E -->|HTTP REST| C
    E -->|HTTP REST| D
```

### Asynchronous Processing

```mermaid
graph TB
    A[User Action] --> B[API Gateway]
    B --> C[Create Automation Rule]
    C --> D[(Database)]
    E[Automation Engine] -->|Poll Every 30s| D
    E -->|When Triggered| F[Execute Action]
    F --> G[Update Database]
    G --> H[Notify User via WebSocket]
```

## Data Storage Patterns

### Caching Strategy

```mermaid
graph TB
    A[Service Request] --> B{Cache Hit?}
    B -->|Yes| C[Return Cached Data]
    B -->|No| D[Fetch from Source]
    D --> E[Store in Redis]
    E --> F[Return Data]
    
    G[Cache Invalidation] -->|TTL: 30s| E
    H[Manual Invalidate] -->|On Updates| E
```

### Database Schema Relationships

```mermaid
erDiagram
    USER ||--o{ PORTFOLIO : has
    USER ||--o{ AUTOMATION_RULE : creates
    USER ||--o{ TRANSACTION : executes
    USER ||--|| SUBSCRIPTION : has
    
    PORTFOLIO ||--o{ POSITION : contains
    PORTFOLIO ||--o{ PORTFOLIO_SNAPSHOT : tracks
    
    AUTOMATION_RULE ||--o{ TRANSACTION : triggers
```

## Error Handling Flow

```mermaid
graph TB
    A[Request] --> B{Valid?}
    B -->|No| C[Return 400 Bad Request]
    B -->|Yes| D[Process Request]
    D --> E{Success?}
    E -->|Yes| F[Return 200 OK]
    E -->|No| G{Error Type?}
    G -->|Network| H[Retry 3x]
    G -->|Validation| I[Return 400]
    G -->|Not Found| J[Return 404]
    G -->|Server| K[Return 500]
    H -->|Still Fails| K
    K --> L[Log Error]
    L --> M[Notify User]
```

