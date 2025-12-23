# Quick Reference Guide

## Service Endpoints

### API Gateway (Port 8080)
- `GET /health` - Health check
- `POST /api/v1/auth/wallet` - Wallet authentication
- `GET /api/v1/portfolios` - Get user portfolios
- `POST /api/v1/automation/rules` - Create automation rule
- `GET /api/v1/ws` - WebSocket connection (authenticated)

### DeFi Service (Port 8081)
- `GET /api/v1/health` - Health check
- `GET /api/v1/protocols` - List protocols
- `GET /api/v1/protocols/:name/apy?asset=USDC&chain=ethereum` - Get APY
- `GET /api/v1/protocols/:name/positions?user_address=0x...` - Get positions
- `GET /api/v1/protocols/:name/health-factor?user_address=0x...` - Get health factor

### ML Service (Port 8001)
- `GET /health` - Health check
- `POST /api/v1/risk/forecast` - Risk prediction
- `POST /api/v1/apy/trend` - APY trend analysis

### Wallet Service (Port 8082)
- `GET /api/v1/health` - Health check
- `POST /api/v1/wallet/connect` - Connect wallet
- `POST /api/v1/wallet/build` - Build transaction

### Automation Engine (Port 8083)
- Runs in background, monitors rules every 30 seconds

## Common Commands

### Docker Management
```bash
# Start all services
docker compose up -d

# Stop all services
docker compose down

# View logs
docker compose logs -f [service-name]

# Restart a service
docker compose restart [service-name]

# Rebuild and restart
docker compose up -d --build [service-name]

# Check service status
docker compose ps
```

### Database Management
```bash
# Connect to PostgreSQL
docker compose exec postgres psql -U defi_user -d defi_optimization

# Run migrations
cd backend/shared && go run cmd/migrate/main.go

# Backup database
docker compose exec postgres pg_dump -U defi_user defi_optimization > backup.sql

# Restore database
docker compose exec -T postgres psql -U defi_user defi_optimization < backup.sql
```

### Development
```bash
# Run API Gateway locally
cd backend/api && go run main.go

# Run ML Service locally
cd ml-service && python app.py

# Run Frontend
cd frontend && npm run dev

# Train ML models
cd ml-service && python training/train_models.py
```

## Environment Variables Quick Reference

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/db` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379` |
| `JWT_SECRET` | Secret for JWT signing | `your-secret-key` |
| `ETH_RPC_URL` | Ethereum RPC endpoint | `https://eth-mainnet.g.alchemy.com/v2/...` |
| `BASE_RPC_URL` | Base RPC endpoint | `https://base-mainnet.g.alchemy.com/v2/...` |
| `WALLETCONNECT_PROJECT_ID` | WalletConnect project ID | `your-project-id` |
| `STRIPE_SECRET_KEY` | Stripe secret key | `sk_test_...` |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook secret | `whsec_...` |

## API Authentication

All protected endpoints require a JWT token in the Authorization header:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/portfolios
```

## WebSocket Events

The platform sends the following WebSocket message types:

- `portfolio_update` - Portfolio data changed
- `risk_alert` - Risk threshold exceeded
- `transaction_status` - Transaction status update
- `automation_triggered` - Automation rule executed

## Automation Rule Examples

### Example 1: APY Drop Rebalancing
```json
{
  "name": "Move to EigenLayer when Aave APY drops",
  "trigger_type": "apy_drop",
  "trigger_config": {
    "protocol": "aave",
    "asset": "USDC",
    "threshold": 4.0
  },
  "action_type": "rebalance",
  "action_config": {
    "from_protocol": "aave",
    "to_protocol": "eigenlayer",
    "asset": "USDC",
    "amount": 1000
  }
}
```

### Example 2: Health Factor Protection
```json
{
  "name": "Emergency withdrawal on low health",
  "trigger_type": "health_factor",
  "trigger_config": {
    "protocol": "aave",
    "threshold": 1.2
  },
  "action_type": "withdraw",
  "action_config": {
    "protocol": "aave",
    "asset": "USDC",
    "amount": 500
  }
}
```

## Troubleshooting Quick Fixes

| Issue | Solution |
|-------|----------|
| Service won't start | Check logs: `docker compose logs [service]` |
| Database connection error | Verify PostgreSQL is healthy: `docker compose ps postgres` |
| CORS errors | Check `FRONTEND_URL` in `.env` matches your frontend URL |
| WebSocket not connecting | Verify JWT token is valid and not expired |
| Automation not running | Check rule is enabled and automation service is running |

## Support Resources

- **Architecture**: See [ARCHITECTURE.md](ARCHITECTURE.md)
- **Usage Guide**: See [USAGE.md](USAGE.md)
- **Data Flows**: See [DATA_FLOW.md](DATA_FLOW.md)
- **Setup Help**: See [SETUP.md](SETUP.md)
- **Troubleshooting**: See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)

