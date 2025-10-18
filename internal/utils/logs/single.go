package logs

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var (
	singletonLogger           *zap.SugaredLogger
	once                      sync.Once
	LoggerNotInitializedError = fmt.Errorf("logger not initialized")
)

func InitSigletonLogger(cfg *TwoPlaceConfig) error {
	once.Do(func() {
		var err error
		singletonLogger, err = InitTwoPlaceLogger(cfg)
		if err != nil {
			panic(err)
		}
	})
	return nil
}

func Debugf(format string, args ...interface{}) {
	once.Do(func() {
		panic(LoggerNotInitializedError)
	})
	singletonLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	once.Do(func() {
		panic(LoggerNotInitializedError)
	})
	singletonLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	once.Do(func() {
		panic(LoggerNotInitializedError)
	})
	singletonLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	once.Do(func() {
		panic(LoggerNotInitializedError)
	})
	singletonLogger.Errorf(format, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	once.Do(func() {
		panic(LoggerNotInitializedError)
	})
	singletonLogger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	once.Do(func() {
		panic(LoggerNotInitializedError)
	})
	singletonLogger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	once.Do(func() {
		panic(LoggerNotInitializedError)
	})
	singletonLogger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	once.Do(func() {
		panic(LoggerNotInitializedError)
	})
	singletonLogger.Errorw(msg, keysAndValues...)
}
