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

// NewLogger creates both standard log and CSV format log
func NewLogger(filename string) (*Logger, error) {
	// Standard log file
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err)
	}

	// CSV log file
	csvFilename := filename[:len(filename)-4] + "_csv.csv"
	csvFile, err := os.OpenFile(csvFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logFile.Close()
		return nil, fmt.Errorf("error opening CSV file: %v", err)
	}

	csvWriter := csv.NewWriter(csvFile)

	// Write CSV header if file is empty
	fileInfo, err := csvFile.Stat()
	if err == nil && fileInfo.Size() == 0 {
		headers := []string{
			"timestamp",
			"level",
			"message",
			"event_type",
			"email",
			"event_timestamp",
			"additional_info",
		}
		if err := csvWriter.Write(headers); err != nil {
			logFile.Close()
			csvFile.Close()
			return nil, fmt.Errorf("error writing CSV headers: %v", err)
		}
		csvWriter.Flush()
	}

	return &Logger{
		logger:    log.New(logFile, "", log.LstdFlags),
		csvWriter: csvWriter,
		csvFile:   csvFile,
	}, nil
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

	// Standard logging
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
		fields)

	l.logger.Println(logMsg)

	// CSV logging
	eventType := ""
	email := ""
	eventTimestamp := ""
	additionalInfo := ""

	// Extract values from keysAndValues
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 >= len(keysAndValues) {
			break
		}
		key := fmt.Sprint(keysAndValues[i])
		value := fmt.Sprint(keysAndValues[i+1])

		switch key {
		case "event":
			eventType = value
		case "email":
			email = value
		case "timestamp":
			eventTimestamp = value
		default:
			if additionalInfo != "" {
				additionalInfo += ";"
			}
			additionalInfo += fmt.Sprintf("%s=%s", key, value)
		}
	}

	// Write to CSV
	record := []string{
		timestamp,
		level,
		msg,
		eventType,
		email,
		eventTimestamp,
		additionalInfo,
	}

	if err := l.csvWriter.Write(record); err != nil {
		l.logger.Printf("Error writing to CSV: %v", err)
	}
	l.csvWriter.Flush()
}

// Close properly closes both log files
func (l *Logger) Close() error {
	l.csvWriter.Flush()
	if err := l.csvFile.Close(); err != nil {
		return err
	}
	return nil
}
