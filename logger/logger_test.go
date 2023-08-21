package logger_test

import (
	"errors"
	"os"
	"time"

	"github.com/catalogfi/wbtc-garden/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("sentry logger", func() {
	Context("when capturing an error", func() {
		It("should send it to sentry", func() {
			dsn := os.Getenv("LOGGER_DSN")
			Expect(dsn).ShouldNot(Equal(""))
			tags := []zap.Field{
				zap.String("upstream", "sentry"),
				zap.String("logger", "zap"),
			}
			core := logger.NewSentryCore(dsn, zapcore.ErrorLevel, tags...)
			lg := zap.New(core)

			err := errors.New("logger is not printing details")
			lg.Error("logger testing", zap.Error(err), zap.String("env", "local"), zap.Int64("unix", time.Now().UnixNano()))
			time.Sleep(time.Second)
		})
	})
})
