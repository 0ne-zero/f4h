package model_function

import (
	"errors"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/0ne-zero/f4h/database"
	"github.com/0ne-zero/f4h/database/model"
	viewmodel "github.com/0ne-zero/f4h/public_struct/view_model"
	general_func "github.com/0ne-zero/f4h/utilities/functions/general"
	"github.com/0ne-zero/f4h/utilities/wrapper_logger"
)

type Model interface {
	model.Forum | model.Discussion | model.CartItem | model.User | model.Product_Tag | model.Product_Category | model.Product_Comment | model.Product | model.Request | model.Discussion_Category | model.BadRequest | model.Topic | model.Topic_Tag
}

func Add[m Model](model *m) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	return db.Create(model).Error
}
func Delete[m Model](model *m, id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	return db.Unscoped().Delete(model, id).Error
}

// Input:
// consider_model = model with its id
// update_model = the model with some change, you wish to apply the consider model
// Returns Changed model (result)
func Update[m Model](consider_model *m, updated_model *m) (*m, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	err := db.Model(consider_model).Updates(updated_model).Error
	if err != nil {
		return nil, err
	}
	//err := db.Save(consider_model).Error

	return consider_model, err
}
func Get[m Model](ref_model *[]m, limit int, orderBy string, orderType string, preloads ...string) error {

	context := database.InitializeOrGetDB()
	if context == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}

	for _, p := range preloads {
		// Include Preload command in db commands chain
		context = context.Preload(p)
	}
	var err error

	if limit < 1 {
		err = context.Order(fmt.Sprintf("%s %s", orderBy, orderType)).Find(ref_model).Error
	} else {
		err = context.Order(fmt.Sprintf("%s %s", orderBy, orderType)).Find(ref_model).Limit(limit).Error
	}
	return err
}
func IsExistsByID[m Model](model *m, id uint) (bool, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var exists bool
	err := db.Model(model).Select("count(*) >0").Where("ID = ?", id).Find(&exists).Error
	return exists, err
}
func GetByID[m Model](model *m, id int) error {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	return db.Where("id = ?", id).First(model).Error
}
func GetUserPassHashByUsername(username string) (string, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var pass_hash string
	err := db.Where("username = ?", username).Select("password_hash").First(&pass_hash).Error
	return pass_hash, err
}
func IsUserExistsByUsername(username string) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var exists bool
	err := db.Model(&model.User{}).Select("count(*) >0").Where("username = ?", username).Find(&exists).Error
	return exists, err
}
func GetUserByUsername(username string) (model.User, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var u model.User
	err := db.Where("username = ?", username).First(&u).Error
	return u, err
}
func GetUsernameByUserID(user_id int) (string, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var username string
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("username").Scan(&username).Error
	return username, err
}

func GetUserAvatarPath(user_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var avatar_path string
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("avatar_path").Scan(&avatar_path).Error
	if err != nil {
		return "", err
	}
	return avatar_path, err
}
func GetFieldsByAnotherFieldValue[m Model](model *m, out_fields_name []string, in_field_name string, in_field_value string) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	err := db.Model(model).Where(fmt.Sprintf("%s = ?", in_field_name), in_field_value).Select(out_fields_name).Scan(model).Error
	return err
}
func GetProductDetailInViewData(p_id int) (*viewmodel.ProductDetailsDetail, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var vm viewmodel.ProductDetailsDetail
	err := db.Model(&model.Product{}).Where("id = ?", p_id).Select("id", "name", "inventory", "price").Scan(&vm).Error
	return &vm, err
}

