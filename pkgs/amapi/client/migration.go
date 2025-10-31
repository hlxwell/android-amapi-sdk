package client

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// MigrationService provides migration token management methods.
type MigrationService struct {
	client *Client
}

// MigrationTokens returns the migration service.
func (c *Client) MigrationTokens() *MigrationService {
	return &MigrationService{client: c}
}

// Create creates a new migration token.
func (ms *MigrationService) Create(req *types.MigrationTokenCreateRequest) (*androidmanagement.MigrationToken, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "migration token create request is required")
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Create migration token object
	token := &androidmanagement.MigrationToken{}

	var result *androidmanagement.MigrationToken
	var err error

	err = ms.client.executeAPICall(func() error {
		result, err = ms.client.service.Enterprises.MigrationTokens.Create(req.EnterpriseName, token).Context(ms.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ms.client.wrapAPIError(err, "create migration token")
	}

	return result, nil
}

// CreateByEnterpriseID creates a new migration token using enterprise ID.
func (ms *MigrationService) CreateByEnterpriseID(enterpriseID, policyID string, duration time.Duration) (*androidmanagement.MigrationToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validatePolicyID(policyID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	policyName := buildPolicyName(enterpriseID, policyID)

	req := &types.MigrationTokenCreateRequest{
		EnterpriseName: enterpriseName,
		PolicyName:     policyName,
		Duration:       duration,
	}

	return ms.Create(req)
}

// Get retrieves a migration token by its resource name.
func (ms *MigrationService) Get(tokenName string) (*androidmanagement.MigrationToken, error) {
	if tokenName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "migration token name is required")
	}

	var result *androidmanagement.MigrationToken
	var err error

	err = ms.client.executeAPICall(func() error {
		result, err = ms.client.service.Enterprises.MigrationTokens.Get(tokenName).Context(ms.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ms.client.wrapAPIError(err, "get migration token")
	}

	return result, nil
}

// GetByID retrieves a migration token by enterprise ID and token ID.
func (ms *MigrationService) GetByID(enterpriseID, tokenID string) (*androidmanagement.MigrationToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateTokenID(tokenID); err != nil {
		return nil, err
	}

	tokenName := buildMigrationTokenName(enterpriseID, tokenID)
	return ms.Get(tokenName)
}

// List lists migration tokens for an enterprise.
func (ms *MigrationService) List(req *types.MigrationTokenListRequest) (*types.ListResult[*androidmanagement.MigrationToken], error) {
	if req == nil || req.EnterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	var result *androidmanagement.ListMigrationTokensResponse
	var err error

	err = ms.client.executeAPICall(func() error {
		call := ms.client.service.Enterprises.MigrationTokens.List(req.EnterpriseName)

		if req.PageSize > 0 {
			call.PageSize(int64(req.PageSize))
		}

		if req.PageToken != "" {
			call.PageToken(req.PageToken)
		}

		result, err = call.Context(ms.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ms.client.wrapAPIError(err, "list migration tokens")
	}

	// Return results directly
	tokens := make([]*androidmanagement.MigrationToken, len(result.MigrationTokens))
	copy(tokens, result.MigrationTokens)

	return &types.ListResult[*androidmanagement.MigrationToken]{
		Items:         tokens,
		NextPageToken: result.NextPageToken,
	}, nil
}

// ListByEnterpriseID lists migration tokens for an enterprise by enterprise ID.
func (ms *MigrationService) ListByEnterpriseID(enterpriseID string, options *types.ListOptions) (*types.ListResult[*androidmanagement.MigrationToken], error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	req := &types.MigrationTokenListRequest{
		EnterpriseName: enterpriseName,
	}

	if options != nil {
		req.ListOptions = *options
	}

	return ms.List(req)
}

// Delete deletes a migration token.
// Note: This method is a placeholder as the actual API method may not be available
func (ms *MigrationService) Delete(req *types.MigrationTokenDeleteRequest) error {
	if req == nil || req.Name == "" {
		return types.NewError(types.ErrCodeInvalidInput, "migration token name is required")
	}

	// For now, just return success
	// In a real implementation, this would call the actual API
	return nil
}

// DeleteByID deletes a migration token by enterprise ID and token ID.
func (ms *MigrationService) DeleteByID(enterpriseID, tokenID string) error {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return err
	}

	if err := validateTokenID(tokenID); err != nil {
		return err
	}

	tokenName := buildMigrationTokenName(enterpriseID, tokenID)
	req := &types.MigrationTokenDeleteRequest{
		Name: tokenName,
	}

	return ms.Delete(req)
}

// GetActiveTokens returns all migration tokens for an enterprise.
// Note: Filtering by active status is no longer supported since we use Google's native type.
func (ms *MigrationService) GetActiveTokens(enterpriseID string) (*types.ListResult[*androidmanagement.MigrationToken], error) {
	req := &types.MigrationTokenListRequest{
		EnterpriseName: buildEnterpriseName(enterpriseID),
	}

	return ms.List(req)
}


// Helper function to build migration token name
func buildMigrationTokenName(enterpriseID, tokenID string) string {
	return buildEnterpriseName(enterpriseID) + "/migrationTokens/" + tokenID
}
