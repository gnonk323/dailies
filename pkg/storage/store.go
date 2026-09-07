package storage

import (
	"os"
	"path/filepath"

	"dailies/pkg/types"
)

type Store interface {
	LoadConfig() (types.DailiesConfig, error)
	SaveConfig(types.DailiesConfig) error
	LoadEntry(date string) (*types.DailyEntry, error)
	SaveEntry(*types.DailyEntry) error
	ListEntries() ([]types.DailyEntry, error)
	Close() error
}

func DefaultConfig() types.DailiesConfig {
	return types.DailiesConfig{
		WordBank: types.WordBank{
			PromotionThreshold: 3,
			MaxWords:           10,
			Moods:              []string{},
			ContextTags:        []string{},
		},
		Integrations: make(map[string]map[string]bool),
	}
}

func DefaultDBPath() string {
	if envPath := os.Getenv("DAILIES_DB"); envPath != "" {
		return envPath
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "dailies.db"
	}
	return filepath.Join(home, ".dailies", "dailies.db")
}
