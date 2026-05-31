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
