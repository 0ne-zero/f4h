package controller_helper

import (
	"net/http"
	"time"

	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	"github.com/0ne-zero/f4h/utilities/functions/general"
	wrapper_logger "github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"github.com/gin-gonic/gin"
)

func ErrorPage(c *gin.Context, log_err error, user_msg string) {

	// Add bad request to database
	AddBadRequest(c)
	// Log Error
	wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: log_err.Error(), Fields: ClientInfoInMap(c), ErrorLocation: general.GetCallerInfo(1)})

	// Return response
	var view_data = gin.H{
		"Title": "Error",
	}
	if user_msg != "" {
		view_data["Error"] = user_msg
	} else {
		view_data["Error"] = "Something bad happened. Come back later"
	}
	c.HTML(http.StatusInternalServerError, "error.html", view_data)
}
func AddBadRequest(c *gin.Context) error {
	var bad_request = model.BadRequest{
		IP:     c.ClientIP(),
		Url:    c.Request.URL.Path,
		Method: c.Request.Method,
		Time:   time.Now().UTC(),
	}
	return model_function.Add(&bad_request)
}
func ClientInfoInMap(c *gin.Context) map[string]string {
	return map[string]string{
		"IP":     c.ClientIP(),
		"URL":    c.Request.URL.Path,
		"METHOD": c.Request.Method,
	}
}