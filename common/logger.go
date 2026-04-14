package common

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	LogError LogLevel = iota
	LogWarning
	LogInfo
	LogDebug
)

type Logger struct {
	name     string
	logLevel LogLevel
}

var (
	loggerMu    sync.RWMutex
	moduleLogs  = map[string]*Logger{}
	moduleLevel = map[string]LogLevel{}
	stdLogger   = log.New(os.Stderr, "", 0)
)

func NewLogger(name string, level LogLevel) *Logger {
	return &Logger{name: name, logLevel: level}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.logLevel = level
}

func (l *Logger) Time() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

func (l *Logger) ModuleName() string {
	return l.name
}

func (l *Logger) shouldLog(level LogLevel) bool {
	return level <= l.logLevel
}

func (l *Logger) Error(args ...any) *Logger {
	stdLogger.Printf("%s %s sdk logger error %v", l.Time(), l.name, args)
	return l
}

func (l *Logger) LogWithError(args ...any) *Logger {
	panic(fmt.Errorf("%v", args))
}

func (l *Logger) Warning(args ...any) *Logger {
	if !l.shouldLog(LogWarning) {
		return l
	}
	stdLogger.Printf("%s %s sdk logger warning %v", l.Time(), l.name, args)
	return l
}

func (l *Logger) Info(args ...any) *Logger {
	if !l.shouldLog(LogInfo) {
		return l
	}
	stdLogger.Printf("%s %s sdk logger info %v", l.Time(), l.name, args)
	return l
}

func (l *Logger) Debug(args ...any) *Logger {
	if !l.shouldLog(LogDebug) {
		return l
	}
	stdLogger.Printf("%s %s sdk logger debug %v", l.Time(), l.name, args)
	return l
}

func CreateLogger(moduleName string) *Logger {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	if existing, ok := moduleLogs[moduleName]; ok {
		return existing
	}
	level, ok := moduleLevel[moduleName]
	if !ok {
		level = LogError
	}
	logger := NewLogger(moduleName, level)
	moduleLogs[moduleName] = logger
	return logger
}

func SetLoggerLevel(moduleName string, level LogLevel) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	moduleLevel[moduleName] = level
	if logger, ok := moduleLogs[moduleName]; ok {
		logger.SetLevel(level)
	}
}
