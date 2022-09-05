package templatefunction

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	viewmodel "github.com/0ne-zero/f4h/public_struct/view_model"
	"github.com/gin-gonic/gin"
)

func remainder(a, b int) int {
	return a % b
}
func sliceLength[T any](slice []T) int {
	return len(slice)
}
func plus(variable int, operand int) int {
	return variable + operand
}
func minus(variable int, operand int) int {
	return variable - operand
}
func getDayMonthYearFromTime(t time.Time) string {
	return t.Format("2006-01-02 15:04 UTC")
}
func toString(i interface{}) string {
	return fmt.Sprint(i)
}

// Get Field of selected element of image view data slice
func getFieldOfIndexOfImageViewData(slice []viewmodel.ImageViewData, index int, field_name string) any {
	e := slice[index]
	switch field_name {
	case "Path":
		return e.Path
	case "Name":
		return e.Name
	default:
		panic("You selected unknown field of ImageViewData")
	}
}
func iterate(end int) []uint {
	var list []uint
	for i := 0; i < end; i++ {
		list = append(list, uint(i))
	}
	return list
}
func titlelizeEachWordFirstLetter(s string) string {
	words := strings.Split(s, " ")
	words_len := len(words)
	var result string
	for i := range words {
		if i == words_len {
			result += strings.Title(words[i])
		} else {
			result += strings.Title(words[i]) + " "
		}
	}
	return result
}
func replaceString(s, o_char, n_char string) string {
	return strings.Replace(s, o_char, n_char, -1)
}
func AddFunctionsToRoute(r *gin.Engine) {
	r.SetFuncMap(
		template.FuncMap{
			"iterate":                 iterate,
			"remainder":               remainder,
			"stringSliceLength":       sliceLength[string],
			"TopicTagsSliceLength":    sliceLength[viewmodel.TopicTagBasicInformation],
			"imagesLength":            sliceLength[viewmodel.ImageViewData],
			"plus":                    plus,
			"minus":                   minus,
			"formatTime":              getDayMonthYearFromTime,
			"titlelizeEachWord":       titlelizeEachWordFirstLetter,
			"replace":                 replaceString,
			"toString":                toString,
			"getFieldOfImageViewData": getFieldOfIndexOfImageViewData,
		},
	)
}
