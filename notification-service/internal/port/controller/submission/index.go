package submission

import (
	logger "go-compiler/common/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SubmissionController struct{}

func NewSubmissionController() *SubmissionController {
	return &SubmissionController{}
}

func (s *SubmissionController) CreateSubmission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		methodName := "CreateSubmission"
		logger.Info("Inside " + methodName)
	}
}
