package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var integrationFlag string

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Pull data from external integration points via the dailies server",
	Run: func(cmd *cobra.Command, args []string) {
		if targetDate == "" {
			targetDate = time.Now().Format("2006-01-02")
		}

		api := mustAPIClient()
		if integrationFlag != "" {
			if _, err := api.FetchIntegration(targetDate, integrationFlag); err != nil {
				fmt.Fprintf(os.Stderr, "Fetch failed: %v\n", err)
				os.Exit(1)
			}
			return
		}
		if err := api.FetchAllIntegrations(targetDate); err != nil {
			fmt.Fprintf(os.Stderr, "Fetch failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	fetchCmd.Flags().StringVarP(&integrationFlag, "integration", "i", "", "Fetch specific named integration exclusively")
	RootCmd.AddCommand(fetchCmd)
}