func GetProductImagesPath(p_id int) ([]string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var images []model.Product_Image
	err := db.Model(&model.Product_Image{}).Where("product_id = ?", p_id).Select("Path").Scan(&images).Error
	if err != nil {
		return nil, err
	}
	var images_path []string
	for i := range images {
		images_path = append(images_path, images[i].Path)
	}
	return images_path, nil
}
func GetProductDetailsImagesInViewData(p_id int) (*viewmodel.ProductDetailsImagesViewData, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Get all images of product
	var images_path []string
	err := db.Model(&model.Product_Image{}).Where("product_id = ?", p_id).Select("path").Scan(&images_path).Error
	if err != nil {
		return nil, err
	}
	// Get product name
	var p_name string
	err = db.Model(&model.Product{}).Where("id = ?", p_id).Select("name").Scan(&p_name).Error
	if err != nil {
		return nil, err
	}
	// Fill view model (vm)
	var vm viewmodel.ProductDetailsImagesViewData
	for i := range images_path {
		if i == 0 {
			vm.MainImage = images_path[0]
			continue
		}
		vm.SubImages = append(vm.SubImages, viewmodel.ImageViewData{Path: images_path[i], Name: p_name})
	}
	// Computing number of slides
	n := float64(len(vm.SubImages)) / float64(3)
	// Is n round or not
	if general_func.IsFloatNumberRound(n) {
		vm.NumberOfSlides = int(n)
	} else {
		vm.NumberOfSlides = int(n) + 1
	}
	if vm.NumberOfSlides == 0 {
		vm.NumberOfSlides = 1
	}
	return &vm, nil
}
func GetUserWishlistInViewmodel(user_id int) ([]viewmodel.ProductBasicViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var w model.Wishlist
	err := db.Model(&model.Wishlist{}).Preload("Products").Preload("Products.Images").Where("user_id = ?", user_id).Find(&w).Error
	if err != nil {
		return nil, err
	}
	if w.Products == nil {
		return []viewmodel.ProductBasicViewModel{}, nil
	}
	var vm = make([]viewmodel.ProductBasicViewModel, len(w.Products))
	for i := range w.Products {
		vm[i].ID = int(w.Products[i].ID)
		vm[i].Name = w.Products[i].Name
		vm[i].Price = w.Products[i].Price
		if w.Products[i].Images != nil {
			vm[i].ImagePath = w.Products[i].Images[0].Path
		}
	}
	return vm, nil
}
func AddProductToCart(p_id, cart_id, quantity int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Check product exist in cart
	exists, err := isProductInCart(cart_id, p_id)
	if err != nil {
		return err
	}
	if exists {
		// Get current quantity of product in cart
		var current_quantity int
		err = db.Model(&model.CartItem{}).Where("cart_id = ? AND product_id = ?", cart_id, p_id).Select("product_quantity").Scan(&current_quantity).Error
		if err != nil {
			return err
		}
		// Update the quantity of product
		final_quantity := current_quantity + quantity
		return db.Model(&model.CartItem{}).Where("cart_id = ? AND product_id = ?", cart_id, p_id).Update("product_quantity", final_quantity).Error
	} else {
		return Add(&model.CartItem{ProductID: uint(p_id), CartID: uint(cart_id), ProductQuantity: uint(quantity)})
	}
}
func GetUserProductsInViewmodel(user_id int) ([]viewmodel.ProductDetailsUserProduct, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var user_products []model.Product
	err := db.Model(&model.Product{}).Where("user_id = ?", user_id).Select("name", "price", "id").Find(&user_products).Error
	if err != nil {
		return nil, err
	}
	user_products_vm := make([]viewmodel.ProductDetailsUserProduct, len(user_products))
	for i := range user_products {
		user_products_vm[i].ID = int(user_products[i].ID)
		user_products_vm[i].Name = user_products[i].Name
		user_products_vm[i].Price = user_products[i].Price
		// Get main image path of product
		main_img_path, err := getMainImagePathOfProduct(user_products_vm[i].ID)
		if err != nil {
			return nil, err
		}
		user_products_vm[i].ImagePath = main_img_path
	}
	return user_products_vm, nil
}
func GetProductCommentsInViewmodel(p_id int) ([]viewmodel.ProductDetailsComment, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var comments []model.Product_Comment
	err := db.Model(&model.Product_Comment{}).Where("product_id = ?", p_id).Select("id", "created_at", "user_id", "text").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	var comments_vm viewmodel.ProductDetailsComments
	for i := range comments {
		// Get username of who written comment
		var c_vm viewmodel.ProductDetailsComment
		c_vm.ID = int(comments[i].ID)
		c_vm.Text = comments[i].Text
		c_vm.Time = &comments[i].CreatedAt
		user_name, err := GetUsernameByUserID(int(comments[i].UserID))
		if err != nil {
			return nil, err
		}
		c_vm.Username = user_name
		comments_vm.Comments = append(comments_vm.Comments, c_vm)
	}
	return comments_vm.Comments, nil
}
func GetProductdetailsTabsContentInViewModel(p_id int) (*viewmodel.ProductDetailsTabs, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Get product needed information
	var p model.Product
	err := db.Model(&model.Product{}).Where("id = ?", p_id).Select("user_id", "name", "description", "created_at").Scan(&p).Error
	if err != nil {
		return nil, err
	}
	var description_vm viewmodel.ProductDetailsDescription
	description_vm.Description = p.Description
	description_vm.Time = p.CreatedAt
	// Get username by id
	user_name, err := GetUsernameByUserID(int(p.UserID))
	if err != nil {
		return nil, err
	}
	description_vm.Username = user_name

	// Get user products
	user_products_vm, err := GetUserProductsInViewmodel(int(p.UserID))
	if err != nil {
		return nil, err
	}
	// Delete taken product from user products list
	var del_element_index int
	for i := range user_products_vm {
		if user_products_vm[i].ID == p_id {
			del_element_index = i
		}
	}
	user_products_vm = general_func.RemoveSliceElement(user_products_vm, del_element_index)

	// Get product comments
	comments_vm, err := GetProductCommentsInViewmodel(p_id)
	if err != nil {
		return nil, err
	}
	return &viewmodel.ProductDetailsTabs{
		DescriptionData: description_vm, UserProductsData: user_products_vm,
		CommentsData:     viewmodel.ProductDetailsComments{ProductID: p_id, Comments: comments_vm},
		NumberOfComments: len(comments_vm)}, nil
}
func GetRecommendedProdcuts(by_product_id int) ([]model.Product, error) {
	return nil, nil
}
func GetRecommendedProdcutsInViewModel(by_product_id int) (*viewmodel.RecommendedViewData, error) {
	return nil, nil
}
func TooManyRequest(ip string, url string, method string) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var req_count int64
	now := time.Now().UTC()
	one_hour_ago := now.Add(time.Duration(-1) * time.Hour)
	err := db.Model(&model.Request{}).Where("ip = ? AND url = ? AND method = ? AND time <= ? ", ip, url, method, one_hour_ago).Count(&req_count).Error
	if err != nil {
		return false, err
	}
	if req_count > 100 {
		return true, nil
	}
	return false, nil
}
func GetProductInProductBasicViewModel(limit int) ([]viewmodel.ProductBasicViewModel, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var products []model.Product
	var err error
	if limit > 0 {
		err = db.Model(&model.Product{}).Limit(limit).Preload("Images").Select("id", "name", "price").Find(&products).Error
	} else {
		err = db.Model(&model.Product{}).Preload("Images").Select("id", "name", "price").Find(&products).Error
	}
	if len(products) < 1 {
		return []viewmodel.ProductBasicViewModel{}, nil
	}
	var vm = make([]viewmodel.ProductBasicViewModel, len(products))
	for i := range products {
		vm[i].ID = int(products[i].ID)
		vm[i].Name = products[i].Name
		vm[i].Price = products[i].Price
		if products[i].Images != nil {
			vm[i].ImagePath = products[i].Images[0].Path
		}
	}
	return vm, err
}
func GetProductByCategoryInViewModel(category_name string, limit int) ([]viewmodel.ProductBasicViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}

	var c model.Product_Category
	err := db.Preload("Products").Preload("Products.Images").Where("name = ?", category_name).Find(&c).Error
	// If error occured or products are nil returns nil
	if err != nil {
		return nil, err
	} else if c.Products == nil {
		return nil, errors.New("products field is empty")
	}

	// Fill view model
	var vm = make([]viewmodel.ProductBasicViewModel, len(c.Products))
	for i := range c.Products {
		vm[i].ID = int(c.Products[i].ID)
		vm[i].Name = c.Products[i].Name
		vm[i].Price = c.Products[i].Price
		if c.Products[i].Images != nil {
			vm[i].ImagePath = c.Products[i].Images[0].Path
		}
	}
	return vm, nil
}
func GetCategoryByOrderingProductsCount(c *[]model.Product_Category) error {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Get categories
	var categories []model.Product_Category
	err := db.Preload("Products").Preload("SubCategories").Find(&categories).Error
	if err != nil {
		return err
	}
	// Order categories by products
	sort.Slice(categories, func(i, j int) bool {
		return len(categories[i].Products) > len(categories[j].Products)
	})
	// Find sub categories
	var sub_categories []model.Product_Category
	for _, cat := range categories {
		cat.Products = nil
		sub_categories = append(sub_categories, cat.SubCategories...)
	}
	// Remove sub categories from parent categories list
	for i, cat := range categories {
		for _, subcat := range sub_categories {
			if cat.IsEqual(&subcat) {
				categories = general_func.RemoveSliceElement(categories, i)
			}
		}
	}
	// Fill input model
	*c = categories
	return nil
}
func GetCategoriesWithRelationsInViewModel(ordering bool) ([]viewmodel.SidebarCategoryViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var categories []model.Product_Category
	var err error
	if !ordering {
		err = db.Preload("SubCategories").Find(&categories).Error
	} else {
		err = db.Preload("Products").Preload("SubCategories").Find(&categories).Error
	}
	if err != nil {
		return nil, err
	}
	if ordering {
		// Order categories by products
		sort.Slice(categories, func(i, j int) bool {
			return len(categories[i].Products) > len(categories[j].Products)
		})
	}
	var result []viewmodel.SidebarCategoryViewModel
	for _, c := range categories {
		var view_cat viewmodel.SidebarCategoryViewModel
		view_cat.Name = c.Name
		// If category is a sub-category skip it
		if c.ParentID != nil {
			continue
		}
		for _, sub := range c.SubCategories {
			var view_cat_sub viewmodel.SidebarCategoryViewModel
			view_cat_sub.Name = sub.Name
			view_cat.SubCategories = append(view_cat.SubCategories, view_cat_sub)
		}

		result = append(result, view_cat)
	}
	return result, err
}
func IncreaseProductCommentPositiveVote(user_id, pc_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Check user already voted
	// Each user can vote one time
	is_voted, err := isUserVotedToProductComment(user_id, pc_id)
	if err != nil {
		return err
	} else if !is_voted {
		var current_positive_vote int
		err = db.Model(&model.Product_Comment_Vote{}).Where("product_comment_id = ?", pc_id).Select("positive").Scan(&current_positive_vote).Error
		if err != nil {
			return err
		}
		return db.Model(&model.Product_Comment_Vote{}).Where("product_comment_id = ?", pc_id).Update("positive", current_positive_vote+1).Error
	}
	return err
}
func DecreaseProductCommentPositiveVote(user_id, pc_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Check user already voted
	// Each user can vote one time
	is_voted, err := isUserVotedToProductComment(user_id, pc_id)
	if err != nil {
		return err
	}
	if !is_voted {
		var current_positive_vote int
		err := db.Model(&model.Product_Comment_Vote{}).Where("product_comment_id = ?", pc_id).Select("positive").Scan(&current_positive_vote).Error
		if err != nil {
			return err
		}
		return db.Model(&model.Product_Comment_Vote{}).Where("product_comment_id = ?", pc_id).Update("positive", current_positive_vote-1).Error
	}
	return err
}
func IncreaseProductCommentNegativeVote(user_id, pc_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Check user already voted
	// Each user can vote one time
	is_voted, err := isUserVotedToProductComment(user_id, pc_id)
	if err != nil {
		return err
	} else if !is_voted {
		var current_negative_vote int
		err := db.Model(&model.Product_Comment_Vote{}).Where("product_comment_id = ?", pc_id).Select("negative").Scan(&current_negative_vote).Error
		if err != nil {
			return err
		}
		return db.Model(&model.Product_Comment_Vote{}).Where("product_comment_id = ?", pc_id).Update("negative", current_negative_vote+1).Error
	}
	return err
}
func DeccreaseProductCommentNegativeVote(user_id, pc_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Check user already voted
	// Each user can vote one time
	is_voted, err := isUserVotedToProductComment(user_id, pc_id)
	if err != nil {
		return err
	} else if !is_voted {
		var current_negative_vote int
		err := db.Model(&model.Product_Comment_Vote{}).Where("product_comment_id = ?", pc_id).Select("negative").Scan(&current_negative_vote).Error
		if err != nil {
			return err
		}
		return db.Model(&model.Product_Comment_Vote{}).Where("product_comment_id = ?", pc_id).Update("negative", current_negative_vote-1).Error
	}
	return err
}

