package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"dailies/pkg/client"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints today's entry data summary",
	Run: func(cmd *cobra.Command, args []string) {
		if targetDate == "" {
			targetDate = time.Now().Format("2006-01-02")
		}

		api := mustAPIClient()
		entry, err := api.GetEntry(targetDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load entry: %v\n", err)
			os.Exit(1)
		}
		if entry == nil {
			fmt.Printf("No entries tracked yet for %s. Fire up 'dailies log' to begin!\n", targetDate)
			return
		}

		fmt.Printf("Entry for %s\n", entry.Date)
		fmt.Println("======================")
		fmt.Printf("Rating:  %d/10\n", entry.DayQuality)
		fmt.Printf("Moods:   [%s]\n", strings.Join(entry.Moods, ", "))
		fmt.Printf("Tags:    [%s]\n", strings.Join(entry.ContextTags, ", "))
		fmt.Println("----------------------")
		fmt.Printf("Journal:\n%s\n", entry.Journal)
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}

func mustAPIClient() *client.APIClient {
	api, err := client.NewFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	return api
}
