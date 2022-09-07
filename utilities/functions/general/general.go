package general

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	html_to_markdown "github.com/JohannesKaufmann/html-to-markdown"
	markdown_to_html "github.com/gomarkdown/markdown"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/public_struct"
	viewmodel "github.com/0ne-zero/f4h/public_struct/view_model"
	"github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"github.com/gin-gonic/gin"
)

func GenerateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	return bytes, err
}
func GenerateRandomHex(length int) (string, error) {
	var byte_size = length
	if length%2 != 0 {
		byte_size += 1
	}
	bytes, err := GenerateRandomBytes(byte_size / 2)
	if err != nil {
		return "", err
	}
	hex := hex.EncodeToString(bytes)
	hex_len := len(hex)
	for hex_len != length {
		hex = hex[:hex_len-1]
		hex_len = len(hex)
	}
	return hex, nil
}
func DeleteFiles(files_path ...string) error {
	for i := range files_path {
		_, err := os.Stat(files_path[i])
		if os.IsNotExist(err) {
			return fmt.Errorf("\"%s\" file isn't exists", files_path[i])
		}
	}
	return nil
}
func GetListofFilesNameInDirectory(dir string) ([]string, error) {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files_name = make([]string, len(files))
	for i := range files {
		files_name[i] = files[i].Name()
	}
	return files_name, nil
}

// category indicates that where folder should be look
func IsImageExists(f_name, category string) (bool, error) {
	switch category {
	case "PRODUCT":
		files_name, err := GetListofFilesNameInDirectory(filepath.Join(constansts.ImagesDirectory, "product"))
		if err != nil {
			return false, err
		}
		for i := range files_name {
			if files_name[i] == f_name {
				return true, nil
			}
		}
		return false, nil
	case "AVATAR":
		files_name, err := GetListofFilesNameInDirectory(filepath.Join(constansts.ImagesDirectory, "avatar"))
		if err != nil {
			return false, err
		}
		for i := range files_name {
			if files_name[i] == f_name {
				return true, nil
			}
		}
		return false, nil
	default:
		panic("You passed unknown category")
	}
}
func Hashing(input string) (string, error) {
	h := sha256.New()
	h.Write([]byte(input))
	r := h.Sum(nil)
	return string(r), nil
}
func ComparePassword(hashed_pass string, pass string) (bool, error) {
	pass_hash, err := Hashing(pass)
	if err != nil {
		return false, nil
	}
	if hashed_pass == pass_hash {
		return true, nil
	} else {
		return false, nil
	}
}
func IsFloatNumberRound(n float64) bool {
	str_n := fmt.Sprint(n)
	if strings.Contains(str_n, ".") {
		return false
	} else {
		return true
	}
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
	return strings.TrimSuffix(strings.TrimPrefix(s, "/"), "/")
}
func ValueExistsInSlice[T comparable](slice *[]T, value T) bool {
	for _, e := range *slice {
		if e == value {
			return true
		}
	}
	return false
}
func MarkdownToHtml(markdown string) string {
	return string(markdown_to_html.ToHTML([]byte(markdown), nil, nil))

	// if _, err := os.Stat(constansts.MarkdownFilePath); err != nil {
	// 	return "", errors.New(fmt.Sprintf("Please put Markdown file in %s path", constansts.MarkdownFilePath))
	// }
	// command := fmt.Sprintf(`echo "%s" | %s`, markdown, constansts.MarkdownFilePath)
	// topic_markdown_html_bytes, err := exec.Command("bash", "-c", command).Output()
	// return string(topic_markdown_html_bytes), err
}
func HtmlToMarkdown(html string) (string, error) {
	converter := html_to_markdown.NewConverter("", true, nil)
	md, err := converter.ConvertString(html)
	return md, err
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
func GetCallerInfo(skip int) wrapper_logger.ErrorLocation {
	// Remove this function from stack
	skip += 1

	pc, path, line, ok := runtime.Caller(skip)
	if !ok {
		log_msg := "Error occurred during get caller info"
		// Get this function info, if possible
		pc, path, line, ok = runtime.Caller(skip)
		// Log location of error and exit
		if !ok {
			// Fill this function info manually
			err_file_info := wrapper_logger.ErrorLocation{
				FilePath: filepath.Join(constansts.ExecutableDirectory, "utilities/functions/general/general.go"),
				FuncName: "GetCallerInfo",
				Line:     0,
			}
			AppendTextToFile(constansts.LogFilePath, AddFieldsToString(log_msg, err_file_info.ToStringMap()))
			os.Exit(1)
		}
		// Fill this function info with call Caller again
		err_file_info := wrapper_logger.ErrorLocation{FilePath: path, Line: line, FuncName: runtime.FuncForPC(pc).Name()}
		AppendTextToFile(constansts.LogFilePath, AddFieldsToString(log_msg, err_file_info.ToStringMap()))
		os.Exit(1)
	}
	return wrapper_logger.ErrorLocation{FilePath: path, Line: line, FuncName: runtime.FuncForPC(pc).Name()}
}

// Case-insensitive strings.Contains
func ContainsI(s string, sub string) bool {
	return strings.Contains(
		strings.ToLower(s),
		strings.ToLower(sub),
	)
}
func AddFieldsToString(s string, fields map[string]string) string {
	s += " | "
	for k, v := range fields {
		s = fmt.Sprintf("%s %s='%s'", s, k, v)
	}
	return s
}
func GetRequestBasicInfo(c *gin.Context) public_struct.RequestBasicInformation {
	return public_struct.RequestBasicInformation{
		IP:     c.ClientIP(),
		Path:   c.Request.URL.Path,
		Method: c.Request.Method,
	}
}
func ConvertIntSliceToStringSlice(int_slice []int) []string {
	var s_slice []string
	for i := range int_slice {
		s_slice = append(s_slice, fmt.Sprint(int_slice[i]))
	}
	return s_slice
}
func SplitEachTagsWithPipe(tags []viewmodel.TopicTagBasicInformation) string {
	if tags == nil {
		return ""
	}
	var res string
	var tags_index_len = len(tags) - 1
	for i := range tags {
		if i != tags_index_len {
			res += tags[i].Name + " | "
		} else {
			res += tags[i].Name
		}
	}
	return res
}

func ExistsStringInStringSlice(s string, slice []string) bool {
	for i := range slice {
		if s == slice[i] {
			return true
		}
	}
	return false
}
func GetMapKeys[key_type comparable, value_type interface{}](m map[key_type]value_type) []key_type {
	keys := make([]key_type, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Checks is root user that runed the program or not
// It will kill program if program runed on windows system
func IsUserRoot() bool {
	usr_id := os.Getuid()
	if usr_id == -1 {
		fmt.Println("This program can only run in unix-like operating systems like linux and other...")
		os.Exit(1)
		return false
	} else if usr_id == 0 {
		return true
	} else {
		return false
	}
}

func DetectFileType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil

}
