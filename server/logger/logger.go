package logger

import (
	"log"
	"os"
)

// Logger - contains the Info, Warning and Error logger
type Logger struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

// Init - Initialises the Loggers with the appropriate files
func (l *Logger) Init() {
	file, err := os.OpenFile("divert.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Panic("Logger: Cannot create divert.log file")
	}

	l.Info = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Warning = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Error = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// WriteInfo - write an info message to the log file as well as stderr
func (l *Logger) WriteInfo(msg string) {
	l.Info.Println(msg)
	log.Println(msg)
}

// WriteWarning - write a warning to the log file as well as stderr
func (l *Logger) WriteWarning(msg string) {
	l.Warning.Println(msg)
	log.Println(msg)
}

// WriteError  - write an error to the log file as well as stderr
func (l *Logger) WriteError(msg string) {
	l.Error.Println(msg)
	log.Fatal(msg)
}
