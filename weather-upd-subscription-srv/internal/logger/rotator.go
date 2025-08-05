package logger

import "gopkg.in/natefinch/lumberjack.v2"

func SetUpRotator(filePath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    5,
		MaxBackups: 2,
		MaxAge:     7,
		Compress:   true,
	}
}
