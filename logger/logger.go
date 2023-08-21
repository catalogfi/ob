package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type sentryCore struct {
	zapcore.LevelEnabler
	tags []zap.Field
}

// NewSentryCore creates a new Zap Core for Sentry.
func NewSentryCore(url string, levelEnable zapcore.LevelEnabler, tags ...zap.Field) zapcore.Core {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: url,
	}); err != nil {
		panic("sentry.Init: " + err.Error())
	}
	return &sentryCore{
		LevelEnabler: levelEnable,
		tags:         tags,
	}
}

// With adds structured context to the core.
func (core *sentryCore) With(fields []zapcore.Field) zapcore.Core {
	core.tags = append(core.tags, fields...)
	return core
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
	if core.LevelEnabler.Enabled(ent.Level) {
		event := sentry.NewEvent()
		event.Level = sentrySeverity(ent.Level)
		event.Message = ent.Message
		event.Tags = core.with(fields)

		sentry.CaptureEvent(event)
	}

	return nil
}

// Sync flushes any buffered logs.
func (core *sentryCore) Sync() error {
	sentry.Flush(2 * time.Second)
	return nil
}

func (core *sentryCore) with(fs []zapcore.Field) map[string]string {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fs {
		f.AddTo(enc)
	}
	for _, f := range core.tags {
		f.AddTo(enc)
	}
	tags := map[string]string{}
	for key, val := range enc.Fields {
		tags[key] = fmt.Sprintf("%v", val)
	}
	return tags
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

func ZapDevelopmentCore() zapcore.Core {
	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	ws := zapcore.AddSync(zapcore.Lock(os.Stdout))
	return zapcore.NewCore(enc, ws, zap.DebugLevel)
}
