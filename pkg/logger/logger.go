package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New builds a zap.Logger based on the level and format strings from config.
func New(level, format string) (*zap.Logger, error) {
	zapLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if format == "console" {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zapLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0)), nil
}

// WithRequestID returns a logger with the request ID field pre-set.
func WithRequestID(l *zap.Logger, requestID string) *zap.Logger {
	return l.With(zap.String("request_id", requestID))
}

// WithTenant returns a logger with the tenant ID field pre-set.
func WithTenant(l *zap.Logger, tenantID string) *zap.Logger {
	return l.With(zap.String("tenant_id", tenantID))
}

// WithUser returns a logger with the user ID field pre-set.
func WithUser(l *zap.Logger, userID string) *zap.Logger {
	return l.With(zap.String("user_id", userID))
}
