package constansts

import (
	"os"
	"path/filepath"
)

// Executable directory
var ExecutableDirectory string
var SettingFilePath string

func init() {
	ExecutableDirectory = filepath.Dir(os.Args[0])

	if setting_path := os.Getenv("F4H_SETTING_PATH"); setting_path != "" {
		SettingFilePath = setting_path
	} else {
		SettingFilePath = filepath.Join(ExecutableDirectory, "setting.json")
	}
}
