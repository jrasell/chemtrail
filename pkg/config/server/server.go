package server

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyBindAddrDefault = "127.0.0.1"
	configKeyBindPortDefault = 8000
	configKeyBindAddr        = "bind-addr"
	configKeyBindPort        = "bind-port"
)

type Config struct {
	Bind string
	Port uint16
}

func (c *Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str(configKeyBindAddr, c.Bind).
		Uint16(configKeyBindPort, c.Port)
}

func GetConfig() Config {
	return Config{
		Bind: viper.GetString(configKeyBindAddr),
		Port: uint16(viper.GetInt(configKeyBindPort)),
	}
}

func RegisterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyBindAddr
			longOpt      = "bind-addr"
			defaultValue = configKeyBindAddrDefault
			description  = "The HTTP server address to bind to"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyBindPort
			longOpt      = "bind-port"
			defaultValue = configKeyBindPortDefault
			description  = "The HTTP server port to bind to"
		)

		flags.Uint16(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
