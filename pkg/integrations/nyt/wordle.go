package nyt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type WordleFetcher struct{}

func (w WordleFetcher) GameKey() string { return "wordle" }

type WordleBasicInfo struct {
	ID       int    `json:"id"`
	Solution string `json:"solution"`
}

type WordleStateResponse struct {
	States []struct {
		PuzzleID string      `json:"puzzle_id"`
		GameData interface{} `json:"game_data"`
	} `json:"states"`
}

func (w WordleFetcher) FetchGame(ctx context.Context, client *NYTClient, dateStr string, cookies string) (interface{}, error) {
	// resolve Date to Puzzle ID
	basicURL := fmt.Sprintf("https://www.nytimes.com/svc/wordle/v2/%s.json", dateStr)
	resp, err := client.DoRequest(ctx, "GET", basicURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("basic lookup failed with status code %d", resp.StatusCode)
	}

	var basicInfo WordleBasicInfo
	if err := json.NewDecoder(resp.Body).Decode(&basicInfo); err != nil {
		return nil, fmt.Errorf("unmarshaling basic metadata failed: %w", err)
	}

	// query user state via resolved puzzle ID
	stateURL := fmt.Sprintf("https://www.nytimes.com/svc/games/state/wordleV2/latests?puzzle_ids=%d", basicInfo.ID)
	stateResp, err := client.DoRequest(ctx, "GET", stateURL)
	if err != nil {
		return nil, err
	}
	defer stateResp.Body.Close()

	if stateResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("state API failed with status code %d", stateResp.StatusCode)
	}

	var stateData WordleStateResponse
	if err := json.NewDecoder(stateResp.Body).Decode(&stateData); err != nil {
		return nil, fmt.Errorf("unmarshaling state response failed: %w", err)
	}

	if len(stateData.States) == 0 {
		return map[string]interface{}{
			"puzzle_id": basicInfo.ID,
			"solution":  basicInfo.Solution,
		}, nil
	}

	return map[string]interface{}{
		"puzzle_id": basicInfo.ID,
		"solution":  basicInfo.Solution,
		"game_data": stateData.States[0].GameData,
	}, nil
}
