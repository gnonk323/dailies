package cmd

import (
	"fmt"
	"strings"
	"time"
	"github.com/spf13/cobra"
	"dailies/pkg/storage"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints today's local entry data summary inside the shell terminal console",
	Run: func(cmd *cobra.Command, args []string) {
		if targetDate == "" {
			targetDate = time.Now().Format("2006-01-02")
		}

		entry, err := storage.LoadEntry(targetDate)
		if err != nil || entry == nil {
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
