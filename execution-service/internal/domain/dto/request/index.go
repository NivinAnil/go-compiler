package request

type NewExecutionRequest struct {
	Id         string `json:"id"`
	Code       string `json:"code"`
	LanguageId int64  `json:"language_id"`
	RequestId  string `json:"request_id"`
	StdIn      string `json:"stdin"`
}
