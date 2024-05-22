package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/oracle/oci-go-sdk/v65/fusionapps"
)

type userBuilder struct {
	client        *fusionapps.FusionApplicationsClient
	environmentID string
	resourceType  *v2.ResourceType
}

func userResource(user *fusionapps.AdminUserSummary) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}
	fullName := fmt.Sprintf("%s %s", *user.FirstName, *user.LastName)

	res, err := resource.NewUserResource(
		fullName,
		userResourceType,
		user.Username,
		[]resource.UserTraitOption{
			resource.WithUserProfile(profile),
			resource.WithEmail(*user.EmailAddress, true),
			resource.WithStatus(v2.UserTrait_Status_STATUS_ENABLED),
		},
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (u *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	users, err := u.client.ListAdminUsers(ctx, fusionapps.ListAdminUsersRequest{
		FusionEnvironmentId: &u.environmentID,
	})
	if err != nil {
		return nil, "", nil, fmt.Errorf("oracle-fusion-cloud-connector: failed to list users: %w", err)
	}

	var rv []*v2.Resource
	for _, user := range users.Items {
		ur, err := userResource(&user) // #nosec G601
		if err != nil {
			return nil, "", nil, fmt.Errorf("oracle-fusion-cloud-connector: failed to create user resource: %w", err)
		}

		rv = append(rv, ur)
	}

	return rv, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (u *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (u *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *fusionapps.FusionApplicationsClient, environmentID string) *userBuilder {
	return &userBuilder{
		client:        client,
		environmentID: environmentID,
		resourceType:  userResourceType,
	}
}
