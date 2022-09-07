package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/database/model_function"
	viewmodel "github.com/0ne-zero/f4h/public_struct/view_model"
	general_func "github.com/0ne-zero/f4h/utilities/functions/general"
	"github.com/0ne-zero/f4h/utilities/wrapper_logger"
	controller_helper "github.com/0ne-zero/f4h/web/controller/controller_helper"
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
		// Log
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Entered empty username or password", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		data := gin.H{
			"Title":      "Login/Signup",
			"LoginError": "Fill all field",
		}
		c.HTML(http.StatusUnprocessableEntity, "login.html", data)
		return
	}
	// Get user
	var user_fields = model.User{}
	err := model_function.GetFieldsByAnotherFieldValue(&user_fields, []string{"password_hash"}, "username", username)
	if err != nil {
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Entered incorrect username", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		data := gin.H{
			"Title":      "Login/Signup",
			"LoginError": "Username or Password is incorrect.",
		}
		c.HTML(http.StatusUnprocessableEntity, "login.html", data)
		return
	}
	// Compare user password with entered password
	status, err := general_func.ComparePassword(user_fields.PasswordHash, password)
	if err != nil {
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Error occurred during compare passwords", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	if !status {
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Entered Incorrect password", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
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
	// Get user ID by username for save it in session
	err = model_function.GetFieldsByAnotherFieldValue(&user_fields, []string{"id"}, "username", username)
	if err != nil {
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Error during get user id by username", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
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
	c.Redirect(http.StatusMovedPermanently, "/Login")
}
func Register_POST(c *gin.Context) {
	// Get username & password from form
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	// Validatae username & password
	if username == "" || password == "" {
		// Log
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Entered username or password is empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})

		view_data := gin.H{}
		view_data["RegisterError"] = "Fill all field"
		c.HTML(http.StatusUnprocessableEntity, "login.html", view_data)
		return
	}
	exists, err := model_function.IsUserExistsByUsername(username)
	if err != nil {
		// Log
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})

		view_data := gin.H{}
		view_data["RegisterError"] = "Unknown error"
		c.HTML(http.StatusUnprocessableEntity, "login.html", view_data)
		return
	}
	if exists {
		// Log
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Entered unavailable username", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		view_data := gin.H{}
		view_data["RegisterError"] = "This username is unavailable"
		c.HTML(http.StatusUnprocessableEntity, "login.html", view_data)
		return
	}
	// So far user not exists, should register client
	// Create password hash
	pass_hash, err := general_func.Hashing(password)
	if err != nil {
		//Log
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})

		view_data := gin.H{}
		view_data["RegisterError"] = "Unknown error,try again."
		c.HTML(http.StatusUnprocessableEntity, "login.html", view_data)
		return
	}
	// Add User
	model_function.Add(&model.User{Username: username, PasswordHash: pass_hash, Email: email})
	// Response
	c.Redirect(http.StatusMovedPermanently, "/")
}
func Index(c *gin.Context) {
	// Get products
	products, err := model_function.GetProductInProductBasicViewModel(15)
	if err != nil {
		// Log Error
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	// Get categories
	categories, err := model_function.GetCategoriesWithRelationsInViewModel(true)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}

	// Everything is ok
	view_data := gin.H{
		"Title": constansts.AppName + " | Home",
	}

	// If user is login, show her/his username. get username from session
	untyped_username := sessions.Default(c).Get("Username")
	if untyped_username != nil {
		username, ok := untyped_username.(string)
		if !ok {
			// Log
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during convert session's username to string ", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			// Return error page
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		view_data["HeaderData"] = gin.H{
			"Username": username,
		}
	}
	if categories != nil {
		view_data["Categories"] = categories
	}
	if products != nil {
		view_data["Products"] = products
	}
	c.HTML(200, "index.html", view_data)
}
func AddProduct_GET(c *gin.Context) {
	view_data := gin.H{
		"Title": controller_helper.SetTitle("Add Product"),
	}
	c.HTML(200, "add-product.html", view_data)
}

// Incomplete (get selected categories)
func AddProduct_POST(c *gin.Context) {
	p_name := strings.TrimSpace(c.PostForm("name"))
	p_price := strings.TrimSpace(c.PostForm("price"))
	p_inventory := strings.TrimSpace(c.PostForm("inventory"))
	p_description := strings.TrimSpace(c.PostForm("description"))
	p_tags := strings.TrimSpace(c.PostForm("tags"))

	// Validate data
	view_data := controller_helper.AddProductValidation(p_name, p_price, p_inventory, p_description, p_tags)
	if view_data != nil {
		c.HTML(200, "add-product.html", view_data)
		return
	}
	user_id := sessions.Default(c).Get("UserID").(int)
	p_price_float, _ := strconv.ParseFloat(p_price, 64)
	p_inventory_int, _ := strconv.Atoi(p_inventory)
	var tags []*model.Product_Tag
	// we have multiple tags
	if strings.Contains(p_tags, "|") {
		splitted_tags := strings.Split(p_tags, "|")
		for i := range splitted_tags {
			tag, err := model_function.FirstOrCreateProductTagByName(strings.TrimSpace(splitted_tags[i]))
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get tag by name", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
			tags = append(tags, tag)
		}
	} else {
		tag, err := model_function.FirstOrCreateProductTagByName(p_tags)
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get tag by name", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		tags = append(tags, tag)
	}
	// Save Images
	form, err := c.MultipartForm()
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred during parse multipart form\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	files := form.File["files[]"]
	p_images := make([]*model.Product_Image, len(files))
	for i := range files {
		var file_name string
		var err error
		var file_name_exists bool = true
		for file_name_exists {
			file_name, err = general_func.GenerateRandomHex(64)
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred generate random hex\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
			// Here loop can be break
			file_name_exists, err = general_func.IsImageExists(file_name, "PRODUCT")
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred during checking image exists\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
		}
		file_path := filepath.Join(constansts.ImagesDirectory, "product", file_name, ".jpeg")
		err = c.SaveUploadedFile(files[i], file_path)
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred during saving uploaded image\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		p_images = append(p_images, &model.Product_Image{Path: file_path})
	}
	//TODO: Selected categories

	// Create product
	product := &model.Product{
		UserID:      uint(user_id),
		Name:        p_name,
		Description: p_description,
		Price:       p_price_float,
		Inventory:   uint(p_inventory_int),
		Tags:        tags,
		Images:      p_images,
	}

	err = model_function.Add(product)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred during add product to database\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}

	c.Redirect(http.StatusMovedPermanently, fmt.Sprint("/ProductDetails/", product.ID))
}
func EditProduct_GET(c *gin.Context) {
	p_id, err := strconv.Atoi(strings.TrimSpace(c.Param("id")))
	if err != nil {
		return
	}
	p, err := model_function.GetProductInViewModel(p_id)
	if err != nil {

	}
	view_data := gin.H{
		"Title":   controller_helper.SetTitle(fmt.Sprintf("Edit %s Product", p.Name)),
		"Product": p,
	}
	c.HTML(200, "edit-product.html", view_data)
}

