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

// logToFile writes a message to a file, prefixed with the current date and time.
// If an error occurs during the write process, it is logged to the console.
func (l Logger) logToFile(filename string, content string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	buf := []byte(now + " " + content + "\n")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	_, err = file.Write(buf)
	_ = file.Close()
	if err != nil {
		log.Println(err)
	}
}
