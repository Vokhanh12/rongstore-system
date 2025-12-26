package logger

// func LogBusinessError(ctx context.Context, err error, extra map[string]interface{}) {
// 	be := errors.GetBusinessError(err)

// 	fields := []zap.Field{
// 		zap.String("error.code", be.Code),
// 		zap.String("error.message", be.Message),
// 		zap.String("severity", be.Severity),
// 	}

// 	fields = append(fields, buildFields(ctx, extra)...)
// 	BusinessLogger.With(fields...).Warn("business error")
// }
