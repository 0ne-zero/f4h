// This package is for open log file before logging and close it after logged
package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/sirupsen/logrus"
)

type logrus_func func(args ...interface{})

var logger *logrus.Logger

func init() {
	logger = logrus.New()

	// Create log file directory
	log_file_directory_path := filepath.Dir(constansts.LogFilePath)
	if _, err := os.Stat(log_file_directory_path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(log_file_directory_path, 0755)
		if err != nil {
			fmt.Printf("MKdirAll: Cannot create directory in %s path", log_file_directory_path)
		}
	}
}
func Log(lf logrus_func, args ...interface{}) {
	file, err := os.OpenFile(constansts.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		fmt.Printf("Cannot open log file in %s path", constansts.LogFilePath)
		os.Exit(1)
	}
	// Close file
	defer file.Close()

	// Set logrus output
	logger.Out = file
	// Log With logrus function
	lf()

}
