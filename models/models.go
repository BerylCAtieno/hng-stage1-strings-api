package models

import "github.com/BerylCAtieno/string-analyzer/analysis"

type AnalysisReponsePayload struct {
	Id         string              `json:"id"`
	Value      string              `json:"value"`
	Properties analysis.Properties `json:"properties"`
	CreatedAt  string              `json:"created_at"`
}
