package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	// "dailies/pkg/client"
	"dailies/pkg/types"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View and interactively edit dailies configuration",
	Run: func(cmd *cobra.Command, args []string) {
		runConfigEditor()
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration with secrets redacted",
	Run: func(cmd *cobra.Command, args []string) {
		runConfigShow()
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	RootCmd.AddCommand(configCmd)
}

func runConfigEditor() {
	api := mustAPIClient()

	cfg, err := api.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	for {
		var section string

		prompt := &survey.Select{
			Message: "What do you want to configure?",
			Options: []string{
				"Word bank",
				"Integrations",
				"GitHub",
				"New York Times",
				"Show current config",
				"Save and exit",
				"Exit without saving",
			},
		}

		if err := survey.AskOne(prompt, &section); err != nil {
			fmt.Fprintf(os.Stderr, "Configuration cancelled: %v\n", err)
			return
		}

		switch section {
		case "Word bank":
			editWordBank(&cfg)

		case "Integrations":
			editIntegrations(&cfg)

		case "GitHub":
			editGitHub(&cfg)

		case "New York Times":
			editNYT(&cfg)

		case "Show current config":
			printConfig(cfg)

		case "Save and exit":
			if err := api.SaveConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save config: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Configuration saved.")
			return

		case "Exit without saving":
			fmt.Println("Configuration changes discarded.")
			return
		}
	}
}

func editWordBank(cfg *types.DailiesConfig) {
	fmt.Println()
	fmt.Println("Word bank configuration")
	fmt.Println("-----------------------")

	promotionThreshold := strconv.Itoa(cfg.WordBank.PromotionThreshold)
	maxWords := strconv.Itoa(cfg.WordBank.MaxWords)

	questions := []*survey.Question{
		{
			Name: "PromotionThreshold",
			Prompt: &survey.Input{
				Message: "Promotion threshold:",
				Default: promotionThreshold,
			},
			Validate: func(val interface{}) error {
				value := strings.TrimSpace(val.(string))
				n, err := strconv.Atoi(value)
				if err != nil || n < 1 {
					return fmt.Errorf("enter a positive integer")
				}
				return nil
			},
		},
		{
			Name: "MaxWords",
			Prompt: &survey.Input{
				Message: "Maximum active words:",
				Default: maxWords,
			},
			Validate: func(val interface{}) error {
				value := strings.TrimSpace(val.(string))
				n, err := strconv.Atoi(value)
				if err != nil || n < 1 {
					return fmt.Errorf("enter a positive integer")
				}
				return nil
			},
		},
		{
			Name: "Moods",
			Prompt: &survey.Input{
				Message: "Moods (comma-separated):",
				Default: strings.Join(cfg.WordBank.Moods, ", "),
			},
		},
		{
			Name: "ContextTags",
			Prompt: &survey.Input{
				Message: "Context tags (comma-separated):",
				Default: strings.Join(cfg.WordBank.ContextTags, ", "),
			},
		},
	}

	answers := struct {
		PromotionThreshold string
		MaxWords           string
		Moods              string
		ContextTags        string
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		fmt.Fprintf(os.Stderr, "Word bank edit cancelled: %v\n", err)
		return
	}

	cfg.WordBank.PromotionThreshold, _ = strconv.Atoi(answers.PromotionThreshold)
	cfg.WordBank.MaxWords, _ = strconv.Atoi(answers.MaxWords)
	cfg.WordBank.Moods = parseList(answers.Moods)
	cfg.WordBank.ContextTags = parseList(answers.ContextTags)

	fmt.Println("Word bank updated.")
}

func editIntegrations(cfg *types.DailiesConfig) {
	if cfg.Integrations == nil {
		cfg.Integrations = make(map[string]map[string]bool)
	}

	// These are the integrations currently registered by the application.
	integrations := []string{
		"github",
		"nyt",
	}

	for _, name := range integrations {
		settings := cfg.Integrations[name]

		enabled := false
		if settings != nil {
			enabled = settings["enabled"]
		}

		var answer bool

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Enable %s integration?", name),
			Default: answerOrDefault(enabled),
		}

		if err := survey.AskOne(prompt, &answer); err != nil {
			fmt.Fprintf(os.Stderr, "Integration configuration cancelled: %v\n", err)
			return
		}

		if settings == nil {
			settings = make(map[string]bool)
		}

		settings["enabled"] = answer
		cfg.Integrations[name] = settings
	}

	fmt.Println("Integrations updated.")
}

