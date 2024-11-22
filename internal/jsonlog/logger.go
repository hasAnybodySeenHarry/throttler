package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	mu       sync.Mutex
	minLevel Level
}

func New(out io.Writer, l Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: l,
	}
}

func (l *Logger) print(level Level, message string, props map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties"`
		Time       string            `json:"time"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Message:    message,
		Properties: props,
		Time:       time.Now().UTC().Format(time.RFC3339),
	}

	if level > l.minLevel {
		aux.Trace = string(debug.Stack())
	}

	var line []byte
	line, err := json.Marshal(&aux)
	if err != nil {
		line = []byte(level.String() + ": unable to marshal the message " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

func (l *Logger) Info(message string, props map[string]string) {
	l.print(LevelInfo, message, props)
}

func (l *Logger) Error(err error, props map[string]string) {
	l.print(LevelError, err.Error(), props)
}

func (l *Logger) Fatal(err error, props map[string]string) {
	l.print(LevelFatal, err.Error(), props)
	os.Exit(1)
}

// server error logger needs a Writer, so we implement Write()
// to make our Logger conform to Writer.
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
