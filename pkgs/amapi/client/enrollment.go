package client

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// EnrollmentService provides enrollment token management methods.
type EnrollmentService struct {
	client *Client
}

// EnrollmentTokens returns the enrollment service.
func (c *Client) EnrollmentTokens() *EnrollmentService {
	return &EnrollmentService{client: c}
}

// Create creates a new enrollment token.
func (es *EnrollmentService) Create(enterpriseName, policyName string, duration time.Duration, allowPersonalUsage, oneTimeOnly bool, user *androidmanagement.User) (*androidmanagement.EnrollmentToken, error) {
	if enterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	if policyName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "policy name is required")
	}

	// Create enrollment token object
	token := &androidmanagement.EnrollmentToken{
		PolicyName:  policyName,
		OneTimeOnly: oneTimeOnly,
	}

	// Set AllowPersonalUsage based on bool value
	if allowPersonalUsage {
		token.AllowPersonalUsage = "PERSONAL_USAGE_ALLOWED"
	} else {
		token.AllowPersonalUsage = "PERSONAL_USAGE_DISALLOWED"
	}

	// Set duration
	if duration > 0 {
		token.Duration = duration.String()
	}

	// Set user information
	if user != nil {
		token.User = user
	}

	var result *androidmanagement.EnrollmentToken
	var err error

	err = es.client.executeAPICall(func() error {
		result, err = es.client.service.Enterprises.EnrollmentTokens.Create(enterpriseName, token).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "create enrollment token")
	}

	return result, nil
}

