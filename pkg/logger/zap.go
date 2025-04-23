package logger

import (
	"github.com/duyanhitbe/go-ecom/internal/config"
	"github.com/duyanhitbe/go-ecom/pkg/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type ZapLogger struct {
	logger *zap.Logger
}

type levelEnabler struct {
}

func (l levelEnabler) Enabled(lvl zapcore.Level) bool {
	if config.Cfg.Server.Mode == constants.DevelopmentMode {
		return true
	}

	return lvl != zapcore.DebugLevel
}

func newEncoderConfig() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	if config.Cfg.Server.Mode == constants.DevelopmentMode {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

func newWriter() zapcore.WriteSyncer {
	syncConsole := zapcore.AddSync(os.Stderr)
	syncFile := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Cfg.Logger.Filename,
		MaxSize:    config.Cfg.Logger.MaxSize, // megabytes
		MaxBackups: config.Cfg.Logger.MaxBackups,
		MaxAge:     config.Cfg.Logger.MaxAge,   //days
		Compress:   config.Cfg.Logger.Compress, // disabled by default
	})
	return zapcore.NewMultiWriteSyncer(syncConsole, syncFile)
}

func NewZapLogger() *ZapLogger {
	encoder := newEncoderConfig()
	writer := newWriter()

	core := zapcore.NewCore(encoder, writer, new(levelEnabler))
	logger := zap.New(core)

	return &ZapLogger{
		logger: logger,
	}
}

func (z ZapLogger) Info(msg string) {
	z.logger.Info(msg)
}

func (z ZapLogger) Debug(msg string) {
	z.logger.Debug(msg)
}

func (z ZapLogger) Warn(msg string) {
	z.logger.Warn(msg)
}

func (z ZapLogger) Error(msg string) {
	z.logger.Error(msg)
}