func IncreaseTopicCommentPositiveVote(user_id, tc_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Check user already voted
	// Each user can vote one time
	is_voted, err := isUserVotedToTopicComment(user_id, tc_id)
	if err != nil {
		return err
	} else if !is_voted {
		var current_positive_vote int
		err := db.Model(&model.Topic_Comment_Vote{}).Where("topic_comment_id = ?", tc_id).Select("positive").Scan(&current_positive_vote).Error
		if err != nil {
			return err
		}
		return db.Model(&model.Topic_Comment_Vote{}).Where("topic_comment_id = ?", tc_id).Update("positive", current_positive_vote+1).Error
	}
	return err
}
func DecreaseTopicCommentPositiveVote(user_id, tc_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Check user already voted
	// Each user can vote one time
	is_voted, err := isUserVotedToTopicComment(user_id, tc_id)
	if err != nil {
		return err
	} else if !is_voted {
		var current_positive_vote int
		err := db.Model(&model.Topic_Comment_Vote{}).Where("topic_comment_id = ?", tc_id).Select("positive").Scan(&current_positive_vote).Error
		if err != nil {
			return err
		}
		return db.Model(&model.Topic_Comment_Vote{}).Where("topic_comment_id = ?", tc_id).Update("positive", current_positive_vote-1).Error
	}
	return err
}
func IncreaseTopicCommentNegativeVote(user_id, tc_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Check user already voted
	// Each user can vote one time
	is_voted, err := isUserVotedToTopicComment(user_id, tc_id)
	if err != nil {
		return err
	} else if !is_voted {
		var current_negative_vote int
		err := db.Model(&model.Topic_Comment_Vote{}).Where("topic_comment_id = ?", tc_id).Select("negative").Scan(&current_negative_vote).Error
		if err != nil {
			return err
		}
		return db.Model(&model.Topic_Comment_Vote{}).Where("topic_comment_id = ?").Update("negative", current_negative_vote+1).Error
	}
	return err
}
func DecreaseTopicCommentNegativeVote(user_id, tc_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Check user already voted
	// Each user can vote one time
	is_voted, err := isUserVotedToTopicComment(user_id, tc_id)
	if err != nil {
		return err
	} else if !is_voted {
		var current_negative_vote int
		err := db.Model(&model.Topic_Comment_Vote{}).Where("topic_comment_id = ?", tc_id).Select("negative").Scan(&current_negative_vote).Error
		if err != nil {
			return err
		}
		return db.Model(&model.Topic_Comment_Vote{}).Where("topic_comment_id = ?").Update("negative", current_negative_vote-1).Error
	}
	return err
}
func GetCategoriesName() ([]string, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var s []string
	err := db.Model(&model.Product_Category{}).Select("name").Scan(&s).Error
	return s, err
}
func GetForumPostsCount(forum_id int) (int, error) {
	var forum_posts_count int
	forum_posts_count = getForumTopicsCount(forum_id)
	forums_topics_comments_count, err := getForumTopicsCommentsCount(forum_id)
	if err != nil {
		return 0, err
	}
	forum_posts_count += forums_topics_comments_count

	return forum_posts_count, nil
}
func GetDiscussionForumsName(d *model.Discussion) ([]string, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var forums_name []string
	err := db.Model(&model.Forum{}).Where("discussion_id = ?", d.ID).Select("name").Find(&forums_name).Error
	return forums_name, err
}
func GetTopicLastPostInViewModel(topic_id int) (viewmodel.LastPost, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Topic comment count
	comment_count, err := getTopicCommentsCount(topic_id)
	if err != nil {
		return viewmodel.LastPost{}, err
	}
	// Return variable
	var lp viewmodel.LastPost

	var temp_lp viewmodel.LastPostViewModelWithUserID
	// Does topic have any comment
	if comment_count > 0 {
		err := db.Model(&model.Topic_Comment{}).Where("topic_id = ?", topic_id).Select("created_at", "user_id").Order("created_at DESC").Limit(1).Scan(&temp_lp).Error
		if err != nil {
			return viewmodel.LastPost{}, err
		}

		lp.CreatedAt = temp_lp.CreatedAt
		username, err := GetUsernameByUserID(temp_lp.UserID)
		if err != nil {
			return viewmodel.LastPost{}, err
		}
		lp.AuthorUsername = username
	} else {
		err := db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("created_at", "user_id").Scan(&temp_lp).Error
		if err != nil {
			return viewmodel.LastPost{}, err
		}
		lp.CreatedAt = temp_lp.CreatedAt
		username, err := GetUsernameByUserID(temp_lp.UserID)
		if err != nil {
			return viewmodel.LastPost{}, err
		}
		lp.AuthorUsername = username
	}
	return lp, nil
}
func GetDiscussionTopicsCount(d *model.Discussion) (int, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var discussion_forums []model.Forum
	err := db.Select("id").Where("discussion_id = ?", d.BasicModel.ID).Find(&discussion_forums).Error
	if err != nil {
		return 0, err
	}
	var forums_ids []int
	for _, f := range discussion_forums {
		forums_ids = append(forums_ids, int(f.ID))
	}
	var discussion_topics_count int64
	err = db.Model(&model.Topic{}).Where("forum_id IN ?", forums_ids).Count(&discussion_topics_count).Error
	if err != nil {
		return 0, err
	}
	return int(discussion_topics_count), nil
}
func GetDiscussionPostsCount(d *model.Discussion) (int, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var discussion_forums []model.Forum
	err := db.Select("id").Where("discussion_id = ?", d.ID).Find(&discussion_forums).Error
	if err != nil {
		return 0, err
	}
	var discussion_forums_posts_count int
	for _, f := range discussion_forums {
		forum_post_count, err := GetForumPostsCount(int(f.ID))
		if err != nil {
			return 0, err
		}
		discussion_forums_posts_count += forum_post_count
	}
	return discussion_forums_posts_count, nil
}
func GetDiscussionForumsInViewModel(discussion_id int) ([]viewmodel.ForumViewModel, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var forums []viewmodel.ForumViewModel
	err := db.Model(&model.Forum{}).Where("discussion_id = ?", discussion_id).Select("name", "description", "id").Scan(&forums).Error
	if err != nil {
		return nil, err
	}
	for i := range forums {
		p_count, err := GetForumPostsCount(int(forums[i].ID))
		if err != nil {
			return nil, err
		}
		forums[i].PostCount = p_count
		forums[i].TopicCount = getForumTopicsCount(int(forums[i].ID))
	}
	return forums, err
}
func GetProductInViewModel(p_id int) (*viewmodel.ProductViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var p model.Product
	err := db.Model(&model.Product{}).Where("id = ?", p_id).Preload("Tags").Preload("Categories").Preload("Images").Find(&p).Error
	if err != nil {
		return nil, err
	}
	vm := viewmodel.ProductViewModel{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Inventory:   int(p.Inventory),
	}
	var i int
	for i = range p.Images {
		vm.ImagesPath = append(vm.ImagesPath, p.Images[i].Path)
	}
	// Reset i
	i = 0

	var tags_str string
	p_tags_len := len(p.Tags)

	for i = range p.Tags {
		if i == p_tags_len {
			tags_str += p.Tags[i].Name
		} else {
			tags_str += p.Tags[i].Name + "|"
		}
	}
	if tags_str != "" && tags_str != "|" {
		vm.Tags = tags_str
	}
	// Reset i
	i = 0
	for i = range p.Categories {
		vm.SelectedCategory = append(vm.SelectedCategory, p.Categories[i].Name)
	}
	return &vm, nil
}

