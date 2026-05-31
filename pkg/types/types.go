package types

type WordBank struct {
	PromotionThreshold int      `json:"promotion_threshold"`
	MaxWords           int      `json:"max_words"`
	Moods              []string `json:"moods"`
	ContextTags        []string `json:"context_tags"`
}

type GitHubConfig struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type DailiesConfig struct {
	Staging      map[string]string          `json:"staging,omitempty"`
	WordBank     WordBank                   `json:"word_bank"`
	Integrations map[string]map[string]bool `json:"integrations"`
	GitHub       GitHubConfig								`json:"github"`
}

type DailyEntry struct {
	Date         string                 `json:"date"`
	DayQuality   int                    `json:"day_quality"`
	Moods        []string               `json:"moods"`
	ContextTags  []string               `json:"context_tags"`
	Journal      string                 `json:"journal"`
	Integrations map[string]interface{} `json:"integrations"`
}
