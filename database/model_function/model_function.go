package model_function

import (
	"errors"
	"fmt"
	"html/template"
	"sort"
	"time"

	"github.com/0ne-zero/f4h/database"
	"github.com/0ne-zero/f4h/database/model"
	viewmodel "github.com/0ne-zero/f4h/public_struct/view_model"
	general_func "github.com/0ne-zero/f4h/utilities/functions/general"
)

type Model interface {
	model.Forum | model.Discussion | model.User | model.Product_Category | model.Product | model.Request | model.Discussion_Category | model.BadRequest | model.Topic | model.Topic_Tag
}

func Add[m Model](model *m) error {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return err
	}
	return db.Create(model).Error
}

// Input:
// consider_model = model with its id
// update_model = the model with some change, you wish to apply the consider model
// Returns Changed model (result)
func Update[m Model](consider_model *m, updated_model *m) (*m, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	err = db.Model(consider_model).Updates(updated_model).Error
	if err != nil {
		return nil, err
	}
	//err = db.Save(consider_model).Error

	return consider_model, err
}
func Get[m Model](model *[]m, limit int, orderBy string, orderType string, preloads ...string) error {

	context, err := database.InitializeOrGetDB()
	if err != nil {
		return err
	}
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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return false, err
	}
	var exists bool
	err = db.Model(model).Select("count(*) >0").Where("ID = ?", id).Find(&exists).Error
	return exists, err
}
func GetByID[m Model](model *m, id int) error {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).First(model).Error
}
func GetUserPassHashByUsername(username string) (string, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return "", err
	}
	var pass_hash string
	err = db.Where("username = ?", username).Select("password_hash").First(&pass_hash).Error
	return pass_hash, err
}
func IsUserExistsByUsername(username string) (bool, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return false, err
	}
	var exists bool
	err = db.Model(&model.User{}).Select("count(*) >0").Where("username = ?", username).Find(&exists).Error
	return exists, err
}
func GetUserByUsername(username string) (model.User, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return model.User{}, err
	}
	var u model.User
	err = db.Where("username = ?", username).First(&u).Error
	return u, err
}
func GetUsernameByUserID(user_id int) (string, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return "", err
	}
	var username string
	err = db.Model(&model.User{}).Where("id = ?", user_id).Select("username").Scan(&username).Error
	return username, err
}
func GetFieldsByAnotherFieldValue[m Model](model *m, out_fields_name []string, in_field_name string, in_field_value string) error {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return err
	}
	err = db.Model(model).Where(fmt.Sprintf("%s = ?", in_field_name), in_field_value).Select(out_fields_name).Scan(model).Error
	return err
}
func TooManyRequest(ip string, url string, method string) (bool, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return false, err
	}
	var req_count int64
	now := time.Now().UTC()
	one_hour_ago := now.Add(time.Duration(-1) * time.Hour)
	err = db.Model(&model.Request{}).Where("ip = ? AND url = ? AND method = ? AND time <= ? ", ip, url, method, one_hour_ago).Count(&req_count).Error
	if err != nil {
		return false, err
	}
	if req_count > 100 {
		return true, nil
	}
	return false, nil
}
func GetProductInViewModel(limit int) ([]viewmodel.ProductBasicViewModel, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var products []viewmodel.ProductBasicViewModel
	if limit > 0 {
		err = db.Model(&model.Product{}).Limit(limit).Select("id", "name", "price").Scan(&products).Error
	} else {
		err = db.Model(&model.Product{}).Select("id", "name", "price").Scan(&products).Error
	}
	return products, err
}
func GetProductByCategoryInViewModel(category_name string, limit int) ([]viewmodel.ProductBasicViewModel, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var products []viewmodel.ProductBasicViewModel
	var c model.Product_Category

	err = db.Preload("Products").Where("name = ?", category_name).Find(&c).Error
	if err != nil {
		return nil, err
	} else if c.Products == nil {
		return nil, errors.New("Products field is empty")
	}
	for _, p := range c.Products {
		products = append(products, viewmodel.ProductBasicViewModel{ID: int(p.ID), Name: p.Name, Price: p.Price})
	}
	return products, nil
}
func GetCategoryByOrderingProductsCount(c *[]model.Product_Category) error {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return err
	}
	// Get categories
	var categories []model.Product_Category
	err = db.Preload("Products").Preload("SubCategories").Find(&categories).Error
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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var categories []model.Product_Category

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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var s []string
	err = db.Model(&model.Product_Category{}).Select("name").Scan(&s).Error
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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var forums_name []string
	err = db.Model(&model.Forum{}).Where("discussion_id = ?", d.ID).Select("name").Find(&forums_name).Error
	return forums_name, err
}
func GetTopicLastPostInViewModel(topic_id int) (viewmodel.LastPost, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return viewmodel.LastPost{}, err
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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return 0, err
	}
	var discussion_forums []model.Forum
	err = db.Select("id").Where("discussion_id = ?", d.BasicModel.ID).Find(&discussion_forums).Error
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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return 0, err
	}
	var discussion_forums []model.Forum
	err = db.Select("id").Where("discussion_id = ?", d.ID).Find(&discussion_forums).Error
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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var forums []viewmodel.ForumViewModel
	err = db.Model(&model.Forum{}).Where("discussion_id = ?", discussion_id).Select("name", "description", "id").Scan(&forums).Error
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
func GetDiscussionTopics(discussion_id int) ([]viewmodel.TopicBriefViewModel, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var forums []model.Forum
	if fields != nil {
		err = db.Where("discussion_id = ?", discussion_id).Select(fields).Find(&forums).Error
	} else {
		err = db.Where("discussion_id = ?", discussion_id).Find(&forums).Error
	}
	return forums, err
}
func GetDiscussionTopicsBasedForums(discussion_id int) ([]viewmodel.TopicBriefViewModel, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var topics []viewmodel.TopicBriefViewModel
	err = db.Where("discussion_id = ?", discussion_id).Select("name", "id", "view_count", "created_at").Find(&topics).Error
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
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	// Temp topic view model
	var temp_topics_view []viewmodel.TopicForShowTopicViewModelWithUserID

	err = db.Model(&model.Topic{}).Where("forum_id = ?", forum_id).Select("id", "name", "view_count", "created_at", "user_id").Scan(&temp_topics_view).Error
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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return viewmodel.TopicForShowTopicViewModel{}, err
	}
	// Get topic basic information
	var t viewmodel.TopicBasicInformation
	err = db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("user_id", "name", "description", "created_at").Scan(&t).Error
	if err != nil {
		return viewmodel.TopicForShowTopicViewModel{}, err
	}
	// Get author information
	u, err := getUserInformationByIDForShowTopicInViewModel(t.UserID)
	if err != nil {
		return viewmodel.TopicForShowTopicViewModel{}, err
	}
	// Get topic tags
	tags_vm, err := getTopicTagsByTopicIDInViewModel(topic_id)

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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	// Get Topic comments
	var tc []viewmodel.TopicCommentBasicInformation
	err = db.Model(&model.Topic_Comment{}).Where("topic_id = ?", topic_id).Select("text", "created_at", "user_id", "reply_id").Scan(&tc).Error
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

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&model).Error
	return err
}
func FirstOrCreateTopicTagByName(name string) (*model.Topic_Tag, error) {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var t model.Topic_Tag
	err = GetFieldsByAnotherFieldValue(&t, []string{"id"}, "name", name)
	if err != nil {
		return nil, err
	}
	// Topic tag is exists
	if t.ID != 0 {
		err = db.Model(&t).Where("id = ?", t.ID).First(&t).Error
		return &t, err
	} else {
		t.Name = name
		err = db.Create(&t).Error
		return &t, err
	}
}
func GetTopicByIDForEdit(topic_id int) (*viewmodel.TopicForEditViewModel, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	// Result
	var t viewmodel.TopicForEditViewModel

	// Get topic name,description
	err = db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("name", "description").Scan(&t).Error
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
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return -1, err
	}
	var f_id int
	err = db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("forum_id").Scan(&f_id).Error
	return f_id, err
}
func GetTopicNameByTopicID(topic_id int) (string, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return "", err
	}
	var topic_name string
	err = db.Model(&model.Topic{}).Where("id = ?", topic_id).Select("name").Scan(&topic_name).Error
	return topic_name, err
}

