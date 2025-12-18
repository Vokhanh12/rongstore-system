
type WarnParams struct {
	LogEntry
}

func LogWarn(ctx context.Context, msg string, extra AccessParams) {
	fields := buildFieldsAccess(ctx, extra)
	AccessLogger.Info(msg, fields...)
}