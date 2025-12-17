package logger

import (
	"fmt"
	"log"
	"time"
)

func Info(msg string, v ...interface{}) {
	log.Printf("[INFO] %s - %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, v...))
}

func Error(msg string, v ...interface{}) {
	log.Printf("[ERROR] %s - %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, v...))
}
