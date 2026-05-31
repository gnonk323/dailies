package cmd

import (
	"time"
	"github.com/spf13/cobra"
	"dailies/pkg/integrations"
)

var integrationFlag string

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Pull data from external integration points",
	Run: func(cmd *cobra.Command, args []string) {
		if targetDate == "" {
			targetDate = time.Now().Format("2006-01-02")
		}

		if integrationFlag != "" {
			integrations.RunManualFetch(integrationFlag, targetDate)
		} else {
			integrations.RunAllManualFetch(targetDate)
		}
	},
}

func init() {
	fetchCmd.Flags().StringVarP(&integrationFlag, "integration", "i", "", "Fetch specific named integration exclusively")
	RootCmd.AddCommand(fetchCmd)
}
