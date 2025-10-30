// Package examples demonstrates basic usage of the amapi client.
package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"amapi-pkg/pkgs/amapi/client"
	"amapi-pkg/pkgs/amapi/config"
	"amapi-pkg/pkgs/amapi/types"
)

// BasicUsageExample demonstrates basic usage patterns of the AMAPI client.
func BasicUsageExample() {
	// Example 1: Basic client setup
	basicSetup()

	// Example 2: Enterprise operations
	enterpriseOperations()

	// Example 3: Policy management
	policyManagement()

	// Example 4: Device management
	deviceManagement()

	// Example 5: Enrollment tokens
	enrollmentTokens()
}

// basicSetup demonstrates basic client configuration and setup.
func basicSetup() {
	fmt.Println("=== Basic Setup ===")

	// Method 1: Auto-load configuration from environment or files
	cfg, err := config.AutoLoadConfig()
	if err != nil {
		log.Printf("Auto-load config failed: %v", err)

		// Method 2: Load from environment variables only
		cfg, err = config.LoadFromEnv()
		if err != nil {
			log.Printf("Environment config failed: %v", err)

			// Method 3: Manual configuration
			cfg = &config.Config{
				ProjectID:       "your-project-id",
				CredentialsFile: "./service-account-key.json",
				CallbackURL:     "https://your-app.com/callback",
				Timeout:         30 * time.Second,
				RetryAttempts:   3,
				EnableRetry:     true,
				LogLevel:        "info",
			}
		}
	}

	// Create client
	c, err := client.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Test connectivity
	if err := c.Health(); err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Println("✓ Client connected successfully")
	}

	// Display client info
	info := c.GetInfo()
	fmt.Printf("✓ Client version: %s\n", info.Version)
	fmt.Printf("✓ Project ID: %s\n", info.ProjectID)
	fmt.Printf("✓ Capabilities: %v\n", info.Capabilities)
}

// enterpriseOperations demonstrates enterprise management operations.
func enterpriseOperations() {
	fmt.Println("\n=== Enterprise Operations ===")

	cfg, _ := config.AutoLoadConfig()
	c, err := client.New(cfg)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer c.Close()

	// Generate signup URL
	signupReq := &types.SignupURLRequest{
		ProjectID:             cfg.ProjectID,
		CallbackURL:           cfg.CallbackURL,
		AdminEmail:            "admin@company.com",
		EnterpriseDisplayName: "Example Company",
	}

	signupURL, err := c.Enterprises().GenerateSignupURL(signupReq)
	if err != nil {
		log.Printf("Failed to generate signup URL: %v", err)
	} else {
		fmt.Printf("✓ Signup URL generated: %s\n", signupURL.URL)
	}

	// List existing enterprises
	enterprises, err := c.Enterprises().List(nil)
	if err != nil {
		log.Printf("Failed to list enterprises: %v", err)
	} else {
		fmt.Printf("✓ Found %d enterprises\n", len(enterprises.Items))
		for _, enterprise := range enterprises.Items {
			fmt.Printf("  - %s (%s)\n", enterprise.DisplayName, enterprise.GetID())
		}
	}

	// Get specific enterprise (if any exist)
	if len(enterprises.Items) > 0 {
		enterprise := enterprises.Items[0]
		detailed, err := c.Enterprises().Get(enterprise.Name)
		if err != nil {
			log.Printf("Failed to get enterprise details: %v", err)
		} else {
			fmt.Printf("✓ Enterprise details: %+v\n", detailed.DisplayName)
		}
	}
}

