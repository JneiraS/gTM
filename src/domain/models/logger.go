package models

import (
	"log"
	"os"
	"time"
)

const FILE_NAME_PATH = "log.txt"

type Logger struct {
}

func (l Logger) LogInfo(message string) {
	l.logToFile(FILE_NAME_PATH, "[INFO] "+message)
}

func (l Logger) LogWarning(message string) {
	l.logToFile(FILE_NAME_PATH, "[WARN] "+message)
}

func (l Logger) LogError(message string) {
	l.logToFile(FILE_NAME_PATH, "[ERROR] "+message)
}

func (l Logger) logToFile(filename string, content string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	now := time.Now()
	_, err = file.WriteString(now.Format("2006-01-02 15:04:05") + " " + content + "\n")
	if err != nil {
		log.Println(err)
	}
}
