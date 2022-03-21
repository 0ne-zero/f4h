package middleware

import (
	"net/http"
	"time"

	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	"github.com/0ne-zero/f4h/utilities/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func TooManyRequest() gin.HandlerFunc {
	return func(c *gin.Context) {

		client_ip := c.ClientIP()
		url := c.Request.URL.Path
		method := c.Request.Method

		yes, err := model_function.TooManyRequest(client_ip, url, method)
		if err != nil {
			log.Log(logrus.Error, err)
		} else if yes == true {
			view_data := make(map[string]interface{})
			view_data["Error"] = "Too many request error.Try later"
			c.HTML(http.StatusTooManyRequests, "error.html", view_data)
			c.Abort()
			return
		} else {
			model_function.Add(&model.Request{IP: client_ip, Url: url, Method: method, Time: time.Now().UTC()})
		}
	}

}
