package github

import (
	"context"
	"dailies/pkg/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GitHubModule struct{}

func (g GitHubModule) GetName() string        { return "github" }
func (g GitHubModule) GetDescription() string { return "Fetches detailed metadata, SHAs, and resource links using the GitHub Search API" }
func (g GitHubModule) GetType() string        { return "manual" }

type GitHubMetadata struct {
	CommitsCount int               `json:"commits_count"`
	Commits      map[string]string `json:"commits"`       // key: commit message, val: sha
	ReposTouched map[string]string `json:"repos_touched"` // eg. gnonk323/dailies -> https://github.com/gnonk323/dailies
	PRsOpened    map[string]string `json:"prs_opened"`    // key: pr title, val: url
	PRsMerged    map[string]string `json:"prs_merged"`    // key: pr title, val: url
}

func (g GitHubModule) Fetch(dateStr string, config types.DailiesConfig) (map[string]interface{}, error) {
	username := config.GitHub.Username
	token := config.GitHub.Token

	if username == "" || token == "" {
		return nil, fmt.Errorf("missing 'github.username' or 'github.token' in config")
	}

	// Validate date format (YYYY-MM-DD)
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	meta := GitHubMetadata{
		Commits:      make(map[string]string),
		ReposTouched: make(map[string]string),
		PRsOpened:    make(map[string]string),
		PRsMerged:    make(map[string]string),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	// FETCH COMMITS via Search API
	commitQuery := fmt.Sprintf("author:%s+committer-date:%s", username, dateStr)
	commitURL := fmt.Sprintf("https://api.github.com/search/commits?q=%s", commitQuery)

	var commitSearchResponse struct {
		Items []struct {
			SHA        string `json:"sha"`
			Commit     struct {
				Message string `json:"message"`
			} `json:"commit"`
			Repository struct {
				FullName string `json:"full_name"`
				HTMLURL  string `json:"html_url"`
			} `json:"repository"`
		} `json:"items"`
	}

	if err := g.makeSearchRequest(ctx, commitURL, token, &commitSearchResponse); err != nil {
		return nil, fmt.Errorf("commit search failed: %w", err)
	}

	meta.CommitsCount = len(commitSearchResponse.Items)
	for _, item := range commitSearchResponse.Items {
		meta.Commits[item.Commit.Message] = item.SHA
		if item.Repository.FullName != "" {
			meta.ReposTouched[item.Repository.FullName] = item.Repository.HTMLURL
		}
	}

	// FETCH PRs OPENED
	prQuery := fmt.Sprintf("author:%s+type:pr+created:%s", username, dateStr)
	prURL := fmt.Sprintf("https://api.github.com/search/issues?q=%s", prQuery)

	var prSearchResponse struct {
		Items []struct {
			Title   string `json:"title"`
			HTMLURL string `json:"html_url"`
		} `json:"items"`
	}

	if err := g.makeSearchRequest(ctx, prURL, token, &prSearchResponse); err != nil {
		return nil, fmt.Errorf("pr search failed: %w", err)
	}

	for _, item := range prSearchResponse.Items {
		meta.PRsOpened[item.Title] = item.HTMLURL
	}

	// FETCH PRs MERGED
	mergedQuery := fmt.Sprintf("author:%s+type:pr+is:merged+merged:%s", username, dateStr)
	mergedURL := fmt.Sprintf("https://api.github.com/search/issues?q=%s", mergedQuery)

	var mergedSearchResponse struct {
		Items []struct {
			Title   string `json:"title"`
			HTMLURL string `json:"html_url"`
		} `json:"items"`
	}

	if err := g.makeSearchRequest(ctx, mergedURL, token, &mergedSearchResponse); err != nil {
		return nil, fmt.Errorf("merged pr search failed: %w", err)
	}

	for _, item := range mergedSearchResponse.Items {
		meta.PRsMerged[item.Title] = item.HTMLURL
	}

	return map[string]interface{}{
		"commits_count": meta.CommitsCount,
		"commits":       meta.Commits,
		"repos":         meta.ReposTouched,
		"prs_opened":    meta.PRsOpened,
		"prs_merged":    meta.PRsMerged,
	}, nil
}

// helper function to handle the API boilerplate
func (g GitHubModule) makeSearchRequest(ctx context.Context, url string, token string, target interface{}) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "dailies-cli")

	// Required custom media type header for the Commit Search API
	req.Header.Set("Accept", "application/vnd.github.cloak-preview+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
