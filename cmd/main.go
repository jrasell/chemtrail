package main

import (
	"fmt"
	"os"

	"github.com/jrasell/chemtrail/cmd/policy"

	"github.com/jrasell/chemtrail/cmd/scale"
	"github.com/jrasell/chemtrail/cmd/server"
	"github.com/jrasell/chemtrail/cmd/system"
	"github.com/jrasell/chemtrail/pkg/build"
	clientCfg "github.com/jrasell/chemtrail/pkg/config/client"
	envCfg "github.com/jrasell/chemtrail/pkg/config/env"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "chemtrail",
		Short: `
Chemtrail is a scaler for HashiCorp Nomad designed to scale workerpool
client nodes based on demand.
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Version: build.GetVersion(),
	}

	envCfg.RegisterCobra(rootCmd)
	clientCfg.RegisterConfig(rootCmd)

	if err := registerCommands(rootCmd); err != nil {
		fmt.Println("error registering commands:", err)
		os.Exit(sysexits.Software)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(sysexits.Software)
	}
}

func registerCommands(rootCmd *cobra.Command) error {
	if err := scale.RegisterCommand(rootCmd); err != nil {
		return err
	}

	if err := policy.RegisterCommand(rootCmd); err != nil {
		return err
	}

	if err := system.RegisterCommand(rootCmd); err != nil {
		return err
	}

	return server.RegisterCommand(rootCmd)
}
