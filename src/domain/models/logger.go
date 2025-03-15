package models

import (
	"log"
)

type Logger struct {
}

func (l Logger) LogInfo(message string) {
	log.Printf("[INFO] %s", message)
}

func (l Logger) LogWarning(message string) {
	log.Printf("[WARN] %s", message)
}

func (l Logger) LogError(message string) {
	log.Printf("[ERROR] %s", message)
}
