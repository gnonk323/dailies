package integrations

import (
	"dailies/pkg/integrations/github"
	"dailies/pkg/integrations/nyt"
	"dailies/pkg/storage"
	"dailies/pkg/types"
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex

var IntegrationRegistry = map[string]Integration{
	"github": github.GitHubModule{},
	"nyt":    nyt.NYTModule{},
}

func RunManualFetch(store storage.Store, integrationName string, dateStr string) {
	config, err := store.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	settings, exists := config.Integrations[integrationName]
	if !exists || !settings["enabled"] {
		fmt.Printf("Integration '%s' is disabled or not declared in configuration.\n", integrationName)
		return
	}

	module, registered := IntegrationRegistry[integrationName]
	if !registered || module.GetType() != "manual" {
		fmt.Printf("Integration '%s' not found or is not configured as a manual type.\n", integrationName)
		return
	}

	manualModule := module.(ManualIntegration)
	fmt.Printf("Fetching data from [%s] for %s...\n", integrationName, dateStr)

	fetchedData, err := manualModule.Fetch(dateStr, config)
	if err != nil {
		fmt.Printf("Error pulling metrics via [%s]: %v\n", integrationName, err)
		return
	}

	payload := IntegrationPayload{
		Data:      fetchedData,
		FetchedAt: time.Now().Format(time.RFC3339),
	}

	mu.Lock()
	defer mu.Unlock()

	entry, err := store.LoadEntry(dateStr)
	if err != nil || entry == nil {
		entry = &types.DailyEntry{
			Date:         dateStr,
			DayQuality:   5,
			Moods:        []string{},
			ContextTags:  []string{},
			Journal:      "",
			Integrations: make(map[string]interface{}),
		}
	}

	if entry.Integrations == nil {
		entry.Integrations = make(map[string]interface{})
	}

	entry.Integrations[integrationName] = payload

	if err := store.SaveEntry(entry); err != nil {
		fmt.Printf("Failed saving updated entry data logs for target key: %v\n", err)
		return
	}
	fmt.Printf("Cleanly merged [%s] metadata into entry %s\n", integrationName, dateStr)
}

func RunAllManualFetch(store storage.Store, dateStr string) {
	config, err := store.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}
	var enabledModules []string

	for name, module := range IntegrationRegistry {
		if module.GetType() == "manual" {
			if settings, ok := config.Integrations[name]; ok && settings["enabled"] {
				enabledModules = append(enabledModules, name)
			}
		}
	}

	if len(enabledModules) == 0 {
		fmt.Println("No manual integrations are currently enabled.")
		return
	}

	fmt.Printf("Executing concurrent fetch for %d manual integrations...\n", len(enabledModules))

	var wg sync.WaitGroup
	for _, name := range enabledModules {
		wg.Add(1)
		go func(integrationName string) {
			defer wg.Done()
			RunManualFetch(store, integrationName, dateStr)
		}(name)
	}

	wg.Wait()
	fmt.Println("Finished all manual integration fetches.")
}