// policyManagement demonstrates policy creation and management.
func policyManagement() {
	fmt.Println("\n=== Policy Management ===")

	cfg, _ := config.AutoLoadConfig()
	c, err := client.New(cfg)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer c.Close()

	// First, get an enterprise to work with
	enterprises, err := c.Enterprises().List(nil)
	if err != nil || len(enterprises.Items) == 0 {
		log.Printf("No enterprises found, skipping policy management example")
		return
	}

	enterpriseID := enterprises.Items[0].GetID()
	fmt.Printf("✓ Using enterprise: %s\n", enterpriseID)

	// Create a basic policy
	policy := &types.Policy{
		StatusBarDisabled:     false,
		KeyguardDisabled:      false,
		AddUserDisabled:       true,
		UninstallAppsDisabled: true,
		CameraDisabled:        false,
		BluetoothDisabled:     false,
		AutoTimeRequired:      true,
	}

	// Validate policy before creating
	if err := policy.Validate(); err != nil {
		log.Printf("Policy validation failed: %v", err)
		return
	}

	// Create policy
	created, err := c.Policies().CreateByEnterpriseID(enterpriseID, "example-policy", policy)
	if err != nil {
		log.Printf("Failed to create policy: %v", err)
	} else {
		fmt.Printf("✓ Policy created: %s\n", created.GetID())
	}

	// List policies
	policies, err := c.Policies().ListByEnterpriseID(enterpriseID, nil)
	if err != nil {
		log.Printf("Failed to list policies: %v", err)
	} else {
		fmt.Printf("✓ Found %d policies\n", len(policies.Items))
		for _, pol := range policies.Items {
			fmt.Printf("  - %s (Mode: %s)\n", pol.GetID(), pol.GetPolicyMode())
		}
	}

	// Update policy
	if len(policies.Items) > 0 {
		policy := policies.Items[0]
		policy.CameraDisabled = true // Disable camera

		updated, err := c.Policies().UpdateByID(enterpriseID, policy.GetID(), &policy)
		if err != nil {
			log.Printf("Failed to update policy: %v", err)
		} else {
			fmt.Printf("✓ Policy updated, camera disabled: %t\n", updated.CameraDisabled)
		}
	}
}

// deviceManagement demonstrates device operations.
func deviceManagement() {
	fmt.Println("\n=== Device Management ===")

	cfg, _ := config.AutoLoadConfig()
	c, err := client.New(cfg)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer c.Close()

	// Get an enterprise to work with
	enterprises, err := c.Enterprises().List(nil)
	if err != nil || len(enterprises.Items) == 0 {
		log.Printf("No enterprises found, skipping device management example")
		return
	}

	enterpriseID := enterprises.Items[0].GetID()

	// List devices
	devices, err := c.Devices().ListByEnterpriseID(enterpriseID, nil)
	if err != nil {
		log.Printf("Failed to list devices: %v", err)
		return
	}

	fmt.Printf("✓ Found %d devices\n", len(devices.Items))

	if len(devices.Items) == 0 {
		fmt.Println("  No devices enrolled yet")
		return
	}

	// Display device information
	for _, device := range devices.Items {
		fmt.Printf("  - Device: %s\n", device.GetID())
		fmt.Printf("    State: %s\n", device.State)
		fmt.Printf("    Compliant: %t\n", device.PolicyCompliant)
		fmt.Printf("    User: %s\n", device.UserName)
		fmt.Printf("    Android Version: %s\n", device.GetAndroidVersion())
		fmt.Printf("    Model: %s\n", device.GetDeviceModel())
		fmt.Printf("    Online: %t\n", device.IsOnline())
	}

	// Get compliance statistics
	compliantDevices, err := c.Devices().GetCompliantDevices(enterpriseID)
	if err != nil {
		log.Printf("Failed to get compliant devices: %v", err)
	} else {
		fmt.Printf("✓ Compliant devices: %d\n", len(compliantDevices.Items))
	}

	nonCompliantDevices, err := c.Devices().GetNonCompliantDevices(enterpriseID)
	if err != nil {
		log.Printf("Failed to get non-compliant devices: %v", err)
	} else {
		fmt.Printf("✓ Non-compliant devices: %d\n", len(nonCompliantDevices.Items))
	}

	// Example device command (commented out to avoid affecting real devices)
	/*
		if len(devices.Items) > 0 {
			deviceID := devices.Items[0].GetID()

			// Lock device for 5 minutes
			err := c.Devices().LockByID(enterpriseID, deviceID, "PT5M")
			if err != nil {
				log.Printf("Failed to lock device: %v", err)
			} else {
				fmt.Printf("✓ Device locked: %s\n", deviceID)
			}
		}
	*/
}

