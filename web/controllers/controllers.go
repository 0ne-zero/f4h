package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	contoller_helper "github.com/0ne-zero/f4h/utilities/functions/controller_helper"
	general_func "github.com/0ne-zero/f4h/utilities/functions/general"
	viewmodel "github.com/0ne-zero/f4h/web/view_model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
			"Title":      "Login/Signup",
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
			"Title":      "Login/Signup",
			"LoginError": "Username or Password is incorrect.",
		}
		c.HTML(http.StatusOK, "login.html", data)
		return
	}
	// Compare user password with entered password
	if err := general_func.ComparePassword(user_fields.PasswordHash, password); err != nil {
		data := gin.H{
			"Title": "Login/Signup",

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

	c.Redirect(http.StatusMovedPermanently, "/")
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
	c.Redirect(http.StatusMovedPermanently, "/")
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
	pass_hash, err := general_func.HashPassword(password)
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
	c.Redirect(http.StatusMovedPermanently, "/")
}
func Index(c *gin.Context) {
	// Get products
	products, err := model_function.GetProductInViewModel(15)
	if err != nil {
		contoller_helper.SomethingBadHappened(c, err, "")
		return
	}
	// Get categories
	categories, err := model_function.GetCategoriesWithRelationsInViewModel(true)
	if err != nil {
		contoller_helper.SomethingBadHappened(c, err, "")
	}

	// Everything is ok
	view_data := gin.H{
		"Title": "Index",
	}

	// If user is login, show her/his username. get username from session
	untyped_username := sessions.Default(c).Get("Username")
	if untyped_username != nil {
		view_data["HeaderData"] = gin.H{
			"Username": untyped_username.(string),
		}
	}
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
		enteredCategory = general_func.RemoveSlashFromBeginAndEnd(enteredCategory)

		// Get all categories name
		categoriesName, err := model_function.GetCategoriesName()
		if err != nil {
			contoller_helper.SomethingBadHappened(c, err, "")
			return
		}
		// Check entered category is exists in database
		if !general_func.ValueExistsInSlice(&categoriesName, enteredCategory) {
			contoller_helper.SomethingBadHappened(c, errors.New("Invalid Category"), "")
			return
		}
		// Get Products by entered category
		products, err := model_function.GetProductByCategoryInViewModel(strings.ToLower(enteredCategory), 15)
		if err != nil {
			contoller_helper.SomethingBadHappened(c, err, "")
			return
		}
		// Get all categories for sidebar
		categories, err := model_function.GetCategoriesWithRelationsInViewModel(true)
		if err != nil {
			contoller_helper.SomethingBadHappened(c, err, "")
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
			contoller_helper.SomethingBadHappened(c, err, "")
			return
		}
		var categories []model.Product_Category
		err = model_function.GetCategoryByOrderingProductsCount(&categories)
		if err != nil {
			contoller_helper.SomethingBadHappened(c, err, "")
			return
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
		contoller_helper.SomethingBadHappened(c, errors.New("Non-Int product id sent to ProductDetails controller"), "")
		return
	}
	// Integer product id
	int_id, err := strconv.ParseInt(str_id, 10, 64)
	// Parse check
	if err != nil {
		contoller_helper.SomethingBadHappened(c, errors.New(fmt.Sprintf("%s (%s)", err.Error(), "Error during parse product id to Int in ProductDetails controller")), "")
		return
	}
	var product model.Product
	err = model_function.GetByID(&product, int(int_id))
	// Check get product error
	if err != nil {
		contoller_helper.SomethingBadHappened(c, err, "")
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
		contoller_helper.SomethingBadHappened(c, err, "")
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
				contoller_helper.SomethingBadHappened(c, err, "")
				return
			}
			topic_count, err := model_function.GetDiscussionTopicsCount(d)
			if err != nil {
				contoller_helper.SomethingBadHappened(c, err, "")
				return
			}
			// Get discussion forums name
			forums_name, err := model_function.GetDiscussionForumsName(d)
			if err != nil {
				contoller_helper.SomethingBadHappened(c, err, "")
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
		contoller_helper.SomethingBadHappened(c, errors.New("Empty discussion name in DiscussionForums controller"), "Discussion name is empty in url parameter")
		return
	}
	var discussion_fields = model.Discussion{}
	err := model_function.GetFieldsByAnotherFieldValue(&discussion_fields, []string{"name"}, "name", discussion_name)
	if err != nil {
		contoller_helper.SomethingBadHappened(c, errors.New(fmt.Sprintf("%s (%s)", err.Error(), "Invalid discussion in DiscussionForums controller (discussion doesn't exists)")), fmt.Sprintf("%s discussion doesn't exists.", discussion_name))
		return
	}
	discussion_forums, err := model_function.GetDiscussionForumsInViewModel(int(discussion_fields.ID))
	if err != nil {
		contoller_helper.SomethingBadHappened(c, errors.New(fmt.Sprintf("%s (%s)", err.Error(), "Error during get discussion forums in Discussion Forums controller")), "")
		return
	}
	discussion_topics, err := model_function.GetDiscussionTopics(int(discussion_fields.ID))
	if err != nil {
		contoller_helper.SomethingBadHappened(c, errors.New(fmt.Sprintf("%s (%s)", err.Error(), "Error during get discussion topics in Discussion Forums controller")), "")
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
		contoller_helper.SomethingBadHappened(c, errors.New("Empty forum name sent to ForumTopics controller"), "Forum name is empty in url parameter")
		return
	}
	var forum_fields = model.Forum{}
	err := model_function.GetFieldsByAnotherFieldValue(&forum_fields, []string{"id"}, "name", forum_name)
	if err != nil {
		contoller_helper.SomethingBadHappened(c, errors.New(fmt.Sprintf("%s (%s)", err, "Invalid forum name in ForumTopics (forum name doesn't exists)")), "")
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
	write_topic_data_view_data := gin.H{}
	write_topic_data_view_data["ForumName"] = c.Param("forum")

	s := sessions.Default(c)
	// Get topic data and convert it in map[string]string
	topic_data, ok := s.Get("TopicData").(map[string]string)
	if !ok {
		contoller_helper.SomethingBadHappened(c, errors.New("Error during convert user's session topic data"), "")
		return
	}
	if topic_data != nil {
		write_topic_data_view_data["TopicSubject"] = topic_data["Subject"]
		write_topic_data_view_data["TopicMarkdown"] = topic_data["Markdown"]
	}

	view_data := gin.H{}
	view_data["WriteTopicData"] = write_topic_data_view_data
	c.HTML(200, "add_topic.html", view_data)
}
func AddTopic_POST(c *gin.Context) {
	// do error handelling

	// User topic Markdown
	topic_markdown := c.Request.FormValue("topic-markdown")
	if topic_markdown == "" {
		return
	}
	topic_markdown = strings.TrimSpace(topic_markdown)
	// Convert Markdown to Html
	topic_html, err := general_func.MarkdownToHtml(topic_markdown)
	if err != nil {
		return
	}
	// Remove space from start and end of topic_html
	topic_html = strings.TrimSpace(topic_html)
	// Prevent XSS Attacks
	topic_html = constansts.XSSPreventor.Sanitize(topic_html)

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
		subject := c.Request.FormValue("subject")
		var topic_data map[string]string
		topic_data["Subject"] = subject
		topic_data["Markdown"] = topic_markdown
		s := sessions.Default(c)
		s.Set("TopicData", topic_data)
		s.Save()
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
			// Login error
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
