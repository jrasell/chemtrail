package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyProviderAWSASGEnabled = "provider-aws-asg-enabled"
)

type ProviderConfig struct {
	AWSASG bool
}

func GetProviderConfig() *ProviderConfig {
	return &ProviderConfig{
		AWSASG: viper.GetBool(configKeyProviderAWSASGEnabled),
	}
}

func RegisterProviderConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyProviderAWSASGEnabled
			longOpt      = "provider-aws-asg-enabled"
			defaultValue = false
			description  = "Enable the AWS AutoScaling Group client provider"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
