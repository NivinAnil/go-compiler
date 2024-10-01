package submission

import (
	"go-compiler/common/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SubmissionController struct{}

func NewSubmissionController() *SubmissionController {
	return &SubmissionController{}
}

func (s *SubmissionController) CreateSubmission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := utils.GetLogger()
		methodName := "CreateSubmission"
		log.Info("Inside " + methodName)
	}
}
