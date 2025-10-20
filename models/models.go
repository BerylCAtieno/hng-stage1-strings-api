package models

import "github.com/BerylCAtieno/string-analyzer/analysis"

type StringReponsePayload struct {
	Id         string              `json:"id"`
	Value      string              `json:"value"`
	Properties analysis.Properties `json:"properties"`
	CreatedAt  string              `json:"created_at"`
}

type FilterPayload struct {
	Data           []analysis.Properties `json:"data"`
	Count          int                   `json:"count"`
	FiltersApplied map[string]any        `json:"filters_applied"`
}

type InterpretedQuery struct {
	Original      string                `json:"original"`
	ParsedFilters []analysis.Properties `json:"parsed_filters"`
}

type NLPFilterPayload struct {
	Data             []analysis.Properties `json:"data"`
	Count            int                   `json:"count"`
	InterpretedQuery InterpretedQuery      `'json:"interpreted_query"`
}