func GetDiscussionTopics(discussion_id int) ([]viewmodel.TopicBriefViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var topics []viewmodel.TopicBriefViewModel
	discussion_forums_ids, err := getDiscussionForumsIDs(discussion_id)
	if err != nil {
		return nil, err
	}
	err = db.Model(&model.Topic{}).Where("forum_id IN ?", discussion_forums_ids).Order("created_at DESC").Scan(&topics).Error
	if err != nil {
		return nil, err
	}
	for i := range topics {
		comments_count, err := getTopicCommentsCount(int(topics[i].ID))
		if err != nil {
			return nil, err
		}
		topics[i].ReplyCount = comments_count
	}
	return topics, err
}
func GetDiscussionForumsByField(discussion_id int, fields []string) ([]model.Forum, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var forums []model.Forum
	var err error
	if fields != nil {
		err = db.Where("discussion_id = ?", discussion_id).Select(fields).Find(&forums).Error
	} else {
		err = db.Where("discussion_id = ?", discussion_id).Find(&forums).Error
	}
	return forums, err
}
func GetDiscussionTopicsBasedForums(discussion_id int) ([]viewmodel.TopicBriefViewModel, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var topics []viewmodel.TopicBriefViewModel
	err := db.Where("discussion_id = ?", discussion_id).Select("name", "id", "view_count", "created_at").Find(&topics).Error
	for _, t := range topics {
		commentCount, err := getTopicCommentsCount(int(t.ID))
		if err != nil {
			return nil, err
		}
		t.ReplyCount = commentCount
	}
	return topics, err
}

