package domain

import audit "github.com/berezovskiydeval/logger-mongo/pkg/domain"

type LogMessage struct {
	Context map[string]interface{} `json:"context"`
	LogItem audit.LogItem          `json:"log_item"`
}
