package errors

import "fmt"

type AppError struct {
	Code         string                 `json:"code"`
	Status       int                    `json:"status"`
	GRPCCode     string                 `json:"grpc_code"`
	Key          string                 `json:"key"`
	Cause        string                 `json:"cause"`
	ClientAction string                 `json:"client_action"`
	ServerAction string                 `json:"server_action"`
	Message      string                 `json:"message"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Severity     string                 `json:"severity,omitempty"`
	Retryable    bool                   `json:"retryable,omitempty"`
}

func (e *AppError) Error() string {
	if e == nil {
		return "<nil AppError>"
	}
	return fmt.Sprintf("%s (%d): %s", e.Code, e.Status, e.Message)
}
