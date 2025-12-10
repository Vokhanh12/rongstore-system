package logger

import "go.uber.org/zap/zapcore"

var SeverityToLevel = map[string]zapcore.Level{
	"S1": zapcore.ErrorLevel,
	"S2": zapcore.WarnLevel,
	"S3": zapcore.InfoLevel,
}
