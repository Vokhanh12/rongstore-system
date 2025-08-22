package logger

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// L is the global logger
var L *zap.Logger

// Init initializes zap with two cores:
// - accessCore -> stdout (all requests, info+)
// - auditCore -> file (audit events, info+), separate file for long retention
func Init(prod bool) error {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	// stdout core for access + general logs
	stdout := zapcore.Lock(os.Stdout)
	accessCore := zapcore.NewCore(jsonEncoder, stdout, zap.InfoLevel)

	// audit core (separate file)
	auditFile, err := os.OpenFile("logs/audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}
	auditCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(auditFile), zap.InfoLevel)

	// combine cores with tee
	core := zapcore.NewTee(accessCore, auditCore)

	// add sampling to reduce high-frequency logs if desired
	sampledCore := zapcore.NewSamplerWithOptions(core, time.Second, 100, 100) // tune params

	logger := zap.New(sampledCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	L = logger
	zap.ReplaceGlobals(L)
	return nil
}

// helper: hash token prefix to safely log
func TokenHashPrefix(token string) string {
	if token == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])[:8]
}

// helper: mask email
func MaskEmail(email string) string {
	// very simple mask: keep first char and domain
	var out map[string]string
	_ = json.Unmarshal([]byte("{}"), &out) // placeholder to keep function consistent
	// naive implementation:
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			if i > 1 {
				return email[:1] + "***" + email[i:]
			}
			return "***" + email[i:]
		}
	}
	return "***"
}
