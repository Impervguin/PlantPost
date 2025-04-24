package main

import "github.com/spf13/viper"

const (
	LogConsoleLevelKey = "console_level"
	LogFileLevelKey    = "file_level"
	LogDirKey          = "dir"
	LogTypeKey         = "type"
	FileFactoryKey     = "file_factory"
	LoggerPrefix       = "log"
	LogFileTypekey     = "file_type"
)

const (
	EveryDayFileFactory = "every_day"
)

const (
	LogFileTypeJson = "json"
)

type LoggerConfig struct {
	LogConsoleLevel string
	LogFileLevel    string
	LogDir          string
	LogType         string
	FileFactory     string
	LogFileType     string
}

func GetLoggerConfig() *LoggerConfig {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	conf := LoggerConfig{
		LogFileLevel:    viper.GetString(Key(LoggerPrefix, LogFileLevelKey)),
		LogConsoleLevel: viper.GetString(Key(LoggerPrefix, LogConsoleLevelKey)),
		LogDir:          viper.GetString(Key(LoggerPrefix, LogDirKey)),
		LogType:         viper.GetString(Key(LoggerPrefix, LogTypeKey)),
		FileFactory:     viper.GetString(Key(LoggerPrefix, FileFactoryKey)),
		LogFileType:     viper.GetString(Key(LoggerPrefix, LogFileTypekey)),
	}
	return &conf
}
