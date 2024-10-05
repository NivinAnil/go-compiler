package submission

import (
	"go-compiler/common/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

type SubmissionController struct{}

func NewSubmissionController() *SubmissionController {
	return &SubmissionController{}
}

func (s *SubmissionController) CreateSubmission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.GetLogger(ctx)
		methodName := "CreateSubmission"
		log.Info("Inside " + methodName)
	}
}