// CreateByEnterpriseID creates a new enrollment token using enterprise ID.
func (es *EnrollmentService) CreateByEnterpriseID(enterpriseID, policyID string, duration time.Duration) (*androidmanagement.EnrollmentToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validatePolicyID(policyID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	policyName := buildPolicyName(enterpriseID, policyID)

	return es.Create(enterpriseName, policyName, duration, false, false, nil)
}

// CreateQuick creates an enrollment token with default settings.
func (es *EnrollmentService) CreateQuick(enterpriseID, policyID string) (*androidmanagement.EnrollmentToken, error) {
	return es.CreateByEnterpriseID(enterpriseID, policyID, 24*time.Hour)
}

// CreateForWorkProfile creates an enrollment token for work profile enrollment.
func (es *EnrollmentService) CreateForWorkProfile(enterpriseID, policyID string, duration time.Duration) (*androidmanagement.EnrollmentToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validatePolicyID(policyID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	policyName := buildPolicyName(enterpriseID, policyID)

	return es.Create(enterpriseName, policyName, duration, true, false, nil)
}

// Get retrieves an enrollment token by its resource name.
func (es *EnrollmentService) Get(tokenName string) (*androidmanagement.EnrollmentToken, error) {
	if tokenName == "" {
		return nil, types.ErrInvalidTokenID
	}

	var result *androidmanagement.EnrollmentToken
	var err error

	err = es.client.executeAPICall(func() error {
		result, err = es.client.service.Enterprises.EnrollmentTokens.Get(tokenName).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "get enrollment token")
	}

	return result, nil
}

// GetByID retrieves an enrollment token by enterprise ID and token ID.
func (es *EnrollmentService) GetByID(enterpriseID, tokenID string) (*androidmanagement.EnrollmentToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateTokenID(tokenID); err != nil {
		return nil, err
	}

	tokenName := buildEnrollmentTokenName(enterpriseID, tokenID)
	return es.Get(tokenName)
}

// List lists enrollment tokens for an enterprise.
func (es *EnrollmentService) List(enterpriseName string, pageSize int, pageToken string, policyName string, includeExpired bool) (*types.ListResult[*androidmanagement.EnrollmentToken], error) {
	if enterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	var result *androidmanagement.ListEnrollmentTokensResponse
	var err error

	err = es.client.executeAPICall(func() error {
		call := es.client.service.Enterprises.EnrollmentTokens.List(enterpriseName)

		if pageSize > 0 {
			call.PageSize(int64(pageSize))
		}

		if pageToken != "" {
			call.PageToken(pageToken)
		}

		result, err = call.Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "list enrollment tokens")
	}

	// Convert results
	tokens := make([]*androidmanagement.EnrollmentToken, len(result.EnrollmentTokens))
	copy(tokens, result.EnrollmentTokens)

	// Apply client-side filtering
	if policyName != "" || !includeExpired {
		filteredTokens := make([]*androidmanagement.EnrollmentToken, 0)
		for _, token := range tokens {
			// Filter by policy name
			if policyName != "" && token.PolicyName != policyName {
				continue
			}

			// Filter expired tokens if requested
			if !includeExpired {
				// Use helper function to check expiration
				if token.ExpirationTimestamp != "" {
					expiration, err := time.Parse(time.RFC3339, token.ExpirationTimestamp)
					if err == nil && time.Now().After(expiration) {
						continue
					}
				}
			}

			filteredTokens = append(filteredTokens, token)
		}
		tokens = filteredTokens
	}

	return &types.ListResult[*androidmanagement.EnrollmentToken]{
		Items:         tokens,
		NextPageToken: result.NextPageToken,
	}, nil
}

// ListByEnterpriseID lists enrollment tokens for an enterprise by enterprise ID.
func (es *EnrollmentService) ListByEnterpriseID(enterpriseID string, pageSize int, pageToken string, policyName string, includeExpired bool) (*types.ListResult[*androidmanagement.EnrollmentToken], error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	return es.List(enterpriseName, pageSize, pageToken, policyName, includeExpired)
}

// Delete deletes an enrollment token.
func (es *EnrollmentService) Delete(tokenName string) error {
	if tokenName == "" {
		return types.ErrInvalidTokenID
	}

	err := es.client.executeAPICall(func() error {
		_, err := es.client.service.Enterprises.EnrollmentTokens.Delete(tokenName).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return es.client.wrapAPIError(err, "delete enrollment token")
	}

	return nil
}

// DeleteByID deletes an enrollment token by enterprise ID and token ID.
func (es *EnrollmentService) DeleteByID(enterpriseID, tokenID string) error {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return err
	}

	if err := validateTokenID(tokenID); err != nil {
		return err
	}

	tokenName := buildEnrollmentTokenName(enterpriseID, tokenID)
	return es.Delete(tokenName)
}

// GenerateQRCode generates QR code data for an enrollment token.
func (es *EnrollmentService) GenerateQRCode(tokenName string, options *types.QRCodeOptions) (*types.QRCodeData, error) {
	// Get the enrollment token
	token, err := es.Get(tokenName)
	if err != nil {
		return nil, err
	}

	// Check if token is expired
	if types.IsEnrollmentTokenExpired(token) {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enrollment token has expired")
	}

	// Generate QR code data
	return types.GenerateQRCodeData(token, options), nil
}

// GenerateQRCodeByID generates QR code data for an enrollment token by IDs.
func (es *EnrollmentService) GenerateQRCodeByID(enterpriseID, tokenID string, options *types.QRCodeOptions) (*types.QRCodeData, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateTokenID(tokenID); err != nil {
		return nil, err
	}

	tokenName := buildEnrollmentTokenName(enterpriseID, tokenID)
	return es.GenerateQRCode(tokenName, options)
}

// GetActiveTokens returns all non-expired enrollment tokens for an enterprise.
func (es *EnrollmentService) GetActiveTokens(enterpriseID string) (*types.ListResult[*androidmanagement.EnrollmentToken], error) {
	enterpriseName := buildEnterpriseName(enterpriseID)
	return es.List(enterpriseName, 0, "", "", false)
}

// GetTokensForPolicy returns all enrollment tokens for a specific policy.
func (es *EnrollmentService) GetTokensForPolicy(enterpriseID, policyID string) (*types.ListResult[*androidmanagement.EnrollmentToken], error) {
	enterpriseName := buildEnterpriseName(enterpriseID)
	policyName := buildPolicyName(enterpriseID, policyID)
	return es.List(enterpriseName, 0, "", policyName, false)
}

// RevokeToken revokes an enrollment token by deleting it.
func (es *EnrollmentService) RevokeToken(tokenName string) error {
	return es.Delete(tokenName)
}

// RevokeTokenByID revokes an enrollment token by enterprise ID and token ID.
func (es *EnrollmentService) RevokeTokenByID(enterpriseID, tokenID string) error {
	return es.DeleteByID(enterpriseID, tokenID)
}

// CreateWithQRCode creates an enrollment token and generates QR code data.
func (es *EnrollmentService) CreateWithQRCode(enterpriseName, policyName string, duration time.Duration, allowPersonalUsage, oneTimeOnly bool, user *androidmanagement.User, qrOptions *types.QRCodeOptions) (*androidmanagement.EnrollmentToken, *types.QRCodeData, error) {
	// Create the enrollment token
	token, err := es.Create(enterpriseName, policyName, duration, allowPersonalUsage, oneTimeOnly, user)
	if err != nil {
		return nil, nil, err
	}

	// Generate QR code data
	qrData := types.GenerateQRCodeData(token, qrOptions)

	return token, qrData, nil
}

// CreateBulkTokens creates multiple enrollment tokens for the same policy.
func (es *EnrollmentService) CreateBulkTokens(enterpriseID, policyID string, count int, duration time.Duration) ([]*androidmanagement.EnrollmentToken, error) {
	if count <= 0 {
		return nil, types.NewError(types.ErrCodeInvalidInput, "count must be positive")
	}

	if count > 100 {
		return nil, types.NewError(types.ErrCodeInvalidInput, "count cannot exceed 100")
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	policyName := buildPolicyName(enterpriseID, policyID)

	tokens := make([]*androidmanagement.EnrollmentToken, 0, count)

	for i := 0; i < count; i++ {
		token, err := es.Create(enterpriseName, policyName, duration, false, false, nil)
		if err != nil {
			return tokens, err // Return partial results with error
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// ExtendTokenExpiration extends the expiration of an enrollment token by creating a new one.
// Note: The API doesn't support extending existing tokens, so this creates a new token.
func (es *EnrollmentService) ExtendTokenExpiration(tokenName string, newDuration time.Duration) (*androidmanagement.EnrollmentToken, error) {
	// Get the existing token
	existingToken, err := es.Get(tokenName)
	if err != nil {
		return nil, err
	}

	// Extract enterprise ID from token name
	enterpriseID, _, err := parseEnrollmentTokenName(tokenName)
	if err != nil {
		return nil, err
	}

	// Create a new token with the same policy but new duration
	enterpriseName := buildEnterpriseName(enterpriseID)
	allowPersonalUsage := types.GetEnrollmentTokenAllowPersonalUsageBool(existingToken)

	// Create new token
	newToken, err := es.Create(enterpriseName, existingToken.PolicyName, newDuration, allowPersonalUsage, existingToken.OneTimeOnly, existingToken.User)
	if err != nil {
		return nil, err
	}

	// Delete the old token
	_ = es.RevokeToken(tokenName) // Ignore error if deletion fails

	return newToken, nil
}

// GetTokenStatistics returns statistics about enrollment tokens for an enterprise.
func (es *EnrollmentService) GetTokenStatistics(enterpriseID string) (map[string]int, error) {
	tokens, err := es.ListByEnterpriseID(enterpriseID, 0, "", "", true)
	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"total":   len(tokens.Items),
		"active":  0,
		"expired": 0,
	}

	for _, token := range tokens.Items {
		if types.IsEnrollmentTokenExpired(token) {
			stats["expired"]++
		} else {
			stats["active"]++
		}
	}

	return stats, nil
}
