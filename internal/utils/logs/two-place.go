package logs

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogFileFactory interface {
	GetLogFile() *os.File
	zapcore.WriteSyncer
}

type TwoPlaceConfig struct {
	Type           string
	FileLevel      string
	ConsoleLevel   string
	LogFileFactory LogFileFactory
}

const (
	TypeDev  = "dev"
	TypeProd = "prod"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

func MapLevels(level string) (zapcore.Level, error) {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel, nil
	case LevelInfo:
		return zapcore.InfoLevel, nil
	case LevelWarn:
		return zapcore.WarnLevel, nil
	case LevelError:
		return zapcore.ErrorLevel, nil
	case LevelFatal:
		return zapcore.FatalLevel, nil
	default:
		return zap.DPanicLevel, fmt.Errorf("unknown logger error")
	}
}

func InitTwoPlaceLogger(tpCfg *TwoPlaceConfig) (*zap.SugaredLogger, error) {
	var cfg zapcore.EncoderConfig
	switch tpCfg.Type {
	case TypeDev:
		cfg = zap.NewDevelopmentEncoderConfig()
	case TypeProd:
		cfg = zap.NewProductionEncoderConfig()
	default:
		return nil, fmt.Errorf("invalid log type: %s", tpCfg.Type)
	}
	fileEncoder := zapcore.NewJSONEncoder(cfg)
	if fileEncoder == nil {
		return nil, fmt.Errorf("failed to create file encoder")
	}
	consoleEncoder := zapcore.NewConsoleEncoder(cfg)
	if consoleEncoder == nil {
		return nil, fmt.Errorf("failed to create console encoder")
	}

	consoleLevel, err := MapLevels(tpCfg.ConsoleLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to map console level: %w", err)
	}
	fileLevel, err := MapLevels(tpCfg.FileLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to map file level: %w", err)
	}

	tee := zapcore.NewTee(
		zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			consoleLevel,
		),
		zapcore.NewCore(
			fileEncoder,
			zapcore.Lock(tpCfg.LogFileFactory),
			fileLevel,
		),
	)

	l := zap.New(tee,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller())
	if l == nil {
		return nil, fmt.Errorf("can't create logger")
	}

	s := l.Sugar()
	if s == nil {
		return nil, fmt.Errorf("can't sugar logger")
	}
	return s, nil
}
