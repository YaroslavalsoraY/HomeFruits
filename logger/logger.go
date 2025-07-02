package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	errorLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags | log.Lshortfile)
	warnLogger = log.New(os.Stdout, "WARN: ", log.LstdFlags)
	infoLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
)

func HaltOnErr(err error, messages ...string) {
	if err == nil {
		return
	}

	message := "An error occured"

	if len(messages) > 0 {
		message = fmt.Sprintf("%s: %s", message, strings.Join(messages, " "))
	}
	errorLogger.Printf("%s: %v", message, err)
	os.Exit(1)
}

func Info(message string) {
	infoLogger.Println(message)
}

func Warn(err error, messages ...string) {
	if err != nil {
		message := "A warning occured"
		if len(messages) > 0 {
			message = fmt.Sprintf("%s: %s", message, strings.Join(messages, " "))
		}
		warnLogger.Printf("%s: %v", message, err)
	}
}