func GetUserDataForUserPanel_Overview_FrontPage(user_id int) (*viewmodel.UserPanel_Overview_Front, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var joined_at time.Time
	err = db.Model(&model.User{}).Where("id = ?", user_id).Select("joined_at").Scan(&joined_at).Error
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
	total_polls, err := getUserPollsCount(user_id)

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
		d_vm.TotalPrice = uint(orders[order_i].Cart.TotalPrice)
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

func GetProductBasicInfoByID(product_id int) (*viewmodel.ProductBasicViewModel, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var p_vm viewmodel.ProductBasicViewModel
	err = db.Model(&model.Product{}).Where("id = ?", product_id).Select("name", "price").Find(&p_vm).Error
	return &p_vm, err
}

// func GetUserDataForUserPanel_Overview_Logins(user_id int) ([]viewmodel.UserPanel_Overview_Login, error) {
// 	db, err := database.InitializeOrGetDB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var activity_logins model.Activity
// 	err = db.Model(&model.Activity{}).Where("user_id = ?", user_id).Select("logins_at").First(&activity_logins).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	var vm_l []viewmodel.UserPanel_Overview_Login
// 	for i := range activity_logins.Logins {
// 		vm_l = append(vm_l, viewmodel.UserPanel_Overview_Login{Login_At: &activity_logins.Logins[i]})
// 	}
// 	return vm_l, nil
// }

func GetOrderStatusByOrderID(order_id int) (string, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return "", err
	}
	var order_status_id int
	err = db.Model(&model.Order{}).Where("id = ?", order_id).Select("order_status_id").Scan(&order_status_id).Error
	if err != nil {
		return "", err
	}
	var s string
	s, err = getOrderStatus_StatusByOrderStatusID(order_status_id)
	return s, err
}
func GetUserDataForUserPanel_Profile_FrontPage(user_id int) {
}
func GetUserDataForUserPanel_Profile_EditAvatar(user_id int) {
}
func GetUserDataForUserPanel_Profile_ManageLogin(user_id int) {
}
func GetUserDataForUserPanel_Product_FrontPage(user_id int) {
}
func GetUserDataForUserPanel_Payment_FrontPage(user_id int) {
}
func GetUserDataForUserPanel_Topic_FrontPage(user_id int) {
}
func GetUserDataForUserPanel_Poll_FrontPage(user_id int) {
}
func GetUserDataForUserPanel_Profile_EditAccount(user_id int) (*viewmodel.UserPanel_Profile_EditAccount, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var email string
	err = db.Model(&model.User{}).Where("id = ?", user_id).Select("email").Scan(&email).Error
	return &viewmodel.UserPanel_Profile_EditAccount{Email: email}, err
}
func GetUserDataForUserPanel_Profile_ManageAddress(user_id int) (*viewmodel.UserPanel_Profile_ManageAddress, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var address model.Address
	err = db.Model(address).Where("user_id = ?", user_id).Find(&address).Error
	if err != nil {
		return nil, err
	}
	return &viewmodel.UserPanel_Profile_ManageAddress{Name: address.Name, Country: address.Country, Province: address.Province, City: address.City, Street: address.Street, BuildingNumber: address.BuildingNumber, PostalCode: address.PostalCode, Description: address.Description}, nil
}
func GetUserDataForUserPanel_Profile_ManageWallet(user_id int) ([]viewmodel.UserPanel_Profile_ManageWallet, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var wallets []model.WalletInfo
	err = db.Model(&model.WalletInfo{}).Where("user_id = ?", user_id).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	var wallets_vm = make([]viewmodel.UserPanel_Profile_ManageWallet, 0, len(wallets))
	for i := range wallets {
		wallets_vm = append(wallets_vm, viewmodel.UserPanel_Profile_ManageWallet{
			Name:      wallets[i].Name,
			Addr:      wallets[i].Addr,
			IsDefault: wallets[i].IsDefault,
		})
	}
	return wallets_vm, nil
}
func GetUserDataForUserPanel_Profile_EditSignature(user_id int) (*viewmodel.UserPanel_Profile_EditSignature, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var signature string
	err = db.Model(&model.User{}).Where("id = ?", user_id).Select("signature").Scan(&signature).Error
	return &viewmodel.UserPanel_Profile_EditSignature{Signature: signature}, err
}
