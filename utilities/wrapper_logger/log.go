// This package is for open log file before logging and close it after logged
package wrapper_logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type customLogrusFormatter struct{}

type LogInfo struct {
	Message       string
	Fields        map[string]string
	ErrorLocation ErrorLocation
}
type ErrorLocation struct {
	FilePath string
	FuncName string
	Line     int
}

func (e *ErrorLocation) ToStringMap() map[string]string {
	var m = map[string]string{
		"Path":     e.FilePath,
		"Line":     fmt.Sprintf("%d", e.Line),
		"FuncName": e.FuncName,
	}
	return m
}

func init() {
	_, err := os.OpenFile(constansts.LogFilePath, os.O_CREATE, 0775)
	if err != nil {
		fmt.Printf("Cannot create log file in %s path", constansts.LogFilePath)
		os.Exit(1)
	}
	// Create logger
	logger = logrus.New()
	// Config logger
	logger.SetFormatter(&customLogrusFormatter{})

	// Create log file directory
	log_file_directory_path := filepath.Dir(constansts.LogFilePath)
	if _, err := os.Stat(log_file_directory_path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(log_file_directory_path, 0755)
		if err != nil {
			fmt.Printf("MKdirAll: Cannot create directory in %s path", log_file_directory_path)
		}
	}

	// Custom fatal log handler
	//logrus.RegisterExitHandler(cleanup.CleanUpResources)

}

func (f *customLogrusFormatter) Format(e *logrus.Entry) ([]byte, error) {
	time := e.Time.UTC().Format(time.RFC3339)
	level := strings.ToUpper(e.Level.String())
	log_msg := e.Message

	log_text := fmt.Sprintf("[%s]-[%s]-[%s]:\nMsg= %s\n%s\n", time, level, constansts.AppName, log_msg, strings.Repeat("-", 70))
	return []byte(log_text), nil
}
func openLogFile() *os.File {
	file, err := os.OpenFile(constansts.LogFilePath, os.O_APPEND|os.O_WRONLY, 0775)
	if err != nil {
		fmt.Printf("Cannot open log file in %s path", constansts.LogFilePath)
		os.Exit(1)
	}
	return file
}

// log_things must be error or string type
func (log_info *LogInfo) log(level string) {
	file := openLogFile()
	// Close file
	defer file.Close()
	// Create log_msg for log
	var fields string
	var log_msg string
	var error_location string
	error_location = fmt.Sprintf("%s:%d %s", log_info.ErrorLocation.FilePath, log_info.ErrorLocation.Line, log_info.ErrorLocation.FuncName)
	if log_info.Fields != nil {
		var counter int
		fields_len := len(log_info.Fields) - 1
		for k, v := range log_info.Fields {
			if counter == fields_len {
				fields += fmt.Sprintf("%s='%s'", k, v)
			} else {
				fields += fmt.Sprintf("%s='%s' | ", k, v)
			}
			counter += 1
		}
		log_msg = fmt.Sprintf("'%s'\nFields= %s\nLocation= %s", log_info.Message, fields, error_location)
	} else {
		log_msg = fmt.Sprintf("'%s'\nLocation= %s", log_info.Message, error_location)
	}

	// Set logrus output
	logger.SetOutput(file)

	// Log With logrus function
	switch level {
	case "INFO":
		logger.Info(log_msg)
	case "DEBUG":
		logger.Debug(log_msg)
	case "WARNING":
		logger.Warning(log_msg)
	case "ERROR":
		logger.Error(log_msg)
	case "PANIC":
		logger.Panic(log_msg)
	case "FATAL":
		// Log and close the program
		logger.Fatal(log_msg)
		// File will be close by os
		//file.Close()
	}
}

func Info(log_info *LogInfo) {
	log_info.log("INFO")
}
func Debug(log_info *LogInfo) {
	log_info.log("DEBUG")
}
func Warning(log_info *LogInfo) {
	log_info.log("WARNING")
}
func Error(log_info *LogInfo) {
	log_info.log("ERROR")
}
func Panic(log_info *LogInfo) {
	log_info.log("PANIC")
}
func Fatal(log_info *LogInfo) {
	log_info.log("FATAL")
}
