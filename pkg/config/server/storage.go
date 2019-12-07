package server

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyStorageConsulEnabled     = "storage-consul-enabled"
	configKeyStorageConsulPath        = "storage-consul-path"
	configKeyStorageConsulPathDefault = "chemtrail/"
)

// StorageConfig is the CLI configuration options for the state storage backend.
type StorageConfig struct {
	ConsulEnabled bool
	ConsulPath    string
}

// GetStorageConfig populates a StorageConfig object with the users CLI parameters.
func GetStorageConfig() *StorageConfig {

	// Check that the path has suffixed with a forward slash, otherwise put this on so we do not
	// need to check in a number of other places.
	path := viper.GetString(configKeyStorageConsulPath)
	if suffix := strings.HasSuffix(path, "/"); !suffix {
		path = path + "/"
	}

	return &StorageConfig{
		ConsulEnabled: viper.GetBool(configKeyStorageConsulEnabled),
		ConsulPath:    path,
	}
}

// RegisterStorageConfig register the storage CLI parameters for the state storage backend.
func RegisterStorageConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyStorageConsulEnabled
			longOpt      = "storage-consul-enabled"
			defaultValue = false
			description  = "Enable the Consul state storage backend"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
	{
		const (
			key          = configKeyStorageConsulPath
			longOpt      = "storage-consul-path"
			defaultValue = configKeyStorageConsulPathDefault
			description  = "The Consul KV base path that will be used to store state"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
