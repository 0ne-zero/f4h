package middleware

import (
	"net/http"

	"github.com/0ne-zero/f4h/database/model_function"
	"github.com/0ne-zero/f4h/utilities/functions/general"
	wrapper_logger "github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"github.com/0ne-zero/f4h/web/controller/controller_helper"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// session := sessions.Default(c)
		// auth := session.Get("authenticated")
		// user_id := session.Get("UserID")
		// if user_id == nil {
		// 	data := gin.H{
		// 		"error": "You are NOT authorized. Go to login page",
		// 	}
		// 	c.HTML(http.StatusUnauthorized, "login.html", data)
		// 	c.Abort()
		// 	return
		// } else if user_id_int, ok := user_id.(int); !ok || user_id_int < 1 {
		// 	// Ban user, they enetered non-int value, how ?; i don't know

		// 	data := gin.H{
		// 		"error": "You are NOT authorized. Go to login page",
		// 	}
		// 	c.HTML(http.StatusUnauthorized, "login.html", data)
		// 	c.Abort()
		// 	return
		// }

		// if auth != true {
		// 	data := gin.H{
		// 		"error": "You are NOT authorized. Go to login page",
		// 	}
		// 	c.HTML(http.StatusUnauthorized, "login.html", data)
		// 	c.Abort()
		// 	return
		// }
	}
}
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Not Found", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general.GetCallerInfo(0)})
		view_data := gin.H{}
		view_data["Title"] = "Not Found"
		view_data["Error"] = "This Page not Found"
		c.HTML(http.StatusNotFound, "error.html", view_data)
		c.Abort()
	}
}
func TooManyRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		client_ip := c.ClientIP()
		url := c.Request.URL.Path
		method := c.Request.Method

		yes, err := model_function.TooManyRequest(client_ip, url, method)
		if err != nil {
			// Log
			wrapper_logger.Error(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general.GetCallerInfo(0)})
		} else if yes == true {
			view_data := gin.H{}
			view_data["Error"] = "Too many request error.Try later"
			c.HTML(http.StatusTooManyRequests, "error.html", view_data)
			c.Abort()
			return
		} else {
			//model_function.Add(&model.Request{IP: client_ip, Url: url, Method: method, Time: time.Now().UTC()})
		}
	}
}
func SetSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		s.Set("UserID", 1)
		s.Set("Username", "admin")
		s.Save()
	}
}
