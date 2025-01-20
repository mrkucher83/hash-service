package logger

import (
	"context"
	formatters "github.com/fabienm/go-logrus-formatters"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	lg *logrus.Logger
}

func NewLogrusLogger() *Logger {
	l := logrus.New()
	l.SetFormatter(formatters.NewGelf("hasher"))
	if env := os.Getenv("ENVIRONMENT"); env == "prod" {
		l.SetLevel(logrus.InfoLevel)
	}

	return &Logger{l}
}

var globalLogger Logger

func InitLogger(logger *Logger) {
	globalLogger = *logger
}

func Info(msg string, args ...interface{}) {
	globalLogger.lg.Infof(msg, args...)
}

func Error(msg string, err error) {
	globalLogger.lg.WithFields(logrus.Fields{
		"stack": errors.WithStack(err),
	}).Errorf(msg, err)
}

func Warn(args ...interface{}) {
	globalLogger.lg.Warn(args...)
}

func Fatal(msg string, err error) {
	globalLogger.lg.WithFields(logrus.Fields{
		"stack": errors.WithStack(err),
	}).Fatalf(msg, err)
}

func InfoCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.lg.WithFields(logrus.Fields{
		"requestID": ctx.Value("requestID").(string),
	}).Infof(msg, args...)
}

func ErrorCtx(ctx context.Context, msg string, err error) {
	globalLogger.lg.WithFields(logrus.Fields{
		"requestID": ctx.Value("requestID").(string),
		"stack":     errors.WithStack(err),
	}).Errorf(msg, err)
}
