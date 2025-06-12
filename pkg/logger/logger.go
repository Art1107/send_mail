package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger(filename string) (*Logger, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err)
	}

	logger := log.New(f, "", log.LstdFlags)
	return &Logger{logger: logger}, nil
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.log("INFO", msg, keysAndValues...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.log("ERROR", msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.log("WARN", msg, keysAndValues...)
}

func (l *Logger) log(level, msg string, keysAndValues ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fields := make([]string, 0, len(keysAndValues)/2)

	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			fields = append(fields, fmt.Sprintf("%v=%v",
				keysAndValues[i], keysAndValues[i+1]))
		}
	}

	logMsg := fmt.Sprintf("%s [%s] %s %s",
		timestamp,
		level,
		msg,
		strings.Join(fields, " "))

	l.logger.Println(logMsg)
}
