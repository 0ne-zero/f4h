package controller_helper

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetTitle(t string) string {
	return constansts.AppName + fmt.Sprintf(" | %s", t)
}

// If returned value be nil means there is no error(data are validated)
func AddProductValidation(name, price, inventory, description, tags string) gin.H {
	// Emptiness check
	if name == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Price":       price,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the title field"},
		}
		return view_data
	}
	if price == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the price field"},
		}
		return view_data
	}
	if inventory == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the inventory field"},
		}
		return view_data
	}
	if description == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Price":     price,
				"Inventory": inventory,
				"Name":      name,
				"Tags":      tags,
			},
			"NotesData": gin.H{"Error": "Please fill the description"},
		}
		return view_data
	}
	if tags == "" {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Inventory":   inventory,
				"Description": description,
			},
			"NotesData": gin.H{"Error": "Please fill the tags field"},
		}
		return view_data
	}
	// Convert to actual type
	_, err := strconv.ParseFloat(price, 64)
	if err != nil {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the price field as number (in Dollor)"},
		}
		return view_data
	}
	_, err = strconv.Atoi(inventory)
	if err != nil {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Please fill the inventory field as number"},
		}
		return view_data
	}
	// Check length
	if len(name) < 5 {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Name field is too short, It should be minimum 5 character"},
		}
		return view_data
	}
	if len(description) < 50 {
		view_data := gin.H{
			"ProductInfo": gin.H{
				"Name":        name,
				"Price":       price,
				"Inventory":   inventory,
				"Description": description,
				"Tags":        tags,
			},
			"NotesData": gin.H{"Error": "Description field is too short, It should be minimum 50 character"},
		}
		return view_data
	}

	if strings.Contains(tags, "|") {
		splitted_tags := strings.Split(tags, "|")
		splitted_tags_len := len(splitted_tags)
		if splitted_tags_len > 5 {
			view_data := gin.H{
				"ProductInfo": gin.H{
					"Name":        name,
					"Price":       price,
					"Inventory":   inventory,
					"Description": description,
					"Tags":        tags,
				},
				"NotesData": gin.H{"Error": fmt.Sprintf("Maximum number of tag is 5, But you inserted %d tags", splitted_tags_len)},
			}
			return view_data
		}
	}
	return nil
}
func EditProductValidation(id, name, price, inventory, description, tags string) gin.H {
	return AddProductValidation(name, price, inventory, description, tags)
}
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