func GetForumTopicsInViewModel(forum_id int) ([]viewmodel.TopicBriefViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Temp topic view model
	var temp_topics_view []viewmodel.TopicForShowTopicViewModelWithUserID

	err := db.Model(&model.Topic{}).Where("forum_id = ?", forum_id).Select("id", "name", "view_count", "created_at", "user_id").Scan(&temp_topics_view).Error
	if err != nil {
		return nil, err
	}
	// Topics view model
	var topics_view []viewmodel.TopicBriefViewModel
	// Fill topics_view variable
	for i := range temp_topics_view {
		var topic viewmodel.TopicBriefViewModel
		topic.ID = temp_topics_view[i].ID
		topic.Name = temp_topics_view[i].Name
		topic.CreatedAt = temp_topics_view[i].CreatedAt
		topic.ViewCount = temp_topics_view[i].ViewCount

		username, err := GetUsernameByUserID(temp_topics_view[i].UserID)
		if err != nil {
			return nil, err
		}
		topic.AuthorUsername = username

		comment_count, err := getTopicCommentsCount(int(topic.ID))
		if err != nil {
			return nil, err
		}
		topic.ReplyCount = comment_count

		last_post, err := GetTopicLastPostInViewModel(int(topic.ID))
		if err != nil {
			return nil, err
		}
		topic.LastPost = last_post

		// Append topic to topics_view
		topics_view = append(topics_view, topic)
	}
	return topics_view, nil
}

