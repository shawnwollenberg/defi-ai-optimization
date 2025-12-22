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
2. Copy `.env.example` to `.env` and configure your environment variables
3. Start services with Docker Compose:

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

## License

MIT

