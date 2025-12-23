# User Guide - How to Use the DeFi Optimization Platform

## Getting Started

### Prerequisites

1. **Docker and Docker Compose** installed
2. **Node.js 18+** and npm (for frontend)
3. **Go 1.21+** (for local development)
4. **Python 3.11+** (for ML service development)

### Initial Setup

1. **Clone and configure environment:**
   ```bash
   cd DefiOptimization
   cp .env.example .env
   # Edit .env with your API keys (see SETUP.md)
   ```

2. **Start all services:**
   ```bash
   docker compose up -d
   ```

3. **Verify services are running:**
   ```bash
   docker compose ps
   # All services should show "Up" status
   ```

4. **Start the frontend:**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

5. **Open your browser:**
   ```
   http://localhost:3000
   ```

## User Workflows

### 1. Connecting Your Wallet

**Step 1:** Open the application in your browser

**Step 2:** Click "Connect Wallet" in the top right

**Step 3:** Select your wallet (MetaMask, WalletConnect, etc.)

**Step 4:** Approve the connection and sign the authentication message

**Step 5:** You're now authenticated! Your wallet address will appear in the header.

### 2. Viewing Your Portfolio

**Step 1:** Navigate to the **Dashboard** (home page)

**Step 2:** The dashboard displays:
   - Total portfolio value in USD
   - Health factor across all positions
   - Liquidation risk assessment
   - Portfolio value trend chart

**Step 3:** Navigate to **Portfolios** to see detailed breakdown:
   - Individual portfolio positions
   - Per-protocol allocations
   - APY for each position
   - Health factors per protocol

### 3. Setting Up Automation Rules

**Step 1:** Navigate to **Automation** page

**Step 2:** Click **"Create Rule"**

**Step 3:** Configure your automation:

   **Example: Auto-Rebalance on APY Drop**
   - **Name**: "Move from Aave to EigenLayer when APY drops"
   - **Trigger Type**: `apy_drop`
   - **Trigger Config**:
     ```json
     {
       "protocol": "aave",
       "asset": "USDC",
       "chain": "ethereum",
       "threshold": 4.0
     }
     ```
   - **Action Type**: `rebalance`
   - **Action Config**:
     ```json
     {
       "from_protocol": "aave",
       "to_protocol": "eigenlayer",
       "asset": "USDC",
       "amount": 1000
     }
     ```

**Step 4:** Toggle the rule **ON** to activate it

**Step 5:** The automation engine will monitor every 30 seconds and execute when conditions are met

### 4. Monitoring Risk

**Step 1:** The **Dashboard** automatically shows your risk assessment

**Step 2:** Risk levels:
   - **Low** (green): Health factor > 1.5, low liquidation risk
   - **Medium** (yellow): Health factor 1.3-1.5, moderate risk
   - **High** (orange): Health factor 1.1-1.3, high risk
   - **Critical** (red): Health factor < 1.1, immediate action needed

**Step 3:** Review recommendations:
   - The ML service provides actionable recommendations
   - Examples: "Consider reducing leverage", "Add collateral"

**Step 4:** Set up automation rules to automatically respond to high risk:
   - Trigger: `risk_threshold` with threshold: `0.5`
   - Action: `withdraw` or `rebalance` to safer positions

### 5. Managing Subscriptions

**Step 1:** Navigate to **Settings**

**Step 2:** View your current subscription tier:
   - **Free**: Basic monitoring only
   - **Basic ($10/month)**: Basic automation, risk monitoring, email alerts
   - **Premium ($50/month)**: Advanced automation, AI risk forecasting, priority support

**Step 3:** To upgrade:
   - Click **"Upgrade Plan"**
   - Select your desired tier
   - Complete payment via Stripe checkout
   - Your subscription activates immediately

## API Usage Examples

### Authentication

```bash
# 1. Connect wallet and get token
curl -X POST http://localhost:8080/api/v1/auth/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0x...",
    "signature": "0x...",
    "message": "Sign in to DeFi Optimizer..."
  }'

# Response:
# {
#   "token": "eyJhbGciOiJIUzI1NiIs...",
#   "user": { "id": 1, "wallet_address": "0x..." }
# }
```

### Get Portfolio Data

```bash
# 2. Get user portfolios
curl -X GET http://localhost:8080/api/v1/portfolios \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response:
# [
#   {
#     "id": 1,
#     "name": "Main Portfolio",
#     "total_value_usd": 50000.00,
#     "health_factor": 1.75,
#     "positions": [...]
#   }
# ]
```

### Create Automation Rule

```bash
# 3. Create automation rule
curl -X POST http://localhost:8080/api/v1/automation/rules \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Auto-rebalance on APY drop",
    "description": "Move funds when Aave APY drops below 4%",
    "enabled": true,
    "trigger_type": "apy_drop",
    "trigger_config": {
      "protocol": "aave",
      "asset": "USDC",
      "chain": "ethereum",
      "threshold": 4.0
    },
    "action_type": "rebalance",
    "action_config": {
      "from_protocol": "aave",
      "to_protocol": "eigenlayer",
      "asset": "USDC",
      "amount": 1000
    }
  }'
```

