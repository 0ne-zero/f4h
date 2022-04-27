package model_function

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/0ne-zero/f4h/constansts"
	"github.com/0ne-zero/f4h/database"
	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/public_struct"
	general_func "github.com/0ne-zero/f4h/utilities/functions/general"
	wrapper_logger "github.com/0ne-zero/f4h/utilities/wrapper_logger"
	viewmodel "github.com/0ne-zero/f4h/web/view_model"
	"gorm.io/gorm"
)

var db *gorm.DB

type Model interface {
	model.Forum | model.Discussion | model.User | model.Product_Category | model.Product | model.Request | model.Discussion_Category | model.BadRequest
}

func init() {
	var err error
	db, err = database.Initialize()
	if err != nil {
		wrapper_logger.Log(&wrapper_logger.FatalLevel{}, "Error when initializing database", &public_struct.ErroredFileInfo{Path: constansts.ExecutableDirectory, Line: 28})
	}
}

func Add[m Model](model *m) error {
	return db.Create(model).Error
}
func Get[m Model](model *[]m, limit int, orderBy string, orderType string, preloads ...string) error {
	var err error
	var context *gorm.DB = db
	if preloads != nil {
		for _, p := range preloads {
			// Include Preload command in db commands chain
			context = context.Preload(p)
		}
	}

	if limit < 1 {
		err = context.Order(fmt.Sprintf("%s %s", orderBy, orderType)).Find(model).Error
	} else {
		err = context.Order(fmt.Sprintf("%s %s", orderBy, orderType)).Find(model).Limit(limit).Error
	}
	return err
}
func IsExistsByID[m Model](model *m, id uint) (bool, error) {
	var exists bool
	err := db.Model(model).Select("count(*) >0").Where("ID = ?", id).Find(&exists).Error
	return exists, err
}
func GetByID[m Model](model *m, id int) error {
	return db.Where("id = ?", id).First(model).Error
}
func GetUserPassHashByUsername(username string) (string, error) {
	var pass_hash string
	err := db.Where("username = ?", username).Select("password_hash").First(&pass_hash).Error
	return pass_hash, err
}
func IsUserExistsByUsername(username string) (bool, error) {
	var exists bool
	err := db.Model(&model.User{}).Select("count(*) >0").Where("username = ?", username).Find(&exists).Error
	return exists, err
}
func GetUserByUsername(username string) (model.User, error) {
	var u model.User
	err := db.Where("username = ?", username).First(&u).Error
	return u, err
}
func GetUsernameByUserID(user_id int) (string, error) {
	var username string
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("username").Scan(&username).Error
	return username, err
}

type x struct {
	password_hash interface{}
}

func GetFieldsByAnotherFieldValue[m Model](model *m, out_fields_name []string, in_field_name string, in_field_value string) error {
	err := db.Model(model).Where(fmt.Sprintf("%s = ?", in_field_name), in_field_value).Select(out_fields_name).Scan(model).Error
	return err
}
func TooManyRequest(ip string, url string, method string) (bool, error) {
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
func GetProductInViewModel(limit int) ([]viewmodel.ProductViewModel, error) {
	var products []viewmodel.ProductViewModel
	var err error
	if limit > 0 {
		err = db.Model(&model.Product{}).Limit(limit).Select("id", "name", "price").Scan(&products).Error
	} else {
		err = db.Model(&model.Product{}).Select("id", "name", "price").Scan(&products).Error
	}
	return products, err
}
func GetProductByCategoryInViewModel(category_name string, limit int) ([]viewmodel.ProductViewModel, error) {
	var products []viewmodel.ProductViewModel
	var c model.Product_Category

	err := db.Preload("Products").Where("name = ?", category_name).Find(&c).Error
	if err != nil {
		return nil, err
	} else if c.Products == nil {
		return nil, errors.New("Products field is empty")
	}
	for _, p := range c.Products {
		products = append(products, viewmodel.ProductViewModel{ID: int(p.ID), Name: p.Name, Price: p.Price})
	}
	return products, nil
}
func GetCategoryByOrderingProductsCount(c *[]model.Product_Category) error {
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

		for _, sub := range c.SubCategories {
			var view_cat_sub viewmodel.SidebarCategoryViewModel
			view_cat_sub.Name = sub.Name
			view_cat.SubCategories = append(view_cat.SubCategories, view_cat_sub)
		}

		result = append(result, view_cat)
	}
	return result, err
}
func GetCategoriesName() ([]string, error) {
	var s []string
	err := db.Model(&model.Product_Category{}).Select("name").Scan(&s).Error
	return s, err
}

func getForumTopicsCount(forum_id int) int {
	return int(db.Model(&model.Forum{BasicModel: model.BasicModel{ID: uint(forum_id)}}).Association("Topics").Count())
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
func getForumTopicsCommentsCount(forum_id int) (int, error) {
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
	// Doesn't works
	//count := db.Model(&model.Topic{}).Where("id = ?", topic_id).Association("Comments").Count()
	var count int
	err := db.Raw("SELECT COUNT(*) FROM topic_comments WHERE topic_id = ?", topic_id).Scan(&count).Error
	return count, err
}
func GetDiscussionForumsName(d *model.Discussion) ([]string, error) {
	var forums_name []string
	err := db.Model(&model.Forum{}).Where("discussion_id = ?", d.ID).Select("name").Find(&forums_name).Error
	return forums_name, err
}

func GetTopicLastPostInViewModel(topic_id int) (viewmodel.LastPost, error) {
	// Temporary data structure
	type last_post_view_model_with_user_id struct {
		UserID int
		viewmodel.LastPost
	}
	// Topic comment count
	comment_count, err := getTopicCommentsCount(topic_id)
	if err != nil {
		return viewmodel.LastPost{}, err
	}
	// Return variable
	var lp viewmodel.LastPost

	var temp_lp last_post_view_model_with_user_id
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
		err = db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("created_at", "user_id").Scan(&temp_lp).Error
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
func GetDiscussionTopics(discussion_id int) ([]viewmodel.TopicViewModel, error) {
	var topics []viewmodel.TopicViewModel
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
	var forums []model.Forum
	var err error
	if fields != nil {
		err = db.Where("discussion_id = ?", discussion_id).Select(fields).Find(&forums).Error
	} else {
		err = db.Where("discussion_id = ?", discussion_id).Find(&forums).Error
	}
	return forums, err
}
func GetDiscussionTopicsBasedForums(discussion_id int) ([]viewmodel.TopicViewModel, error) {
	var topics []viewmodel.TopicViewModel
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
func getDiscussionForumsIDs(discussion_id int) ([]int, error) {
	var IDs []int
	err := db.Model(&model.Forum{}).Where("discussion_id = ?", discussion_id).Select("id").Find(&IDs).Error
	return IDs, err
}

func GetForumTopicsInViewModel(forum_id int) ([]viewmodel.TopicViewModel, error) {
	// Temp struct for send just one request and get all data i need
	type topic_view_model_with_user_id struct {
		UserID int
		viewmodel.TopicViewModel
	}
	// Temp topic view model
	var temp_topics_view []topic_view_model_with_user_id

	err := db.Model(&model.Topic{}).Where("forum_id = ?", forum_id).Select("id", "name", "view_count", "created_at", "user_id").Scan(&temp_topics_view).Error
	if err != nil {
		return nil, err
	}
	// Topics view model
	var topics_view []viewmodel.TopicViewModel
	// Fill topics_view variable
	for i := range temp_topics_view {
		var topic viewmodel.TopicViewModel
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
