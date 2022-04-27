package middleware

import (
	"net/http"
	"time"

	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	"github.com/0ne-zero/f4h/utilities/functions/general"
	wrapper_logger "github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		auth := session.Get("authenticated")

		if auth == true {
			data := gin.H{
				"error": "You are NOT authorized. Go to login page",
			}
			c.HTML(http.StatusUnauthorized, "login.html", data)
			c.Abort()
			return
		}
	}
}
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		view_data := gin.H{}
		view_data["Title"] = "Not Found"
		view_data["Error"] = "This Page not Found"
		c.HTML(http.StatusNotFound, "error.html", view_data)
	}
}
func TooManyRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		client_ip := c.ClientIP()
		url := c.Request.URL.Path
		method := c.Request.Method

		yes, err := model_function.TooManyRequest(client_ip, url, method)
		if err != nil {
			err_file_info, err := general.GetCallerInfo(1)
			if err != nil {
				err_file_info, err = general.GetCallerInfo(0)
				wrapper_logger.Log(&wrapper_logger.FatalLevel{}, "Error occurred during get caller info", &err_file_info)
			}
			wrapper_logger.Log(&wrapper_logger.ErrorLevel{}, err, &err_file_info)
		} else if yes == true {
			view_data := gin.H{}
			view_data["Error"] = "Too many request error.Try later"
			c.HTML(http.StatusTooManyRequests, "error.html", view_data)
			c.Abort()
			return
		} else {
			model_function.Add(&model.Request{IP: client_ip, Url: url, Method: method, Time: time.Now().UTC()})
		}
	}
}
