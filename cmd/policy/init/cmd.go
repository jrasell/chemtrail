package init

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates an example scaling policy",
		Run: func(cmd *cobra.Command, args []string) {
			runInit(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runInit(_ *cobra.Command, _ []string) {
	fmt.Print(initPolicy)
}

const initPolicy = `{
  "Enabled": true,
  "MinCount": 2,
  "MaxCount": 4,
  "ScaleOutCount": 1,
  "ScaleInCount": 1,
  "Provider": "aws-autoscaling",
  "ProviderConfig": {
    "asg-name": "chemtrail-test"
  },
  "Checks": {
    "cpu-in": {
      "Enabled": true,
      "Resource": "cpu",
      "ComparisonOperator": "less-than",
      "ComparisonPercentage": 25,
      "Action": "scale-in"
    },
    "cpu-out": {
      "Enabled": true,
      "Resource": "cpu",
      "ComparisonOperator": "greater-than",
      "ComparisonPercentage": 80,
      "Action": "scale-out"
    },
    "memory-in": {
      "Enabled": true,
      "Resource": "memory",
      "ComparisonOperator": "less-than",
      "ComparisonPercentage": 25,
      "Action": "scale-in"
    },
    "memory-out": {
      "Enabled": true,
      "Resource": "memory",
      "ComparisonOperator": "greater-than",
      "ComparisonPercentage": 80,
      "Action": "scale-out"
    }
  }
}
`