func GetTopicByIDForShowTopicInViewModel(topic_id int) (viewmodel.TopicForShowTopicViewModel, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Get topic basic information
	var t viewmodel.TopicBasicInformation
	err := db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("user_id", "name", "description", "created_at").Scan(&t).Error
	if err != nil {
		return viewmodel.TopicForShowTopicViewModel{}, err
	}
	// Get author information
	u, err := getUserInformationByIDForShowTopicInViewModel(t.UserID)
	if err != nil {
		return viewmodel.TopicForShowTopicViewModel{}, err
	}
	// Get topic tags
	var tags_vm []viewmodel.TopicTagBasicInformation
	tags_vm, err = getTopicTagsByTopicIDInViewModel(topic_id)
	if err != nil {
		return viewmodel.TopicForShowTopicViewModel{}, err
	}

	// Fill view model and return it
	var topic_vm = viewmodel.TopicForShowTopicViewModel{
		Title:       t.Name,
		Description: template.HTML(t.Description),
		CreatedAt:   t.CreatedAt,
		UserInfo:    u,
		Tags:        tags_vm,
	}
	return topic_vm, nil
}

func GetTopicCommentsByIDForShowTopicInViewModel(topic_id int) ([]viewmodel.TopicCommentViewModel, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Get Topic comments
	var tc []viewmodel.TopicCommentBasicInformation
	err := db.Model(&model.Topic_Comment{}).Where("topic_id = ?", topic_id).Select("text", "created_at", "user_id", "reply_id").Scan(&tc).Error
	if err != nil {
		return nil, err
	}

	// Fill view model and return it
	var tc_vm = make([]viewmodel.TopicCommentViewModel, len(tc))
	for i := range tc {
		tc_vm[i].Text = tc[i].Text
		tc_vm[i].CreatedAt = tc[i].CreatedAt

		// If topic comment is a reply to another topic comment get that topic comment
		if tc[i].ReplyID != 0 {
			tc_vm[i].Reply, err = getTopicCommentByIDInViewModel(tc[i].ReplyID)
			if err != nil {
				return nil, err
			}
		}
		// Get comment's author information
		u, err := getUserInformationByIDForShowTopicInViewModel(tc[i].UserID)
		if err != nil {
			return nil, err
		}
		tc_vm[i].UserInfo = u
	}
	return tc_vm, nil
}

