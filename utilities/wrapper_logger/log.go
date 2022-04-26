// This package is for open log file before logging and close it after logged
package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/utilities/functions/setting"
	"github.com/sirupsen/logrus"
)

type logrus_func func(args ...interface{})

var logger *logrus.Logger

type CustomLogrusFormatter struct{}

func (f *CustomLogrusFormatter) Format(e *logrus.Entry) ([]byte, error) {

	app_name, err := setting.ReadFieldInSettingData("APP_NAME")
	if err != nil {
		return nil, err
	}
	time := e.Time.UTC().Format(time.RFC3339)
	level := e.Level.String()
	msg := e.Message
	file_loc := fmt.Sprintf("%s:%d", e.Caller.File, e.Caller.Line)

	log_text := fmt.Sprintf(`[%s]-[%s] %s: msg='%s' file='%s'`, time, level, app_name, msg, file_loc)
	return []byte(log_text), nil
}

func init() {
	// Create logger
	logger = logrus.New()
	// Config logger
	logger.SetFormatter(&CustomLogrusFormatter{})

	// Create log file directory
	log_file_directory_path := filepath.Dir(constansts.LogFilePath)
	if _, err := os.Stat(log_file_directory_path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(log_file_directory_path, 0755)
		if err != nil {
			fmt.Printf("MKdirAll: Cannot create directory in %s path", log_file_directory_path)
		}

	}
	_, err := os.OpenFile(constansts.LogFilePath, os.O_CREATE, 0775)
	if err != nil {
		fmt.Printf("Cannot create log file in %s path", constansts.LogFilePath)
		os.Exit(1)
	}
}
func Log(lf logrus_func, args ...interface{}) {
	file, err := os.OpenFile(constansts.LogFilePath, os.O_APPEND|os.O_WRONLY, 0775)
	if err != nil {
		fmt.Printf("Cannot open log file in %s path", constansts.LogFilePath)
		os.Exit(1)
	}
	// Close file
	defer file.Close()

	// Set logrus output
	logger.SetOutput(file)
	// Log With logrus function
	lf(args)

}
