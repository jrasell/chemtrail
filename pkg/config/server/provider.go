package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyProviderAWSASGEnabled = "provider-aws-asg-enabled"
	configKeyProviderNoOpEnabled   = "provider-noop-enabled"
)

type ProviderConfig struct {
	AWSASG bool
	NoOp   bool
}

func GetProviderConfig() *ProviderConfig {
	return &ProviderConfig{
		AWSASG: viper.GetBool(configKeyProviderAWSASGEnabled),
		NoOp:   viper.GetBool(configKeyProviderNoOpEnabled),
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
	{
		const (
			key          = configKeyProviderNoOpEnabled
			longOpt      = "provider-noop-enabled"
			defaultValue = true
			description  = "Enable the NoOp client provider"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
