package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

func InitLog() (*os.File, error) {
	// Open a log file
	timeNow := time.Now().Format("02-01-2006")

	logFile, err := os.OpenFile(fmt.Sprintf("logs/%s_serverlogs.log", timeNow), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set the log output to the file
	log.SetOutput(logFile)

	return logFile, nil
}

func CloseLogFile(logFile *os.File) {
	if logFile != nil {
		logFile.Close()
	}
}

func LogAndPrint(msg string, a ...interface{}) {
	fmt.Printf(msg+"\n", a...)
	log.Printf(msg, a...)
}
