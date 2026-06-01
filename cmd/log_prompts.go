package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"dailies/pkg/storage"
	"dailies/pkg/types"
	"github.com/AlecAivazis/survey/v2"
)

func runInteractiveLog(existing *types.DailyEntry) (*types.DailyEntry, error) {
	activeMoods := GetQuickPicks("moods")
	activeTags := GetQuickPicks("context_tags")

	var defaultQuality string
	var defaultMoods []string
	var defaultTags []string
	var defaultJournal string
	var leftoverMoods []string
	var leftoverTags []string

	if existing != nil {
		defaultQuality = strconv.Itoa(existing.DayQuality)
		defaultJournal = existing.Journal

		// Sort existing moods into multi-select defaults vs custom text input fallbacks
		for _, m := range existing.Moods {
			matched := false
			for _, am := range activeMoods {
				if m == am {
					defaultMoods = append(defaultMoods, m)
					matched = true
					break
				}
			}
			if !matched {
				leftoverMoods = append(leftoverMoods, m)
			}
		}

		// Sort existing tags into multi-select defaults vs custom text input fallbacks
		for _, t := range existing.ContextTags {
			matched := false
			for _, at := range activeTags {
				if t == at {
					defaultTags = append(defaultTags, t)
					matched = true
					break
				}
			}
			if !matched {
				leftoverTags = append(leftoverTags, t)
			}
		}
	}

	var qs []*survey.Question

	qs = append(qs, &survey.Question{
		Name: "DayQuality",
		Prompt: &survey.Input{
			Message: "Rate your day quality (1-10):",
			Default: defaultQuality,
		},
		Validate: func(val interface{}) error {
			str, _ := val.(string)
			num, err := strconv.Atoi(str)
			if err != nil || num < 1 || num > 10 {
				return fmt.Errorf("please enter an integer between 1 and 10")
			}
			return nil
		},
	})

	if len(activeMoods) > 0 {
		qs = append(qs, &survey.Question{
			Name: "ChosenMoods",
			Prompt: &survey.MultiSelect{
				Message: "Select active moods for today:",
				Options: activeMoods,
				Default: defaultMoods,
			},
		})
	}

	qs = append(qs, &survey.Question{
		Name: "CustomMoods",
		Prompt: &survey.Input{
			Message: "Enter any custom or additional moods (comma-separated, or leave blank):",
			Default: strings.Join(leftoverMoods, ", "),
		},
	})

	if len(activeTags) > 0 {
		qs = append(qs, &survey.Question{
			Name: "ChosenTags",
			Prompt: &survey.MultiSelect{
				Message: "Select active context tags:",
				Options: activeTags,
				Default: defaultTags,
			},
		})
	}

	qs = append(qs, &survey.Question{
		Name: "CustomTags",
		Prompt: &survey.Input{
			Message: "Enter context tags for today (comma-separated, or leave blank):",
			Default: strings.Join(leftoverTags, ", "),
		},
	})

	qs = append(qs, &survey.Question{
		Name: "Journal",
		Prompt: &survey.Input{
			Message: "Write your daily reflection journal:",
			Default: defaultJournal,
		},
	})

	type FormAnswers struct {
		DayQuality  string
		ChosenMoods []string
		CustomMoods string
		ChosenTags  []string
		CustomTags  string
		Journal     string
	}

	var answers FormAnswers
	if err := survey.Ask(qs, &answers); err != nil {
		return nil, err
	}

	dq, _ := strconv.Atoi(answers.DayQuality)

	// Combine dynamic and static items cleanly
	moodMap := make(map[string]bool)
	var finalMoods []string
	for _, m := range answers.ChosenMoods {
		if !moodMap[m] {
			moodMap[m] = true
			finalMoods = append(finalMoods, m)
		}
	}
	if answers.CustomMoods != "" {
		for _, m := range strings.Split(answers.CustomMoods, ",") {
			clean := strings.ToLower(strings.TrimSpace(m))
			if clean != "" && !moodMap[clean] {
				moodMap[clean] = true
				finalMoods = append(finalMoods, clean)
			}
		}
	}

	tagMap := make(map[string]bool)
	var finalTags []string
	for _, t := range answers.ChosenTags {
		if !tagMap[t] {
			tagMap[t] = true
			finalTags = append(finalTags, t)
		}
	}
	if answers.CustomTags != "" {
		for _, t := range strings.Split(answers.CustomTags, ",") {
			clean := strings.ToLower(strings.TrimSpace(t))
			if clean != "" && !tagMap[clean] {
				tagMap[clean] = true
				finalTags = append(finalTags, clean)
			}
		}
	}

	return &types.DailyEntry{
		DayQuality:  dq,
		Moods:       finalMoods,
		ContextTags: finalTags,
		Journal:     answers.Journal,
	}, nil
}

func handleInteractivePromotion(field string, todaysWords []string) {
	candidates := GetPromotionCandidates(field, todaysWords)
	if len(candidates) == 0 {
		return
	}

	label := "context tag(s)"
	if field == "moods" {
		label = "mood(s)"
	}

	var chosen []string
	prompt := &survey.MultiSelect{
		Message: fmt.Sprintf("The following %s hit your threshold and were used today. Add to permanent word bank?", label),
		Options: candidates,
	}

	if err := survey.AskOne(prompt, &chosen); err == nil && len(chosen) > 0 {
		config := storage.LoadConfig()
		if field == "moods" {
			config.WordBank.Moods = append(config.WordBank.Moods, chosen...)
		} else {
			config.WordBank.ContextTags = append(config.WordBank.ContextTags, chosen...)
		}
		_ = storage.SaveConfig(config)
		fmt.Printf("Permanently added to config word_bank.%s: [%s]\n", field, strings.Join(chosen, ", "))
	}
}
