// Package examples demonstrates complete enterprise setup workflow.
package main

import (
	"fmt"
	"log"
	"time"

	"amapi-pkg/pkgs/amapi/client"
	"amapi-pkg/pkgs/amapi/config"
	"amapi-pkg/pkgs/amapi/presets"
	"amapi-pkg/pkgs/amapi/types"
)

func main() {
	// Complete enterprise setup workflow
	enterpriseSetupWorkflow()
}

// enterpriseSetupWorkflow demonstrates the complete process of setting up a new enterprise.
func enterpriseSetupWorkflow() {
	fmt.Println("=== Enterprise Setup Workflow ===")

	// Step 1: Initialize client
	cfg, err := config.AutoLoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	fmt.Println("✓ Client initialized")

	// Step 2: Generate enterprise signup URL
	signupURL := generateSignupURL(c, cfg)

	// Step 3: Simulate enterprise creation (in real world, this happens after admin completes signup)
	fmt.Println("\n--- Simulating Enterprise Creation ---")
	fmt.Printf("In a real scenario, the admin would visit: %s\n", signupURL.URL)
	fmt.Println("After completing signup, they would be redirected to your callback URL with a token.")
	fmt.Println("You would then use that token to create the enterprise.")

	// For demo purposes, let's work with an existing enterprise if available
	enterprises, err := c.Enterprises().List(nil)
	if err != nil {
		log.Printf("Failed to list enterprises: %v", err)
		return
	}

	var enterprise *types.Enterprise
	if len(enterprises.Items) > 0 {
		enterprise = &enterprises.Items[0]
		fmt.Printf("✓ Using existing enterprise: %s\n", enterprise.DisplayName)
	} else {
		fmt.Println("No existing enterprises found. In a real scenario, create one using the signup flow.")
		return
	}

	// Step 4: Configure enterprise settings
	configureEnterprise(c, enterprise)

	// Step 5: Create policies for different device types
	policies := createPolicies(c, enterprise.GetID())

	// Step 6: Create enrollment tokens
	tokens := createEnrollmentTokens(c, enterprise.GetID(), policies)

	// Step 7: Display setup summary
	displaySetupSummary(enterprise, policies, tokens)
}

// generateSignupURL creates a signup URL for enterprise registration.
func generateSignupURL(c *client.Client, cfg *config.Config) *types.EnterpriseSignupURL {
	fmt.Println("\n--- Step 1: Generate Signup URL ---")

	req := &types.SignupURLRequest{
		ProjectID:             cfg.ProjectID,
		CallbackURL:           cfg.CallbackURL,
		AdminEmail:            "admin@example.com",
		EnterpriseDisplayName: "Example Corporation",
		Locale:                "en_US",
	}

	signupURL, err := c.Enterprises().GenerateSignupURL(req)
	if err != nil {
		log.Printf("Failed to generate signup URL: %v", err)
		return nil
	}

	fmt.Printf("✓ Signup URL generated: %s\n", signupURL.URL)
	fmt.Printf("✓ Callback URL: %s\n", signupURL.CallbackURL)
	fmt.Printf("✓ Project ID: %s\n", signupURL.ProjectID)

	return signupURL
}

