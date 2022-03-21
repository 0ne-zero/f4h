// This package is for open log file before logging and close it after logged
package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0ne-zero/f4h/utilities"
	"github.com/sirupsen/logrus"
)

type logrus_func func(args ...interface{})

var logger *logrus.Logger
var mainDirectory string
var appName string
var logFileDirectory string
var logFileParentDirectory string

func init() {
	logger = logrus.New()
	mainDirectory = filepath.Dir(os.Args[0])
	setting_file_path := filepath.Join(mainDirectory, "setting.json")
	// Read setting file and get appName and logFileParentDirectory
	appName, err := utilities.ReadFieldInSettingFile(setting_file_path, "APP_NAME")
	if err != nil {
		fmt.Println("APP_NAME doesn't exists in setting file")
		os.Exit(1)
	} else if appName == "" {
		fmt.Println("APP_NAME is empty in setting file")
		os.Exit(1)
	}
	logFileParentDirectory, err := utilities.ReadFieldInSettingFile(setting_file_path, "LOG_FILE_PARENT_DIRECTORY")
	if err != nil {
		fmt.Println("LOG_FILE_PARENT_DIRECTORY doesn't exists in setting file")
		os.Exit(1)
	} else if logFileParentDirectory == "" {
		fmt.Println("LOG_FILE_PARENT_DIRECTORY is empty in setting file")
		os.Exit(1)
	}
	logFileDirectory = filepath.Join(logFileParentDirectory, appName)
	// Create log file directory
	if _, err := os.Stat(logFileDirectory); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(logFileDirectory, 0755)
		if err != nil {
			fmt.Printf("MKdirAll: Cannot create directory in %s path", logFileDirectory)
		}
	}
}
func Log(lf logrus_func, args ...interface{}) {
	log_file_path := filepath.Join(logFileParentDirectory, appName, "log.txt")
	file, err := os.OpenFile(log_file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		fmt.Printf("Cannot open log file in %s path", log_file_path)
		os.Exit(1)
	}
	// Close file
	defer file.Close()

	// Set logrus output
	logger.Out = file
	// Log With logrus function
	lf()

}
