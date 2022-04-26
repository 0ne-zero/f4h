package controllerhelper

import (
	"net/http"

	wrapper_logger "github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SomethingBadHappened(c *gin.Context, log_err error, user_msg string) {
	wrapper_logger.Log(logrus.Error, log_err)
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
