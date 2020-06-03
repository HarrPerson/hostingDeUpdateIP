package models

type FindRecords struct {
	AuthToken string `json:"authToken"`
	Filter    struct {
		Field string `json:"field"`
		Value string `json:"value"`
	} `json:"filter"`
}
