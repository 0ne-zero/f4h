package database

import (
	"time"

	"github.com/0ne-zero/f4h/config/constansts"
	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/utilities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Initialize() (*gorm.DB, error) {

	// Get DSN from setting file
	// DSN = Data source name (like connection string for database)
	dsn, err := utilities.ReadFieldInSettingFile(constansts.SettingFilePath, "DSN")
	if err != nil {
		return nil, err
	}

	// Open connection to database
	db, err := gorm.Open(
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
		return nil, err
	}
	db.Set("gorm:auto_preload", true)
	return db, nil
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
		&model.OrderItem{},
		&model.OrderStatus{},
		&model.Cart{},
		&model.Product_Comment_Vote{},
		&model.Topic_Comment_Vote{},
		&model.Poll_Vote{},
		&model.Topic_Vote{},
		&model.Forum{},
		&model.Topic{},
		&model.Topic_Tag{},
		&model.Forum_Category{},
		&model.Topic_Comment{},
		&model.Topic_Category{},
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
		{Name: "Art", Description: "Selling Arts"},
		{Name: "Drugs", Description: "Selling Drugs"},
		{Name: "Cannabis", Description: "Selling Cannabis", ParentCategoryID: 2},
		{Name: "Dissociative", Description: "Selling Dissociative"},
		{Name: "Ecstasy", Description: "Selling Ecstasy"},
		{Name: "Opioid", Description: "Selling Opioids"},
	}
	roles = []model.Role{
		{Name: "Member", Description: "Member users"},
		{Name: "Maintainer", Description: "Maintainers users"},
		{Name: "Admin", Description: "Administrators users"},
	}
	db.Create(&user)
	db.Create(&order_statuses)
	db.Create(&categories)
	db.Create(&roles)
}
func CreateTestData(db *gorm.DB) {
	// user
	// role
	db.Create(&model.User{Username: "x", Email: "x", PasswordHash: "x", IsSeller: false, Roles: []*model.Role{{Name: "xxxxx", Description: "xssxx"}, {Name: "xxssssssssxxx", Description: "xssxssssssx"}}})
	// activity
	utc_now_time := time.Now().UTC()
	db.Create(&model.Activity{LastLoginAt: &utc_now_time, LastBuyAt: nil, LastChatAt: nil, LastChangePasswordAt: nil, UserID: 1})
	// product
	// inventory
	// tag
	// image path
	db.Create(&model.Product{Name: "product", Description: "xxx", Price: 1111, UserID: 1, Categories: []*model.Product_Category{{Name: "XXXX", Description: "xxxx"}, {Name: "XXxxxxXX", Description: "xxxxxxx"}}, Tags: []*model.Product_Tag{{Name: "sssssssssssss", Description: "ssssssssssss"}, {Name: "sssss", Description: "sssss"}}, Inventory: 21232223, Images: []*model.Product_Image{{Path: "xsfsdfsadfd"}}})
	// address
	db.Create(&model.Address{Name: "ssasdgss", Country: "ssss", City: "sdfasdfd", Street: "sasdgs", BuildingNumber: "sdfasdf", PostalCode: "sfaffas", Description: "sadfsdds", UserID: 1})
	db.Create(&model.Address{Name: "ssgsgsasss", Country: "ssdfsdfxcxcvbdfshsdfgsadfsss", City: "sasdfasdfasdfasdfd", Street: "sasdfsdsafasfdsafaasdgs", BuildingNumber: "sadsfsafsadfasdf", PostalCode: "ssadfasdffaffas", Description: "sadfsssadsadfasffdds", UserID: 1})
	// wallet
	db.Create(&model.WalletInfo{Name: "xxx", Addr: "xxx", IsDefault: true, UserID: 1})
	db.Create(&model.WalletInfo{Name: "xxxx", Addr: "xxxx", IsDefault: false, UserID: 1, OrderID: 1})
	db.Create(&model.WalletInfo{Name: "xxxx", Addr: "xxxx", IsDefault: true, UserID: 1, OrderID: 1})
	// cart
	db.Create(&model.Cart{TotalPrice: 2322, OrderItemQuantity: 3, IsOrdered: true, UserID: 1})
	// order status
	db.Create(&model.OrderStatus{Status: "xxxxxxxxxx"})
	db.Create(&model.OrderStatus{Status: "xxxxxxxxxxx"})
	db.Create(&model.OrderStatus{Status: "xxxxxxxxxxxx"})
	// order
	db.Create(&model.Order{SenderWalletInfoID: 1, UserID: 1, OrderStatusID: 1, CartID: 1})
	// category
	db.Create(&model.Forum_Category{Name: "xxxx", Description: "sssss"})
	db.Create(&model.Topic_Category{Name: "xxxxsssss", Description: "sssssxxxxx"})
	// order item
	db.Create(&model.OrderItem{ProductID: 1, CartID: 1})
	db.Create(&model.OrderItem{ProductID: 1, CartID: 1})
	db.Create(&model.OrderItem{ProductID: 1, CartID: 1})
	// poll
	db.Create(&model.Poll{Name: "xxxx", Description: "xxxx", UserID: 1})
	// comment
	db.Create(&model.Topic_Comment{Text: "xxxxxxxx", UserID: 1, TopicID: 1})
	// vote
	db.Create(&model.Poll_Vote{UserID: 1, PollID: 1})
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
