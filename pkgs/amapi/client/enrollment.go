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
func (es *EnrollmentService) Create(req *types.EnrollmentTokenCreateRequest) (*types.EnrollmentToken, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enrollment token create request is required")
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Create enrollment token object
	token := &androidmanagement.EnrollmentToken{
		PolicyName:  req.PolicyName,
		OneTimeOnly: req.OneTimeOnly,
	}

	// Set AllowPersonalUsage based on bool value
	if req.AllowPersonalUsage {
		token.AllowPersonalUsage = "PERSONAL_USAGE_ALLOWED"
	} else {
		token.AllowPersonalUsage = "PERSONAL_USAGE_DISALLOWED"
	}

	// Set duration
	if req.Duration > 0 {
		token.Duration = req.Duration.String()
	}

	// Set user information
	if req.User != nil {
		token.User = req.User.ToAMAPIUser()
	}

	var result *androidmanagement.EnrollmentToken
	var err error

	err = es.client.executeAPICall(func() error {
		result, err = es.client.service.Enterprises.EnrollmentTokens.Create(req.EnterpriseName, token).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "create enrollment token")
	}

	return types.FromAMAPIEnrollmentToken(result), nil
}

// CreateByEnterpriseID creates a new enrollment token using enterprise ID.
func (es *EnrollmentService) CreateByEnterpriseID(enterpriseID, policyID string, duration time.Duration) (*types.EnrollmentToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validatePolicyID(policyID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	policyName := buildPolicyName(enterpriseID, policyID)

	req := &types.EnrollmentTokenCreateRequest{
		EnterpriseName: enterpriseName,
		PolicyName:     policyName,
		Duration:       duration,
	}

	return es.Create(req)
}

// CreateQuick creates an enrollment token with default settings.
func (es *EnrollmentService) CreateQuick(enterpriseID, policyID string) (*types.EnrollmentToken, error) {
	return es.CreateByEnterpriseID(enterpriseID, policyID, 24*time.Hour)
}

// CreateForWorkProfile creates an enrollment token for work profile enrollment.
func (es *EnrollmentService) CreateForWorkProfile(enterpriseID, policyID string, duration time.Duration) (*types.EnrollmentToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validatePolicyID(policyID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	policyName := buildPolicyName(enterpriseID, policyID)

	req := &types.EnrollmentTokenCreateRequest{
		EnterpriseName:     enterpriseName,
		PolicyName:         policyName,
		Duration:           duration,
		AllowPersonalUsage: true, // Enable personal usage for work profile
	}

	return es.Create(req)
}

// Get retrieves an enrollment token by its resource name.
func (es *EnrollmentService) Get(tokenName string) (*types.EnrollmentToken, error) {
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

	return types.FromAMAPIEnrollmentToken(result), nil
}

// GetByID retrieves an enrollment token by enterprise ID and token ID.
func (es *EnrollmentService) GetByID(enterpriseID, tokenID string) (*types.EnrollmentToken, error) {
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
func (es *EnrollmentService) List(req *types.EnrollmentTokenListRequest) (*types.ListResult[types.EnrollmentToken], error) {
	if req == nil || req.EnterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	var result *androidmanagement.ListEnrollmentTokensResponse
	var err error

	err = es.client.executeAPICall(func() error {
		call := es.client.service.Enterprises.EnrollmentTokens.List(req.EnterpriseName)

		if req.PageSize > 0 {
			call.PageSize(int64(req.PageSize))
		}

		if req.PageToken != "" {
			call.PageToken(req.PageToken)
		}

		result, err = call.Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "list enrollment tokens")
	}

	// Convert results
	tokens := make([]types.EnrollmentToken, len(result.EnrollmentTokens))
	for i, token := range result.EnrollmentTokens {
		tokens[i] = *types.FromAMAPIEnrollmentToken(token)
	}

	// Apply client-side filtering
	if req.PolicyName != "" || !req.IncludeExpired {
		filteredTokens := make([]types.EnrollmentToken, 0)
		for _, token := range tokens {
			// Filter by policy name
			if req.PolicyName != "" && token.PolicyName != req.PolicyName {
				continue
			}

			// Filter expired tokens if requested
			if !req.IncludeExpired && token.IsExpired() {
				continue
			}

			filteredTokens = append(filteredTokens, token)
		}
		tokens = filteredTokens
	}

	return &types.ListResult[types.EnrollmentToken]{
		Items:         tokens,
		NextPageToken: result.NextPageToken,
	}, nil
}

// ListByEnterpriseID lists enrollment tokens for an enterprise by enterprise ID.
func (es *EnrollmentService) ListByEnterpriseID(enterpriseID string, options *types.ListOptions) (*types.ListResult[types.EnrollmentToken], error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	req := &types.EnrollmentTokenListRequest{
		EnterpriseName: enterpriseName,
	}

	if options != nil {
		req.ListOptions = *options
	}

	return es.List(req)
}

// Delete deletes an enrollment token.
func (es *EnrollmentService) Delete(req *types.EnrollmentTokenDeleteRequest) error {
	if req == nil || req.Name == "" {
		return types.ErrInvalidTokenID
	}

	err := es.client.executeAPICall(func() error {
		_, err := es.client.service.Enterprises.EnrollmentTokens.Delete(req.Name).Context(es.client.ctx).Do()
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
	req := &types.EnrollmentTokenDeleteRequest{
		Name: tokenName,
	}

	return es.Delete(req)
}

// GenerateQRCode generates QR code data for an enrollment token.
func (es *EnrollmentService) GenerateQRCode(tokenName string, options *types.QRCodeOptions) (*types.QRCodeData, error) {
	// Get the enrollment token
	token, err := es.Get(tokenName)
	if err != nil {
		return nil, err
	}

	// Check if token is expired
	if token.IsExpired() {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enrollment token has expired")
	}

	// Generate QR code data
	return token.GenerateQRCodeData(options), nil
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
func (es *EnrollmentService) GetActiveTokens(enterpriseID string) (*types.ListResult[types.EnrollmentToken], error) {
	req := &types.EnrollmentTokenListRequest{
		EnterpriseName: buildEnterpriseName(enterpriseID),
		IncludeExpired: false,
	}

	return es.List(req)
}

// GetTokensForPolicy returns all enrollment tokens for a specific policy.
func (es *EnrollmentService) GetTokensForPolicy(enterpriseID, policyID string) (*types.ListResult[types.EnrollmentToken], error) {
	policyName := buildPolicyName(enterpriseID, policyID)
	req := &types.EnrollmentTokenListRequest{
		EnterpriseName: buildEnterpriseName(enterpriseID),
		PolicyName:     policyName,
		IncludeExpired: false,
	}

	return es.List(req)
}

// RevokeToken revokes an enrollment token by deleting it.
func (es *EnrollmentService) RevokeToken(tokenName string) error {
	req := &types.EnrollmentTokenDeleteRequest{
		Name: tokenName,
	}

	return es.Delete(req)
}

// RevokeTokenByID revokes an enrollment token by enterprise ID and token ID.
func (es *EnrollmentService) RevokeTokenByID(enterpriseID, tokenID string) error {
	return es.DeleteByID(enterpriseID, tokenID)
}

// CreateWithQRCode creates an enrollment token and generates QR code data.
func (es *EnrollmentService) CreateWithQRCode(req *types.EnrollmentTokenCreateRequest, qrOptions *types.QRCodeOptions) (*types.EnrollmentToken, *types.QRCodeData, error) {
	// Create the enrollment token
	token, err := es.Create(req)
	if err != nil {
		return nil, nil, err
	}

	// Generate QR code data
	qrData := token.GenerateQRCodeData(qrOptions)

	return token, qrData, nil
}

// CreateBulkTokens creates multiple enrollment tokens for the same policy.
func (es *EnrollmentService) CreateBulkTokens(enterpriseID, policyID string, count int, duration time.Duration) ([]*types.EnrollmentToken, error) {
	if count <= 0 {
		return nil, types.NewError(types.ErrCodeInvalidInput, "count must be positive")
	}

	if count > 100 {
		return nil, types.NewError(types.ErrCodeInvalidInput, "count cannot exceed 100")
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	policyName := buildPolicyName(enterpriseID, policyID)

	tokens := make([]*types.EnrollmentToken, 0, count)

	for i := 0; i < count; i++ {
		req := &types.EnrollmentTokenCreateRequest{
			EnterpriseName: enterpriseName,
			PolicyName:     policyName,
			Duration:       duration,
		}

		token, err := es.Create(req)
		if err != nil {
			return tokens, err // Return partial results with error
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// ExtendTokenExpiration extends the expiration of an enrollment token by creating a new one.
// Note: The API doesn't support extending existing tokens, so this creates a new token.
func (es *EnrollmentService) ExtendTokenExpiration(tokenName string, newDuration time.Duration) (*types.EnrollmentToken, error) {
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
	req := &types.EnrollmentTokenCreateRequest{
		EnterpriseName:     buildEnterpriseName(enterpriseID),
		PolicyName:         existingToken.PolicyName,
		Duration:           newDuration,
		AllowPersonalUsage: existingToken.AllowPersonalUsage,
		OneTimeOnly:        existingToken.OneTimeOnly,
	}

	// Create new token
	newToken, err := es.Create(req)
	if err != nil {
		return nil, err
	}

	// Delete the old token
	_ = es.RevokeToken(tokenName) // Ignore error if deletion fails

	return newToken, nil
}

// GetTokenStatistics returns statistics about enrollment tokens for an enterprise.
func (es *EnrollmentService) GetTokenStatistics(enterpriseID string) (map[string]int, error) {
	tokens, err := es.ListByEnterpriseID(enterpriseID, nil)
	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"total":   len(tokens.Items),
		"active":  0,
		"expired": 0,
	}

	for _, token := range tokens.Items {
		if token.IsExpired() {
			stats["expired"]++
		} else {
			stats["active"]++
		}
	}

	return stats, nil
}