// enrollmentTokens demonstrates enrollment token management.
func enrollmentTokens() {
	fmt.Println("\n=== Enrollment Tokens ===")

	cfg, _ := config.AutoLoadConfig()
	c, err := client.New(cfg)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer c.Close()

	// Get enterprise and policy to work with
	enterprises, err := c.Enterprises().List(nil)
	if err != nil || len(enterprises.Items) == 0 {
		log.Printf("No enterprises found, skipping enrollment token example")
		return
	}

	enterpriseID := enterprises.Items[0].GetID()

	policies, err := c.Policies().ListByEnterpriseID(enterpriseID, nil)
	if err != nil || len(policies.Items) == 0 {
		log.Printf("No policies found, skipping enrollment token example")
		return
	}

	policyID := policies.Items[0].GetID()

	// Create enrollment token
	token, err := c.EnrollmentTokens().CreateByEnterpriseID(
		enterpriseID,
		policyID,
		24*time.Hour, // Valid for 24 hours
	)
	if err != nil {
		log.Printf("Failed to create enrollment token: %v", err)
		return
	}

	fmt.Printf("✓ Enrollment token created: %s\n", token.GetID())
	fmt.Printf("  Token value: %s\n", token.Value)
	fmt.Printf("  Policy: %s\n", token.GetPolicyID())
	fmt.Printf("  Expires: %s\n", token.ExpirationTimestamp)

	// Generate QR code
	qrOptions := &types.QRCodeOptions{
		WiFiSSID:         "CompanyWiFi",
		WiFiPassword:     "password123",
		WiFiSecurityType: types.WiFiSecurityTypeWPA2,
		SkipSetupWizard:  true,
		Locale:           "en_US",
	}

	qrData, err := c.EnrollmentTokens().GenerateQRCodeByID(
		enterpriseID,
		token.GetID(),
		qrOptions,
	)
	if err != nil {
		log.Printf("Failed to generate QR code: %v", err)
	} else {
		fmt.Printf("✓ QR code data generated\n")
		qrJSON, _ := qrData.ToJSON()
		fmt.Printf("  QR JSON: %s\n", qrJSON)
	}

	// List all tokens
	tokens, err := c.EnrollmentTokens().GetActiveTokens(enterpriseID)
	if err != nil {
		log.Printf("Failed to list tokens: %v", err)
	} else {
		fmt.Printf("✓ Found %d active enrollment tokens\n", len(tokens.Items))
		for _, tok := range tokens.Items {
			fmt.Printf("  - Token: %s (expires in %v)\n",
				tok.GetID(),
				tok.TimeUntilExpiration())
		}
	}

	// Get token statistics
	stats, err := c.EnrollmentTokens().GetTokenStatistics(enterpriseID)
	if err != nil {
		log.Printf("Failed to get token statistics: %v", err)
	} else {
		fmt.Printf("✓ Token statistics: %+v\n", stats)
	}
}

// errorHandlingExample demonstrates proper error handling.
func errorHandlingExample() {
	fmt.Println("\n=== Error Handling ===")

	cfg, _ := config.AutoLoadConfig()
	c, err := client.New(cfg)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
	defer c.Close()

	// Try to get a non-existent enterprise
	_, err = c.Enterprises().GetByID("non-existent-enterprise")
	if err != nil {
		if apiErr, ok := err.(*types.Error); ok {
			fmt.Printf("✓ API Error caught:\n")
			fmt.Printf("  Code: %d\n", apiErr.Code)
			fmt.Printf("  Message: %s\n", apiErr.Message)
			fmt.Printf("  Type: %s\n", apiErr.GetErrorType())
			fmt.Printf("  Retryable: %t\n", apiErr.IsRetryable())

			if apiErr.IsRetryable() {
				delay := apiErr.RetryDelay(1, time.Second)
				fmt.Printf("  Retry delay: %v\n", delay)
			}
		} else {
			fmt.Printf("✓ Generic error: %v\n", err)
		}
	}
}

// contextExample demonstrates context usage.
func contextExample() {
	fmt.Println("\n=== Context Usage ===")

	cfg, _ := config.AutoLoadConfig()

	// Create client with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := client.NewWithContext(ctx, cfg)
	if err != nil {
		log.Printf("Failed to create client with context: %v", err)
		return
	}
	defer c.Close()

	// Operations will respect the context timeout
	enterprises, err := c.Enterprises().List(nil)
	if err != nil {
		log.Printf("Context-aware operation failed: %v", err)
	} else {
		fmt.Printf("✓ Context-aware operation succeeded, found %d enterprises\n", len(enterprises.Items))
	}
}