package nyt

import "context"

type MidiCrosswordFetcher struct{}
func (m MidiCrosswordFetcher) GameKey() string { return "crossword_midi" }
func (m MidiCrosswordFetcher) FetchGame(ctx context.Context, client *NYTClient, dateStr string, cookies string) (interface{}, error) {
	return FetchCrosswordSummary(ctx, client, "midi", dateStr)
}