func editGitHub(cfg *types.DailiesConfig) {
	fmt.Println()
	fmt.Println("GitHub configuration")
	fmt.Println("--------------------")

	fmt.Printf("Current username: %s\n", emptyAsUnset(cfg.GitHub.Username))
	fmt.Printf("Token: %s\n", secretStatus(cfg.GitHub.Token))

	var username string
	var changeToken bool

	questions := []*survey.Question{
		{
			Name: "Username",
			Prompt: &survey.Input{
				Message: "GitHub username:",
				Default: cfg.GitHub.Username,
			},
		},
		{
			Name: "ChangeToken",
			Prompt: &survey.Confirm{
				Message: "Change GitHub token?",
				Default: false,
			},
		},
	}

	answers := struct {
		Username    string
		ChangeToken bool
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		fmt.Fprintf(os.Stderr, "GitHub configuration cancelled: %v\n", err)
		return
	}

	username = strings.TrimSpace(answers.Username)
	cfg.GitHub.Username = username
	changeToken = answers.ChangeToken

	if changeToken {
		var token string

		prompt := &survey.Password{
			Message: "GitHub token:",
		}

		if err := survey.AskOne(prompt, &token); err != nil {
			fmt.Fprintf(os.Stderr, "GitHub token update cancelled: %v\n", err)
			return
		}

		cfg.GitHub.Token = strings.TrimSpace(token)
	}

	fmt.Println("GitHub configuration updated.")
}

func editNYT(cfg *types.DailiesConfig) {
	fmt.Println()
	fmt.Println("New York Times configuration")
	fmt.Println("----------------------------")
	fmt.Printf("Cookies: %s\n", secretStatus(cfg.NYT.Cookies))

	var changeCookies bool

	prompt := &survey.Confirm{
		Message: "Change NYT cookies?",
		Default: false,
	}

	if err := survey.AskOne(prompt, &changeCookies); err != nil {
		fmt.Fprintf(os.Stderr, "NYT configuration cancelled: %v\n", err)
		return
	}

	if changeCookies {
		var cookies string

		prompt := &survey.Password{
			Message: "NYT cookies:",
		}

		if err := survey.AskOne(prompt, &cookies); err != nil {
			fmt.Fprintf(os.Stderr, "NYT cookie update cancelled: %v\n", err)
			return
		}

		cfg.NYT.Cookies = strings.TrimSpace(cookies)
	}

	fmt.Println("NYT configuration updated.")
}

func runConfigShow() {
	api := mustAPIClient()

	cfg, err := api.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	printConfig(cfg)
}

func printConfig(cfg types.DailiesConfig) {
	fmt.Println()
	fmt.Println("Dailies configuration")
	fmt.Println("=====================")

	fmt.Println()
	fmt.Println("Word bank")
	fmt.Printf("  Promotion threshold: %d\n", cfg.WordBank.PromotionThreshold)
	fmt.Printf("  Max words:           %d\n", cfg.WordBank.MaxWords)
	fmt.Printf("  Moods:               [%s]\n", strings.Join(cfg.WordBank.Moods, ", "))
	fmt.Printf("  Context tags:        [%s]\n", strings.Join(cfg.WordBank.ContextTags, ", "))

	fmt.Println()
	fmt.Println("Integrations")

	integrationNames := []string{"github", "nyt"}

	for _, name := range integrationNames {
		enabled := false

		if settings, ok := cfg.Integrations[name]; ok {
			enabled = settings["enabled"]
		}

		fmt.Printf("  %-8s enabled=%t\n", name, enabled)
	}

	fmt.Println()
	fmt.Println("GitHub")
	fmt.Printf("  Username: %s\n", emptyAsUnset(cfg.GitHub.Username))
	fmt.Printf("  Token:    %s\n", secretStatus(cfg.GitHub.Token))

	fmt.Println()
	fmt.Println("New York Times")
	fmt.Printf("  Cookies:  %s\n", secretStatus(cfg.NYT.Cookies))

	fmt.Println()
}

func parseList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}

	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, item := range strings.Split(value, ",") {
		item = strings.ToLower(strings.TrimSpace(item))

		if item == "" || seen[item] {
			continue
		}

		seen[item] = true
		result = append(result, item)
	}

	return result
}

func secretStatus(value string) string {
	if strings.TrimSpace(value) == "" {
		return "<not configured>"
	}

	return "<configured>"
}

func emptyAsUnset(value string) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return "<not configured>"
	}

	return value
}

func answerOrDefault(value bool) bool {
	return value
}
