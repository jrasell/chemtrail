package health

import (
	"fmt"
	"os"

	"github.com/jrasell/chemtrail/pkg/api"
	"github.com/jrasell/chemtrail/pkg/config/client"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Retrieve health information of a Chemtrail server",
		Run: func(cmd *cobra.Command, args []string) {
			runDelete(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runDelete(_ *cobra.Command, _ []string) {
	clientConfig := client.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	chemtrailClient, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Chemtrail client:", err)
		os.Exit(sysexits.Software)
	}

	health, err := chemtrailClient.System().Health()
	if err != nil {
		fmt.Println("Error calling server health:", err)
		os.Exit(sysexits.Software)
	}

	fmt.Println("Chemtrail server status:", health.Status)
}
