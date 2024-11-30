package config

import (
	"log"
	"os"
)

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	RESET  = "\033[0m"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger(flags int) *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "", flags),
		errorLogger: log.New(os.Stderr, "", flags),
	}
}

func (l *Logger) Info(msg any) {
	l.infoLogger.Println(msg)
}

func (l *Logger) Success(msg any) {
	l.infoLogger.Printf("%s%v%s\n", GREEN, msg, RESET)
}

func (l *Logger) Warning(msg any) {
	l.errorLogger.Printf("%s%v%s\n", YELLOW, msg, RESET)
}

func (l *Logger) Error(msg any) {
	l.errorLogger.Printf("%s%v%s\n", RED, msg, RESET)
}
