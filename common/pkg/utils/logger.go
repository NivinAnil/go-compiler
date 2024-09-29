package utils

import (
	"context"
	"go-compiler/common/pkg/constants"

	"github.com/kodnest/go-common-libraries/logger"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// GetLogger : Return SugaredLogger instance from library according to log level set in config
// kept variadic fn as context is optional parameter
func GetLogger(ctx ...context.Context) *zap.SugaredLogger {
	logLevel := logger.LogLevel(viper.GetInt(constants.LoggerLevelKey))
	ctxCount := len(ctx)
	if ctxCount > 0 {
		log := logger.New(logLevel)
		if ctxCount > 1 {
			log.Warn("Illegal State : two contexts were found for the logger; initiating the logger using the first context")
		}
		return log
	}
	return logger.New(logLevel)
}
