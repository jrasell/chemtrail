package scale

import (
	"fmt"
	"os"

	"github.com/jrasell/chemtrail/cmd/scale/in"
	"github.com/jrasell/chemtrail/cmd/scale/out"
	"github.com/jrasell/chemtrail/cmd/scale/status"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "scale",
		Short: "Perform scaling actions against a Nomad client class and check status",
		Run: func(cmd *cobra.Command, args []string) {
			runScale(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	if err := registerCommands(cmd); err != nil {
		fmt.Println("Error registering commands:", err)
		os.Exit(sysexits.Software)
	}
	return nil
}

func runScale(cmd *cobra.Command, _ []string) {
	_ = cmd.Usage()
}

func registerCommands(cmd *cobra.Command) error {
	if err := in.RegisterCommand(cmd); err != nil {
		return err
	}
	if err := out.RegisterCommand(cmd); err != nil {
		return err
	}
	return status.RegisterCommand(cmd)
}
