package cmd

import (
	"fmt"
	"os"
	"time"
	"github.com/spf13/cobra"
	"dailies/pkg/storage"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Interactive prompt to create or append to today's data entry",
	Run: func(cmd *cobra.Command, args []string) {
		if targetDate == "" {
			targetDate = time.Now().Format("2006-01-02")
		}

		fmt.Printf("Logging for %s\n", targetDate)
		existingEntry, _ := storage.LoadEntry(targetDate)

		if existingEntry != nil {
			fmt.Printf("Found existing data log for %s. Pre-filling options...\n", targetDate)
		}

		coreAnswers, err := runInteractiveLog(existingEntry)
		if err != nil {
			if err.Error() == "interrupt" {
				fmt.Println("\nOperation cancelled. Exiting...")
				os.Exit(0)
			}
			fmt.Println("Operation aborted or unexpected prompt failure:", err)
			return
		}

		coreAnswers.Date = targetDate
		coreAnswers.Integrations = make(map[string]interface{})
		if existingEntry != nil {
			coreAnswers.Integrations = existingEntry.Integrations
		}

		if err := storage.SaveEntry(coreAnswers); err != nil {
			fmt.Println("Error writing entry to disk:", err)
			return
		}
		fmt.Printf("\nLog written to disk at data/%s.json\n", targetDate)

		handleInteractivePromotion("moods", coreAnswers.Moods)
		handleInteractivePromotion("context_tags", coreAnswers.ContextTags)
	},
}

func init() {
	RootCmd.AddCommand(logCmd)
}
