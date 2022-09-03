package model_function

import (
	"github.com/0ne-zero/f4h/database"
	"github.com/0ne-zero/f4h/database/model"
	viewmodel "github.com/0ne-zero/f4h/public_struct/view_model"
	general_func "github.com/0ne-zero/f4h/utilities/functions/general"
	"github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"gorm.io/gorm/clause"
)

func getOrderStatus_StatusByOrderStatusID(order_status_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var s string
	err := db.Model(&model.OrderStatus{}).Where("id = ?", order_status_id).Select("status").Scan(&s).Error
	return s, err
}
func getForumTopicsCount(forum_id int) int {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	return int(db.Model(&model.Forum{BasicModel: model.BasicModel{ID: uint(forum_id)}}).Association("Topics").Count())
}
func getForumTopicsCommentsCount(forum_id int) (int, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var forums_topics []model.Topic
	err := db.Where("forum_id = ?", forum_id).Find(&forums_topics).Error
	if err != nil {
		return -1, err
	}
	var forums_topics_comments_count int
	for _, t := range forums_topics {
		comment_count, err := getTopicCommentsCount(int(t.ID))
		if err != nil {
			return -1, err
		}
		forums_topics_comments_count += comment_count
	}
	return forums_topics_comments_count, nil
}
func getTopicCommentsCount(topic_id int) (int, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	// Doesn't works
	//count := db.Model(&model.Topic{}).Where("id = ?", topic_id).Association("Comments").Count()
	var count int
	err := db.Raw("SELECT COUNT(*) FROM topic_comments WHERE topic_id = ?", topic_id).Scan(&count).Error
	return count, err
}

func getUserPostCount(user_id int) (int, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var tc_count int
	err := db.Raw("SELECT COUNT(*) FROM topic_comments WHERE user_id = ?", user_id).Scan(&tc_count).Error
	if err != nil {
		return 0, err
	}
	var t_count int
	err = db.Raw("SELECT COUNT(*) FROM topics WHERE user_id = ?", user_id).Scan(&t_count).Error
	if err != nil {
		return 0, err
	}

	return (tc_count + t_count), err
}
func getTopicTagsByTopicIDInViewModel(topic_id int) ([]viewmodel.TopicTagBasicInformation, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var tags_id []int
	err := db.Table("topic_tag_m2m").Where("topic_id = ?", topic_id).Select("topic_tag_id").Scan(&tags_id).Error
	if err != nil {
		return nil, err
	}
	var tags []viewmodel.TopicTagBasicInformation
	err = db.Model(&model.Topic_Tag{}).Where("id IN ?", tags_id).Select("name").Scan(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func getTopicCommentByIDInViewModel(comment_id int) (*viewmodel.TopicCommentViewModel, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var tc viewmodel.TopicCommentBasicInformation
	// Get Topic Basic information
	err := db.Model(&model.Topic_Comment{}).Where("id = ?", comment_id).Select("text", "created_at", "user_id", "reply_id").Scan(&tc).Error
	if err != nil {
		return nil, err
	}
	// If topic comment is a reply to another topic comment get that topic comment
	var tc_reply *viewmodel.TopicCommentViewModel
	if tc.ReplyID != 0 {
		tc_reply, err = getTopicCommentByIDInViewModel(tc.ReplyID)
		if err != nil {
			return nil, err
		}
	}

	// Get author information
	u, err := getUserInformationByIDForShowTopicInViewModel(tc.UserID)
	if err != nil {
		return nil, err
	}

	// Fill view model and finally return it
	tc_vm := viewmodel.TopicCommentViewModel{
		Text:      tc.Text,
		CreatedAt: tc.CreatedAt,
		UserInfo:  u,
		Reply:     tc_reply,
	}
	return &tc_vm, nil
}

func getTopicForumNameByTopicID(topic_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var forum_id string
	err := db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("forum_id").Scan(&forum_id).Error
	if err != nil {
		return "", err
	}
	var forum_name string
	err = db.Model(&model.Forum{}).Where("id = ?", forum_id).Select("name").Scan(&forum_name).Error
	return forum_name, err
}

func getUserPollsCount(user_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var count int64
	err := db.Model(&model.Poll{}).Where("user_id", user_id).Count(&count).Error
	return int(count), err
}
func getUserOrders(user_id int) ([]model.Order, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var orders []model.Order
	err := db.Model(&model.Order{}).Where("user_id = ?", user_id).Preload(clause.Associations).Preload("Cart.CartItems.Product").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
func getUserProductsCount(user_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var count int64
	err := db.Model(&model.Product{}).Where("user_id = ?", user_id).Count(&count).Error

	return int(count), err
}

func convertProductToBasicProductBasicViewModel(product *model.Product) viewmodel.ProductBasicViewModel {
	if product.Images == nil {
		return viewmodel.ProductBasicViewModel{
			ID:    int(product.ID),
			Name:  product.Name,
			Price: product.Price,
		}
	} else {
		return viewmodel.ProductBasicViewModel{
			ID:        int(product.ID),
			Name:      product.Name,
			Price:     product.Price,
			ImagePath: product.Images[0].Path,
		}
	}
}

func getUserInformationByIDForShowTopicInViewModel(user_id int) (*viewmodel.TopicUserViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var u viewmodel.UserBasicInformation
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("username", "created_at").First(&u).Error
	if err != nil {
		return nil, err
	}
	u_posts_count, err := getUserPostCount(user_id)
	if err != nil {
		return nil, err
	}
	var u_vm = viewmodel.TopicUserViewModel{
		Username:  u.Username,
		JoinedAt:  &u.CreatedAt,
		PostCount: uint(u_posts_count),
	}
	return &u_vm, nil
}

func getDiscussionForumsIDs(discussion_id int) ([]int, error) {

	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var IDs []int
	err := db.Model(&model.Forum{}).Where("discussion_id = ?", discussion_id).Select("id").Find(&IDs).Error
	return IDs, err
}
func getUserActivity(user_id int) (*model.Activity, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "InitializeOrGetDB returns nil db", ErrorLocation: general_func.GetCallerInfo(1)})
	}
	var a model.Activity
	err := db.Model(&model.Activity{}).Where("id = ?", user_id).First(&a).Error
	return &a, err
}
