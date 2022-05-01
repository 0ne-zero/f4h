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
	log_fields := map[string]string{
		"ip":             c.ClientIP(),
		"request_url":    c.Request.URL.Path,
		"request_method": c.Request.Method,
	}
	log_msg := wrapper_logger.AddFieldsToString(log_err.Error()+" | ", log_fields)
	err_file_info, err := general.GetCallerInfo(1)
	if err != nil {
		// Log and exit
		err_file_info, err = general.GetCallerInfo(0)
		wrapper_logger.Log(&wrapper_logger.FatalLevel{}, "Error occurred during get caller info", &err_file_info)
	}
	wrapper_logger.Log(&wrapper_logger.WarningLevel{}, log_msg, &err_file_info)

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