// Incomplete (get selected categories)
func EditProduct_POST(c *gin.Context) {
	p_id := strings.TrimSpace(c.PostForm("id"))
	p_name := strings.TrimSpace(c.PostForm("name"))
	p_price := strings.TrimSpace(c.PostForm("price"))
	p_inventory := strings.TrimSpace(c.PostForm("inventory"))
	p_description := strings.TrimSpace(c.PostForm("description"))
	p_tags := strings.TrimSpace(c.PostForm("tags"))

	// Validate data
	view_data := controller_helper.AddProductValidation(p_name, p_price, p_inventory, p_description, p_tags)
	if view_data != nil {
		c.HTML(200, "add-product.html", view_data)
		return
	}
	if p_id == "" {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Product id is Empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	p_id_int, err := strconv.Atoi(p_id)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Product id is non-int", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	user_id := sessions.Default(c).Get("UserID").(int)
	p_price_float, _ := strconv.ParseFloat(p_price, 64)
	p_inventory_int, _ := strconv.Atoi(p_inventory)

	var tags []*model.Product_Tag
	// we have multiple tags
	if strings.Contains(p_tags, "|") {
		splitted_tags := strings.Split(p_tags, "|")
		for i := range splitted_tags {
			tag, err := model_function.FirstOrCreateProductTagByName(strings.TrimSpace(splitted_tags[i]))
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get tag by name", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
			tags = append(tags, tag)
		}
	} else {
		tag, err := model_function.FirstOrCreateProductTagByName(p_tags)
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get tag by name", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		tags = append(tags, tag)
	}
	// Save Images
	form, err := c.MultipartForm()
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred during parse multipart form\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	files := form.File["files[]"]
	p_images := make([]*model.Product_Image, len(files))
	// Save uploaded images
	for i := range files {
		var file_name string
		var err error
		var file_name_exists bool = true
		for file_name_exists {
			file_name, err = general_func.GenerateRandomHex(64)
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred generate random hex\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
			// Here loop can be break
			file_name_exists, err = general_func.IsImageExists(file_name, "PRODUCT")
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred during checking image exists\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
		}
		file_path := filepath.Join(constansts.ImagesDirectory, "product", file_name, ".jpeg")
		err = c.SaveUploadedFile(files[i], file_path)
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred during saving uploaded image\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		p_images = append(p_images, &model.Product_Image{Path: file_path})
	}
	// Delete old images
	old_images_path, err := model_function.GetProductImagesPath(p_id_int)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred during get product old images\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	err = general_func.DeleteFiles(old_images_path...)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: fmt.Sprintf("Error occurred during delete product old images\n%s", err.Error()), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	//TODO: Get selected categories

	product := &model.Product{
		BasicModel:  model.BasicModel{ID: uint(p_id_int)},
		UserID:      uint(user_id),
		Name:        p_name,
		Description: p_description,
		Price:       p_price_float,
		Inventory:   uint(p_inventory_int),
		Tags:        tags,
		Images:      p_images,
	}
	// Edit Product
	_, err = model_function.Update(&model.Product{BasicModel: model.BasicModel{ID: uint(p_id_int)}}, product)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during edit product", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}

	c.Redirect(http.StatusMovedPermanently, fmt.Sprint("ProductDetails/", p_id_int))
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
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		// Check entered category is exists in database
		if !general_func.ValueExistsInSlice(&categoriesName, enteredCategory) {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Entered category name doesn't exists in database", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, "Category name entered in the url is invalid.")
			return
		}
		// Get Products by entered category
		products, err := model_function.GetProductByCategoryInViewModel(strings.ToLower(enteredCategory), 15)
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		// Get all categories for sidebar
		categories, err := model_function.GetCategoriesWithRelationsInViewModel(true)
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
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
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		var categories []model.Product_Category
		err = model_function.GetCategoryByOrderingProductsCount(&categories)
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
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
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Entered product id is empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Product id entered in the url is empty")
		return
	}
	// Integer product id
	int_id, err := strconv.Atoi(str_id)
	// Parse check
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Product id entered in the url is invalid")
		return
	}
	var product model.Product
	err = model_function.GetByID(&product, int_id)
	// Check get product error
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	// Get categories
	categories, err := model_function.GetCategoriesWithRelationsInViewModel(true)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	images_view_data, err := model_function.GetProductDetailsImagesInViewData(int_id)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	tabs_view_data, err := model_function.GetProductdetailsTabsContentInViewModel(int_id)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	recommended_viewdata, err := model_function.GetRecommendedProdcutsInViewModel(int_id)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}

	detail_view_data := viewmodel.ProductDetailsDetail{
		ID:        int(product.ID),
		Name:      product.Name,
		Price:     product.Price,
		Inventory: int(product.Inventory),
	}
	// Everything is ok
	view_data := gin.H{}
	view_data["Title"] = product.Name
	view_data["Product"] = product
	view_data["ProductImages"] = images_view_data
	view_data["ProductDetail"] = detail_view_data
	view_data["ProductTabs"] = tabs_view_data
	view_data["RecommendedProducts"] = recommended_viewdata
	view_data["Categories"] = categories
	c.HTML(http.StatusOK, "product-details.html", view_data)
}
func AddToCart(c *gin.Context) {
	user_id := sessions.Default(c).Get("UserID").(int)
	// Get product id
	p_id, err := strconv.Atoi(c.PostForm("p_id"))
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Non_int product id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	// Get quantity of product
	p_quantity, err := strconv.Atoi(c.PostForm("p_quantity"))
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Non-int quantity", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}

	// Get user cart id by user id
	cart_id, err := model_function.GetCartIDByUserID(user_id)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get cart id by user id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	err = model_function.AddProductToCart(p_id, cart_id, p_quantity)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during add product to user cart", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/ProductDetails/%d", p_id))
}
func Wishlist(c *gin.Context) {
	user_id := sessions.Default(c).Get("UserID").(int)
	// Get wishlist products
	wishlist_products, err := model_function.GetUserWishlistInViewmodel(user_id)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get wishlist products", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	view_data := gin.H{}
	view_data["Title"] = "Wishlist"
	view_data["Products"] = wishlist_products
	c.HTML(200, "wishlist.html", view_data)
}
func AddToWishlist(c *gin.Context) {
	user_id := sessions.Default(c).Get("UserID").(int)
	p_id, err := strconv.Atoi(c.Param("p_id"))
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Non-int product id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	err = model_function.AddToWishlist(user_id, p_id)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during add product to wishlist", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	// Redirect to referer page
	c.Redirect(http.StatusMovedPermanently, c.Request.Referer())
}
func AddProductComment(c *gin.Context) {
	user_id := sessions.Default(c).Get("UserID").(int)
	product_id_str := strings.TrimSpace(c.PostForm("product_id"))
	if product_id_str == "" {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Empty product id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	product_id, err := strconv.Atoi(product_id_str)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Non-int product id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	comment_text := strings.TrimSpace(c.PostForm("text"))
	// Sanitize
	comment_text = constansts.XSSPreventor.Sanitize(comment_text)
	if comment_text == "" {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Empty comment", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}

	err = model_function.AddProductComment(user_id, product_id, comment_text)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during add product comment", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/ProductDetails/%d", product_id))
}

