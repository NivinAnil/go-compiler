package request

type NewExecutionRequest struct {
	Id          string `json:"id"`
	Code        string `json:"code"`
	LanguageId  int64  `json:"language_id"`
	ConnctionId string `json:"connection_id"`
	StdIn       string `json:"stdin"`
}
