package nyt

import (
	"context"
	"dailies/pkg/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NYTModule struct{}

func (n NYTModule) GetName()        string { return "nyt" }
func (n NYTModule) GetDescription() string { return "Aggregates structural logs and life stats from NYT Games" }
func (n NYTModule) GetType()        string { return "manual" }

// Struct for parsing the general /latests stats schema
type PlayerStatsPayload struct {
	Player struct {
		Stats struct {
			Wordle        interface{} `json:"wordle"`
			Connections   interface{} `json:"connections"`
			CrosswordMini interface{} `json:"crossword_mini"`
			CrosswordMidi interface{} `json:"crossword_midi"`
		} `json:"stats"`
	} `json:"player"`
}

func (n NYTModule) Fetch(dateStr string, config types.DailiesConfig) (map[string]interface{}, error) {
	cookies := config.NYT.Cookies 
	if cookies == "" {
		return nil, fmt.Errorf("missing 'nyt.cookies' block within active configurations")
	}

	client := NewNYTClient(cookies)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	outputPayload := make(map[string]interface{})

	// add user identity key
	if userId, err := client.ExtractRegiID(); err == nil {
		outputPayload["user_id"] = userId
	}

	gamesToFetch := []GameFetcher{
		WordleFetcher{},
		ConnectionsFetcher{},
		MiniCrosswordFetcher{},
		MidiCrosswordFetcher{},
	}

	for _, game := range gamesToFetch {
		gameData, err := game.FetchGame(ctx, client, dateStr, cookies)
		if err != nil {
			fmt.Printf(" [NYT Module Warning] Failed retrieving %s data logs: %v\n", game.GameKey(), err)
			continue
		}
		outputPayload[game.GameKey()] = gameData
	}

	summaryURL := "https://www.nytimes.com/svc/games/state/wordleV2/latests"
	summaryResp, err := client.DoRequest(ctx, "GET", summaryURL)
	if err == nil {
		defer summaryResp.Body.Close()
		if summaryResp.StatusCode == http.StatusOK {
			var rawStats PlayerStatsPayload
			if err := json.NewDecoder(summaryResp.Body).Decode(&rawStats); err == nil {
				outputPayload["summary"] = map[string]interface{}{
					"wordle":         rawStats.Player.Stats.Wordle,
					"connections":    rawStats.Player.Stats.Connections,
					"crossword_mini": rawStats.Player.Stats.CrosswordMini,
					"crossword_midi": rawStats.Player.Stats.CrosswordMidi,
				}
			}
		}
	} else {
		fmt.Printf("[NYT Module Warning] Failed retrieving life user summary section: %v\n", err)
	}

	return outputPayload, nil
}