// configureEnterprise sets up enterprise-level settings.
func configureEnterprise(c *client.Client, enterprise *types.Enterprise) {
	fmt.Println("\n--- Step 2: Configure Enterprise Settings ---")

	// Enable notifications
	notificationTypes := []string{
		types.NotificationTypeEnrollment,
		types.NotificationTypeComplianceReport,
		types.NotificationTypeStatusReport,
		types.NotificationTypeCommand,
	}

	updated, err := c.Enterprises().EnableNotifications(enterprise.Name, notificationTypes)
	if err != nil {
		log.Printf("Failed to enable notifications: %v", err)
	} else {
		fmt.Printf("✓ Enabled %d notification types\n", len(updated.EnabledNotificationTypes))
	}

	// Configure Pub/Sub topic (if available)
	if enterprise.PubsubTopic == "" {
		topicName := fmt.Sprintf("projects/%s/topics/amapi-events", enterprise.GetEnterpriseID())
		_, err := c.Enterprises().SetPubSubTopic(enterprise.Name, topicName)
		if err != nil {
			log.Printf("Failed to set Pub/Sub topic: %v", err)
		} else {
			fmt.Printf("✓ Pub/Sub topic configured: %s\n", topicName)
		}
	}

	// Update enterprise display information
	updateReq := &types.EnterpriseUpdateRequest{
		DisplayName: "Example Corporation - Updated",
		ContactInfo: &types.ContactInfo{
			ContactEmail:                   "admin@example.com",
			DataProtectionOfficerName:      "John Doe",
			DataProtectionOfficerEmail:     "dpo@example.com",
			DataProtectionOfficerPhone:     "+1-555-0123",
		},
	}

	updated, err = c.Enterprises().Update(enterprise.Name, updateReq)
	if err != nil {
		log.Printf("Failed to update enterprise: %v", err)
	} else {
		fmt.Printf("✓ Enterprise information updated: %s\n", updated.DisplayName)
	}
}

// createPolicies creates different policy types for various use cases.
func createPolicies(c *client.Client, enterpriseID string) []*types.Policy {
	fmt.Println("\n--- Step 3: Create Policies ---")

	var createdPolicies []*types.Policy

	// Policy 1: Fully Managed Corporate Devices
	fullyManagedPreset := presets.GetFullyManagedPreset()
	fullyManagedPolicy, err := c.Policies().CreateByEnterpriseID(
		enterpriseID,
		"corporate-fully-managed",
		fullyManagedPreset.Policy,
	)
	if err != nil {
		log.Printf("Failed to create fully managed policy: %v", err)
	} else {
		fmt.Printf("✓ Created fully managed policy: %s\n", fullyManagedPolicy.GetID())
		createdPolicies = append(createdPolicies, fullyManagedPolicy)
	}

	// Policy 2: BYOD Work Profile
	workProfilePreset := presets.GetWorkProfilePreset()
	workProfilePolicy, err := c.Policies().CreateByEnterpriseID(
		enterpriseID,
		"byod-work-profile",
		workProfilePreset.Policy,
	)
	if err != nil {
		log.Printf("Failed to create work profile policy: %v", err)
	} else {
		fmt.Printf("✓ Created work profile policy: %s\n", workProfilePolicy.GetID())
		createdPolicies = append(createdPolicies, workProfilePolicy)
	}

	// Policy 3: Kiosk Mode for Retail
	kioskPreset := presets.GetRetailKioskPreset()
	// Add a sample retail app
	retailApp := types.NewKioskApp("com.example.retailapp")
	kioskPreset.Policy.AddApplication(retailApp)

	kioskPolicy, err := c.Policies().CreateByEnterpriseID(
		enterpriseID,
		"retail-kiosk",
		kioskPreset.Policy,
	)
	if err != nil {
		log.Printf("Failed to create kiosk policy: %v", err)
	} else {
		fmt.Printf("✓ Created kiosk policy: %s\n", kioskPolicy.GetID())
		createdPolicies = append(createdPolicies, kioskPolicy)
	}

	// Policy 4: Secure Workstation
	securePreset := presets.GetSecureWorkstationPreset()
	// Add required corporate applications
	securePreset.Policy.AddApplication(types.NewRequiredApp("com.company.vpn"))
	securePreset.Policy.AddApplication(types.NewRequiredApp("com.company.security"))

	securePolicy, err := c.Policies().CreateByEnterpriseID(
		enterpriseID,
		"secure-workstation",
		securePreset.Policy,
	)
	if err != nil {
		log.Printf("Failed to create secure workstation policy: %v", err)
	} else {
		fmt.Printf("✓ Created secure workstation policy: %s\n", securePolicy.GetID())
		createdPolicies = append(createdPolicies, securePolicy)
	}

	// Policy 5: Custom Education Policy
	educationPolicy := createCustomEducationPolicy(c, enterpriseID)
	if educationPolicy != nil {
		createdPolicies = append(createdPolicies, educationPolicy)
	}

	return createdPolicies
}