func FirstOrCreate[m Model](model *m) error {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	err := db.FirstOrCreate(&model).Error
	return err
}
func FirstOrCreateProductTagByName(name string) (*model.Product_Tag, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var t model.Product_Tag
	err := GetFieldsByAnotherFieldValue(&t, []string{"id"}, "name", name)
	if err != nil {
		return nil, err
	}
	if t.ID != 0 {
		err := db.Model(&t).Where("id = ?", t.ID).First(&t).Error
		if err != nil {
			return nil, err
		}
		return &t, err
	} else {
		t.Name = name
		err = Add(&t)
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
}
func FirstOrCreateTopicTagByName(name string) (*model.Topic_Tag, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var t model.Topic_Tag
	err := GetFieldsByAnotherFieldValue(&t, []string{"id"}, "name", name)
	if err != nil {
		return nil, err
	}
	// Topic tag is exists
	if t.ID != 0 {
		err := db.Model(&t).Where("id = ?", t.ID).First(&t).Error
		if err != nil {
			return nil, err
		}
		return &t, err
	} else {
		t.Name = name
		err = Add(&t)
		if err != nil {
			return nil, err
		}
		return &t, err
	}
}
func GetTopicByIDForEdit(topic_id int) (*viewmodel.TopicForEditViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Result
	var t viewmodel.TopicForEditViewModel

	// Get topic name,description
	err := db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("name", "description").Scan(&t).Error
	if err != nil {
		return nil, err
	}
	// Get topic tags
	t.Tags, err = getTopicTagsByTopicIDInViewModel(topic_id)
	if err != nil {
		return nil, err
	}
	// Get topic forum name
	forum_name, err := getTopicForumNameByTopicID(topic_id)
	if err != nil {
		return nil, err
	}
	t.ForumName = forum_name
	return &t, err
}

func GetTopicForumIDByTopicID(topic_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var f_id int
	err := db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("forum_id").Scan(&f_id).Error
	return f_id, err
}
func GetTopicNameByTopicID(topic_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var topic_name string
	err := db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("name").Scan(&topic_name).Error
	return topic_name, err
}

func GetOrderStatusByOrderID(order_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var order_status_id int
	err := db.Model(&model.Order{}).Where("id = ?", order_id).Select("order_status_id").Scan(&order_status_id).Error
	if err != nil {
		return "", err
	}
	var s string
	s, err = getOrderStatus_StatusByOrderStatusID(order_status_id)
	return s, err
}

func GetProductBasicInfoByID(product_id int) (*viewmodel.ProductBasicViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var p model.Product
	err := db.Model(&model.Product{}).Where("id = ?", product_id).Preload("Images").Select("name", "price").Find(&p).Error
	if err != nil {
		return nil, err
	}
	var p_vm viewmodel.ProductBasicViewModel
	p_vm.ID = product_id
	p_vm.Name = p.Name
	p_vm.Price = p.Price
	if p.Images != nil {
		p_vm.ImagePath = p.Images[0].Path
	}
	return &p_vm, nil
}
func AddProductComment(user_id int, p_id int, text string) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	return Add(&model.Product_Comment{Text: text, UserID: uint(user_id), ProductID: uint(p_id)})

}
func IncreaseCartItemQuantity(cart_item_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var current_quantity int
	err := db.Model(&model.CartItem{}).Where("id = ?", cart_item_id).Select("product_quantity").Scan(&current_quantity).Error
	if err != nil {
		return err
	}
	return db.Model(&model.CartItem{}).Where("id = ?", cart_item_id).Update("product_quantity", current_quantity+1).Error
}
func DecreaseCartItemQuantity(cart_item_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var current_quantity int
	err := db.Model(&model.CartItem{}).Where("id = ?", cart_item_id).Select("product_quantity").Scan(&current_quantity).Error
	if err != nil {
		return err
	}
	return db.Model(&model.CartItem{}).Where("id = ?", cart_item_id).Update("product_quantity", current_quantity-1).Error
}
func DeleteCartItem(cart_item_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	err := db.Unscoped().Delete(&model.CartItem{}, cart_item_id).Error
	return err
}
func GetCartIDByUserID(user_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var cart_id int
	err := db.Model(&model.Cart{}).Where("user_id = ? AND is_ordered = FALSE", user_id).Select("id").Scan(&cart_id).Error
	return cart_id, err
}
func AddToWishlist(user_id, p_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	wishlist_id, err := getWishlistIDByUserID(user_id)
	if err != nil {
		return err
	}
	already_exists, err := isProductInUserWishlist(wishlist_id, p_id)
	if err != nil {
		return err
	}
	if already_exists {
		return nil
	}
	// Insert into product-wishlist relation table (many to many)
	p_w_data := struct {
		WishlistID int
		ProductID  int
	}{WishlistID: wishlist_id, ProductID: p_id}
	err = db.Table("product_wishlist_m2m").Create(&p_w_data).Error
	return err
}
func GetUserCart(user_id int) (*viewmodel.Cart, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Get cart id for getting cart items
	cart_id, err := GetCartIDByUserID(user_id)
	if err != nil {
		return nil, err
	}
	// Get cart items
	var cart_items []model.CartItem
	err = db.Model(&model.CartItem{}).Where("cart_id = ?", cart_id).Find(&cart_items).Error
	if err != nil {
		return nil, err
	}
	// Create result(viewdata) and return it
	var cart_vm viewmodel.Cart
	for i := range cart_items {
		// Get item(product info)
		p, err := getProductInfoForCartItem(int(cart_items[i].ProductID))
		if err != nil {
			return nil, err
		}
		p.Quantity = int(cart_items[i].ProductQuantity)
		// Calculate total price of cart item
		item_total_price := float64(p.Price) * float64(p.Quantity)
		cart_vm.CartItems = append(cart_vm.CartItems, viewmodel.CartItem{ID: int(cart_items[i].ID), ProductID: p.ID, Name: p.Name, Price: p.Price, ImagePath: p.ImagePath, Quantity: p.Quantity, TotalPrice: item_total_price})
	}
	// Calculate total price of cart
	for i := range cart_vm.CartItems {
		cart_vm.TotalPrice += cart_vm.CartItems[i].TotalPrice
	}
	return &cart_vm, err
}

