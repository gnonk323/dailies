package nyt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ConnectionsFetcher struct{}

func (c ConnectionsFetcher) GameKey() string { return "connections" }

type ConnectionsBasicInfo struct {
	ID int `json:"id"`
}

type ConnectionsStateResponse struct {
	States []struct {
		PuzzleID string      `json:"puzzle_id"`
		GameData interface{} `json:"game_data"`
	} `json:"states"`
}

func (c ConnectionsFetcher) FetchGame(ctx context.Context, client *NYTClient, dateStr string, cookies string) (interface{}, error) {
	// lookup internal puzzle identifier by date string
	lookupURL := fmt.Sprintf("https://www.nytimes.com/svc/connections/v2/%s.json", dateStr)
	resp, err := client.DoRequest(ctx, "GET", lookupURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("connections base lookup failed with code %d", resp.StatusCode)
	}

	var basic ConnectionsBasicInfo
	if err := json.NewDecoder(resp.Body).Decode(&basic); err != nil {
		return nil, fmt.Errorf("failed decoding connections basic info: %w", err)
	}

	// grab the user's specific state for this puzzle
	stateURL := "https://www.nytimes.com/svc/games/state/connections/latests"
	stateResp, err := client.DoRequest(ctx, "GET", stateURL)
	if err != nil {
		return nil, err
	}
	defer stateResp.Body.Close()

	if stateResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("connections state endpoint failed with code %d", stateResp.StatusCode)
	}

	var statePayload ConnectionsStateResponse
	if err := json.NewDecoder(stateResp.Body).Decode(&statePayload); err != nil {
		return nil, fmt.Errorf("failed decoding connections state body: %w", err)
	}

	if len(statePayload.States) == 0 {
		return map[string]interface{}{
			"puzzle_id": basic.ID,
		}, nil
	}

	return map[string]interface{}{
		"puzzle_id": basic.ID,
		"game_data": statePayload.States[0].GameData,
	}, nil
}
