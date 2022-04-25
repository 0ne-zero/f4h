package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/0ne-zero/f4h/config/constansts"
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
	view_data := gin.H{
		"Title": "Login/Signup",
	}
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
	var user_fields = model.User{}
	err := model_function.GetFieldsByAnotherFieldValue(&user_fields, []string{"password_hash"}, "username", username)
	if err != nil {
		data := gin.H{
			"LoginError": "Username or Password is incorrect.",
		}
		c.HTML(http.StatusOK, "login.html", data)
		return
	}
	// Compare user password with entered password
	if err := utilities.ComparePassword(user_fields.PasswordHash, password); err != nil {
		data := gin.H{
			"LoginError": "Username or Password is incorrect.",
		}
		c.HTML(http.StatusOK, "login.html", data)
		return
	}

	// Create session for user
	session := sessions.Default(c)
	session.Set("Username", username)
	// Get user ID by username
	err = model_function.GetFieldsByAnotherFieldValue(&user_fields, []string{"id"}, "username", username)
	if err != nil {
		return
	}
	session.Set("UserID", user_fields.ID)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
func Logout(c *gin.Context) {
	// Get session
	session := sessions.Default(c)
	// Remove session in server-side
	session.Clear()
	// Remove session in client-side
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	// Remove
	session.Save()
	// Redirect to home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
func Register_POST(c *gin.Context) {
	// Get username & password from form
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	// Validatae username & password
	if username == "" || password == "" {
		view_data := gin.H{}
		view_data["RegisterError"] = "Fill all field"
		c.HTML(http.StatusOK, "login.html", view_data)
		return
	}
	exists, err := model_function.IsUserExistsByUsername(username)
	if err != nil {
		view_data := gin.H{}
		view_data["RegisterError"] = "Unknown error"
		c.HTML(http.StatusOK, "login.html", view_data)
		return
	}
	if exists == true {
		view_data := gin.H{}
		view_data["RegisterError"] = "User is exists"
		c.HTML(http.StatusOK, "login.html", view_data)
		return
	}
	// So far user not exists, i should register client
	// Create password hash
	pass_hash, err := utilities.HashPassword(password)
	if err != nil {
		//Log
		view_data := gin.H{}
		view_data["RegisterError"] = "Unknown error"
		c.HTML(http.StatusOK, "login.html", view_data)
		return
	}
	// Add User
	model_function.Add(&model.User{Username: username, PasswordHash: pass_hash, Email: email})
	// Response
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
func Index(c *gin.Context) {
	// Get products
	products, err := model_function.GetProductInViewModel(15)
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)
	}
	// Get categories
	categories, err := model_function.GetCategoriesWithRelationsInViewModel(true)
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)
		return
	}
	// Everything is ok
	view_data := gin.H{}
	view_data["Title"] = "Index"
	view_data["Username"] = sessions.Default(c).Get("Username").(string)
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
			view_data := gin.H{}
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
			return
		}
		// Check entered category is exists in database
		if !utilities.ValueExistsInSlice(&categoriesName, enteredCategory) {
			log.Log(logrus.Error, errors.New("Invalid Category"))
			view_data := gin.H{}
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
			return
		}
		// Get Products by entered category
		products, err := model_function.GetProductByCategoryInViewModel(strings.ToLower(enteredCategory), 15)
		if err != nil {
			log.Log(logrus.Error, err)
			view_data := gin.H{}
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
			view_data := gin.H{}
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
			return
		}
		// Everything is ok
		view_data := gin.H{}
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
			view_data := gin.H{}
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
		}
		var categories []model.Product_Category
		err = model_function.GetCategoryByOrderingProductsCount(&categories)
		if err != nil {
			log.Log(logrus.Error, err)
			view_data := gin.H{}
			view_data["Title"] = "Error"
			view_data["Error"] = "Something bad happened. Come back later"
			c.HTML(http.StatusInternalServerError, "error.html", view_data)
		}
		// Everything is ok
		view_data := gin.H{}
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
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = "Go and select product again."
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	// Integer product id
	int_id, err := strconv.ParseInt(str_id, 10, 64)
	// Parse check
	if err != nil {
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = "Go and select product again."
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	var product model.Product
	err = model_function.GetByID(&product, int(int_id))
	// Check get product error
	if err != nil {
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = "Go and select product again. If this page comes up again, there must be a problem, So come later ):"
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	// Everything is ok
	view_data := gin.H{}
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
		view_data := gin.H{}
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
				view_data := gin.H{}
				view_data["Title"] = "Error"
				view_data["Error"] = "Something bad happened. Come back later"
				c.HTML(http.StatusInternalServerError, "error.html", view_data)
				return
			}
			topic_count, err := model_function.GetDiscussionTopicsCount(d)
			if err != nil {
				log.Log(logrus.Error, err)
				view_data := gin.H{}
				view_data["Title"] = "Error"
				view_data["Error"] = "Something bad happened. Come back later"
				c.HTML(http.StatusInternalServerError, "error.html", view_data)
				return
			}
			// Get discussion forums name
			forums_name, err := model_function.GetDiscussionForumsName(d)
			if err != nil {
				log.Log(logrus.Error, err)
				view_data := gin.H{}
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
	view_data := gin.H{}
	view_data["Title"] = "Discussions"
	view_data["DiscussionsCategories"] = discussion_categories_view_model
	c.HTML(http.StatusOK, "discussions.html", view_data)
}

func DiscussionForums(c *gin.Context) {
	discussion_name := c.Param("discussion")
	if discussion_name == "" {
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = "Enter a discussion name"
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	var discussion_fields = model.Discussion{}
	err := model_function.GetFieldsByAnotherFieldValue(&discussion_fields, []string{"name"}, "name", discussion_name)
	if err != nil {
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = fmt.Sprintf("%s discussion doesn't exists. Please enter a valid discussion name", discussion_name)
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	discussion_forums, err := model_function.GetDiscussionForumsInViewModel(int(discussion_fields.ID))
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)
		return
	}
	discussion_topics, err := model_function.GetDiscussionTopics(int(discussion_fields.ID))
	if err != nil {
		log.Log(logrus.Error, err)
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = "Something bad happened. Come back later"
		c.HTML(http.StatusInternalServerError, "error.html", view_data)
		return
	}

	view_data := gin.H{}
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
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = "Enter a discussion name"
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	var forum_fields = model.Forum{}
	err := model_function.GetFieldsByAnotherFieldValue(&forum_fields, []string{"id"}, "name", forum_name)
	if err != nil {
		view_data := gin.H{}
		view_data["Title"] = "Error"
		view_data["Error"] = fmt.Sprintf("%s forum doesn't exists. Please enter a valid discussion name", forum_name)
		c.HTML(http.StatusBadRequest, "error.html", view_data)
		return
	}
	forum_topics, err := model_function.GetForumTopicsInViewModel(int(forum_fields.ID))

	view_data := gin.H{}
	var _topics_list_template_view_model = map[string]interface{}{
		"Topics":    forum_topics,
		"ForumName": forum_name,
	}
	view_data["Title"] = fmt.Sprintf("%s Topics", forum_name)
	view_data["Topics"] = _topics_list_template_view_model
	c.HTML(http.StatusBadRequest, "forum_topics.html", view_data)
}

func AddTopic_GET(c *gin.Context) {
	view_data := gin.H{
		"WriteTopicData": gin.H{
			"ForumName": c.Param("forum"),
		},
	}
	c.HTML(200, "add_topic.html", view_data)
}
func AddTopic_POST(c *gin.Context) {
	// User topic Markdown
	topic_markdown := c.Request.FormValue("topic-markdown")
	if topic_markdown == "" {
		// return error
	}
	topic_markdown = strings.TrimSpace(topic_markdown)

	// Convert Markdown to Html
	topic_html, err := utilities.MarkdownToHtml(topic_markdown)
	if err != nil {
		return
	}
	// Remove space from start and end of topic_html
	topic_html = strings.TrimSpace(topic_html)
	// Prevent XSS Attacks
	topic_html = constansts.XSS_Preventor.Sanitize(topic_html)

	var view_data = gin.H{}

	// User wants preview of her/his topic markdown
	if is_preview := c.Request.FormValue("preview"); is_preview != "" {
		view_data = gin.H{
			"WriteTopicData": gin.H{
				"ForumName":     c.Param("forum"),
				"Preview":       template.HTML(topic_html),
				"TopicMarkdown": topic_markdown,
			},
			"Title": "Post new topic",
		}

		c.HTML(200, "add_topic.html", view_data)
		return
		// User wants save her/his topic markdown
	} else if is_save := c.Request.FormValue("save"); is_save != "" {
		// User wants submit her/his text
	} else if is_submit := c.Request.FormValue("submit"); is_submit != "" {
		forum_name := c.Param("forum")
		if forum_name == "" {
			return
		}
		var forum_fields = model.Forum{}
		err := model_function.GetFieldsByAnotherFieldValue(&forum_fields, []string{"id"}, "name", forum_name)
		if err != nil {
			return
		}
		subject := c.Request.FormValue("subject")
		if subject == "" {
			return
		}

		var topic model.Topic
		topic.Name = subject
		topic.Description = topic_html
		topic.ForumID = uint(forum_fields.ID)
		//topic.Tags =

		untyped_user_id := sessions.Default(c).Get("UserID")
		if untyped_user_id == nil {
			return
		}
		user_id, ok := untyped_user_id.(uint)
		if !ok {
			return
		}

		topic.UserID = user_id

	}
}

//endregion

// func About()

// func TagProducts()

// func CategoryProducts()

// func Categories()

// func Dashboard()
