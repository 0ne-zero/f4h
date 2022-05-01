package general

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/public_struct"
	"github.com/0ne-zero/f4h/utilities/functions/setting"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	// Generate bcrypt hash from password with 17 cost
	// Get hash cost number from settings file
	hash_cost_number_string, err := setting.ReadFieldInSettingData("HASH_COST_NUMBER")
	if err != nil {
		return "", err
	}
	// convert hash_cost_number_string to int
	hash_cost_number, err := strconv.ParseInt(hash_cost_number_string, 10, 64)
	if err != nil {
		return "", err
	}
	hash_bytes, err := bcrypt.GenerateFromPassword([]byte(pass), int(hash_cost_number))
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
	if _, err := os.Stat(constansts.MarkdownFilePath); err != nil {
		return "", errors.New(fmt.Sprintf("Please put Markdown file in %s path", constansts.MarkdownFilePath))
	}
	command := fmt.Sprintf(`echo "%s" | %s`, markdown, constansts.MarkdownFilePath)
	topic_markdown_html_bytes, err := exec.Command("bash", "-c", command).Output()
	return string(topic_markdown_html_bytes), err
}

func AppendTextToFile(path string, text string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(text))
	if err != nil {
		return err
	}
	return nil
}

func GetCallerInfo(skip int) (public_struct.ErroredFileInfo, error) {
	// fix file infos
	if constansts.GetCallerInfoError >= 10 {
		log_msg := fmt.Sprintf("%s file='%s:%d'", "More than 9 Errors occurred during get caller info ", "/utilities/functions/general/general.GetCallerInfo", 88)
		AppendTextToFile(constansts.LogFilePath, log_msg)
		os.Exit(1)
	}
	// Remove this function from stack
	skip += 1

	_, path, line, ok := runtime.Caller(skip)
	if !ok {
		constansts.GetCallerInfoError += 1
		return public_struct.ErroredFileInfo{}, errors.New("Error occurred during get caller information")
	}
	if constansts.GetCallerInfoError > 0 {
		constansts.GetCallerInfoError = 0
	}

	return public_struct.ErroredFileInfo{Path: path, Line: line}, nil
}

// Case-insensitive strings.Contains
func ContainsI(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}
