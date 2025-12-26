package logger

import (
	"context"

	"go.uber.org/zap"
)

func buildFieldsAccess(ctx context.Context, p AccessParams) []zap.Field {
	fields := []zap.Field{
		zap.String("path", p.Path),
		zap.String("method", p.Method),
		zap.Int("http_code", p.HTTPCode),
		zap.String("ip", p.IP),
		zap.String("user_agent", p.UserAgent),
	}

	for k, v := range p.Extra {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}

func buildFieldsError(ctx context.Context, p ErrorParams) []zap.Field {

	fields := []zap.Field{
		zap.String("service", p.LogEntry.ServiceInfo.Name),
		zap.String("code", p.LogEntry.Code),
		zap.String("service", p.LogEntry.Key),
		zap.String("http_status", p.LogEntry.HTTPStatus),
		zap.String("grpc_code", p.LogEntry.GRPCCode),
		zap.String("message", p.LogEntry.Message),
		zap.String("cause", p.LogEntry.Cause),
		zap.String("cause_detail", p.LogEntry.CauseDetail),
		zap.String("client_action", p.LogEntry.ClientAction),
		zap.String("server_action", p.LogEntry.ServerAction),
	}

	if !p.LogEntry.RequestContext.IsEmpty() {
		fields = append(fields, zap.String("trace_id", p.LogEntry.RequestContext.TraceId))
		fields = append(fields, zap.String("user_id", p.LogEntry.RequestContext.UserId))
		fields = append(fields, zap.String("realm_id", p.LogEntry.RequestContext.RealmId))
		fields = append(fields, zap.String("client_id", p.LogEntry.RequestContext.ClientId))
	}

	for k, v := range p.Extra {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}

func buildFieldsInfo(ctx context.Context, p InfoParams) []zap.Field {

	fields := []zap.Field{
		zap.String("service", p.LogEntry.ServiceInfo.Name),
		zap.String("code", p.LogEntry.Code),
		zap.String("service", p.LogEntry.Key),
		zap.String("http_status", p.LogEntry.HTTPStatus),
		zap.String("grpc_code", p.LogEntry.GRPCCode),
		zap.String("message", p.LogEntry.Message),
		zap.String("cause", p.LogEntry.Cause),
		zap.String("cause_detail", p.LogEntry.CauseDetail),
		zap.String("client_action", p.LogEntry.ClientAction),
		zap.String("server_action", p.LogEntry.ServerAction),
	}

	if !p.LogEntry.RequestContext.IsEmpty() {
		fields = append(fields, zap.String("trace_id", p.LogEntry.RequestContext.TraceId))
		fields = append(fields, zap.String("user_id", p.LogEntry.RequestContext.UserId))
		fields = append(fields, zap.String("realm_id", p.LogEntry.RequestContext.RealmId))
		fields = append(fields, zap.String("client_id", p.LogEntry.RequestContext.ClientId))
	}

	for k, v := range p.Extra {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}

func buildFieldsWarn(ctx context.Context, p WarnParams) []zap.Field {

	fields := []zap.Field{
		zap.String("service", p.LogEntry.ServiceInfo.Name),
		zap.String("code", p.LogEntry.Code),
		zap.String("service", p.LogEntry.Key),
		zap.String("http_status", p.LogEntry.HTTPStatus),
		zap.String("grpc_code", p.LogEntry.GRPCCode),
		zap.String("message", p.LogEntry.Message),
		zap.String("cause", p.LogEntry.Cause),
		zap.String("cause_detail", p.LogEntry.CauseDetail),
		zap.String("client_action", p.LogEntry.ClientAction),
		zap.String("server_action", p.LogEntry.ServerAction),
	}

	if !p.LogEntry.RequestContext.IsEmpty() {
		fields = append(fields, zap.String("trace_id", p.LogEntry.RequestContext.TraceId))
		fields = append(fields, zap.String("user_id", p.LogEntry.RequestContext.UserId))
		fields = append(fields, zap.String("realm_id", p.LogEntry.RequestContext.RealmId))
		fields = append(fields, zap.String("client_id", p.LogEntry.RequestContext.ClientId))
	}

	for k, v := range p.Extra {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}

func buildFields(ctx context.Context, extra map[string]interface{}) []zap.Field {
	var fields []zap.Field

	for k, v := range extra {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}
