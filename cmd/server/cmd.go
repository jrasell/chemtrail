package server

import (
	"fmt"
	"os"

	logCfg "github.com/jrasell/chemtrail/pkg/config/log"
	serverCfg "github.com/jrasell/chemtrail/pkg/config/server"
	"github.com/jrasell/chemtrail/pkg/logger"
	"github.com/jrasell/chemtrail/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start a Chemtrail server",
		Run: func(cmd *cobra.Command, args []string) {
			runServer(cmd, args)
		},
	}
	serverCfg.RegisterConfig(cmd)
	serverCfg.RegisterTLSConfig(cmd)
	serverCfg.RegisterTelemetryConfig(cmd)
	serverCfg.RegisterProviderConfig(cmd)
	serverCfg.RegisterAutoscalerConfig(cmd)
	serverCfg.RegisterStorageConfig(cmd)
	logCfg.RegisterConfig(cmd)
	rootCmd.AddCommand(cmd)

	return nil
}

func runServer(_ *cobra.Command, _ []string) {
	autoscaleConfig := serverCfg.GetAutoscalerConfig()
	providerConfig := serverCfg.GetProviderConfig()
	storageConfig := serverCfg.GetStorageConfig()
	serverConfig := serverCfg.GetConfig()
	tlsConfig := serverCfg.GetTLSConfig()
	telemetryConfig := serverCfg.GetTelemetryConfig()

	// Setup the server logging.
	logConfig := logCfg.GetConfig()
	if err := logger.Setup(logConfig); err != nil {
		fmt.Println(err)
		os.Exit(sysexits.Software)
	}

	cfg := &server.Config{
		Autoscale: autoscaleConfig,
		Provider:  providerConfig,
		Server:    &serverConfig,
		Storage:   storageConfig,
		TLS:       &tlsConfig,
		Telemetry: &telemetryConfig,
	}
	srv := server.New(log.Logger, cfg)

	if err := srv.Start(); err != nil {
		os.Exit(sysexits.Software)
	}
}