// Incomplete
func Profile_GET(c *gin.Context) {
	// Available tabs and their modes
	tabs_modes := map[string][]string{
		"overview": {"front_page", "logins", "orders"},
		"profile":  {"edit_account", "edit_signature", "edit_avatar", "manage_login", "manage_address", "manage_wallets"},
		"products": {"front_page"},
		"payments": {"front_page"},
		"topics":   {"front_page"},
		"polls":    {"front_page"},
	}
	// Get user id
	user_id := sessions.Default(c).Get("UserID").(int)
	// Sort tabs and modes
	sorted_tabs := general_func.GetMapKeys(tabs_modes)
	sort.Strings(sorted_tabs)

	// Finally view data
	view_data := gin.H{
		"Title": fmt.Sprintf("%s Profile | %s", sessions.Default(c).Get("Username"), constansts.AppName),
		"Tabs":  sorted_tabs,
	}

	// Check is user selected tab
	if tab := c.Query("tab"); tab != "" {
		// Check selected tab is available
		if !general_func.ExistsStringInStringSlice(tab, general_func.GetMapKeys(tabs_modes)) {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Entered unavailable tab", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		// Insert tab to information that will send to template
		view_data["Tab"] = tab
		// Check user selected tab mode
		if mode := c.Query("mode"); mode != "" {
			// Check entered mode available
			if !general_func.ExistsStringInStringSlice(mode, tabs_modes[tab]) {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Entered unavailable mode", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
			// Insert mode to information that will send to template
			view_data["Mode"] = mode

			// Get panel data
			switch tab {
			case "overview":
				switch mode {
				case "front_page":
					// panel_data, err := model_function.GetUserDataForUserPanel_Overview_FrontPage(user_id)
					// if err != nil {
					// 	wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
					// 	controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
					// 	return
					// }
					// view_data["PanelData"] = panel_data
				case "logins":
					//panel_data,err := model_function.GetUserDataForUserPanel_Overview_FrontPage()
				case "orders":
				}
			case "profile":
				switch mode {
				case "edit_account":
				case "edit_signature":
				case "edit_avatar":
				case "manage_login":
				case "manage_address":
				case "manage_wallets":
				}
			case "products":
				switch mode {
				case "front_page":
				}
			case "payments":
				switch mode {
				case "front_page":
				}
			case "topics":
				switch mode {
				case "front_page":
				}
			case "polls":
				switch mode {
				case "front_page":
				}

			}

			c.HTML(200, "profile.html", view_data)
			return
		}
		// Tab's mode isn't selected so select default mode (first element)
		view_data["Mode"] = tabs_modes[tab][0]

		// Get panel data
		switch tab {
		case "overview":
		case "profile":
		case "products":
		case "payments":
		case "topics":
		case "polls":
		}
		c.HTML(200, "profile.html", view_data)
		return
	}

	// Neither tab nor tab mode is selected, so select default tab and default tab mode, which that means overview tab and its first mode
	view_data["Tab"] = "overview"
	view_data["Mode"] = "front_page"
	view_data["TabModes"] = tabs_modes["overview"]

	// Get panel data
	panel_data, err := model_function.GetUserDataForUserPanel_Overview_FrontPage(user_id)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	view_data["PanelData"] = panel_data
	c.HTML(200, "profile.html", view_data)
}

func Cart(c *gin.Context) {
	user_id := sessions.Default(c).Get("UserID").(int)
	// Get users cart information
	cart, err := model_function.GetUserCart(user_id)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get users cart information", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	view_data := gin.H{}
	view_data["Title"] = "Cart"
	view_data["CartItems"] = cart.CartItems
	view_data["CartInfo"] = cart
	c.HTML(200, "cart.html", view_data)
}

func DeleteCartItem(c *gin.Context) {
	id_str := c.Param("id")
	if id_str != "" {
		if id_int, err := strconv.Atoi(id_str); err == nil {
			err = model_function.DeleteCartItem(id_int)
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during delete cart item", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/Cart")
}
func DecreaseCartItemQuantity(c *gin.Context) {
	id_str := c.Param("id")
	if id_str != "" {
		if id_int, err := strconv.Atoi(id_str); err == nil {
			err = model_function.DecreaseCartItemQuantity(id_int)
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during decrease cart item quantity", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/Cart")
}
func IncreaseCartItemQuantity(c *gin.Context) {
	id_str := c.Param("id")
	if id_str != "" {
		if id_int, err := strconv.Atoi(id_str); err == nil {
			err = model_function.IncreaseCartItemQuantity(id_int)
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during increase cart item quantity", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/Cart")
}

//region Forum
func Discussions(c *gin.Context) {
	// Discussions categories and Discussion (preload Discussion)
	var Discussion_categories []model.Discussion_Category
	err := model_function.Get(&Discussion_categories, -1, "created_at", "ASC", "Discussions")
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
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
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
			topic_count, err := model_function.GetDiscussionTopicsCount(d)
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
			// Get discussion forums name
			forums_name, err := model_function.GetDiscussionForumsName(d)
			if err != nil {
				wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
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
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Entered discussion name is empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Discussion name entered in the url is empty")
		return
	}
	var discussion_fields = model.Discussion{}
	err := model_function.GetFieldsByAnotherFieldValue(&discussion_fields, []string{"name"}, "name", discussion_name)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, fmt.Sprintf("%s discussion doesn't exists.", discussion_name))
		return
	}
	discussion_forums, err := model_function.GetDiscussionForumsInViewModel(int(discussion_fields.ID))
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	discussion_topics, err := model_function.GetDiscussionTopics(int(discussion_fields.ID))
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}

	view_data := gin.H{}
	_topics_list_template_view_model := map[string]interface{}{
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
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Entered forum name is empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Forum name entered in the url is empty.")
		return
	}
	var forum_fields = model.Forum{}
	err := model_function.GetFieldsByAnotherFieldValue(&forum_fields, []string{"id"}, "name", forum_name)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Forum name entered in the url is invalid.")
		return
	}
	forum_topics, err := model_function.GetForumTopicsInViewModel(int(forum_fields.ID))
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Forum name entered in the url is invalid.")
		return
	}
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
	// Get topic data from session if any exists
	// Check sent forum name is exists
	forum_name := c.Param("forum")
	if forum_name == "" {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Entered forum name is empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Forum name entered in the url is empty.")
		return
	}
	var forum_fields model.Forum
	err := model_function.GetFieldsByAnotherFieldValue(&forum_fields, []string{"id"}, "name", forum_name)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error() + "Entered forum name is invalid", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Forum name entered in the url is invalid.")
		return
	}

	s := sessions.Default(c)
	view_data := gin.H{
		"Title": "Post new topic",
		"MSG":   fmt.Sprintf("You are posting topic in %s forum", forum_name),
		"WriteTopicData": gin.H{
			"Mode":          "Add",
			"ForumName":     forum_name,
			"TopicSubject":  s.Get("TopicSubject"),
			"TopicMarkdown": s.Get("TopicMarkdown"),
			"TopicTags":     s.Get("TopicTags"),
		},
	}
	c.HTML(200, "edit_add_topic.html", view_data)
}
func AddTopic_POST(c *gin.Context) {
	// Get forum name from the url
	forum_name := c.Param("forum")
	// Check sent forum name is exists
	if forum_name == "" {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Entered forum name is empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Forum name entered in the url is empty.")
		return
	}
	var forum_fields = model.Forum{}
	err := model_function.GetFieldsByAnotherFieldValue(&forum_fields, []string{"id"}, "name", forum_name)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error() + "Entered forum name is invalid", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Forum name entered in the url is invalid.")
		return
	}

	// User topic Markdown
	topic_markdown := c.Request.FormValue("topic-markdown")
	if topic_markdown == "" {
		// Log
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Entered topic markdown is empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		// Send error in add topic page for user
		view_data := gin.H{
			"Title": "Post new topic",
			"MSG":   fmt.Sprintf("You are posting topic in %s forum", forum_name),
			"WriteTopicData": gin.H{
				"Mode":          "Add",
				"TopicSubject":  c.Request.FormValue("subject"),
				"TopicMarkdown": topic_markdown,
				"TopicTags":     c.Request.FormValue("tags"),
				"Error":         "Topic Markdown field must be filled",
			},
		}
		c.HTML(http.StatusUnprocessableEntity, "edit_add_topic.html", view_data)
		return
	}
	topic_markdown = strings.TrimSpace(topic_markdown)
	// Convert Markdown to Html
	topic_html := general_func.MarkdownToHtml(topic_markdown)
	// Remove space from start and end of topic_html
	topic_html = strings.TrimSpace(topic_html)
	// Prevent XSS Attacks
	topic_html = constansts.XSSPreventor.Sanitize(topic_html)

	// User wants preview of her/his topic markdown
	if is_preview := c.Request.FormValue("preview"); is_preview != "" {
		view_data := gin.H{
			"Title": "Post new topic",
			"MSG":   fmt.Sprintf("You are posting topic in %s forum", forum_name),
			"WriteTopicData": gin.H{
				"Mode":          "Add",
				"ForumName":     forum_name,
				"TopicPreview":  template.HTML(topic_html),
				"TopicSubject":  c.Request.FormValue("subject"),
				"TopicMarkdown": topic_markdown,
				"TopicTags":     c.Request.FormValue("tags"),
			},
		}
		c.HTML(200, "edit_add_topic.html", view_data)
		return
		// User wants save her/his topic markdown (draft)
	} else if is_save := c.Request.FormValue("save"); is_save != "" {
		s := sessions.Default(c)
		s.Set("TopicSubject", c.Request.FormValue("subject"))
		s.Set("TopicMarkdown", topic_markdown)
		s.Set("TopicTags", c.Request.FormValue("tags"))
		s.Save()
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/AddTopic/%s", forum_name))
		// User wants delete her/his topic draft
	} else if is_delete := c.Request.FormValue("delete"); is_delete != "" {
		controller_helper.DeleteUserTopicDraftFromSession(c)
		view_data := gin.H{
			"Title": "Post new topic",
			"MSG":   fmt.Sprintf("You are posting topic in %s forum", forum_name),
			"WriteTopicData": gin.H{
				"Mode":          "Add",
				"ForumName":     forum_name,
				"TopicSubject":  c.Request.FormValue("subject"),
				"TopicMarkdown": topic_markdown,
				"TopicTags":     c.Request.FormValue("tags"),
			},
		}
		c.HTML(200, "edit_add_topic.html", view_data)
		return
		// User wants submit her/his text
	} else if is_submit := c.Request.FormValue("submit"); is_submit != "" {
		subject := c.Request.FormValue("subject")
		if subject == "" {
			// Log
			wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Empty topic subject sent", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			// Send error in add topic page for user
			view_data := gin.H{
				"Title": "Post new topic",
				"MSG":   fmt.Sprintf("You are posting topic in %s forum", forum_name),
				"WriteTopicData": gin.H{
					"Mode":          "Add",
					"TopicSubject":  c.Request.FormValue("subject"),
					"TopicMarkdown": topic_markdown,
					"TopicTags":     c.Request.FormValue("tags"),
					"Error":         "Topic subject field must be filled",
				},
			}
			// 442 = Unprocessable
			c.HTML(442, "edit_add_topic.html", view_data)
			return
		}
		// Fill topic and insert it in database
		var topic model.Topic
		topic.Name = subject
		topic.Description = topic_html
		topic.ForumID = uint(forum_fields.ID)
		// Add topic tags if user entered any
		if sloppy_tags := c.Request.FormValue("tags"); sloppy_tags != "" {
			splited_tags := strings.Split(sloppy_tags, "|")
			//Return error response if the number of topic tags is greater than five
			if len(splited_tags) > 5 {
				// Log
				wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Empty topic subject sent", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				// Send error in add topic page for user
				view_data := gin.H{
					"Title": "Post new topic",
					"MSG":   fmt.Sprintf("You are posting topic in %s forum", forum_name),
					"WriteTopicData": gin.H{
						"TopicSubject":  c.Request.FormValue("subject"),
						"TopicMarkdown": topic_markdown,
						"TopicTags":     c.Request.FormValue("tags"),
						"Error":         "Number of topic tags must less than five",
					},
				}
				// 442 = Unprocessable
				c.HTML(442, "edit_add_topic.html", view_data)
				return
			}
			for i := range splited_tags {
				// Get or create tag by its name
				tag, err := model_function.FirstOrCreateTopicTagByName(splited_tags[i])
				if err != nil {
					wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during check tag existance", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
					controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
					return
				}
				topic.Tags = append(topic.Tags, tag)
			}
		}
		// Add user id
		untyped_user_id := sessions.Default(c).Get("UserID")
		if untyped_user_id == nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "User hasn't user_id in her/his session values", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, "You are not logged in,go to login page.")
			return
		}
		user_id, ok := untyped_user_id.(int)
		if !ok {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Invalid user_id in user's session", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, "Your user_id in your session is invalid.")
			return
		}
		topic.UserID = uint(user_id)

		// Insert topic
		// When Add Function executed topic.ID automatically set by dataase, and so we have topic's id
		err := model_function.Add(&topic)
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}

		// Delete topic information in user session
		controller_helper.DeleteUserTopicDraftFromSession(c)

		// Send user to topic page
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/Topic/%d", topic.ID))
	}
}
func EditTopic_Get(c *gin.Context) {
	// Get topic id
	topic_id_str := c.Param("topic_id")
	if topic_id_str == "" {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Empty toipc id entered", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Entered topic id in the url is empty.")
		return
	}
	// Convert topic id to int
	topic_id, err := strconv.ParseInt(topic_id_str, 10, 64)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Invalid toipc id entered", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Entered topic id in the url is Invalid.")
		return
	}

	// Check user has saved edit information
	s := sessions.Default(c)
	if topic_subject := s.Get("TopicSubject"); topic_subject != nil {
		topic_name, err := model_function.GetTopicNameByTopicID(int(topic_id))
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get topic name" + err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		view_data := gin.H{
			"Title": fmt.Sprintf("Edit %s topic", topic_name),
			"MSG":   fmt.Sprintf("You are editing %s topic", topic_name),
			"WriteTopicData": gin.H{
				"Mode":          "Edit",
				"TopicID":       topic_id,
				"TopicSubject":  topic_subject,
				"TopicMarkdown": s.Get("TopicMarkdown"),
				"TopicTags":     s.Get("TopicTags"),
			},
		}
		c.HTML(200, "edit_add_topic.html", view_data)
		return
	}

	// User hasn't saved edit information, so we show topic information saved in database
	// Get topic informations for edit
	topic, err := model_function.GetTopicByIDForEdit(int(topic_id))
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get topic for edit " + err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	if topic == nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "User entered topic id which that topic id dosen't exists" + err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, fmt.Sprintf("This topic with %d id doesn't exists", topic_id))
		return
	}
	// Convert topic html to markdown
	topic.Description, err = general_func.HtmlToMarkdown(topic.Description)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during convert topic html to markdown " + err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	view_data := gin.H{
		"Title": fmt.Sprintf("Edit %s topic", topic.Name),
		"MSG":   fmt.Sprintf("You are editing %s topic", topic.Name),
		"WriteTopicData": gin.H{
			"Mode":          "Edit",
			"TopicID":       topic_id,
			"TopicSubject":  topic.Name,
			"TopicMarkdown": topic.Description,
			"TopicTags":     general_func.SplitEachTagsWithPipe(topic.Tags),
		},
	}
	c.HTML(200, "edit_add_topic.html", view_data)
}
func EditTopic_POST(c *gin.Context) {
	// Get topic id
	topic_id_str := strings.TrimSpace(c.PostForm("id"))
	if topic_id_str == "" {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Empty toipc id entered", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Entered topic id in the url is empty.")
		return
	}
	// Convert topic id to int
	topic_id, err := strconv.ParseInt(topic_id_str, 10, 64)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Invalid toipc id entered", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Entered topic id in the url is Invalid.")
		return
	}

	// topic Markdown
	topic_markdown := c.Request.FormValue("topic-markdown")
	if topic_markdown == "" {
		// Log
		wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Entered topic markdown is empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		// Send error in add topic page for user
		topic_name, err := model_function.GetTopicNameByTopicID(int(topic_id))
		if err != nil {
			// Log
			wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Error occurred during get topic name form its id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			// Return error page
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		view_data := gin.H{
			"Title": "Edit a topic",
			"MSG":   fmt.Sprintf("You are editing %s topic", topic_name),
			"WriteTopicData": gin.H{
				"Mode":          "Edit",
				"TopicID":       topic_id,
				"TopicSubject":  c.Request.FormValue("subject"),
				"TopicMarkdown": topic_markdown,
				"TopicTags":     c.Request.FormValue("tags"),
				"Error":         "Topic Markdown field must be filled",
			},
		}
		c.HTML(http.StatusUnprocessableEntity, "edit_add_topic.html", view_data)
		return
	}
	topic_markdown = strings.TrimSpace(topic_markdown)
	// Convert Markdown to Html
	topic_html := general_func.MarkdownToHtml(topic_markdown)
	// Remove space from start and end of topic_html
	topic_html = strings.TrimSpace(topic_html)
	// Prevent XSS Attacks
	topic_html = constansts.XSSPreventor.Sanitize(topic_html)

	// User wants preview of her/his topic markdown
	if is_preview := c.Request.FormValue("preview"); is_preview != "" {
		topic_name, err := model_function.GetTopicNameByTopicID(int(topic_id))
		if err != nil {
			// Log
			wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Error occurred during get topic name form its id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			// Return error page
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		view_data := gin.H{
			"Title": "Edit a topic",
			"MSG":   fmt.Sprintf("You are editing %s topic", topic_name),
			"WriteTopicData": gin.H{
				"Mode":          "Edit",
				"TopicID":       topic_id,
				"TopicPreview":  template.HTML(topic_html),
				"TopicSubject":  c.Request.FormValue("subject"),
				"TopicMarkdown": topic_markdown,
				"TopicTags":     c.Request.FormValue("tags"),
			},
		}
		c.HTML(200, "edit_add_topic.html", view_data)
		return
		// User wants save her/his topic markdown
	} else if is_save := c.Request.FormValue("save"); is_save != "" {
		s := sessions.Default(c)
		s.Set("TopicSubject", c.Request.FormValue("subject"))
		s.Set("TopicMarkdown", topic_markdown)
		s.Set("TopicTags", c.Request.FormValue("tags"))
		s.Save()
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/EditTopic/%d", topic_id))
		// User wants delete her/his topic draft
	} else if is_delete := c.Request.FormValue("delete"); is_delete != "" {
		controller_helper.DeleteUserTopicDraftFromSession(c)
		topic_name, err := model_function.GetTopicNameByTopicID(int(topic_id))
		if err != nil {
			// Log
			wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Error occurred during get topic name form its id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			// Return error page
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		view_data := gin.H{
			"Title": "Edit a topic",
			"MSG":   fmt.Sprintf("You are editing %s topic", topic_name),
			"WriteTopicData": gin.H{
				"Mode":          "Edit",
				"TopicID":       topic_id,
				"TopicSubject":  c.Request.FormValue("subject"),
				"TopicMarkdown": topic_markdown,
				"TopicTags":     c.Request.FormValue("tags"),
			},
		}
		c.HTML(200, "edit_add_topic.html", view_data)
		return
		// User wants submit her/his text
	} else if is_submit := c.Request.FormValue("submit"); is_submit != "" {
		subject := c.Request.FormValue("subject")
		if subject == "" {
			// Log
			wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Empty topic subject sent", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			// Send error in add topic page for user
			topic_name, err := model_function.GetTopicNameByTopicID(int(topic_id))
			if err != nil {
				// Log
				wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Error occurred during get topic name form its id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				// Return error page
				controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
				return
			}
			view_data := gin.H{
				"Title": "Post new topic",
				"MSG":   fmt.Sprintf("You are editing %s topic", topic_name),
				"WriteTopicData": gin.H{
					"Mode":          "Edit",
					"TopicID":       topic_id,
					"TopicSubject":  c.Request.FormValue("subject"),
					"TopicMarkdown": topic_markdown,
					"TopicTags":     c.Request.FormValue("tags"),
					"Error":         "Topic subject field must be filled",
				},
			}
			// 442 = Unprocessable
			c.HTML(442, "edit_add_topic.html", view_data)
			return
		}
		// Fill topic and insert it in database
		var updated_topic model.Topic
		updated_topic.Name = subject
		updated_topic.Description = topic_html
		forum_id, err := model_function.GetTopicForumIDByTopicID(int(topic_id))
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during get topic forum id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}
		updated_topic.ForumID = uint(forum_id)
		// Add topic tags if user entered any
		if sloppy_tags := c.Request.FormValue("tags"); sloppy_tags != "" {
			splited_tags := strings.Split(sloppy_tags, "|")
			//Return error response if the number of topic tags is greater than five
			if len(splited_tags) > 5 {
				// Log
				wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Empty topic subject sent", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
				// Send error in add topic page for user
				topic_name, err := model_function.GetTopicNameByTopicID(int(topic_id))
				if err != nil {
					// Log
					wrapper_logger.Debug(&wrapper_logger.LogInfo{Message: "Error occurred during get topic name form its id", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
					// Return error page
					controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
					return
				}
				view_data := gin.H{
					"Title": "Post new topic",
					"MSG":   fmt.Sprintf("You are editing %s topic", topic_name),
					"WriteTopicData": gin.H{
						"TopicSubject":  c.Request.FormValue("subject"),
						"TopicID":       topic_id,
						"TopicMarkdown": topic_markdown,
						"TopicTags":     c.Request.FormValue("tags"),
						"Error":         "Number of topic tags must less than five",
					},
				}

				// 442 = Unprocessable
				c.HTML(442, "edit_add_topic.html", view_data)
				return
			}
			for i := range splited_tags {
				// Get or create tag by its name
				tag, err := model_function.FirstOrCreateTopicTagByName(splited_tags[i])
				if err != nil {
					wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Error occurred during check tag existance", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
					controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
					return
				}
				updated_topic.Tags = append(updated_topic.Tags, tag)
			}
		}
		// Add user id
		untyped_user_id := sessions.Default(c).Get("UserID")
		if untyped_user_id == nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "User hasn't user_id in her/his session values", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, "You are not logged in,go to login page.")
			return
		}
		user_id, ok := untyped_user_id.(int)
		if !ok {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Invalid user_id in user's session", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, "Your user_id in your session is invalid.")
			return
		}
		updated_topic.UserID = uint(user_id)

		// Update topic
		// Create a topic with id for detect which topic should be change
		result_model, err := model_function.Update(&model.Topic{BasicModel: model.BasicModel{ID: uint(topic_id)}}, &updated_topic)
		if err != nil {
			wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
			controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
			return
		}

		// Send user to topic page
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/Topic/%d", result_model.ID))
	}
}
func ShowTopic(c *gin.Context) {
	topic_id_string := c.Param("topic_id")
	if topic_id_string == "" {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: "Entered topic id in the url is empty", Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Topic id entered in the url is empty.")
		return
	}
	topic_id, err := strconv.ParseInt(topic_id_string, 10, 64)
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, "Topic id entered in the url is invalid.")
		return
	}
	// Get Topic and topic comments
	topic, err := model_function.GetTopicByIDForShowTopicInViewModel(int(topic_id))
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}
	topic_comments, err := model_function.GetTopicCommentsByIDForShowTopicInViewModel(int(topic_id))
	if err != nil {
		wrapper_logger.Warning(&wrapper_logger.LogInfo{Message: err.Error(), Fields: controller_helper.ClientInfoInMap(c), ErrorLocation: general_func.GetCallerInfo(0)})
		controller_helper.ErrorPage(c, constansts.SomethingBadHappenedError)
		return
	}

	// Add topic title to each topic comment
	for i := range topic_comments {
		topic_comments[i].Title = topic.Title
	}

	view_data := gin.H{}
	// Topic and topic comments
	view_data["Topic"] = topic
	view_data["TopicComments"] = topic_comments
	c.HTML(200, "view_topic.html", view_data)
}
func Admin_Index(c *gin.Context) {

	c.HTML(200, "admin_index.html", nil)
}

//endregion

// func About()

// func TagProducts()

// func CategoryProducts()

// func Categories()

// func Dashboard()
