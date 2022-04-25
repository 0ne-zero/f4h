package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/0ne-zero/f4h/config/constansts"
	"golang.org/x/crypto/bcrypt"
)

func ReadSettingFile(setting_path string) (map[string]interface{}, error) {
	var data map[string]interface{}
	file, err := ioutil.ReadFile(setting_path)
	if err != nil {
		return nil, errors.New("Error when opening setting file")
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, errors.New("Error during unmarshal setting file")
	}
	return data, nil
}
func ReadFieldInSettingFile(setting_path string, field_name string) (string, error) {
	settings, err := ReadSettingFile(setting_path)
	if err != nil {
		return "", err
	}

	field_value, exists := settings[field_name]
	if !exists {
		return "", errors.New("Field not exists in setting file")
	}

	// Convert type interface to string and return value of field
	return fmt.Sprint(field_value), nil

}
func HashPassword(pass string) (string, error) {
	// Generate bcrypt hash from password with 17 cost
	hash_bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 17)
	if err != nil {
		return "", err
	}
	return string(hash_bytes), nil
}
func ComparePassword(hashed_pass string, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed_pass), []byte(pass))
}

// Remove slice element and keeping the order
// Just shift
func RemoveSliceElement[T any](slice []T, i int) []T {
	// If index is out of range remove last element
	if i >= len(slice) {
		return slice[:len(slice)-1]

	} else {
		return append(slice[:i], slice[i+1:]...)
	}
}
func RemoveSlashFromBeginAndEnd(s string) string {
	if strings.HasPrefix(s, "/") {
		s = s[1:]
	}
	if strings.HasSuffix(s, "/") {
		s = s[:len(s)-1]
	}
	return s
}
func ValueExistsInSlice[T comparable](slice *[]T, value T) bool {
	for _, e := range *slice {
		if e == value {
			return true
		}
	}
	return false
}
func MarkdownToHtml(markdown string) (string, error) {
	command := fmt.Sprintf(`echo "%s" | %s`, markdown, constansts.MarkdownFilePath)
	topic_markdown_html_bytes, err := exec.Command("bash", "-c", command).Output()
	return string(topic_markdown_html_bytes), err
}
