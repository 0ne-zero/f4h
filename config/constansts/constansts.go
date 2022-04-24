package constansts

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// Paths
var ExecutableDirectory string
var SettingFilePath string
var UtilitiesDirectory string
var MarkdownFilePath string

// XSS preventation things
var XSS_Preventor *bluemonday.Policy

// All gin routes that created
// This field is filled by "github.com/0ne-zero/f4h/web/route.MakeRoute" function
var Routes gin.RoutesInfo

func init() {
	// Paths
	ExecutableDirectory = filepath.Dir(os.Args[0])
	if setting_path := os.Getenv("F4H_SETTING_PATH"); setting_path != "" {
		SettingFilePath = setting_path
	} else {
		SettingFilePath = filepath.Join(ExecutableDirectory, "setting.json")
	}
	UtilitiesDirectory = filepath.Join(ExecutableDirectory, "utilities")
	MarkdownFilePath = filepath.Join(UtilitiesDirectory, "Markdown.pl")

	// XSS preventation
	XSS_Preventor = bluemonday.UGCPolicy()
}
