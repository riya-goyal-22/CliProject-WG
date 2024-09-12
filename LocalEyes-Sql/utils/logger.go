package utils

import (
	"fmt"
	"log"
	"os"
)

var Logger *log.Logger
var file *os.File

func InitLoggerFile() {
	var err error
	file, err = os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	// Create a new logger that writes to the file
	Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func CloseLoggerFile() {
	err := file.Close()
	if err != nil {
		fmt.Println("Failed to close log file")
		return
	}
}
