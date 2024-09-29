package request

type NewExecutionRequest struct {
	Id           string `json:"id"`
	Code         string `json:"code"`
	LanguageId   string `json:"language_id"`
	ConnectionId string `json:"connection_id"`
	StdIn        string `json:"stdin"`
}
