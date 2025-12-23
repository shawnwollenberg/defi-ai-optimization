# Stripe Webhooks Setup Guide

Setting up Stripe webhooks can be tricky, especially for local development. This guide covers both local and production setups.

## Option 1: Local Development with Stripe CLI (Recommended)

The easiest way to test webhooks locally is using the Stripe CLI, which forwards webhook events directly to your local server.

### Step 1: Install Stripe CLI

**If Homebrew works:**
```bash
brew install stripe/stripe-cli/stripe
```

**If Homebrew doesn't work (e.g., unsupported macOS version):**

See [INSTALL_STRIPE_CLI.md](INSTALL_STRIPE_CLI.md) for detailed alternative installation methods.

**Quick manual install for macOS:**
1. Go to https://github.com/stripe/stripe-cli/releases/latest
2. Download `stripe_X.X.X_darwin_arm64.tar.gz` (for Apple Silicon) or `stripe_X.X.X_darwin_amd64.tar.gz` (for Intel)
3. Extract: `tar -xzf stripe_*.tar.gz`
4. Move to PATH: `sudo mv stripe /usr/local/bin/` or `mv stripe ~/bin/` (and add ~/bin to PATH)
5. Verify: `stripe --version`

### Step 2: Login to Stripe CLI

```bash
stripe login
```

This will open your browser to authenticate. Make sure you're logged into the same Stripe account you're using for your project.

### Step 3: Get Your Webhook Signing Secret

For local development, you'll use a special signing secret from the CLI:

```bash
stripe listen --forward-to localhost:8080/api/v1/webhooks/stripe
```

This command will:
- Start listening for webhook events
- Forward them to your local server
- Display a webhook signing secret (starts with `whsec_`)

**Copy this signing secret** - you'll need it for your `.env` file.

### Step 4: Update Your .env File

```bash
STRIPE_WEBHOOK_SECRET=whsec_xxxxx  # The secret from step 3
```

### Step 5: Start Your API Server

Make sure your API gateway is running:

```bash
cd backend/api
go run main.go
```

Or with Docker:
```bash
docker-compose up api
```

### Step 6: Test Webhooks

In another terminal, trigger test events:

```bash
# Test checkout completion
stripe trigger checkout.session.completed

# Test subscription update
stripe trigger customer.subscription.updated

# Test subscription cancellation
stripe trigger customer.subscription.deleted
```

You should see the events being received by your server!

---

## Option 2: Local Development with ngrok

If you prefer to use the Stripe Dashboard webhook interface, you can expose your local server using ngrok.

### Step 1: Install ngrok

```bash
# macOS
brew install ngrok

# Or download from https://ngrok.com/download
```

### Step 2: Start Your API Server

```bash
cd backend/api
go run main.go
# Server should be running on localhost:8080
```

### Step 3: Expose Your Server with ngrok

```bash
ngrok http 8080
```

This will give you a public URL like: `https://abc123.ngrok.io`

### Step 4: Set Up Webhook in Stripe Dashboard

