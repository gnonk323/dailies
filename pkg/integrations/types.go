package integrations

import "dailies/pkg/types"

type Integration interface {
	GetName() string
	GetDescription() string
	GetType() string // "manual" or "polling"
}

type ManualIntegration interface {
	Integration
	Fetch(date string, config types.DailiesConfig) (map[string]interface{}, error)
}

type IntegrationPayload struct {
	Data      map[string]interface{} `json:"data"`
	FetchedAt string								 `json:"fetched_at"`
}
