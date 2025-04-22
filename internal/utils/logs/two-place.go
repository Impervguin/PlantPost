package logs

import (
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

func MapLevels(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func InitTwoPlaceLogger(tpCfg *TwoPlaceConfig) *zap.SugaredLogger {
	var cfg zapcore.EncoderConfig
	switch tpCfg.Type {
	case TypeDev:
		cfg = zap.NewDevelopmentEncoderConfig()
	case TypeProd:
		cfg = zap.NewProductionEncoderConfig()
	default:
		cfg = zap.NewDevelopmentEncoderConfig()
	}
	fileEncoder := zapcore.NewJSONEncoder(cfg)
	consoleEncoder := zapcore.NewConsoleEncoder(cfg)
	tee := zapcore.NewTee(
		zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			MapLevels(tpCfg.ConsoleLevel),
		),
		zapcore.NewCore(
			fileEncoder,
			zapcore.Lock(tpCfg.LogFileFactory),
			MapLevels(tpCfg.FileLevel),
		),
	)

	l := zap.New(tee,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller())

	s := l.Sugar()
	return s
}
