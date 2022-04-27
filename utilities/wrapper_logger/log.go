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
	"github.com/0ne-zero/f4h/public_struct"
	"github.com/0ne-zero/f4h/utilities/functions/general"
	"github.com/0ne-zero/f4h/utilities/functions/setting"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type customLogrusFormatter struct{}

// Erorr levels
type InfoLevel struct{}
type DebugLevel struct{}
type WarningLevel struct{}
type ErrorLevel struct{}
type PanicLevel struct{}
type FatalLevel struct{}
type ErrorLevels interface {
	InfoLevel | DebugLevel | WarningLevel | ErrorLevel | PanicLevel | FatalLevel
}

func init() {
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
	_, err := os.OpenFile(constansts.LogFilePath, os.O_CREATE, 0775)
	if err != nil {
		fmt.Printf("Cannot create log file in %s path", constansts.LogFilePath)
		os.Exit(1)
	}
}

func (f *customLogrusFormatter) Format(e *logrus.Entry) ([]byte, error) {

	app_name, err := setting.ReadFieldInSettingData("APP_NAME")
	if err != nil {
		return nil, err
	}
	time := e.Time.UTC().Format(time.RFC3339)
	level := strings.ToUpper(e.Level.String())
	log_msg := e.Message

	if e.Caller != nil {
		file_loc := fmt.Sprintf("%s:%d", e.Caller.File, e.Caller.Line)
		log_text := fmt.Sprintf(`[%s]-[%s]-[%s]: msg='%s' file='%s'`, time, level, app_name, log_msg, file_loc)
		return []byte(log_text), nil
	}
	log_text := fmt.Sprintf(`[%s]-[%s]-[%s]: log_msg='%s'`, time, level, app_name, log_msg)
	return []byte(log_text), nil
}

func Log[error_level ErrorLevels](e_level *error_level, log_things interface{}, err_file *public_struct.ErroredFileInfo) {
	file, err := os.OpenFile(constansts.LogFilePath, os.O_APPEND|os.O_WRONLY, 0775)
	if err != nil {
		fmt.Printf("Cannot open log file in %s path", constansts.LogFilePath)
		os.Exit(1)
	}
	// Close file
	defer file.Close()

	// Convert interface to string
	log_msg, ok := log_things.(string)
	if !ok {
		log_error, ok := log_things.(error)
		if !ok {
			// Log and exit from program
			err_file_info, err := general.GetCallerInfo(1)
			if err != nil {
				err_file_info, err = general.GetCallerInfo(0)
				Log(&FatalLevel{}, "Error occurred during get caller info", &err_file_info)
			}
			Log(&FatalLevel{}, "Error occurred during convert interface{} to string in Log function", &err_file_info)
		}
		log_msg = log_error.Error()
	}

	// Add errored file information to log_msg
	log_msg = fmt.Sprintf("%s file='%s:%d'", log_msg, err_file.Path, err_file.Line)
	// Add new line to log_msg
	log_msg = fmt.Sprintf("%s%s", log_msg, "\n")

	// Set logrus output
	logger.SetOutput(file)

	type_of_error := fmt.Sprintf("%T", e_level)[16:]
	// Log With logrus function
	switch type_of_error {
	case "InfoLevel":
		logger.Info(log_msg)
	case "DebugLevel":
		logger.Debug(log_msg)
	case "WarningLevel":
		logger.Warning(log_msg)
	case "ErrorLevel":
		logger.Error(log_msg)
	case "PanicLevel":
		logger.Panic(log_msg)
	case "FatalLevel":
		file.Close()
		logger.Fatal(log_msg)
	}
}

func AddFieldsToString(s string, fields map[string]string) string {
	for k, v := range fields {
		s = fmt.Sprintf("%s %s='%s'", s, k, v)
	}
	return s
}
