package client

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyChemtrailAddr        = "addr"
	configKeyChemtrailAddrDefault = "http://127.0.0.1:8000"

	configKeyChemtrailClientCertPath    = "client-cert-path"
	configKeyChemtrailClientCertKeyPath = "client-cert-key-path"
	configKeyChemtrailCAPath            = "client-ca-path"
)

type Config struct {
	Addr        string
	CertPath    string
	CertKeyPath string
	CAPath      string
}

func GetConfig() Config {
	return Config{
		Addr:        viper.GetString(configKeyChemtrailAddr),
		CertPath:    viper.GetString(configKeyChemtrailClientCertPath),
		CertKeyPath: viper.GetString(configKeyChemtrailClientCertKeyPath),
		CAPath:      viper.GetString(configKeyChemtrailCAPath),
	}
}

func RegisterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyChemtrailAddr
			longOpt      = "addr"
			defaultValue = configKeyChemtrailAddrDefault
			description  = "The HTTP(S) address of the Chemtrail server"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyChemtrailClientCertPath
			longOpt      = "client-cert-path"
			defaultValue = ""
			description  = "Path to a PEM encoded client certificate for TLS authentication to the Chemtrail server"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyChemtrailClientCertKeyPath
			longOpt      = "client-cert-key-path"
			defaultValue = ""
			description  = "Path to an unencrypted PEM encoded private key matching the client certificate"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyChemtrailCAPath
			longOpt      = "client-ca-path"
			defaultValue = ""
			description  = "Path to a PEM encoded CA cert file to use to verify the Chemtrail server SSL certificate"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
