package main

import (
	"context"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	FusionEnvironmentID  string `mapstructure:"fusion-environment-id"`
	PathToConfigFile     string `mapstructure:"path-to-config-file"`
	Profile              string `mapstructure:"profile"`
	Tenancy              string `mapstructure:"tenancy"`
	User                 string `mapstructure:"user"`
	Region               string `mapstructure:"region"`
	Fingerprint          string `mapstructure:"fingerprint"`
	PrivateKey           string `mapstructure:"private-key"`
	PrivateKeyPassphrase string `mapstructure:"private-key-passphrase"`
}

// getConfigProvider returns a configuration provider based on the configuration provided.
func (cfg *config) GetConfigProvider() (*common.ConfigurationProvider, error) {
	configProvider := common.DefaultConfigProvider()

	if cfg.PathToConfigFile != "" {
		var err error
		switch {
		// Attempt to use a more detailed configuration when both profile and private key are provided.
		case cfg.Profile != "" && cfg.PrivateKeyPassphrase != "":
			configProvider, err = common.ConfigurationProviderFromFileWithProfile(cfg.PathToConfigFile, cfg.Profile, cfg.PrivateKeyPassphrase)
			if err != nil {
				return nil, err
			}

		// If only private key is specified, use a basic file configuration.
		case cfg.PrivateKeyPassphrase != "":
			configProvider, err = common.ConfigurationProviderFromFile(cfg.PathToConfigFile, cfg.PrivateKeyPassphrase)
			if err != nil {
				return nil, err
			}

		// If only profile is specified, use the custom profile configuration.
		case cfg.Profile != "":
			configProvider = common.CustomProfileConfigProvider(cfg.PathToConfigFile, cfg.Profile)

		// If only the path to the config file is specified, try to pass empty passphrase.
		default:
			configProvider, err = common.ConfigurationProviderFromFile(cfg.PathToConfigFile, "")
			if err != nil {
				return nil, err
			}
		}
	}

	// Check if all parameters needed for raw configuration are present.
	if cfg.Tenancy != "" && cfg.User != "" && cfg.Region != "" && cfg.Fingerprint != "" && cfg.PrivateKey != "" && cfg.PrivateKeyPassphrase != "" {
		configProvider = common.NewRawConfigurationProvider(cfg.Tenancy, cfg.User, cfg.Region, cfg.Fingerprint, cfg.PrivateKey, &cfg.PrivateKeyPassphrase)
	}

	return &configProvider, nil
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.FusionEnvironmentID == "" {
		return status.Error(codes.InvalidArgument, "environment ID of the Fusion Cloud environment must be provided, use --help for more information")
	}

	return nil
}

func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		"fusion-environment-id",
		"",
		"The environment ID for the Fusion Cloud environment. For example: 'ocid1.test.oc1..<unique_ID>EXAMPLE-fusionEnvironmentId-Value'. ($BATON_FUSION_ENVIRONMENT_ID)",
	)

	// flags for setting up the config provider - more information can be found at https://docs.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm
	cmd.PersistentFlags().String("path-to-config-file", "", "The path to the OCI config file. ($BATON_PATH_TO_CONFIG_FILE)")
	cmd.PersistentFlags().String("profile", "", "The profile in the OCI config file. ($BATON_PROFILE)")

	cmd.PersistentFlags().String("tenancy", "", "The tenancy of the OCI config file. ($BATON_TENANCY)")
	cmd.PersistentFlags().String("user", "", "The user of the OCI config file. ($BATON_USER)")
	cmd.PersistentFlags().String("region", "", "The region of the OCI config file. ($BATON_REGION)")
	cmd.PersistentFlags().String("fingerprint", "", "The fingerprint of the OCI config file. ($BATON_FINGERPRINT)")

	cmd.PersistentFlags().String("private-key", "", "The private key in the OCI config file. ($BATON_PRIVATE_KEY)")
	cmd.PersistentFlags().String("private-key-passphrase", "", "The passphrase for the private key in the OCI config file. ($BATON_PRIVATE_KEY_PASSPHRASE)")
}
