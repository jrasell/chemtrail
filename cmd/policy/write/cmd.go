package write

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jrasell/chemtrail/pkg/api"
	"github.com/jrasell/chemtrail/pkg/config/client"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "Uploads a policy from file",
		Run: func(cmd *cobra.Command, args []string) {
			runWrite(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runWrite(_ *cobra.Command, args []string) {
	switch {
	case len(args) < 2:
		fmt.Println("Not enough arguments, expected 2 args got", len(args))
		os.Exit(sysexits.Usage)
	case len(args) > 2:
		fmt.Println("Too many arguments, expected 2 args got", len(args))
		os.Exit(sysexits.Usage)
	}

	path := strings.TrimSpace(args[1])

	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading scaling policy file:", err)
		os.Exit(sysexits.Software)
	}

	clientConfig := client.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	chemtrailClient, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Chemtrail client:", err)
		os.Exit(sysexits.Software)
	}

	var policy *api.ScalingPolicy
	if err = json.Unmarshal(b, &policy); err != nil {
		fmt.Println("Error parsing scaling policy file:", err)
		os.Exit(sysexits.Software)
	}

	if err := chemtrailClient.Policy().Write(args[0], policy); err != nil {
		fmt.Println("Error writing class scaling policy:", err)
		os.Exit(sysexits.Software)
	}

	fmt.Println("Successfully wrote client scaling policy")
	os.Exit(sysexits.OK)
}