### Get Risk Forecast

```bash
# 4. Get risk assessment (via API Gateway)
curl -X POST http://localhost:8080/api/v1/risk/forecast \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_address": "0x...",
    "positions": [
      {
        "protocol": "aave",
        "asset": "USDC",
        "amount": 10000,
        "apy": 4.5
      }
    ],
    "health_factor": 1.5,
    "total_collateral": 15000,
    "total_debt": 5000
  }'

# Response:
# {
#   "liquidation_risk": 0.15,
#   "risk_level": "medium",
#   "recommendations": [
#     "Monitor health factor regularly",
#     "Consider rebalancing if APY drops"
#   ],
#   "confidence": 0.85
# }
```

## Automation Rule Types

### 1. APY Drop Trigger

Triggers when APY falls below a threshold:

```json
{
  "trigger_type": "apy_drop",
  "trigger_config": {
    "protocol": "aave",
    "asset": "USDC",
    "chain": "ethereum",
    "threshold": 4.0
  }
}
```

### 2. Health Factor Trigger

Triggers when health factor drops below threshold:

```json
{
  "trigger_type": "health_factor",
  "trigger_config": {
    "protocol": "aave",
    "chain": "ethereum",
    "threshold": 1.3
  }
}
```

### 3. Risk Threshold Trigger

Triggers when ML-predicted risk exceeds threshold:

```json
{
  "trigger_type": "risk_threshold",
  "trigger_config": {
    "threshold": 0.5
  }
}
```

## Action Types

### 1. Rebalance Action

Moves funds between protocols:

```json
{
  "action_type": "rebalance",
  "action_config": {
    "from_protocol": "aave",
    "to_protocol": "eigenlayer",
    "asset": "USDC",
    "amount": 1000
  }
}
```

### 2. Withdraw Action

Withdraws funds from a protocol:

```json
{
  "action_type": "withdraw",
  "action_config": {
    "protocol": "aave",
    "asset": "USDC",
    "amount": 500
  }
}
```

### 3. Deposit Action

Deposits funds into a protocol:

```json
{
  "action_type": "deposit",
  "action_config": {
    "protocol": "eigenlayer",
    "asset": "ETH",
    "amount": 1.0
  }
}
```

## Real-time Updates

The platform uses WebSockets for real-time updates:

```javascript
// Frontend WebSocket connection
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?token=YOUR_JWT');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch(message.type) {
    case 'portfolio_update':
      // Update portfolio display
      break;
    case 'risk_alert':
      // Show risk warning
      break;
    case 'transaction_status':
      // Update transaction status
      break;
  }
};
```

## Best Practices

### 1. Start Conservative
   - Begin with high thresholds for automation rules
   - Monitor for a few days before lowering thresholds
   - Test with small amounts first

### 2. Monitor Regularly
   - Check dashboard daily
   - Review automation rule executions
   - Adjust rules based on market conditions

### 3. Risk Management
   - Set up health factor alerts
   - Use risk threshold triggers for safety
   - Diversify across multiple protocols

### 4. Gas Optimization
   - Batch operations when possible
   - Use automation during low gas periods
   - Consider Layer 2 (Base) for lower fees

## Troubleshooting

### Services Not Starting

```bash
# Check logs
docker compose logs api
docker compose logs defi-service

# Restart specific service
docker compose restart api

# Rebuild and restart
docker compose up -d --build api
```

### Database Connection Issues

```bash
# Check PostgreSQL is healthy
docker compose ps postgres

# View database logs
docker compose logs postgres

# Restart database
docker compose restart postgres
```

### Frontend Not Connecting

1. Verify API Gateway is running: `curl http://localhost:8080/health`
2. Check CORS settings in `.env` (FRONTEND_URL)
3. Verify JWT token is being sent in requests
4. Check browser console for errors

### Automation Not Executing

1. Verify rule is enabled (toggle ON)
2. Check automation engine logs: `docker compose logs automation`
3. Verify trigger conditions are met
4. Check that user has sufficient funds for the action

## Advanced Usage

### Custom ML Models

Train your own models:

```bash
cd ml-service
python training/train_models.py
```

Models are saved to `ml-service/models/` and automatically loaded on service start.

### Direct Service Access

Access services directly (bypassing API Gateway):

```bash
# DeFi Service
curl http://localhost:8081/api/v1/protocols/aave/apy?asset=USDC&chain=ethereum

# ML Service
curl -X POST http://localhost:8001/api/v1/risk/forecast \
  -H "Content-Type: application/json" \
  -d '{...}'
```

### Database Queries

Connect to PostgreSQL:

```bash
docker compose exec postgres psql -U defi_user -d defi_optimization

# Example queries
SELECT * FROM users;
SELECT * FROM automation_rules WHERE enabled = true;
SELECT * FROM portfolios ORDER BY created_at DESC;
```

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review logs: `docker compose logs [service-name]`
3. See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues

