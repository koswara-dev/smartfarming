package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var LogFile *os.File

func InitLogger() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	updateLogFile()

	// Initial clean up of logs older than 30 days
	CleanupOldLogs(30)

	// Start a ticker to check for rotation and cleanup old logs hourly
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			updateLogFile()
			CleanupOldLogs(30)
		}
	}()
}

func updateLogFile() {
	currentDate := time.Now().Format("2006-01-02")
	logPath := filepath.Join("logs", fmt.Sprintf("app-%s.log", currentDate))

	// Check if file is already open for today
	if LogFile != nil {
		stat, err := LogFile.Stat()
		if err == nil && stat.Name() == fmt.Sprintf("app-%s.log", currentDate) {
			return
		}
		LogFile.Close()
	}

	var err error
	LogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open log file %s: %v", logPath, err)
		return
	}

	// Output to both console and log file
	mw := io.MultiWriter(os.Stdout, LogFile)
	log.SetOutput(mw)

	// Set Gin default and error writer to the MultiWriter
	gin.DefaultWriter = mw
	gin.DefaultErrorWriter = mw
}

func CleanupOldLogs(days int) {
	files, err := os.ReadDir("logs")
	if err != nil {
		log.Printf("Failed to read logs directory: %v", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -days)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		info, err := f.Info()
		if err != nil {
			continue
		}

		// Check if file modification time is older than the cutoff
		if info.ModTime().Before(cutoff) {
			path := filepath.Join("logs", f.Name())
			if err := os.Remove(path); err != nil {
				log.Printf("Failed to remove old log file %s: %v", path, err)
			} else {
				log.Printf("Removed old log file: %s", path)
			}
		}
	}
}
