package models

type ServiceQueryResult struct {
	ServiceQueryData []QueryResult `json:"serviceQueryData"`
}

type QueryResult struct {
	Provider SystemDefinition `json:"provider"`
}
