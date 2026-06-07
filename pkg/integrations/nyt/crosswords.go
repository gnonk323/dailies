package nyt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CrosswordMetadata struct {
	ID              int      `json:"id"`
	PublicationDate string   `json:"publicationDate"`
	Title           string   `json:"title,omitempty"`
	Editor          string   `json:"editor,omitempty"`
	Constructors    []string `json:"constructors"`
	Dimensions      struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"dimensions"`
	CompletionFraction int `json:"completionFraction"`
	FirstSolve         int `json:"firstSolve"`
	PlayTime           int `json:"playTime"`
}

type RawCrosswordPayload struct {
	ID              int      `json:"id"`
	PublicationDate string   `json:"publicationDate"`
	Title           string   `json:"title,omitempty"`
	Editor          string   `json:"editor,omitempty"`
	Constructors    []string `json:"constructors"`
	Body            []struct {
		Dimensions struct {
			Height int `json:"height"`
			Width  int `json:"width"`
		} `json:"dimensions"`
	} `json:"body"`
}

// Struct for parsing modern /svc/games/state JSON layouts (Midi)
type ModernStateResponse struct {
	States []struct {
		GameData struct {
			CompletionFraction int `json:"completionFraction"`
			FirstSolve         int `json:"FirstSolve"`
			PlayTimeSeconds    int `json:"playTimeSeconds"`
		} `json:"game_data"`
	} `json:"states"`
}

func FetchCrosswordSummary(ctx context.Context, client *NYTClient, gameType, dateStr string) (interface{}, error) {
	url := fmt.Sprintf("https://www.nytimes.com/svc/crosswords/v6/puzzle/%s/%s.json", gameType, dateStr)
	
	resp, err := client.DoRequest(ctx, "GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return map[string]interface{}{"status": "not_published"}, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("extraction failed with status: %d", resp.StatusCode)
	}

	var raw RawCrosswordPayload
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("unmarshaling payload error: %w", err)
	}

	var height, width int
	if len(raw.Body) > 0 {
		height = raw.Body[0].Dimensions.Height
		width = raw.Body[0].Dimensions.Width
	}

	summary := CrosswordMetadata{
		ID:              raw.ID,
		PublicationDate: raw.PublicationDate,
		Title:           raw.Title,
		Editor:          raw.Editor,
		Constructors:    raw.Constructors,
	}
	summary.Dimensions.Height = height
	summary.Dimensions.Width = width

	stateURL := fmt.Sprintf("https://www.nytimes.com/svc/games/state/crossword_%s/latests?puzzle_ids=%d", gameType, raw.ID)
	if stateResp, err := client.DoRequest(ctx, "GET", stateURL); err == nil {
		defer stateResp.Body.Close()
		if stateResp.StatusCode == http.StatusOK {
			var stateData ModernStateResponse
			if err := json.NewDecoder(stateResp.Body).Decode(&stateData); err == nil && len(stateData.States) > 0 {
				summary.CompletionFraction = stateData.States[0].GameData.CompletionFraction
				summary.FirstSolve = stateData.States[0].GameData.FirstSolve
				summary.PlayTime = stateData.States[0].GameData.PlayTimeSeconds
			}
		}
	}

	return summary, nil
}
