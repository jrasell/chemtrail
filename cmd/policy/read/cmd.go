package read

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
	checkHeader = "Name|Enabled|Resource|Operator|Value|Action"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Details the scaling policy",
		Run: func(cmd *cobra.Command, args []string) {
			runRead(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runRead(_ *cobra.Command, args []string) {
	switch {
	case len(args) < 1:
		fmt.Println("Not enough arguments, expected 1 args got", len(args))
		os.Exit(sysexits.Usage)
	case len(args) > 1:
		fmt.Println("Too many arguments, expected 1 args got", len(args))
		os.Exit(sysexits.Usage)
	}

	clientConfig := client.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	chemtrailClient, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Chemtrail client:", err)
		os.Exit(sysexits.Software)
	}

	policy, err := chemtrailClient.Policy().Info(args[0])
	if err != nil {
		fmt.Println("Error getting policy info:", err)
		os.Exit(sysexits.Software)
	}

	out := []string{
		fmt.Sprintf("Class|%s", policy.Class),
		fmt.Sprintf("Enabled|%v", policy.Enabled),
		fmt.Sprintf("MaxCount|%v", policy.MaxCount),
		fmt.Sprintf("MinCount|%v", policy.MinCount),
		fmt.Sprintf("ScaleInCount|%v", policy.ScaleInCount),
		fmt.Sprintf("ScaleOutCount|%v", policy.ScaleOutCount),
		fmt.Sprintf("Provider|%v", policy.Provider),
		fmt.Sprintf("ProviderConfig|%v", policy.ProviderConfig),
	}

	var checks []string

	if policy.Checks != nil {
		checks = append(checks, checkHeader)

		for name, check := range policy.Checks {
			checks = append(checks, fmt.Sprintf("%s|%v|%s|%s|%v|%s",
				name, check.Enabled, check.Resource, check.ComparisonOperator, check.ComparisonPercentage, check.Action))
		}
	}

	fmt.Println(helper.FormatKV(out))
	fmt.Println("")
	if len(checks) > 0 {
		fmt.Println(helper.FormatList(checks))
	}
}
