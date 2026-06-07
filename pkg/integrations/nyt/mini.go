package nyt

import "context"

type MiniCrosswordFetcher struct{}
func (m MiniCrosswordFetcher) GameKey() string { return "crossword_mini" }
func (m MiniCrosswordFetcher) FetchGame(ctx context.Context, client *NYTClient, dateStr string, cookies string) (interface{}, error) {
	return FetchCrosswordSummary(ctx, client, "mini", dateStr)
}
