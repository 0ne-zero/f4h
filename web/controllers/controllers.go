package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	"github.com/0ne-zero/f4h/utilities"
	"github.com/0ne-zero/f4h/utilities/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Main Directory (direcotry of executable file)
var MainDirectory = filepath.Dir(os.Args[0])

func Login_GET(c *gin.Context) {
	view_data := make(map[string]interface{})
	c.HTML(http.StatusOK, "login.html", view_data)
}
func Login_POST(c *gin.Context) {
	// Get username & password from form
	username := c.PostForm("username")
	password := c.PostForm("password")
	// Validatae username & password
	if username == "" || password == "" {
		data := gin.H{
			"LoginError": "Fill all field",
		}
		c.HTML(http.StatusOK, "login.html", data)
		return
	}
	// Get user
	user, err := model_function.GetUserByUsername(username)
	if err != nil {
		data := gin.H{
			"LoginError": "Username or Password is incorrect.",
		}
		c.HTML(http.StatusOK, "login.html", data)
		return
	}
	// Compare user password with entered password
	if err := utilities.ComparePassword(user.PasswordHash, password); err != nil {
		data := gin.H{
			"LoginError": "Username or Password is incorrect.",
		}
		c.HTML(http.StatusOK, "login.html", data)
		return
	}

	// Set Cookie
	session := sessions.Default(c)
	session.Set("authenticated", true)
	session.Save()

	data := gin.H{
		"Title":    "Index",
		"Username": user.Username,
	}
	c.HTML(http.StatusOK, "index.html", data)
}
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("authenticated", false)
	data := gin.H{
		"Title": "Index",
	}
	c.HTML(http.StatusOK, "index.html", data)
}
func Register_POST(c *gin.Context) {
	// Get username & password from form
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	// Validatae username & password
	if username == "" || password == "" {
		view_data := make(map[string]interface{})
		view_data["RegisterError"] = "Fill all field"
		c.HTML(http.StatusOK, "login.html", view_data)
		return
	}
	exists, err := model_function.IsUserExistsByUsername(username)
	if err != nil {
		view_data := make(map[string]interface{})
		view_data["RegisterError"] = "Unknown error"
		c.HTML(http.StatusOK, "login.html", view_data)
		return
	}
	if exists == true {
		view_data := make(map[string]interface{})
		view_data["RegisterError"] = "User is exists"
		c.HTML(http.StatusOK, "login.html", view_data)
		return
	}
	// So far user not exists, i should register client
	// Create password hash
	pass_hash, err := utilities.HashPassword(password)
	if err != nil {
		//Log
		view_data := make(map[string]interface{})
		view_data["RegisterError"] = "Unknown error"
		c.HTML(http.StatusOK, "login.html", view_data)
		return
	}
	// Add User
	model_function.Add(&model.User{Username: username, PasswordHash: pass_hash, Email: email})
	// Response
	view_data := make(map[string]interface{})
	view_data["Username"] = username
	c.HTML(http.StatusOK, "index.html", view_data)
}
func Index(c *gin.Context) {
	products := []model.Product{}
	// Get products
	err := model_function.Get(&products, 15, "created_at", "desc")
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)

	}
	categories := []model.Product_Category{}
	// Get category
	err = model_function.Get(&categories, -1, "created_at", "desc", "SubCategories")
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)
		return
	}
	// Everything is ok
	view_data := make(map[string]interface{})
	view_data["Title"] = "Index"
	if categories != nil {
		view_data["Categories"] = categories
	}
	if products != nil {
		view_data["Products"] = products
	}
	c.HTML(http.StatusOK, "index.html", view_data)
}
func ProductList(c *gin.Context) {
	enteredCategory := c.Param("category")
	// Specified category
	if enteredCategory != "" && enteredCategory != "/" {
		// Remove slashes and Make Title from input category
		enteredCategory = utilities.RemoveSlashFromBeginAndEnd(enteredCategory)
		enteredCategory = strings.Title(enteredCategory)

		// Get all categories name
		var categoriesName []string
		err := model_function.GetCategoriesName(&categoriesName)
		if err != nil {
			log.Log(logrus.Error, err)
			view_data := make(map[string]interface{})
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
			return
		}
		// Check entered category is exists in database
		if !utilities.ValueExistsInSlice(&categoriesName, strings.Title(enteredCategory)) {
			log.Log(logrus.Error, errors.New("Invalid Category"))
			view_data := make(map[string]interface{})
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
			return
		}
		// Get Products by entered category
		var products []model.Product
		err = model_function.GetProductsByCategory(&products, strings.ToLower(enteredCategory))
		if err != nil {
			log.Log(logrus.Error, err)
			view_data := make(map[string]interface{})
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
			return
		} else if products == nil {

		}
		// Get all categories for sidebar
		var categories []model.Product_Category
		err = model_function.GetCategoryByOrderingProductsCount(&categories)
		if err != nil {
			log.Log(logrus.Error, err)
			view_data := make(map[string]interface{})
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
			return
		}
		// Everything is ok
		view_data := make(map[string]interface{})
		view_data["Title"] = fmt.Sprintf("%s Products List", strings.ToTitle(enteredCategory))
		if categories != nil {
			view_data["Categories"] = categories
		}
		if products != nil {
			view_data["Products"] = products
		}
		c.HTML(http.StatusOK, "product-list.html", view_data)
	} else {
		// Unspecified category

		var products []model.Product
		err := model_function.Get(&products, -1, "created_at", "desc", "Images")
		if err != nil {
			log.Log(logrus.Error, err)
			view_data := make(map[string]interface{})
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
		}
		var categories []model.Product_Category
		err = model_function.GetCategoryByOrderingProductsCount(&categories)
		if err != nil {
			log.Log(logrus.Error, err)
			view_data := make(map[string]interface{})
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
		}
		// Everything is ok
		view_data := make(map[string]interface{})
		view_data["Title"] = "Products List"
		if categories != nil {
			view_data["Categories"] = categories
		}
		if products != nil {
			view_data["Products"] = products
		}
		c.HTML(http.StatusOK, "product-list.html", view_data)
	}
}
func ProductDetails(c *gin.Context) {
	// String product id
	str_id := c.Param("id")
	// check str_id is Empty
	if str_id == "" {
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Go and select product again."
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	// Integer product id
	int_id, err := strconv.ParseInt(str_id, 10, 64)
	// Parse check
	if err != nil {
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Go and select product again."
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	var product model.Product
	err = model_function.GetByID(&product, int(int_id))
	// Check get product error
	if err != nil {
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Go and select product again. If this page comes up again, there must be a problem, So come later ):"
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	// Everything is ok
	view_data := make(map[string]interface{})
	view_data["Title"] = product.Name
	view_data["Product"] = product
	c.HTML(http.StatusOK, "product-details.html", view_data)
}

// func About()

// func TagProducts()

// func CategoryProducts()

// func Categories()

// func Dashboard()
