# DeFi AI Optimization Platform

An AI-powered DeFi optimization and risk management platform that automates portfolio rebalancing, forecasts risks with machine learning, and integrates with wallets for one-click execution on Ethereum and Base.

## Features

- **Real-time DeFi Data Aggregation**: Monitor APYs across protocols (Aave, Compound, EigenLayer)
- **AI-Powered Risk Forecasting**: ML models for liquidation risk prediction and APY trend analysis
- **Automated Rebalancing**: Trigger-based actions for optimal portfolio management
- **Wallet Integration**: WalletConnect/MetaMask support for seamless transactions
- **Subscription Management**: Tiered pricing with performance-based fees

## Architecture

The platform is built as a microservices architecture:

- **Backend Services (Go)**: API Gateway, DeFi Data Service, Automation Engine, Wallet Service, User Service
- **ML Service (Python)**: FastAPI service for risk forecasting and APY predictions
- **Frontend (React/TypeScript)**: Dashboard for portfolio monitoring and automation controls
- **Database**: PostgreSQL for structured data, Redis for caching

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+
- Node.js 18+
- Python 3.11+

### Setup

1. Clone the repository
2. **Set up third-party services** (see [SETUP.md](docs/SETUP.md) for detailed instructions):
   - Sign up for blockchain RPC provider (Alchemy/Infura) → Get Ethereum & Base RPC URLs
   - Sign up for WalletConnect → Get Project ID
   - Sign up for Stripe → Get API keys and create subscription products
3. Copy `.env.example` to `.env` and configure your environment variables
4. Start services with Docker Compose:

```bash
docker-compose up -d
```

4. Install frontend dependencies:

```bash
cd frontend
npm install
npm run dev
```

## Development

### Backend Services

Each backend service is a Go module:

```bash
cd backend/api
go run main.go
```

### ML Service

```bash
cd ml-service
pip install -r requirements.txt
python app.py
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Project Structure

```
DefiOptimization/
├── backend/          # Go microservices
├── ml-service/       # Python ML service
├── frontend/         # React/TypeScript dashboard
├── contracts/        # Solidity smart contracts
├── docker/           # Docker configurations
└── docs/             # Documentation
```

## Documentation

### Getting Started
- [Setup Guide](docs/SETUP.md) - Third-party service setup instructions
- [Quick Reference](docs/QUICK_REFERENCE.md) - Quick command and API reference

### Architecture & Design
- [System Overview](docs/SYSTEM_OVERVIEW.md) - High-level system architecture
- [Architecture Documentation](docs/ARCHITECTURE.md) - Detailed architecture and component diagrams
- [Data Flow Documentation](docs/DATA_FLOW.md) - Detailed data flow and sequence diagrams

### Usage & Operations
- [User Guide](docs/USAGE.md) - How to use the platform
- [Troubleshooting](docs/TROUBLESHOOTING.md) - Common issues and solutions

### Integration Guides
- [Stripe Webhooks Guide](docs/STRIPE_WEBHOOKS.md) - Detailed webhook setup for local and production
- [Stripe CLI Installation](docs/INSTALL_STRIPE_CLI.md) - Alternative installation methods

## Quick Start

1. **Set up environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your API keys (see docs/SETUP.md)
   ```

2. **Start services:**
   ```bash
   docker compose up -d
   ```

3. **Start frontend:**
   ```bash
   cd frontend && npm install && npm run dev
   ```

4. **Open browser:**
   ```
   http://localhost:3000
   ```

5. **Connect wallet and start optimizing!**

See [USAGE.md](docs/USAGE.md) for detailed usage instructions.

## License

MIT

