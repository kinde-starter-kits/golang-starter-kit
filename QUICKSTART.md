# 🚀 Quick Start - Run in 2 Minutes

## You Need Kinde Credentials First!

The app requires valid Kinde credentials. Here's how to get them:

### Step 1: Get Your Kinde Credentials (1 minute)

1. Go to [https://app.kinde.com](https://app.kinde.com)
2. Sign in to your Kinde account
3. Go to **Settings** → **Applications** → Select or create an application
4. Copy these values:
   - **Domain** (e.g., `https://yourbusiness.kinde.com`)
   - **Client ID** (looks like: `abc123def456...`)
   - **Client Secret** (click "Show" to reveal)

### Step 2: Update Callback URLs in Kinde (30 seconds)

In your Kinde application settings, add:
- **Allowed callback URLs**: `http://localhost:3000/callback`
- **Allowed logout redirect URLs**: `http://localhost:3000`

Click **Save**

### Step 3: Configure the `.env` File (30 seconds)

Edit `.env`:

```bash
KINDE_DOMAIN=https://yourbusiness.kinde.com  # ← Your actual domain
KINDE_CLIENT_ID=your_client_id_here          # ← Your actual client ID
KINDE_CLIENT_SECRET=your_client_secret_here  # ← Your actual client secret
KINDE_REDIRECT_URI=http://localhost:3000/callback
KINDE_LOGOUT_REDIRECT_URI=http://localhost:3000
PORT=3000
SESSION_SECRET=any-random-string-you-want-here-make-it-long
```

### Step 4: Run the App! ⚡

```bash
go run main.go
```

You should see:
```
🚀 Kinde Golang Starter Kit running on http://localhost:3000
📝 Make sure you've configured your Kinde app with redirect URI: http://localhost:3000/callback
```

### Step 5: Test It! 🎉

1. Open browser: **http://localhost:3000**
2. Click **"Sign In"** or **"Sign Up"**
3. Login with Kinde
4. You'll be redirected to the dashboard!

---

## Without Kinde Credentials?

The app **cannot run** without valid Kinde credentials because:
- It validates the configuration on startup
- Authentication requires communication with Kinde's OAuth2 endpoints

If you don't have credentials yet:
1. Sign up for free at [https://kinde.com](https://kinde.com)
2. Follow Step 1 above to get your credentials
3. It takes less than 5 minutes!

---

## Troubleshooting

**App exits immediately?**
- Check your `.env` file has real values (not placeholders)
- Make sure `KINDE_DOMAIN` starts with `https://`

**"Invalid state parameter"?**
- Clear browser cookies
- Make sure callback URL matches exactly in Kinde

**Can't authenticate?**
- Verify callback URLs are added in Kinde dashboard
- Check Client ID and Secret are correct

---

Ready to customize? Check out the [README.md](README.md) for full documentation!



