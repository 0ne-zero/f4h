package controller_helper

import (
	"net/http"
	"time"

	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ErrorPage(c *gin.Context, user_msg string) {
	// Add bad request to database
	AddBadRequest(c)

	// Return response
	var view_data = gin.H{}
	view_data["Title"] = "Error"
	view_data["Error"] = user_msg
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
func DeleteUserSession(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
}
func DeleteUserTopicDraftFromSession(c *gin.Context) {
	s := sessions.Default(c)
	s.Delete("TopicSubject")
	s.Delete("TopicMarkdown")
	s.Delete("TopicTags")
	s.Save()
}
