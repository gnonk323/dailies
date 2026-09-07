package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"dailies/pkg/types"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Interactive prompt to create or append to today's data entry",
	Long: `Create or update daily entries. If called without subcommands, it defaults to interactive mode.
Available subcommands: add, remove, set.`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetDate()
		api := mustAPIClient()

		fmt.Printf("Logging for %s\n", targetDate)
		existingEntry, err := api.GetEntry(targetDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load entry: %v\n", err)
			os.Exit(1)
		}

		if existingEntry != nil {
			fmt.Printf("Found existing data log for %s. Pre-filling options...\n", targetDate)
		}

		cfg, err := api.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}
		entries, err := api.ListEntries()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load entries: %v\n", err)
			os.Exit(1)
		}

		coreAnswers, err := runInteractiveLog(cfg, entries, existingEntry)
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

		if err := api.SaveEntry(coreAnswers); err != nil {
			fmt.Println("Error saving entry:", err)
			return
		}
		fmt.Printf("Log saved for %s\n", targetDate)

		handleInteractivePromotion(api, cfg, entries, "moods", coreAnswers.Moods)
		handleInteractivePromotion(api, cfg, entries, "context_tags", coreAnswers.ContextTags)
	},
}

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
	logCmd.AddCommand(logAddCmd)
	logCmd.AddCommand(logRemoveCmd)
	logCmd.AddCommand(logSetCmd)
	RootCmd.AddCommand(logCmd)
}

func ensureTargetDate() {
	if targetDate == "" {
		targetDate = time.Now().Format("2006-01-02")
	}
}

func getOrCreateEntry() (*types.DailyEntry, error) {
	ensureTargetDate()
	existingEntry, err := mustAPIClient().GetEntry(targetDate)
	if err != nil {
		return nil, err
	}
	if existingEntry != nil {
		return existingEntry, nil
	}

	return &types.DailyEntry{
		Date:         targetDate,
		Integrations: make(map[string]interface{}),
	}, nil
}

func saveHeadless(entry *types.DailyEntry) error {
	if err := mustAPIClient().SaveEntry(entry); err != nil {
		return fmt.Errorf("error saving entry: %w", err)
	}
	fmt.Printf("Log updated for %s\n", targetDate)
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