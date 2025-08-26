/**
 * Author : ngdangkietswe
 * Since  : 8/15/2025
 */

package logger

import (
	"github.com/ngdangkietswe/go-rabbitmq/pkg/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewAppLogger(env string) *zap.Logger {
	var (
		logger *zap.Logger
		err    error
	)

	if env == string(constants.EnvProduction) {
		logger, err = zap.NewProduction()
	} else {
		logger, err = NewLoggerDevelopment()
	}

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return logger
}

func NewLoggerDevelopment() (*zap.Logger, error) {
	loggerCfg := zap.NewDevelopmentConfig()
	loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	loggerCfg.Encoding = "console"
	loggerCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return loggerCfg.Build()
}
