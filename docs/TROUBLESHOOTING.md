# Troubleshooting Guide

## Stripe CLI - macOS Security Warning

If you see: *"stripe Not opened. Apple could not verify stripe is free of malware..."*

This is macOS Gatekeeper blocking the unsigned binary. Fix it by removing the quarantine attribute:

```bash
# Find where stripe is installed
which stripe

# Remove quarantine attribute (use the path from above)
sudo xattr -d com.apple.quarantine /usr/local/bin/stripe

# Or if installed in your home directory:
xattr -d com.apple.quarantine ~/bin/stripe

# Verify it works now
stripe --version
```

**Alternative:** If the above doesn't work, you can allow it in System Settings:
1. Go to System Settings → Privacy & Security
2. Scroll down to find the blocked app
3. Click "Open Anyway"

## Homebrew - Unsupported macOS Version

If you see: `unknown or unsupported macOS version: "26.1"`

This means your macOS version is too new for your Homebrew installation. Solutions:

1. **Update Homebrew:**
   ```bash
   brew update
   ```

2. **Reinstall Homebrew:**
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/uninstall.sh)"
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

3. **Use manual installation instead** (see [INSTALL_STRIPE_CLI.md](INSTALL_STRIPE_CLI.md))

## Other Common Issues

### Port Already in Use

If you get "port already in use" errors:

```bash
# Find what's using the port
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Database Connection Errors

If PostgreSQL connection fails:

```bash
# Check if Docker containers are running
docker-compose ps

# Restart services
docker-compose restart postgres

# Check logs
docker-compose logs postgres
```

### Go Module Errors

If you see "cannot find module" errors:

```bash
# In each backend service directory
cd backend/api
go mod tidy
go mod download
```

### Frontend Build Errors

If frontend won't start:

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