// createCustomEducationPolicy creates a custom policy for educational devices.
func createCustomEducationPolicy(c *client.Client, enterpriseID string) *types.Policy {
	// Start with education preset
	educationPreset := presets.GetEducationTabletPreset()
	policy := educationPreset.Policy.Clone()

	// Add educational applications
	educationalApps := []string{
		"com.google.android.apps.classroom",
		"com.google.android.apps.docs.editors.docs",
		"com.google.android.apps.docs.editors.sheets",
		"com.google.android.apps.docs.editors.slides",
		"com.google.android.youtube",
		"com.microsoft.office.word",
		"com.adobe.reader",
	}

	for _, appPackage := range educationalApps {
		policy.AddApplication(types.NewRequiredApp(appPackage))
	}

	// Block social media and games
	blockedApps := []string{
		"com.facebook.katana",
		"com.instagram.android",
		"com.snapchat.android",
		"com.twitter.android",
	}

	for _, appPackage := range blockedApps {
		policy.AddApplication(types.NewBlockedApp(appPackage))
	}

	// Create the policy
	educationPolicy, err := c.Policies().CreateByEnterpriseID(
		enterpriseID,
		"education-tablet-custom",
		policy,
	)
	if err != nil {
		log.Printf("Failed to create education policy: %v", err)
		return nil
	}

	fmt.Printf("✓ Created custom education policy: %s\n", educationPolicy.GetID())
	return educationPolicy
}

// createEnrollmentTokens creates enrollment tokens for each policy.
func createEnrollmentTokens(c *client.Client, enterpriseID string, policies []*types.Policy) []*types.EnrollmentToken {
	fmt.Println("\n--- Step 4: Create Enrollment Tokens ---")

	var tokens []*types.EnrollmentToken

	for _, policy := range policies {
		// Create different token types based on policy
		var duration time.Duration
		var allowPersonalUsage bool

		switch policy.GetPolicyMode() {
		case types.PolicyModeWorkProfile:
			duration = 7 * 24 * time.Hour // 1 week for BYOD
			allowPersonalUsage = true
		case types.PolicyModeDedicated:
			duration = 1 * time.Hour // 1 hour for kiosk (short-lived)
			allowPersonalUsage = false
		default:
			duration = 24 * time.Hour // 1 day for corporate devices
			allowPersonalUsage = false
		}

		// Create enrollment token
		req := &types.EnrollmentTokenCreateRequest{
			EnterpriseName:     fmt.Sprintf("enterprises/%s", enterpriseID),
			PolicyName:         policy.Name,
			Duration:           duration,
			AllowPersonalUsage: allowPersonalUsage,
			OneTimeOnly:        false,
		}

		token, err := c.EnrollmentTokens().Create(req)
		if err != nil {
			log.Printf("Failed to create token for policy %s: %v", policy.GetID(), err)
			continue
		}

		fmt.Printf("✓ Created enrollment token for %s (valid for %v)\n",
			policy.GetID(), duration)

		// Generate QR code for each token
		qrOptions := &types.QRCodeOptions{
			SkipSetupWizard: true,
			Locale:          "en_US",
		}

		// Add WiFi configuration for non-kiosk devices
		if policy.GetPolicyMode() != types.PolicyModeDedicated {
			qrOptions.WiFiSSID = "CorpWiFi"
			qrOptions.WiFiPassword = "CompanyPassword123"
			qrOptions.WiFiSecurityType = types.WiFiSecurityTypeWPA2
		}

		qrData, err := c.EnrollmentTokens().GenerateQRCode(token.Name, qrOptions)
		if err != nil {
			log.Printf("Failed to generate QR code: %v", err)
		} else {
			fmt.Printf("  ✓ QR code generated for %s\n", token.GetID())
		}

		tokens = append(tokens, token)
	}

	return tokens
}

