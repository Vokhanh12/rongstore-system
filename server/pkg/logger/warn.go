
type WarnParams struct {
	BaseLogLevel
}

func LogWarn(ctx context.Context, msg string, extra AccessParams) {
	fields := buildFieldsAccess(ctx, extra)
	AccessLogger.Info(msg, fields...)
}