package submission

type NewSubmissionRequest struct {
	Code       string `json:"code" validate:"required"`
	LanguageId string `json:"language_id" validate:"required"`
	SessionId  string `json:"session_id"`
}
