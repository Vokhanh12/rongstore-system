package logger

type Option func(*ServiceInfo)

func WithService(name string) Option {
	return func(s *ServiceInfo) {
		s.Name = name
	}
}
