package constansts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// App name
var AppName string

// Paths
var ExecutableDirectory string
var SettingFilePath string
var UtilitiesDirectory string
var MarkdownFilePath string
var LogFilePath string
var ImagesDirectory string
var DefaultAvatarPath string

// DSN
var DSN string

// XSS preventation things
var XSSPreventor *bluemonday.Policy

// All gin routes that created
// This field is filled by "github.com/0ne-zero/f4h/web/route.MakeRoute" function
var Routes gin.RoutesInfo

// Setting file data
var SettingData map[string]string

// Number of get caller info (/utilities/functions/general.GetCallerInfo) error
var GetCallerInfoErrorCount int

// Errors
var SomethingBadHappenedError string

// File name length
var FileNameLength int

// Is already loaded
var loaded bool

func init() {
	if !loaded {
		// Errors
		SomethingBadHappenedError = "Something Bad Happened, Please back later."
		// Paths
		ExecutableDirectory = filepath.Dir(os.Args[0])
		if setting_path := os.Getenv("F4H_SETTING_PATH"); setting_path != "" {
			SettingFilePath = setting_path
		} else {
			SettingFilePath = filepath.Join(ExecutableDirectory, "config", "setting.json")
		}
		UtilitiesDirectory = filepath.Join(ExecutableDirectory, "utilities")
		MarkdownFilePath = filepath.Join(UtilitiesDirectory, "Markdown.pl")

		// Load Setting File
		var err error
		SettingData, err = readSettingFile(SettingFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		AppName = SettingData["APP_NAME"]
		LogFilePath = filepath.Join(SettingData["LOG_FILE_PARENT_DIRECTORY"], SettingData["APP_NAME"], "log.txt")
		ImagesDirectory = SettingData["IMAGES_DIRECTORY"]
		file_name_length_int, err := strconv.Atoi(SettingData["FILE_NAME_LENGTH"])
		if err != nil {
			fmt.Println("FILE_NAME_LENGTH field is non-int in setting file\n%s", err.Error())
			os.Exit(1)
		}
		FileNameLength = file_name_length_int

		// Initial XSS preventation
		XSSPreventor = bluemonday.UGCPolicy()
		// Set loaded to true, so next time it won't load again and will use them immediately
		loaded = true
	}
}

func readSettingFile(setting_path string) (map[string]string, error) {
	file_bytes, err := ioutil.ReadFile(setting_path)
	if err != nil {
		return nil, errors.New("Error when opening setting file")
	}
	var data map[string]string
	err = json.Unmarshal(file_bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("error occurred during unmarshal setting file")
	}
	err = validateSettingData(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func validateSettingData(data map[string]string) error {
	var expect_fields_name = []string{
		"APP_NAME",
		"LOG_FILE_PARENT_DIRECTORY",
		"DSN",
		"FILE_NAME_LENGTH",
		"CONTACT_EMAIL",
		"IMAGES_DIRECTORY",
		"DEFAULT_AVATAR_PATH",
	}

	var exists bool
	var data_value string
	for _, data_name := range expect_fields_name {
		data_value, exists = data[data_name]
		if !exists {
			return fmt.Errorf("%s doesn't exists in setting file", data_name)
		}
		if data_value == "" {
			return fmt.Errorf("%s is empty in setting file", data_name)
		}

	}
	return nil
}
