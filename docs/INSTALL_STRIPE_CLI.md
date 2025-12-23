# Installing Stripe CLI (Alternative Methods)

If Homebrew isn't working (e.g., unsupported macOS version), use one of these methods:

## Method 1: Direct Download (Recommended)

1. **Go to the Stripe CLI releases page:**
   https://github.com/stripe/stripe-cli/releases/latest

2. **Download the correct file for your system:**
   - **macOS (Apple Silicon/M1/M2):** `stripe_X.X.X_darwin_arm64.tar.gz`
   - **macOS (Intel):** `stripe_X.X.X_darwin_amd64.tar.gz`
   - **Linux:** `stripe_X.X.X_linux_amd64.tar.gz`
   - **Windows:** `stripe_X.X.X_windows_amd64.zip`

3. **Extract and install:**

   **macOS/Linux:**
   ```bash
   # Extract
   tar -xzf stripe_X.X.X_darwin_arm64.tar.gz
   
   # Move to a directory in your PATH
   # Option A: System-wide (requires sudo)
   sudo mv stripe /usr/local/bin/
   
   # Option B: User directory (no sudo needed)
   mkdir -p ~/bin
   mv stripe ~/bin/
   echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
   source ~/.zshrc
   
   # IMPORTANT: Remove macOS quarantine attribute (if you get "not verified" error)
   # For system-wide installation:
   sudo xattr -d com.apple.quarantine /usr/local/bin/stripe
   
   # For user directory installation:
   xattr -d com.apple.quarantine ~/bin/stripe
   
   # Verify
   stripe --version
   ```

   **Windows:**
   - Extract the ZIP file
   - Move `stripe.exe` to a directory in your PATH (e.g., `C:\Program Files\stripe\`)
   - Or add the extraction directory to your PATH environment variable

## Method 2: Using Go (if you have Go installed)

```bash
go install github.com/stripe/stripe-cli@latest
```

This installs to `~/go/bin/` (or `$GOPATH/bin`). Make sure that's in your PATH.

## Method 3: Using npm (if you have Node.js)

```bash
npm install -g @stripe/stripe-cli
```

## Method 4: Fix Homebrew (if you want to use Homebrew)

If your Homebrew is outdated or doesn't recognize your macOS version:

```bash
# Update Homebrew
brew update

# If that fails, try updating Homebrew manually
cd /usr/local/Homebrew
git pull

# Or reinstall Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/uninstall.sh)"
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

## Verify Installation

After installation, verify it works:

```bash
stripe --version
```

You should see something like: `stripe version X.X.X`

## Next Steps

Once Stripe CLI is installed, continue with the webhook setup:

1. Login: `stripe login`
2. Forward webhooks: `stripe listen --forward-to localhost:8080/api/v1/webhooks/stripe`
3. Copy the webhook secret to your `.env` file

See [STRIPE_WEBHOOKS.md](STRIPE_WEBHOOKS.md) for complete webhook setup instructions.

