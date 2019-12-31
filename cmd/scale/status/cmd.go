package status

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jrasell/chemtrail/cmd/helper"
	"github.com/jrasell/chemtrail/pkg/api"
	clientCfg "github.com/jrasell/chemtrail/pkg/config/client"
	"github.com/ryanuber/columnize"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

const (
	listOutputHeader   = "ID|Status|Provider|LastUpdate"
	eventsOutputHeader = "Time|Source|Event"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display the status output for scaling activities",
		Run: func(cmd *cobra.Command, args []string) {
			runStatus(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runStatus(_ *cobra.Command, args []string) {
	if len(args) > 1 {
		fmt.Println("Too many arguments, expected maximum 1 args got", len(args))
		os.Exit(sysexits.Usage)
	}

	clientConfig := clientCfg.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	client, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Chemtrail client:", err)
		os.Exit(sysexits.Software)
	}

	switch {
	case len(args) == 1:
		if err := runStatusInfo(client, args[0]); err != nil {
			fmt.Println("Error querying scale status info:", err)
			os.Exit(sysexits.Software)
		}
	case len(args) == 0:
		if err := runStatusList(client); err != nil {
			fmt.Println("Error querying scale status:", err)
			os.Exit(sysexits.Software)
		}
	}
}

func runStatusList(c *api.Client) error {
	resp, err := c.Scale().StatusList()
	if err != nil {
		return err
	}

	out := []string{listOutputHeader}

	for k, v := range *resp {
		out = append(out, fmt.Sprintf("%s|%s|%s|%v",
			k, v.Status, v.Provider.String(), time.Unix(0, v.LastUpdate).UTC()))
	}

	if len(out) > 1 {
		fmt.Println(formatList(out))
	}
	return nil
}

func runStatusInfo(c *api.Client, id string) error {
	resp, err := c.Scale().StatusInfo(id)
	if err != nil {
		return err
	}

	header := []string{
		fmt.Sprintf("ID|%s", id),
		fmt.Sprintf("Status|%s", resp.Status),
		fmt.Sprintf("LastUpdate|%v", time.Unix(0, resp.LastUpdate).UTC()),
		fmt.Sprintf("Direction|%s", resp.Direction),
		fmt.Sprintf("Provider|%s", resp.Provider.String()),
		fmt.Sprintf("ProviderConfig|%s", strings.Join(helper.MapStringsToSliceString(resp.ProviderCfg, ":"), ",")),
	}

	events := []string{eventsOutputHeader}
	for _, v := range resp.Events {
		events = append(events, fmt.Sprintf("%v|%s|%s", time.Unix(0, v.Timestamp).UTC(), v.Source, v.Message))
	}

	fmt.Println(formatKV(header))
	fmt.Println("")
	fmt.Println(formatList(events))

	return nil
}

func formatList(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	return columnize.Format(in, columnConf)
}

func formatKV(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	columnConf.Glue = " = "
	return columnize.Format(in, columnConf)
}
