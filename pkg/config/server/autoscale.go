package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyAutoscalerThreadNumberDefault       = 3
	configKeyAutoscalerEvaluationIntervalDefault = 180

	configKeyAutoscalerEnabled            = "autoscaler-enabled"
	configKeyAutoscalerEvaluationInterval = "autoscaler-evaluation-interval"
	configKeyAutoscalerThreadNumber       = "autoscaler-num-threads"
)

type AutoscalerConfig struct {
	Enabled  bool
	Interval int
	Threads  int
}

func GetAutoscalerConfig() *AutoscalerConfig {
	return &AutoscalerConfig{
		Enabled:  viper.GetBool(configKeyAutoscalerEnabled),
		Interval: viper.GetInt(configKeyAutoscalerEvaluationInterval),
		Threads:  viper.GetInt(configKeyAutoscalerThreadNumber),
	}
}

func RegisterAutoscalerConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyAutoscalerEnabled
			longOpt      = "autoscaler-enabled"
			defaultValue = false
			description  = "Enable the internal autoscaling engine"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyAutoscalerEvaluationInterval
			longOpt      = "autoscaler-evaluation-interval"
			defaultValue = configKeyAutoscalerEvaluationIntervalDefault
			description  = "The time period in seconds between autoscaling evaluation runs"
		)

		flags.Int(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyAutoscalerThreadNumber
			longOpt      = "autoscaler-num-threads"
			defaultValue = configKeyAutoscalerThreadNumberDefault
			description  = "Specifies the number of parallel autoscaler threads to run"
		)

		flags.Int(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
