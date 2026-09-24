# Quick Setup Guide

## Step 1: Configure Kinde

1. Go to [https://kinde.com](https://kinde.com) and sign in to your account
2. Create a new application (or use an existing one)
3. In your Kinde application settings, configure:
   - **Allowed callback URLs**: `http://localhost:3000/callback`
   - **Allowed logout redirect URLs**: `http://localhost:3000`
4. Note down your:
   - Domain (e.g., `https://yourbusiness.kinde.com`)
   - Client ID
   - Client Secret

## Step 2: Configure Environment Variables

Edit the `.env` file in the project root and replace the placeholder values:

```env
KINDE_DOMAIN=https://yourbusiness.kinde.com
KINDE_CLIENT_ID=your_actual_client_id
KINDE_CLIENT_SECRET=your_actual_client_secret
KINDE_REDIRECT_URI=http://localhost:3000/callback
KINDE_LOGOUT_REDIRECT_URI=http://localhost:3000
PORT=3000
SESSION_SECRET=generate-a-random-secure-string-here
```

## Step 3: Run the Application

```bash
# Install dependencies (if not done already)
go mod download

# Run the application
go run main.go
```

You should see:
```
🚀 Kinde Golang Starter Kit running on http://localhost:3000
📝 Make sure you've configured your Kinde app with redirect URI: http://localhost:3000/callback
```

## Step 4: Test It Out

1. Open your browser and go to [http://localhost:3000](http://localhost:3000)
2. Click "Sign In" or "Sign Up"
3. You'll be redirected to Kinde for authentication
4. After logging in, you'll be redirected back to the dashboard

## Troubleshooting

### Error: "KINDE_DOMAIN is required"
- Make sure your `.env` file exists and has all required values
- Check that you're running from the correct directory

### Error: "Failed to load configuration"
- Verify your `.env` file format
- Ensure no extra spaces in the values

### Authentication fails
- Check that your callback URLs match exactly in Kinde
- Verify your Client ID and Client Secret are correct
- Make sure your Kinde domain includes `https://`

## Next Steps

- Customize the templates in `templates/`
- Add your own routes and handlers
- Style the app by editing `static/style.css`
- Deploy to production (see README.md for deployment guides)



