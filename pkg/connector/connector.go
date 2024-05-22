package connector

import (
	"context"
	"io"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/fusionapps"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FusionCloud struct {
	client        *fusionapps.FusionApplicationsClient
	environmentID string
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (fc *FusionCloud) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(fc.client, fc.environmentID),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (fc *FusionCloud) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (fc *FusionCloud) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "OracleFusionCloud",
		Description: "Connector syncing Oracle Fusion Cloud admin users to Baton",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (fc *FusionCloud) Validate(ctx context.Context) (annotations.Annotations, error) {
	_, err := fc.client.ListAdminUsers(ctx, fusionapps.ListAdminUsersRequest{
		FusionEnvironmentId: &fc.environmentID,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to authenticate with FusionCloud")
	}

	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, configProvider *common.ConfigurationProvider, environmentID string) (*FusionCloud, error) {
	client, err := fusionapps.NewFusionApplicationsClientWithConfigurationProvider(*configProvider)
	if err != nil {
		return nil, err
	}

	return &FusionCloud{
		client:        &client,
		environmentID: environmentID,
	}, nil
}
