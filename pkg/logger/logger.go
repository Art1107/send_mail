package logger

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	logger    *log.Logger
	csvWriter *csv.Writer
	csvFile   *os.File
}

var csvHeaders = []string{
	"timestamp",
	"level",
	"message",
	"event_type",
	"email",
	"event_timestamp",
	"additional_info",
}

func NewLogger(filename string) (*Logger, error) {
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	csvFilename := filename[:len(filename)-4] + ".csv"
	csvFile, err := os.OpenFile(csvFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logFile.Close()
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}

	csvWriter := csv.NewWriter(csvFile)

	if fileInfo, err := csvFile.Stat(); err == nil && fileInfo.Size() == 0 {
		if err := csvWriter.Write(csvHeaders); err != nil {
			logFile.Close()
			csvFile.Close()
			return nil, fmt.Errorf("failed to write CSV headers: %w", err)
		}
		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			logFile.Close()
			csvFile.Close()
			return nil, fmt.Errorf("failed to flush CSV headers: %w", err)
		}
	}

	return &Logger{
		logger:    log.New(logFile, "", log.LstdFlags),
		csvWriter: csvWriter,
		csvFile:   csvFile,
	}, nil
}

type LogEntry struct {
	Timestamp      string
	Level          string
	Message        string
	EventType      string
	Email          string
	EventTimestamp string
	AdditionalInfo string
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
	entry := LogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Level:     level,
		Message:   msg,
	}

	fields := make([]string, 0, len(keysAndValues)/2)
	additionalInfo := make([]string, 0)

	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 >= len(keysAndValues) {
			break
		}
		key := fmt.Sprint(keysAndValues[i])
		value := fmt.Sprint(keysAndValues[i+1])

		fields = append(fields, fmt.Sprintf("%s=%s", key, value))

		switch key {
		case "event":
			entry.EventType = value
		case "email":
			entry.Email = value
		case "timestamp":
			entry.EventTimestamp = value
		default:
			additionalInfo = append(additionalInfo, fmt.Sprintf("%s=%s", key, value))
		}
	}

	l.logger.Printf("%s [%s] %s %v",
		entry.Timestamp,
		entry.Level,
		entry.Message,
		fields,
	)

	if len(additionalInfo) > 0 {
		entry.AdditionalInfo = joinStrings(additionalInfo, ";")
	}

	record := []string{
		entry.Timestamp,
		entry.Level,
		entry.Message,
		entry.EventType,
		entry.Email,
		entry.EventTimestamp,
		entry.AdditionalInfo,
	}

	if err := l.csvWriter.Write(record); err != nil {
		l.logger.Printf("Error writing to CSV: %v", err)
		return
	}

	l.csvWriter.Flush()
	if err := l.csvWriter.Error(); err != nil {
		l.logger.Printf("Error flushing CSV: %v", err)
	}
}

func joinStrings(items []string, sep string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += sep
		}
		result += item
	}
	return result
}

func (l *Logger) Close() error {
	l.csvWriter.Flush()
	if err := l.csvFile.Close(); err != nil {
		return err
	}
	return nil
}
