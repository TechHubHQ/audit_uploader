package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// fileLogger receives every level (debug and up); consoleLogger only info and up,
// so the terminal stays readable while the log file keeps the full trace.
var fileLogger *log.Logger
var consoleLogger *log.Logger

func InitLogger() error {
	// Create logs directory if it doesn't exist
	dirName := time.Now().Format(time.DateOnly)
	logDir := fmt.Sprintf("logs/%s", dirName)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp in filename
	timestamp := time.Now().Format("15-04-05")
	logFile := filepath.Join(logDir, fmt.Sprintf("auditUploader_%s.log", timestamp))

	// Open or create the log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	fileLogger = log.New(file, "", log.LstdFlags)
	consoleLogger = log.New(colorWriter{writer: os.Stdout}, "", log.LstdFlags)

	return nil
}

func Info(msg string, args ...any) {
	record := formatRecord("INFO", msg, args...)
	fileLogger.Print(record)
	consoleLogger.Print(record)
}

func Error(msg string, args ...any) {
	record := formatRecord("ERROR", msg, args...)
	fileLogger.Print(record)
	consoleLogger.Print(record)
}

func Warn(msg string, args ...any) {
	record := formatRecord("WARN", msg, args...)
	fileLogger.Print(record)
	consoleLogger.Print(record)
}

func Debug(msg string, args ...any) {
	// debug only goes to the file, not the terminal
	fileLogger.Print(formatRecord("DEBUG", msg, args...))
}

func Infof(format string, v ...any) {
	Info(fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...any) {
	Error(fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...any) {
	Warn(fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...any) {
	Debug(fmt.Sprintf(format, v...))
}

func formatRecord(level, msg string, args ...any) string {
	var b strings.Builder
	b.WriteString("[")
	b.WriteString(level)
	b.WriteString("] ")
	b.WriteString(msg)

	var multilineFields strings.Builder
	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprintf("arg_%d", i)
		if k, ok := args[i].(string); ok && k != "" {
			key = k
		}

		var val any = "<missing>"
		if i+1 < len(args) {
			val = args[i+1]
		}

		formatted, multiline := formatValue(val)
		if multiline {
			multilineFields.WriteString("\n  ")
			multilineFields.WriteString(key)
			multilineFields.WriteString(":\n")
			multilineFields.WriteString(indentLines(formatted, "    "))
		} else {
			b.WriteString(" ")
			b.WriteString(key)
			b.WriteString("=")
			b.WriteString(formatted)
		}
	}

	b.WriteString(multilineFields.String())

	return b.String()
}

func formatValue(v any) (string, bool) {
	if v == nil {
		return "<nil>", false
	}

	switch t := v.(type) {
	case error:
		return t.Error(), false
	case string:
		return t, false
	case []byte:
		return string(t), false
	}

	k := reflect.TypeOf(v).Kind()
	if k == reflect.Struct || k == reflect.Map || k == reflect.Slice || k == reflect.Array {
		pretty, err := json.MarshalIndent(v, "", "  ")
		if err == nil {
			return string(pretty), true
		}
	}

	return fmt.Sprint(v), false
}

func indentLines(s, indent string) string {
	if s == "" {
		return indent
	}

	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = indent + lines[i]
	}
	return strings.Join(lines, "\n")
}

type colorWriter struct {
	writer io.Writer
}

func (w colorWriter) Write(p []byte) (int, error) {
	if strings.Contains(string(p), "[ERROR]") {
		colored := "\033[31m" + string(p) + "\033[0m"
		if _, err := io.WriteString(w.writer, colored); err != nil {
			return 0, err
		}
		return len(p), nil
	}

	return w.writer.Write(p)
}
