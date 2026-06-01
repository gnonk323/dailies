package cmd

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"
	"dailies/pkg/storage"
	"dailies/pkg/sync"
	"dailies/pkg/types"
	"github.com/spf13/cobra"
)

// Base log command (defaults to interactive mode if no subcommand is called)
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Interactive prompt to create or append to today's data entry",
	Long: `Create or update daily entries. If called without subcommands, it defaults to interactive mode.
Available subcommands: add, remove, set.`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetDate()

		if err := sync.SyncPull(); err != nil {
			fmt.Printf("Sync Warning: %v\nProceeding with local files anyway...\n", err)
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
		fmt.Printf("Log written to disk at data/%s.json\n", targetDate)

		if err := sync.SyncPush(); err != nil {
			fmt.Printf("Sync Error: %v\nYour data is saved locally but not pushed.\n", err)
		}

		handleInteractivePromotion("moods", coreAnswers.Moods)
		handleInteractivePromotion("context_tags", coreAnswers.ContextTags)
	},
}

// 2. Subcommand: log add <field> <value>
var logAddCmd = &cobra.Command{
	Use:   "add [field] [value]",
	Short: "Add an item to an array field (mood or context)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		field := strings.ToLower(args[0])
		value := strings.ToLower(strings.TrimSpace(args[1]))

		entry, err := getOrCreateEntry()
		if err != nil {
			return err
		}

		switch field {
		case "mood":
			if !sliceContains(entry.Moods, value) {
				entry.Moods = append(entry.Moods, value)
			}
		case "context", "tag":
			if !sliceContains(entry.ContextTags, value) {
				entry.ContextTags = append(entry.ContextTags, value)
			}
		default:
			return fmt.Errorf("invalid field '%s': only 'mood' or 'context' can be added to", field)
		}

		return saveHeadless(entry)
	},
}

// 3. Subcommand: log remove <field> <value>
var logRemoveCmd = &cobra.Command{
	Use:   "remove [field] [value]",
	Short: "Remove an item from an array field (mood or context)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		field := strings.ToLower(args[0])
		value := strings.ToLower(strings.TrimSpace(args[1]))

		entry, err := getOrCreateEntry()
		if err != nil {
			return err
		}

		switch field {
		case "mood":
			entry.Moods = removeFromSlice(entry.Moods, value)
		case "context", "tag":
			entry.ContextTags = removeFromSlice(entry.ContextTags, value)
		default:
			return fmt.Errorf("invalid field '%s': only 'mood' or 'context' can be removed from", field)
		}

		return saveHeadless(entry)
	},
}

// 4. Subcommand: log set <field> <value>
var logSetCmd = &cobra.Command{
	Use:   "set [field] [value]",
	Short: "Set scalar text or rating fields (rating or journal)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		field := strings.ToLower(args[0])
		value := strings.TrimSpace(args[1])

		entry, err := getOrCreateEntry()
		if err != nil {
			return err
		}

		switch field {
		case "rating", "quality":
			num, err := strconv.Atoi(value)
			if err != nil || num < 1 || num > 10 {
				return fmt.Errorf("rating value must be an integer between 1 and 10")
			}
			entry.DayQuality = num
		case "journal", "reflection":
			entry.Journal = value
		default:
			return fmt.Errorf("invalid field '%s': only 'rating' or 'journal' can be set", field)
		}

		return saveHeadless(entry)
	},
}

func init() {
	// Nest your subcommands under the base logCmd
	logCmd.AddCommand(logAddCmd)
	logCmd.AddCommand(logRemoveCmd)
	logCmd.AddCommand(logSetCmd)

	// Add logCmd to RootCmd
	RootCmd.AddCommand(logCmd)
}


func ensureTargetDate() {
	if targetDate == "" {
		targetDate = time.Now().Format("2006-01-02")
	}
}

func getOrCreateEntry() (*types.DailyEntry, error) {
	ensureTargetDate()

	if err := sync.SyncPull(); err != nil {
		fmt.Printf("Sync Warning: %v\nProceeding with local files anyway...\n", err)
	}

	existingEntry, _ := storage.LoadEntry(targetDate)
	if existingEntry != nil {
		return existingEntry, nil
	}

	return &types.DailyEntry{
		Date:         targetDate,
		Integrations: make(map[string]interface{}),
	}, nil
}

func saveHeadless(entry *types.DailyEntry) error {
	if err := storage.SaveEntry(entry); err != nil {
		return fmt.Errorf("error writing entry to disk: %w", err)
	}
	fmt.Printf("Log updated and written to disk at data/%s.json\n", targetDate)

	if err := sync.SyncPush(); err != nil {
		fmt.Printf("Sync Error: %v\nYour data is saved locally but not pushed.\n", err)
	}
	
	return nil
}

func sliceContains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func removeFromSlice(slice []string, val string) []string {
	var result []string
	for _, item := range slice {
		if item != val {
			result = append(result, item)
		}
	}
	return result
}
