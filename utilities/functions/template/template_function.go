package templatefunction

import (
	"html/template"
	"time"

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
func AddFunctionsToRoute(r *gin.Engine) {
	r.SetFuncMap(
		template.FuncMap{
			"remainder":         remainder,
			"stringSliceLength": sliceLength[string],
			"plus":              plus,
			"minus":             minus,
			"formatTime":        getDayMonthYearFromTime,
		},
	)
}
