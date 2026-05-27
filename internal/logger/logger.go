package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func Init() {

	InfoLogger = log.New(
		os.Stdout,
		"[INFO] ",
		log.Ldate|log.Ltime,
	)

	ErrorLogger = log.New(
		os.Stdout,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	DebugLogger = log.New(
		os.Stdout,
		"[DEBUG] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
}