// Overview tab
func GetUserDataForUserPanel_Overview_FrontPage(user_id int) (*viewmodel.UserPanel_Overview_Front, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var joined_at time.Time
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("created_at").Scan(&joined_at).Error
	if err != nil {
		return nil, err
	}

	activity, err := getUserActivity(user_id)
	if err != nil {
		return nil, err
	}
	total_posts, err := getUserPostCount(user_id)
	if err != nil {
		return nil, err
	}
	total_products, err := getUserProductsCount(user_id)
	if err != nil {
		return nil, err
	}
	var total_polls int
	total_polls, err = getUserPollsCount(user_id)
	if err != nil {
		return nil, err
	}

	return &viewmodel.UserPanel_Overview_Front{
		JoinedAt:              &joined_at,
		LastLoginAt:           activity.LastLoginAt,
		LastPasswordChangedAt: activity.LastChangePasswordAt,
		LastBuyAt:             activity.LastBuyAt,
		TotalPosts:            total_posts,
		TotalProducts:         total_products,
		TotalPolls:            total_polls,
	}, nil
}
func GetUserDataForUserPanel_Overview_Orders(user_id int) ([]viewmodel.UserPanel_Overview_LastBuy, error) {
	orders, err := getUserOrders(user_id)
	if err != nil {
		return nil, err
	}
	var last_buy_vm []viewmodel.UserPanel_Overview_LastBuy
	for order_i := range orders {
		var d_vm viewmodel.UserPanel_Overview_LastBuy
		d_vm.Time = &orders[order_i].CreatedAt
		d_vm.TotalPrice = (orders[order_i].Cart.TotalPrice)
		order_status, err := GetOrderStatusByOrderID(int(orders[order_i].ID))
		if err != nil {
			return nil, err
		}
		d_vm.OrderStatus = order_status
		d_vm.OrderID = int(orders[order_i].ID)
		last_buy_vm = append(last_buy_vm, d_vm)
	}
	return last_buy_vm, nil
}
func GetUserDataForUserPanel_Overview_Logins(user_id int) ([]viewmodel.UserPanel_Overview_Login, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var activity_logins model.Activity
	err := db.Model(&model.Activity{}).Where("user_id = ?", user_id).Select("logins_at").First(&activity_logins).Error
	if err != nil {
		return nil, err
	}
	// Check if there is more than one logintime in database (logins time seperated with "|")
	if strings.Contains(activity_logins.LoginsAt, "|") {
		// Seperate times and get a list of them in string
		str_times := strings.Split(activity_logins.LoginsAt, "|")
		vm_l := make([]viewmodel.UserPanel_Overview_Login, len(str_times))
		for i := range str_times {
			t, err := time.Parse(time.RFC3339, str_times[i])
			if err != nil {
				return nil, err
			}
			vm_l[i].Login_At = &t
		}
		return vm_l, nil
	} else {
		t, err := time.Parse(time.RFC3339, activity_logins.LoginsAt)
		if err != nil {
			return nil, err
		}
		return []viewmodel.UserPanel_Overview_Login{{Login_At: &t}}, nil
	}
}

// Profile tab
// Should be complete
func GetUserDataForUserPanel_Profile_EditAvatar(user_id int) (*viewmodel.UserPanel_Overview_EditAvatar, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var a_path string
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("avatar_path").Scan(&a_path).Error
	return &viewmodel.UserPanel_Overview_EditAvatar{AvatarPath: a_path}, err
}
func GetUserDataForUserPanel_Profile_EditAccount(user_id int) (*viewmodel.UserPanel_Profile_EditAccount, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var email string
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("email").Scan(&email).Error
	return &viewmodel.UserPanel_Profile_EditAccount{Email: email}, err
}
func GetUserDataForUserPanel_Profile_ManageAddress(user_id int) (*viewmodel.UserPanel_Profile_ManageAddress, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var address model.Address
	err := db.Model(address).Where("user_id = ?", user_id).Find(&address).Error
	if err != nil {
		return nil, err
	}
	return &viewmodel.UserPanel_Profile_ManageAddress{Name: address.Name, Country: address.Country, Province: address.Province, City: address.City, Street: address.Street, BuildingNumber: address.BuildingNumber, PostalCode: address.PostalCode, Description: address.Description}, nil
}
func GetUserDataForUserPanel_Profile_ManageWallet(user_id int) ([]viewmodel.UserPanel_Profile_ManageWallet, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var wallets []model.WalletInfo
	err := db.Model(&model.WalletInfo{}).Where("user_id = ?", user_id).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	var wallets_vm = make([]viewmodel.UserPanel_Profile_ManageWallet, 0, len(wallets))
	for i := range wallets {
		wallets_vm = append(wallets_vm, viewmodel.UserPanel_Profile_ManageWallet{
			WalletAddr: wallets[i].Addr,
		})
	}
	return wallets_vm, nil
}
func GetUserDataForUserPanel_Profile_EditSignature(user_id int) (*viewmodel.UserPanel_Profile_EditSignature, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var signature string
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("signature").Scan(&signature).Error
	return &viewmodel.UserPanel_Profile_EditSignature{Signature: signature}, err
}

// Product tab
func GetUserDataForUserPanel_Product_FrontPage(user_id int) {
}

// Payment tab
func GetUserDataForUserPanel_Payment_FrontPage(user_id int) {
}

// Topic tab
func GetUserDataForUserPanel_Topic_FrontPage(user_id int) {
}

// Poll tab
func GetUserDataForUserPanel_Poll_FrontPage(user_id int) {
}
