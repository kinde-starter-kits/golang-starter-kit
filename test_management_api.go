package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kinde-oss/kinde-go/kinde"
	"github.com/kinde-oss/kinde-go/kinde/management_api"
	"github.com/kinde-oss/kinde-go/oauth2/client_credentials"
)

// TestManagementAPI tests the Kinde Management API with null handling
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	kindeDomain := os.Getenv("KINDE_DOMAIN")
	clientID := os.Getenv("KINDE_CLIENT_ID")
	clientSecret := os.Getenv("KINDE_CLIENT_SECRET")

	if kindeDomain == "" || clientID == "" || clientSecret == "" {
		log.Fatal("KINDE_DOMAIN, KINDE_CLIENT_ID, and KINDE_CLIENT_SECRET must be set")
	}

	fmt.Println("🔧 Testing Kinde Management API with local SDK...")
	fmt.Println("This will test the null value handling fix for permissions")
	fmt.Println()

	// Create Client Credentials Flow
	ctx := context.Background()
	flow, err := client_credentials.NewClientCredentialsFlow(
		kindeDomain,
		clientID,
		clientSecret,
		client_credentials.WithKindeManagementAPI(kindeDomain),
	)
	if err != nil {
		log.Fatalf("Failed to create client credentials flow: %v", err)
	}

	// Create Management API client
	client, err := kinde.NewManagementAPI(ctx, kindeDomain, flow)
	if err != nil {
		log.Fatalf("Failed to create Management API client: %v", err)
	}

	fmt.Println("✅ Management API client created successfully")
	fmt.Println()

	// Test: Get permissions (this will test null description handling)
	fmt.Println("📋 Fetching permissions...")
	response, err := client.GetPermissions(ctx, management_api.GetPermissionsParams{
		PageSize: management_api.NewOptNilInt(10),
	})
	if err != nil {
		log.Fatalf("❌ Failed to get permissions: %v", err)
	}

	switch res := response.(type) {
	case *management_api.GetPermissionsResponse:
		fmt.Printf("✅ Successfully fetched %d permissions\n\n", len(res.Permissions))
		
		// Display permissions and their descriptions
		for i, perm := range res.Permissions {
			fmt.Printf("Permission %d:\n", i+1)
			
			if id, ok := perm.ID.Get(); ok {
				fmt.Printf("  ID: %s\n", id)
			}
			
			if key, ok := perm.Key.Get(); ok {
				fmt.Printf("  Key: %s\n", key)
			}
			
			if name, ok := perm.Name.Get(); ok {
				fmt.Printf("  Name: %s\n", name)
			}
			
			// This is the critical test - descriptions can be null
			if perm.Description.IsSet() {
				if desc, ok := perm.Description.Get(); ok {
					fmt.Printf("  Description: %s\n", desc)
				}
			} else {
				fmt.Printf("  Description: (null/not set) ✅ Handled correctly!\n")
			}
			
			fmt.Println()
		}
		
		fmt.Println("🎉 SUCCESS! Null descriptions were handled correctly!")
		fmt.Println("The fix is working - no 'unexpected byte 110' errors!")
		
	default:
		log.Fatalf("Unexpected response type: %T", response)
	}
}

