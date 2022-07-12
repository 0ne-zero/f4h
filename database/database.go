package database

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/utilities/functions/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitializeOrGetDB() (*gorm.DB, error) {
	if db == nil {
		// Get DSN from setting file
		// DSN = Data source name (like connection string for database)
		dsn, err := setting.ReadFieldInSettingData("DSN")
		if err != nil {
			return nil, err
		}

		// Open connection to database
		var db *gorm.DB
		try_again := true
		for try_again {
			// Connect to database with gorm
			db, err = gorm.Open(
				// Open Databse
				mysql.New(mysql.Config{DSN: dsn}),
				// Config GORM
				&gorm.Config{
					// Allow create tables with null foreignkey
					DisableForeignKeyConstraintWhenMigrating: true,
					// All Datetime in database is in UTC
					NowFunc:              func() time.Time { return time.Now().UTC() },
					FullSaveAssociations: true,
				})

			if err != nil {
				// If databse doesn't exists, so we have to create the database
				if strings.Contains(err.Error(), "Unknown database") {
					err = CreateDatabaseFromDSN(dsn)
					if err != nil {
						fmt.Println(fmt.Sprintf("Mentioned database in dsn isn't created,program tried to create that database but it can't do that.\nError: %s", err.Error()))
						os.Exit(1)
					}
					// We don't need to set try_again to True, its default value
					// try_again = true
				} else {
					return nil, err
				}
			}
			// We don't need to try again to connect to database because we are connected
			try_again = false
		}
		db.Set("gorm:auto_preload", true)
		return db, nil
	} else {
		return db, nil
	}

}
func MigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Address{},
		&model.WalletInfo{},
		&model.Activity{},
		&model.Role{},
		&model.Product{},
		&model.Product_Category{},
		&model.Product_Tag{},
		&model.Product_Image{},
		&model.Product_Comment{},
		&model.Order{},
		&model.CartItem{},
		&model.OrderStatus{},
		&model.Cart{},
		&model.Product_Comment_Vote{},
		&model.Topic_Comment_Vote{},
		&model.Poll_Vote{},
		&model.Topic_Vote{},
		&model.Forum{},
		&model.Topic{},
		&model.Topic_Tag{},
		&model.Discussion_Category{},
		&model.Topic_Comment{},
		&model.Discussion{},
		&model.Poll{},
		&model.Request{},
	)
}
func CreateEssentialData(db *gorm.DB) {
	var (
		user           []model.User
		order_statuses []model.OrderStatus
		categories     []model.Product_Category
		roles          []model.Role
		products       []model.Product
	)
	user = []model.User{
		{Username: "Unknown", Email: "Unknown", PasswordHash: "Unknown", IsSeller: false},
	}

	order_statuses = []model.OrderStatus{
		{Status: "Waiting for payment"},
		{Status: "Paid out"},
		{Status: "Waiting for send"},
		{Status: "Sending"},
		{Status: "Delivered"},
	}

	categories = []model.Product_Category{
		{Name: "Art", Description: "Selling Arts", SubCategories: []model.Product_Category{{Name: "Music", UserID: 1, Description: "msuics"}, {Name: "Theater", Description: "Theater"}}},
		{Name: "Drugs", Description: "Selling Drugs"},
		{Name: "Cannabis", Description: "Selling Cannabis"},
		{Name: "Dissociative", Description: "Selling Dissociative"},
		{Name: "Ecstasy", Description: "Selling Ecstasy"},
		{Name: "Opioid", Description: "Selling Opioids"},
	}
	products = []model.Product{
		{Name: "xx", Price: 234, UserID: 1, Description: "ssdf", Inventory: 434324},
		{Name: "xx11111", Price: 234, UserID: 1, Description: "ssdf", Inventory: 43},
	}
	products[1].Categories = append(products[1].Categories, &categories[0])
	products[0].Categories = append(products[0].Categories, &categories[0])
	products[0].Categories = append(products[0].Categories, &categories[2])
	products[0].Categories = append(products[0].Categories, &categories[3])
	roles = []model.Role{
		{Name: "Member", Description: "Member users"},
		{Name: "Maintainer", Description: "Maintainers users"},
		{Name: "Admin", Description: "Administrators users"},
	}
	db.Create(&user)
	db.Create(&order_statuses)
	db.Create(&categories)
	db.Create(&products)
	db.Create(&roles)
}
func CreateTestData(db *gorm.DB) {
	// user
	// role
	db.Create(&model.User{Username: "admin", Email: "admin@gmail.com", PasswordHash: "x", IsSeller: false, Roles: []*model.Role{{Name: "xxxx", Description: "xssxx"}, {Name: "xxssssssssxxx", Description: ""}}})
	// activity
	utc_now_time := time.Now().UTC()
	db.Create(&model.Activity{LastLoginAt: &utc_now_time, LastBuyAt: nil, LastChangePasswordAt: nil, UserID: 1})
	// product
	// inventory
	// tag
	// image path
	db.Create(&model.Product{Name: "drug", Description: "x", Price: 114325, UserID: 1, Categories: []*model.Product_Category{{Name: "XXXX", Description: "xxxx"}, {Name: "XXxxxxXX", Description: "xxxxxxx"}}, Tags: []*model.Product_Tag{{Name: "sssssssssssss", Description: "ssssssssssss"}, {Name: "sssss", Description: "sssss"}}, Inventory: 21232223, Images: []*model.Product_Image{{Path: "xsfsdfsadfd"}}})
	// address
	db.Create(&model.Address{Name: "ssasdgss", Country: "ssss", City: "sdfasdfd", Street: "sasdgs", BuildingNumber: "sdfasdf", PostalCode: "sfaffas", Description: "sadfsdds", UserID: 1})
	db.Create(&model.Address{Name: "ssgsgsasss", Country: "ssdfsdfxcxcvbdfshsdfgsadfsss", City: "sasdfasdfasdfasdfd", Street: "sasdfsdsafasfdsafaasdgs", BuildingNumber: "sadsfsafsadfasdf", PostalCode: "ssadfasdffaffas", Description: "sadfsssadsadfasffdds", UserID: 1})
	// wallet
	db.Create(&model.WalletInfo{Name: "xxx", Addr: "xxx", IsDefault: true, UserID: 1})
	db.Create(&model.WalletInfo{Name: "xxxx", Addr: "xxxx", IsDefault: false, UserID: 1, OrderID: 1})
	db.Create(&model.WalletInfo{Name: "xxxx", Addr: "xxxx", IsDefault: true, UserID: 1, OrderID: 1})

	// Discussion category
	db.Create(&model.Discussion_Category{Name: "Bitcoin", Description: "Bitcoin Forum general discussion about the Bitcoin ecosystem that doesn't fit better elsewhere. The Bitcoin community, innovations, the general environment, etc. Discussion of specific Bitcoin-related services usually belongs in other sections.", UserID: 1})
	db.Create(&model.Discussion_Category{Name: "Bitcoin", Description: "Bitcoin Forum general discussion about the Bitcoin ecosystem that doesn't fit better elsewhere. The Bitcoin community, innovations, the general environment, etc. Discussion of specific Bitcoin-related services usually belongs in other sections.", UserID: 1})

	db.Create(&model.Discussion{Name: "Drug", Description: "drugs things", UserID: 1, Categories: []*model.Discussion_Category{{BasicModel: model.BasicModel{ID: 1}}}})
	db.Create(&model.Discussion{Name: "Art", Description: "Art things", UserID: 1, Categories: []*model.Discussion_Category{{BasicModel: model.BasicModel{ID: 1}}}})
	db.Create(&model.Discussion{Name: "Other", Description: "Other things", UserID: 1, Categories: []*model.Discussion_Category{{BasicModel: model.BasicModel{ID: 1}}}})

	db.Create(&model.Forum{Name: "Bitcoin", Description: "Bitcoin Forum general discussion about the Bitcoin ecosystem that doesn't fit better elsewhere. The Bitcoin community, innovations, the general environment, etc. Discussion of specific Bitcoin-related services usually belongs in other sections.", DiscussionID: 1, UserID: 1})
	db.Create(&model.Forum{Name: "Bitcoin", Description: "Bitcoin Forum general discussion about the Bitcoin ecosystem that doesn't fit better elsewhere. The Bitcoin community, innovations, the general environment, etc. Discussion of specific Bitcoin-related services usually belongs in other sections.", DiscussionID: 1, UserID: 1})

	// forum
	db.Create(&model.Forum{Name: "blockchain", Description: "some description for blockchain", UserID: 1, DiscussionID: 1})
	db.Create(&model.Forum{Name: "bitcoin", Description: "some description", UserID: 1, DiscussionID: 1})
	db.Create(&model.Forum{Name: "sub", Description: "sub description", DiscussionID: 2, UserID: 1})
	db.Create(&model.Forum{Name: "sub", Description: "sub description", DiscussionID: 3, UserID: 1})
	db.Create(&model.Forum{Name: "sub", Description: "sub description", DiscussionID: 3, UserID: 1})

	// topics

	db.Create(&model.Topic{Name: "topic", Description: "some description", ForumID: 1, UserID: 1})
	db.Create(&model.Topic{Name: "topic", Description: "some description", ForumID: 1, UserID: 1})
	db.Create(&model.Topic{Name: "topic", Description: "some description", ForumID: 2, UserID: 1})
	db.Create(&model.Topic{Name: "topic", Description: "some description", ForumID: 2, UserID: 1})
	db.Create(&model.Topic{Name: "topic", Description: "some description", ForumID: 2, UserID: 1})
	db.Create(&model.Topic{Name: "topic", Description: "some description", ForumID: 3, UserID: 1})
	db.Create(&model.Topic{Name: "topic", Description: "some description", ForumID: 3, UserID: 1})
	db.Create(&model.Topic{Name: "topic", Description: "some description", ForumID: 3, UserID: 1})
	db.Create(&model.Topic{Name: "topic", Description: "some description", ForumID: 3, UserID: 1})

	// cart
	db.Create(&model.Cart{TotalPrice: 2322, IsOrdered: true, UserID: 1})
	// order status
	db.Create(&model.OrderStatus{Status: "xxxxxxxxxx"})
	db.Create(&model.OrderStatus{Status: "xxxxxxxxxxx"})
	db.Create(&model.OrderStatus{Status: "xxxxxxxxxxxx"})
	// order
	db.Create(&model.Order{SenderWalletInfoID: 1, UserID: 1, OrderStatusID: 1, CartID: 1})
	// order item
	db.Create(&model.CartItem{ProductID: 1, CartID: 1, ProductQuantity: 2})
	db.Create(&model.CartItem{ProductID: 1, CartID: 1, ProductQuantity: 2})
	db.Create(&model.CartItem{ProductID: 1, CartID: 1, ProductQuantity: 2})
	// poll
	db.Create(&model.Poll{Name: "xxxx", Description: "xxxx", UserID: 1})
	// comment
	db.Create(&model.Topic_Comment{Text: "xxxxxxxx", UserID: 1, TopicID: 1})
	db.Create(&model.Topic_Comment{Text: "xxxxxxxx", UserID: 1, TopicID: 1})
	db.Create(&model.Topic_Comment{Text: "xxxxxxxx", UserID: 1, TopicID: 1})
	// vote
	db.Create(&model.Topic_Comment_Vote{UserID: 1, Positive: 12121, Negative: 121, Topic_CommentID: 1})
	db.Create(&model.Topic_Comment_Vote{UserID: 1, Positive: 12121, Negative: 121, Topic_CommentID: 2})
	db.Create(&model.Topic_Comment_Vote{UserID: 1, Positive: 12121, Negative: 121, Topic_CommentID: 3})
	db.Create(&model.Poll_Vote{UserID: 1, PollID: 1, Positive: 20, Negative: 75})
	db.Create(&model.Topic_Vote{UserID: 1, TopicID: 1, Positive: 20523, Negative: 7523})
	db.Create(&model.Topic_Vote{UserID: 1, TopicID: 2, Positive: 2120, Negative: 7545})
	db.Create(&model.Topic_Vote{UserID: 1, TopicID: 3, Positive: 205342, Negative: 7534})
	db.Create(&model.Product_Comment_Vote{UserID: 1, Product_CommentID: 1})
	// unknown user
	db.Create(&model.User{Username: "Unknown", Email: "ss", PasswordHash: "sda", IsSeller: false})
}
func AnonymizeUser(db *gorm.DB, user *model.User) {
	var unknown_user_ID uint
	db.Unscoped().Select("ID").Where("username = ?", "Unknown").First(&unknown_user_ID)

	//region Anonymize data
	// Anonymize user.Comments
	for _, comment := range user.Comments {
		comment.UserID = unknown_user_ID
		db.Save(&comment)
	}
	// Anonymize user.Polls
	for _, poll := range user.Polls {
		poll.UserID = unknown_user_ID
		db.Save(&poll)
	}
	// Anonymize user.Votes
	// for _, vote := range user.Votes {
	// 	vote.UserID = unknown_user_ID
	// 	db.Save(&vote)
	// }
	// Anonymize user.Orders
	for _, order := range user.Orders {
		order.UserID = unknown_user_ID
		db.Save((&order))
	}
	// Anonymize user.Carts
	for _, cart := range user.Carts {
		cart.UserID = unknown_user_ID
		db.Save(&cart)
	}
	//endregion

	//region Delete data
	for _, wallet := range user.WalletInfos {
		db.Unscoped().Delete(&wallet)
	}
	// If user.Activity exists, delete it
	if user.Activity.ID != 0 {
		db.Unscoped().Delete(&user.Activity)
	}

	for _, product := range user.Products {
		db.Unscoped().Delete(&product)
	}
	for _, role := range user.Roles {
		db.Unscoped().Delete(&role)
	}
	for _, address := range user.Addresses {
		db.Unscoped().Delete(&address)
	}
	db.Unscoped().Delete(&user, user.ID)
	//endregion
}
func GetDatabaseNameFromDSN(dsn string) string {
	before, _, _ := strings.Cut(strings.Split(dsn, "/")[1], "?")
	return before
}

func CreateDatabaseFromDSN(dsn string) error {
	// Create database
	dsn_without_database := strings.Split(dsn, "/")[0] + "/"
	db, err := sql.Open("mysql", dsn_without_database)
	if err != nil {
		if !StartMySqlService() {
			fmt.Println(fmt.Sprintf("We can't connect to mysql and we can't even start mysql.service\nError: %s", err.Error()))
			os.Exit(1)
		}
		db, err = sql.Open("mysql", dsn_without_database)
		if err != nil {
			fmt.Println(fmt.Sprintf("mysql.service is in start mode, but for any reason we can't connect to database\nError: %s", err.Error()))
			os.Exit(1)
		}
	}
	db_name := GetDatabaseNameFromDSN(dsn)
	_, err = db.Exec("CREATE DATABASE " + db_name)
	return err
}
func StartMySqlService() bool {
	command := fmt.Sprintf("systemctl start mysql.service")
	_, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return false
	}
	return true
}
