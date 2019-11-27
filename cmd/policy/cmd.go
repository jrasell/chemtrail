package policy

import (
	"fmt"
	"os"

	"github.com/jrasell/chemtrail/cmd/policy/delete"
	initcmd "github.com/jrasell/chemtrail/cmd/policy/init"
	"github.com/jrasell/chemtrail/cmd/policy/list"
	"github.com/jrasell/chemtrail/cmd/policy/read"
	"github.com/jrasell/chemtrail/cmd/policy/write"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Interact with scaling policies",
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
	if err := delete.RegisterCommand(cmd); err != nil {
		return err
	}

	if err := list.RegisterCommand(cmd); err != nil {
		return err
	}

	if err := read.RegisterCommand(cmd); err != nil {
		return err
	}

	if err := initcmd.RegisterCommand(cmd); err != nil {
		return err
	}

	return write.RegisterCommand(cmd)
}
