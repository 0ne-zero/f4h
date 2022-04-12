package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	"github.com/0ne-zero/f4h/utilities"
	"github.com/0ne-zero/f4h/utilities/log"
	viewmodel "github.com/0ne-zero/f4h/web/view_model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
	user_pass_hash, err := model_function.GetUserPassHashByUsername(username)
	if err != nil {
		data := gin.H{
			"LoginError": "Username or Password is incorrect.",
		}
		c.HTML(http.StatusOK, "login.html", data)
		return
	}
	// Compare user password with entered password
	if err := utilities.ComparePassword(user_pass_hash, password); err != nil {
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
		"Username": username,
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
	// Get products
	products, err := model_function.GetProductInViewModel(15)
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)
	}
	// Get categories
	categories, err := model_function.GetCategoriesWithRelationsInViewModel(true)
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

		// Get all categories name
		categoriesName, err := model_function.GetCategoriesName()
		if err != nil {
			log.Log(logrus.Error, err)
			view_data := make(map[string]interface{})
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
			return
		}
		// Check entered category is exists in database
		if !utilities.ValueExistsInSlice(&categoriesName, enteredCategory) {
			log.Log(logrus.Error, errors.New("Invalid Category"))
			view_data := make(map[string]interface{})
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
			return
		}
		// Get Products by entered category
		products, err := model_function.GetProductByCategoryInViewModel(strings.ToLower(enteredCategory), 15)
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
		categories, err := model_function.GetCategoriesWithRelationsInViewModel(true)
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

//region Forum
func Discussions(c *gin.Context) {
	// Discussions categories and Discussion (preload Discussion)
	var Discussion_categories []model.Discussion_Category
	err := model_function.Get(&Discussion_categories, -1, "created_at", "ASC", "Discussions")
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)
		return
	}

	// View model
	// Main view model
	discussion_categories_view_model := []viewmodel.DiscussionCategoryViewModel{}

	// For on Discussion categories for create a view data to send template
	for _, d_c := range Discussion_categories {
		// Data that should append to view model
		// Create discussion category view model
		discussion_category := viewmodel.DiscussionCategoryViewModel{
			ID: d_c.ID, Name: d_c.Name, Description: d_c.Description,
			CreatedAt: d_c.CreatedAt,
			UpdatedAt: d_c.UpdatedAt, UserID: d_c.UserID,
		}
		// For on Discussions
		for _, d := range d_c.Discussions {
			// Get discussion posts count
			post_count, err := model_function.GetDiscussionPostsCount(d)
			if err != nil {
				log.Log(logrus.Error, err)
				view_data := make(map[string]interface{})
				view_data["Title"] = "Error"
				view_data["Error"] = "Something bad happened. Come back later"
				c.HTML(http.StatusInternalServerError, "error.html", view_data)
				return
			}
			topic_count, err := model_function.GetDiscussionTopicsCount(d)
			if err != nil {
				log.Log(logrus.Error, err)
				view_data := make(map[string]interface{})
				view_data["Title"] = "Error"
				view_data["Error"] = "Something bad happened. Come back later"
				c.HTML(http.StatusInternalServerError, "error.html", view_data)
				return
			}
			// Get discussion forums name
			forums_name, err := model_function.GetDiscussionForumsName(d)
			if err != nil {
				log.Log(logrus.Error, err)
				view_data := make(map[string]interface{})
				view_data["Title"] = "Error"
				view_data["Error"] = "Something bad happened. Come back later"
				c.HTML(http.StatusInternalServerError, "error.html", view_data)
				return
			}
			// Create discussion view model
			discussion := viewmodel.DiscussionViewModel{
				Discussion: *d,
				// Get discussion topics count
				TopicCount: topic_count,
				PostCount:  post_count,
				ForumsName: forums_name,
			}
			// Append discussion view model to main view model
			discussion_category.Discussions = append(discussion_category.Discussions, &discussion)
		}
		// Append discussion_category to main view model
		discussion_categories_view_model = append(discussion_categories_view_model, discussion_category)
	}
	view_data := make(map[string]interface{})
	view_data["Title"] = "Discussions"
	view_data["DiscussionsCategories"] = discussion_categories_view_model
	c.HTML(http.StatusOK, "discussions.html", view_data)
}

func DiscussionForums(c *gin.Context) {
	discussion_name := c.Param("discussion")
	if discussion_name == "" {
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Enter a discussion name"
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	discussion_id, err := model_function.GetModelIDByFieldValue(&model.Discussion{}, "name", discussion_name)
	if err != nil {
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = fmt.Sprintf("%s discussion doesn't exists. Please enter a valid discussion name", discussion_name)
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	discussion_forums, err := model_function.GetDiscussionForumsInViewModel(discussion_id)
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)
		return
	}
	discussion_topics, err := model_function.GetDiscussionTopics(discussion_id)
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)
		return
	}

	view_data := make(map[string]interface{})
	var _topics_list_template_view_model = map[string]interface{}{
		"Topics":         discussion_topics,
		"DiscussionName": discussion_name,
	}
	view_data["Topics"] = _topics_list_template_view_model
	view_data["Forums"] = discussion_forums
	c.HTML(http.StatusOK, "discussion_forums.html", view_data)
}

func ForumTopics(c *gin.Context) {
	forum_name := c.Param("forum")
	if forum_name == "" {
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = "Enter a discussion name"
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	forum_id, err := model_function.GetModelIDByFieldValue(&model.Forum{}, "name", forum_name)
	if err != nil {
		view_data := make(map[string]interface{})
		view_data["Title"] = "Error"
		view_data["Error"] = fmt.Sprintf("%s forum doesn't exists. Please enter a valid discussion name", forum_name)
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	forum_topics, err := model_function.GetForumTopicsInViewModel(forum_id)

	view_data := make(map[string]interface{})
	var _topics_list_template_view_model = map[string]interface{}{
		"Topics":    forum_topics,
		"ForumName": forum_name,
	}
	view_data["Title"] = fmt.Sprintf("%s Topics", forum_name)
	view_data["Topics"] = _topics_list_template_view_model
	c.HTML(http.StatusBadRequest, "forum_topics.html", view_data)

}

//endregion

// func About()

// func TagProducts()

// func CategoryProducts()

// func Categories()

// func Dashboard()
