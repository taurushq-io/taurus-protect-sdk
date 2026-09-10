package service

import "context"

// Field is one structured key/value pair on a log line.
//
// SECURITY: metadata only — an identifier, a resource name, a reason. Never a
// payload, a hash preimage, a token or key material. The SDK's own call sites are
// pinned by TestNoLoggerCallLeaksPayload; a consumer's adapter must forward these
// verbatim and add nothing.
type Field struct {
	Key   string
	Value any
}

// Logger receives SDK diagnostic events.
//
// ctx is the first parameter so a consumer whose logger derives fields from it —
// a correlation id, a tenant — keeps them. A logger shaped like log/slog
// (msg string, args ...any) has nowhere to put ctx, which would make those fields
// structurally unrecoverable; the SDK holds a live ctx at every call site, and the
// events it reports are exactly the ones worth correlating to a request.
//
// The SDK deliberately depends on no logging library. Adapting a concrete logger
// is a few lines on the consumer's side.
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
}

// NoopLogger discards everything. It is the default so that call sites never need
// a nil check — a nil interface would panic on the first method call, and the
// first call is on an error path that is rare in production and easy to miss.
type NoopLogger struct{}

func (NoopLogger) Debug(context.Context, string, ...Field) {}
func (NoopLogger) Info(context.Context, string, ...Field)  {}
func (NoopLogger) Warn(context.Context, string, ...Field)  {}
func (NoopLogger) Error(context.Context, string, ...Field) {}

// ServiceOption configures a service at construction. It exists for the logger;
// keep it that way rather than growing an options framework.
type ServiceOption func(*serviceOptions)

type serviceOptions struct {
	logger Logger
}

// WithServiceLogger sets the logger for a service.
func WithServiceLogger(l Logger) ServiceOption {
	return func(o *serviceOptions) {
		if l != nil {
			o.logger = l
		}
	}
}

// applyServiceOptions resolves options to a usable config, defaulting the logger
// to the no-op rather than nil.
func applyServiceOptions(opts []ServiceOption) serviceOptions {
	resolved := serviceOptions{logger: NoopLogger{}}
	for _, opt := range opts {
		if opt != nil {
			opt(&resolved)
		}
	}
	return resolved
}