1. Go to [Stripe Dashboard → Developers → Webhooks](https://dashboard.stripe.com/test/webhooks)
2. Click **"Add endpoint"**
3. Enter endpoint URL: `https://abc123.ngrok.io/api/v1/webhooks/stripe`
4. Select events to listen to:
   - `checkout.session.completed`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
5. Click **"Add endpoint"**
6. Copy the **"Signing secret"** (starts with `whsec_`)

### Step 5: Update Your .env File

```bash
STRIPE_WEBHOOK_SECRET=whsec_xxxxx  # From step 4
```

### Step 6: Test Webhooks

1. In Stripe Dashboard, go to your webhook endpoint
2. Click **"Send test webhook"**
3. Select an event type and click **"Send test webhook"**
4. Check your server logs to see if it received the event

**Note:** The ngrok URL changes every time you restart ngrok (unless you have a paid plan). You'll need to update the webhook URL in Stripe Dashboard each time.

---

## Option 3: Production Setup

For production, you'll set up webhooks directly in the Stripe Dashboard.

### Step 1: Deploy Your API

Make sure your API is deployed and accessible at a public URL, e.g.:
- `https://api.yourdomain.com/api/v1/webhooks/stripe`

### Step 2: Set Up Webhook in Stripe Dashboard

1. Go to [Stripe Dashboard → Developers → Webhooks](https://dashboard.stripe.com/webhooks)
2. Click **"Add endpoint"**
3. Enter your production endpoint URL
4. Select events:
   - `checkout.session.completed`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
5. Click **"Add endpoint"**
6. Copy the **"Signing secret"**

### Step 3: Update Production Environment Variables

Set `STRIPE_WEBHOOK_SECRET` in your production environment with the signing secret from step 2.

---

## Troubleshooting

### Webhook Not Receiving Events

1. **Check your server is running:**
   ```bash
   curl http://localhost:8080/health
   ```

2. **Verify webhook secret matches:**
   - Make sure `STRIPE_WEBHOOK_SECRET` in `.env` matches the secret from Stripe
   - For Stripe CLI: Use the secret shown when running `stripe listen`
   - For Dashboard: Use the secret from the webhook endpoint settings

3. **Check webhook endpoint URL:**
   - Should be: `http://localhost:8080/api/v1/webhooks/stripe` (local)
   - Or: `https://your-domain.com/api/v1/webhooks/stripe` (production)

4. **Verify webhook handler is registered:**
   Check `backend/api/main.go` - should have:
   ```go
   r.POST("/api/v1/webhooks/stripe", handlers.HandleStripeWebhook)
   ```

### "Invalid signature" Error

This means the webhook secret doesn't match. Solutions:

1. **For Stripe CLI:** Make sure you're using the secret from the `stripe listen` command
2. **For Dashboard:** Make sure you copied the correct signing secret
3. **Check .env file:** Verify `STRIPE_WEBHOOK_SECRET` is set correctly
4. **Restart server:** After changing `.env`, restart your API server

### Webhook Events Not Being Processed

1. **Check server logs** for errors
2. **Verify event types** are selected in Stripe Dashboard
3. **Test with Stripe CLI:**
   ```bash
   stripe trigger checkout.session.completed
   ```

### ngrok URL Changed

If using ngrok and your URL changed:
1. Get the new ngrok URL
2. Update the webhook endpoint URL in Stripe Dashboard
3. Or use Stripe CLI instead (easier for development)

---

## Quick Reference

### Stripe CLI Commands

```bash
# Login to Stripe
stripe login

# Forward webhooks to local server
stripe listen --forward-to localhost:8080/api/v1/webhooks/stripe

# Trigger test events
stripe trigger checkout.session.completed
stripe trigger customer.subscription.updated
stripe trigger customer.subscription.deleted

# View webhook events
stripe events list
```

### Required Environment Variables

```bash
STRIPE_SECRET_KEY=sk_test_xxxxx
STRIPE_PUBLISHABLE_KEY=pk_test_xxxxx
STRIPE_WEBHOOK_SECRET=whsec_xxxxx  # Different for CLI vs Dashboard
STRIPE_BASIC_PRICE_ID=price_xxxxx
STRIPE_PREMIUM_PRICE_ID=price_xxxxx
```

### Webhook Endpoint

- **Local:** `http://localhost:8080/api/v1/webhooks/stripe`
- **Production:** `https://your-domain.com/api/v1/webhooks/stripe`

---

## Recommended Approach

**For Development:**
- Use **Stripe CLI** (Option 1) - it's the easiest and most reliable

**For Production:**
- Use **Stripe Dashboard** (Option 3) with your deployed API

**Avoid for Development:**
- ngrok (Option 2) - only use if you specifically need to test the Dashboard webhook interface

