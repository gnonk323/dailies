package storage

import (
	"encoding/json"
	"errors"
	"os"
	"io/fs"
	"path/filepath"
	"strings"
	"dailies/pkg/types"
)

func GetConfigPath() string {
	envPath := os.Getenv("DAILIES_CONFIG")
	if envPath != "" {
		if strings.HasPrefix(envPath, "~") {
			home, _ := os.UserHomeDir()
			return filepath.Join(home, envPath[1:])
		}
		return filepath.Clean(envPath)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".dailies", "config.json")
}

func LoadConfig() types.DailiesConfig {
	// default fallbacks
	defaultCfg := types.DailiesConfig{
		WordBank: types.WordBank{
			PromotionThreshold: 3,
			MaxWords:           10,
			Moods:              []string{},
			ContextTags:        []string{},
		},
		Integrations: make(map[string]map[string]bool),
	}

	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); errors.Is(err, fs.ErrNotExist) {
		return defaultCfg
	}

	raw, err := os.ReadFile(configPath)
	if err != nil {
		return defaultCfg
	}

	var cfg types.DailiesConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return defaultCfg
	}
	return cfg
}

func SaveConfig(cfg types.DailiesConfig) error {
	configPath := GetConfigPath()
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}
	
	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, bytes, 0644)
}

func GetDataDirectory() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "data")
}

func GetEntryPath(dateStr string) string {
	return filepath.Join(GetDataDirectory(), dateStr+".json")
}

func LoadEntry(dateStr string) (*types.DailyEntry, error) {
	filePath := GetEntryPath(dateStr)
	if _, err := os.Stat(filePath); errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}

	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var entry types.DailyEntry
	if err := json.Unmarshal(raw, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func SaveEntry(entry *types.DailyEntry) error {
	dir := GetDataDirectory()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := GetEntryPath(entry.Date)
	bytes, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, bytes, 0644)
}
