package logger

import (
	"bytes"
	"github.com/duyanhitbe/go-ecom/internal/config"
	"github.com/duyanhitbe/go-ecom/pkg/constants"
	"go.uber.org/zap"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func NewZapLoggerWithWriter(buffer *bytes.Buffer) *ZapLogger {
	encoder := newEncoderConfig()
	writer := zapcore.AddSync(buffer)
	core := zapcore.NewCore(encoder, writer, new(levelEnabler))
	zapLogger := zap.New(core)
	return &ZapLogger{
		logger: zapLogger,
	}
}

func TestNewZapLogger(t *testing.T) {
	type testcase struct {
		name string
		mode string
	}

	testcases := []testcase{
		{
			name: "development mode",
			mode: constants.DevelopmentMode,
		},
		{
			name: "production mode",
			mode: constants.ProductionMode,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			config.Cfg.Server.Mode = tc.mode
			logger := NewZapLogger()
			require.NotNil(t, logger)
		})
	}
}

func TestZapLogger_Info(t *testing.T) {
	var buffer bytes.Buffer
	logger := NewZapLoggerWithWriter(&buffer)

	logger.Info("info message")

	require.Contains(t, buffer.String(), "info message")
	require.Contains(t, buffer.String(), "INFO")
}

func TestZapLogger_Debug(t *testing.T) {
	var buffer bytes.Buffer
	logger := NewZapLoggerWithWriter(&buffer)

	logger.Debug("debug message")

	if config.Cfg.Server.Mode == constants.DevelopmentMode {
		require.Contains(t, buffer.String(), "debug message")
		require.Contains(t, buffer.String(), "DEBUG")
	} else {
		assert.NotContains(t, buffer.String(), "debug message")
	}
}

func TestZapLogger_Warn(t *testing.T) {
	var buffer bytes.Buffer
	logger := NewZapLoggerWithWriter(&buffer)

	logger.Warn("warn message")

	require.Contains(t, buffer.String(), "warn message")
	require.Contains(t, buffer.String(), "WARN")
}

func TestZapLogger_Error(t *testing.T) {
	var buffer bytes.Buffer
	logger := NewZapLoggerWithWriter(&buffer)

	logger.Error("error message")

	require.Contains(t, buffer.String(), "error message")
	require.Contains(t, buffer.String(), "ERROR")
}
