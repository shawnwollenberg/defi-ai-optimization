# Setup Guide - Third-Party Services

This guide walks you through setting up all the required third-party services for the DeFi Optimization Platform.

## Required Third-Party Services

### 1. Blockchain RPC Providers

You need RPC endpoints to interact with Ethereum and Base blockchains. You can use any of these providers:

#### Option A: Alchemy (Recommended)
1. Go to [https://www.alchemy.com/](https://www.alchemy.com/)
2. Sign up for a free account
3. Create a new app:
   - Select "Ethereum" network → Get your API key
   - Select "Base" network → Get your API key
4. Copy your API keys to `.env`:
   ```
   ETH_RPC_URL=https://eth-mainnet.g.alchemy.com/v2/YOUR_ETH_API_KEY
   BASE_RPC_URL=https://base-mainnet.g.alchemy.com/v2/YOUR_BASE_API_KEY
   ```

#### Option B: Infura
1. Go to [https://www.infura.io/](https://www.infura.io/)
2. Sign up for a free account
3. Create a new project
4. Get your API keys for Ethereum and Base networks
5. Copy to `.env`:
   ```
   ETH_RPC_URL=https://mainnet.infura.io/v3/YOUR_PROJECT_ID
   BASE_RPC_URL=https://base-mainnet.infura.io/v3/YOUR_PROJECT_ID
   ```

#### Option C: QuickNode
1. Go to [https://www.quicknode.com/](https://www.quicknode.com/)
2. Sign up and create endpoints for Ethereum and Base
3. Copy the HTTP URLs to `.env`

**Note:** Free tiers typically have rate limits. For production, consider paid plans.

---

### 2. WalletConnect

WalletConnect enables wallet connections in your frontend.

1. Go to [https://cloud.walletconnect.com/](https://cloud.walletconnect.com/)
2. Sign up for a free account
3. Create a new project
4. Copy your Project ID to `.env`:
   ```
   WALLETCONNECT_PROJECT_ID=your-project-id-here
   ```

**Note:** The free tier includes 1 million requests/month, which is sufficient for development and small-scale production.

---

### 3. Stripe (For Subscriptions)

Stripe handles subscription payments and billing.

1. Go to [https://dashboard.stripe.com/register](https://dashboard.stripe.com/register)
2. Sign up for a Stripe account
3. Get your API keys from the Dashboard:
   - Go to Developers → API keys
   - Copy your **Secret key** (starts with `sk_test_` for test mode)
   - Copy your **Publishable key** (starts with `pk_test_` for test mode)
4. Create products and prices:
   - Go to Products → Add product
   - Create "Basic" subscription ($10/month) → Copy the Price ID
   - Create "Premium" subscription ($50/month) → Copy the Price ID
5. Set up webhooks:
   - **For local development:** Use Stripe CLI (see [STRIPE_WEBHOOKS.md](STRIPE_WEBHOOKS.md) for detailed instructions)
   - **For production:** Go to Developers → Webhooks → Add endpoint
   - Select events: `checkout.session.completed`, `customer.subscription.updated`, `customer.subscription.deleted`
   - Copy the webhook signing secret
6. Add to `.env`:
   ```
   STRIPE_SECRET_KEY=sk_test_your_secret_key
   STRIPE_PUBLISHABLE_KEY=pk_test_your_publishable_key
   STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
   STRIPE_BASIC_PRICE_ID=price_your_basic_price_id
   STRIPE_PREMIUM_PRICE_ID=price_your_premium_price_id
   ```

**Note:** Use test mode keys for development. Switch to live mode keys for production.

---

## Optional Services

### Database & Caching

**PostgreSQL** and **Redis** are included in Docker Compose, so no external signup is needed. However, for production, you may want to use managed services:

- **PostgreSQL**: AWS RDS, Google Cloud SQL, or Supabase
- **Redis**: AWS ElastiCache, Redis Cloud, or Upstash

---

## Environment Variables Summary

After signing up for the services above, your `.env` file should look like this:

```bash
# Database (Docker - no signup needed)
POSTGRES_USER=defi_user
POSTGRES_PASSWORD=defi_password
POSTGRES_DB=defi_optimization
DATABASE_URL=postgres://defi_user:defi_password@localhost:5432/defi_optimization

# Redis (Docker - no signup needed)
REDIS_URL=redis://localhost:6379

# JWT (Generate your own)
JWT_SECRET=your-secret-key-change-in-production

# Blockchain RPC (Sign up: Alchemy/Infura/QuickNode)
ETH_RPC_URL=https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY
BASE_RPC_URL=https://base-mainnet.g.alchemy.com/v2/YOUR_API_KEY

# WalletConnect (Sign up: https://cloud.walletconnect.com/)
WALLETCONNECT_PROJECT_ID=your-walletconnect-project-id

# Stripe (Sign up: https://stripe.com/)
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
STRIPE_BASIC_PRICE_ID=price_your_basic_price_id
STRIPE_PREMIUM_PRICE_ID=price_your_premium_price_id

# Service URLs (Local development)
API_URL=http://localhost:8080
DEFI_SERVICE_URL=http://localhost:8081
ML_SERVICE_URL=http://localhost:8001
WALLET_SERVICE_URL=http://localhost:8082
AUTOMATION_SERVICE_URL=http://localhost:8083

# Frontend
VITE_API_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000
```

---

## Quick Setup Checklist

- [ ] Sign up for Alchemy/Infura/QuickNode → Get Ethereum RPC URL
- [ ] Sign up for Alchemy/Infura/QuickNode → Get Base RPC URL
- [ ] Sign up for WalletConnect → Get Project ID
- [ ] Sign up for Stripe → Get API keys and create products
- [ ] Generate a secure JWT_SECRET (use `openssl rand -hex 32`)
- [ ] Copy `.env.example` to `.env` and fill in all values
- [ ] Start services: `docker-compose up -d`

---

## Cost Estimates (Free Tiers)

- **Alchemy/Infura**: Free tier includes ~300K requests/month
- **WalletConnect**: Free tier includes 1M requests/month
- **Stripe**: No monthly fee, only transaction fees (2.9% + $0.30 per transaction)
- **Docker Services**: Free (runs locally)

**Total Monthly Cost for Development**: $0 (using free tiers)

For production, expect to pay:
- RPC Provider: $50-200/month (depending on usage)
- WalletConnect: Free up to 1M requests, then pay-as-you-go
- Stripe: Transaction fees only
- Managed Database/Redis: $20-100/month (optional)

---

## Troubleshooting

### RPC Rate Limits
If you hit rate limits, you can:
- Upgrade to a paid plan
- Use multiple RPC providers and rotate between them
- Implement request caching

### WalletConnect Connection Issues
- Ensure your Project ID is correct
- Check that your frontend URL is whitelisted in WalletConnect dashboard
- Verify CORS settings in your API gateway

### Stripe Webhook Not Working
- **Recommended:** Use Stripe CLI for local testing (see [STRIPE_WEBHOOKS.md](STRIPE_WEBHOOKS.md))
- Alternative: Use [ngrok](https://ngrok.com/) to expose your local server
- Verify the webhook secret matches (different for CLI vs Dashboard)
- Check that webhook events are properly configured in Stripe dashboard
- See [STRIPE_WEBHOOKS.md](STRIPE_WEBHOOKS.md) for detailed troubleshooting