// displaySetupSummary shows a summary of the completed setup.
func displaySetupSummary(enterprise *types.Enterprise, policies []*types.Policy, tokens []*types.EnrollmentToken) {
	fmt.Println("\n=== Setup Summary ===")

	fmt.Printf("Enterprise: %s (%s)\n", enterprise.DisplayName, enterprise.GetID())
	fmt.Printf("Notification Types: %d enabled\n", len(enterprise.EnabledNotificationTypes))

	fmt.Printf("\nPolicies Created: %d\n", len(policies))
	for _, policy := range policies {
		fmt.Printf("  - %s (%s mode)\n", policy.GetID(), policy.GetPolicyMode())
		fmt.Printf("    Applications: %d configured\n", len(policy.Applications))
		fmt.Printf("    Compliance Rules: %d defined\n", len(policy.ComplianceRules))
	}

	fmt.Printf("\nEnrollment Tokens: %d created\n", len(tokens))
	for _, token := range tokens {
		fmt.Printf("  - %s\n", token.GetID())
		fmt.Printf("    Policy: %s\n", token.GetPolicyID())
		fmt.Printf("    Expires: %v\n", token.TimeUntilExpiration())
		fmt.Printf("    Personal Usage: %t\n", token.AllowPersonalUsage)
	}

	fmt.Println("\n=== Next Steps ===")
	fmt.Println("1. Distribute enrollment tokens/QR codes to device administrators")
	fmt.Println("2. Monitor device enrollment in the console")
	fmt.Println("3. Review compliance reports as devices check in")
	fmt.Println("4. Adjust policies based on organizational needs")
	fmt.Println("5. Set up monitoring and alerting for policy violations")

	fmt.Println("\n=== Device Enrollment Instructions ===")
	fmt.Println("For new device enrollment:")
	fmt.Println("1. Factory reset the device")
	fmt.Println("2. During setup, scan the QR code or enter enrollment token")
	fmt.Println("3. Follow the guided setup process")
	fmt.Println("4. Device will automatically apply the assigned policy")

	fmt.Println("\n=== Monitoring and Management ===")
	fmt.Println("- Use the device list API to monitor enrolled devices")
	fmt.Println("- Set up webhook endpoints to receive real-time notifications")
	fmt.Println("- Regularly review compliance reports")
	fmt.Println("- Update policies as business requirements change")
}

// demonstrateOngoingManagement shows examples of ongoing enterprise management tasks.
func demonstrateOngoingManagement() {
	fmt.Println("\n=== Ongoing Management Examples ===")

	cfg, _ := config.AutoLoadConfig()
	c, _ := client.New(cfg)
	defer c.Close()

	enterprises, _ := c.Enterprises().List(nil)
	if len(enterprises.Items) == 0 {
		return
	}

	enterpriseID := enterprises.Items[0].GetID()

	// Example 1: Monitor device compliance
	fmt.Println("\n--- Compliance Monitoring ---")
	compliantDevices, _ := c.Devices().GetCompliantDevices(enterpriseID)
	nonCompliantDevices, _ := c.Devices().GetNonCompliantDevices(enterpriseID)

	fmt.Printf("✓ Compliant devices: %d\n", len(compliantDevices.Items))
	fmt.Printf("✓ Non-compliant devices: %d\n", len(nonCompliantDevices.Items))

	// Example 2: Token management
	fmt.Println("\n--- Token Management ---")
	stats, _ := c.EnrollmentTokens().GetTokenStatistics(enterpriseID)
	fmt.Printf("✓ Token statistics: %+v\n", stats)

	// Example 3: Policy updates
	fmt.Println("\n--- Policy Updates ---")
	policies, _ := c.Policies().ListByEnterpriseID(enterpriseID, nil)
	for _, policy := range policies.Items {
		devicesUsingPolicy, _ := c.Policies().GetDevicesUsingPolicy(policy.Name)
		fmt.Printf("✓ Policy %s: %d devices\n", policy.GetID(), len(devicesUsingPolicy.Items))
	}
}