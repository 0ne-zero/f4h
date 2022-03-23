package constansts

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Executable directory
var ExecutableDirectory string
var SettingFilePath string

// All gin routes that created
// This field is filled by "github.com/0ne-zero/f4h/web/route.MakeRoute" function
var Routes gin.RoutesInfo

func init() {
	ExecutableDirectory = filepath.Dir(os.Args[0])

	if setting_path := os.Getenv("F4H_SETTING_PATH"); setting_path != "" {
		SettingFilePath = setting_path
	} else {
		SettingFilePath = filepath.Join(ExecutableDirectory, "setting.json")
	}
}
