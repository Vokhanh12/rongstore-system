package logger

var (
	BASE_LOGGING ServiceInfo
)

func WithService(name string) func(*ServiceInfo) {
	return func(l *ServiceInfo) {
		l.Name = name
	}
}
