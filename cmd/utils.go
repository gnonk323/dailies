package cmd

import (
	"sort"
	"strings"

	"dailies/pkg/types"
)

func GetQuickPicks(config types.DailiesConfig, entries []types.DailyEntry, field string) []string {
	maxWords := config.WordBank.MaxWords

	var historicalBank []string
	if field == "moods" {
		historicalBank = config.WordBank.Moods
	} else {
		historicalBank = config.WordBank.ContextTags
	}

	if len(historicalBank) <= maxWords {
		return historicalBank
	}

	frequencyMap := make(map[string]int)
	recencyMap := make(map[string]string)

	for _, word := range historicalBank {
		frequencyMap[word] = 0
		recencyMap[word] = "0000-00-00"
	}

	for _, entry := range entries {
		var words []string
		if field == "moods" {
			words = entry.Moods
		} else {
			words = entry.ContextTags
		}

		for _, w := range words {
			clean := strings.ToLower(strings.TrimSpace(w))
			if _, exists := frequencyMap[clean]; exists {
				frequencyMap[clean]++
				if entry.Date > recencyMap[clean] {
					recencyMap[clean] = entry.Date
				}
			}
		}
	}

	keys := make([]string, 0, len(frequencyMap))
	for k := range frequencyMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if frequencyMap[keys[j]] != frequencyMap[keys[i]] {
			return frequencyMap[keys[i]] > frequencyMap[keys[j]]
		}
		return recencyMap[keys[i]] > recencyMap[keys[j]]
	})

	if len(keys) > maxWords {
		return keys[:maxWords]
	}
	return keys
}

func GetPromotionCandidates(config types.DailiesConfig, entries []types.DailyEntry, field string, todaysWords []string) []string {
	threshold := config.WordBank.PromotionThreshold

	var historicalBank []string
	if field == "moods" {
		historicalBank = config.WordBank.Moods
	} else {
		historicalBank = config.WordBank.ContextTags
	}

	if len(todaysWords) == 0 {
		return []string{}
	}

	historicalMap := make(map[string]bool)
	for _, w := range historicalBank {
		historicalMap[w] = true
	}

	frequencyMap := make(map[string]int)
	for _, w := range todaysWords {
		clean := strings.ToLower(strings.TrimSpace(w))
		if clean != "" && !historicalMap[clean] {
			frequencyMap[clean] = 0
		}
	}

	if len(frequencyMap) == 0 {
		return []string{}
	}

	for _, entry := range entries {
		var words []string
		if field == "moods" {
			words = entry.Moods
		} else {
			words = entry.ContextTags
		}
		for _, w := range words {
			clean := strings.ToLower(strings.TrimSpace(w))
			if _, exists := frequencyMap[clean]; exists {
				frequencyMap[clean]++
			}
		}
	}

	var candidates []string
	for word, freq := range frequencyMap {
		if freq >= threshold {
			candidates = append(candidates, word)
		}
	}
	return candidates
}
