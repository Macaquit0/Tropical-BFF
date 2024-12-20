package logger

import (
	"context"
	"os"

	"github.com/backend/bff-cognito/pkg/trace"
	"github.com/rs/zerolog"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(build, applicationName string) *Logger {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("build", build).
		Str("application-name", applicationName).
		Logger()
	return &Logger{
		logger: &logger,
	}
}

func (l *Logger) log(ctx context.Context, logFunc *zerolog.Event) *zerolog.Event {
	traceId, ok := ctx.Value(trace.TraceId).(string)
	if ok {
		logFunc = logFunc.Str("trace-id", traceId)
	}
	return logFunc

}

type Event struct {
	event *zerolog.Event
}

func (e *Event) Msg(message string, args ...any) {
	e.event.Msgf(message, args...)
}

func (l *Logger) Debug(ctx context.Context) *Event {
	return &Event{l.log(ctx, l.logger.Debug())}
}

func (l *Logger) Info(ctx context.Context) *Event {
	return &Event{l.log(ctx, l.logger.Info())}
}

func (l *Logger) Warn(ctx context.Context) *Event {
	return &Event{l.log(ctx, l.logger.Warn())}
}

func (l *Logger) Error(ctx context.Context) *Event {
	return &Event{l.log(ctx, l.logger.Error())}
}
