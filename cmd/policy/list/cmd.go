package list

import (
	"fmt"
	"os"

	"github.com/jrasell/chemtrail/cmd/helper"
	"github.com/jrasell/chemtrail/pkg/api"
	"github.com/jrasell/chemtrail/pkg/config/client"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

const (
	outputHeader = "Class|Enabled|MinCount|MaxCount|Provider"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all scaling policies",
		Run: func(cmd *cobra.Command, args []string) {
			runList(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runList(_ *cobra.Command, _ []string) {
	clientConfig := client.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	chemtrailClient, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Chemtrail client:", err)
		os.Exit(sysexits.Software)
	}

	resp, err := chemtrailClient.Policy().List()
	if err != nil {
		fmt.Println("Error listing scaling policies:", err)
		os.Exit(sysexits.Software)
	}

	if len(*resp) == 0 {
		os.Exit(sysexits.OK)
	}
	out := []string{outputHeader}

	for class, pol := range *resp {
		out = append(out, fmt.Sprintf("%s|%v|%v|%v|%s",
			class, pol.Enabled, pol.MinCount, pol.MaxCount, pol.Provider))
	}
	fmt.Println(helper.FormatList(out))
}
