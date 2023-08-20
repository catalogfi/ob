package logger

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type sentryCore struct {
	zapcore.LevelEnabler
	fields map[string]interface{}
}

// NewSentryCore creates a new Zap Core for Sentry.
func NewSentryCore(url string, enab zapcore.LevelEnabler) zapcore.Core {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: url,
	}); err != nil {
		panic("sentry.Init: " + err.Error())
	}
	return zapcore.NewTee(zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(zapcore.Lock(os.Stdout)), zap.DebugLevel), &sentryCore{
		LevelEnabler: zapcore.ErrorLevel,
		fields:       make(map[string]interface{}),
	})
}

// With adds structured context to the core.
func (core *sentryCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *core
	clone.fields = make(map[string]interface{}, len(core.fields)+len(fields))
	for k, v := range core.fields {
		clone.fields[k] = v
	}
	return &clone
}

// Check determines whether the logger should log at the given level.
func (core *sentryCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if core.Enabled(ent.Level) {
		return ce.AddCore(ent, core)
	}
	return ce
}

// Write logs the entry and fields supplied at the log site.
func (core *sentryCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	event := sentry.NewEvent()
	event.Level = sentrySeverity(ent.Level)
	event.Message = ent.Message
	event.Extra = core.fields

	sentry.CaptureEvent(event)
	return nil
}

// Sync flushes any buffered logs.
func (core *sentryCore) Sync() error {
	sentry.Flush(2 * time.Second)
	return nil
}

// sentrySeverity converts Zap's log level to Sentry's severity.
func sentrySeverity(level zapcore.Level) sentry.Level {
	switch level {
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.InfoLevel:
		return sentry.LevelInfo
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		return sentry.LevelFatal
	case zapcore.FatalLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelInfo
	}
}
