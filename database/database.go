package database

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/0ne-zero/f4h/database/model"
	general_func "github.com/0ne-zero/f4h/utilities/functions/general"
	"github.com/0ne-zero/f4h/utilities/functions/setting"
	"github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// If it couldn't connect to database, and also it didn't close program, returns nil
func InitializeOrGetDB() *gorm.DB {
	if db == nil {
		// DSN = Data source name (like connection string for database)
		dsn, err := setting.ReadFieldInSettingData("DSN")
		if err != nil {
			return nil
		}

		// For error handling
		var connect_again = true
		for connect_again {
			db, err = connectDB(dsn)
			if err != nil {
				// Specific error handling

				//Databse doesn't exists, we have to create the database
				if strings.Contains(err.Error(), "Unknown database") {
					err = CreateDatabaseFromDSN(dsn)
					if err != nil {
						// Database isn't exists
						// Also we can't create database from dsn
						fmt.Println(fmt.Sprintf("Mentioned database in dsn isn't created,program tried to create that database but it can't do that.\nError: %s", err.Error()))
						os.Exit(1)
					}
					// Database created in mysql
					// Don't check rest of possible errors and try to connect again
					continue
				}
				// Error handling with error type detection
				switch err.(type) {
				case *net.OpError:
					op_err := err.(*net.OpError)
					// Get TCPAddr if exists
					if tcp_addr, ok := op_err.Addr.(*net.TCPAddr); ok {
						// Check error occurred when we trired to connect to mysql
						if tcp_addr.Port == 3306 {
							// Try to start mysql service
							connect_again = StartMySqlService()
						}
					}
				default:
					wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "Cannot connect to database " + err.Error(), Fields: map[string]string{"Hint": "Maybe you should start database service(deamon)"}, ErrorLocation: general_func.GetCallerInfo(0)})
				}
			} else {
				// We don't need to try again to connect to database because we are connected
				connect_again = false
			}
		}

		db.Set("gorm:auto_preload", true)
		return getDB()
	} else {
		return getDB()
	}
}
func connectDB(dsn string) (*gorm.DB, error) {
	// Connect to database with gorm
	return gorm.Open(
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
}

func getDB() *gorm.DB {
	return db
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
		&model.BadRequest{},
		&model.CartItem{},
		&model.OrderStatus{},
		&model.Cart{},
		&model.Product_Comment_Vote{},
		&model.Topic_Comment_Vote{},
		&model.Wishlist{},
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
func CreateInitialData(db *gorm.DB) {
	utc_now_time := time.Now().UTC()
	// user
	// role
	db.Create(
		&model.User{
			Username: "admin", Email: "admin@gmail.com",
			PasswordHash: "d83f900baedb967c3b4e5fb5411abb4ea1986b018804bfea9887345626e9a36efb711872a76fa212ab83e634f96bb7c711808bc8231db849777361eb7ae409db",
			IsSeller:     false,
			Roles:        []*model.Role{{Name: "admin", Description: "admin group"}},
			Signature:    "The great admin",
			AvatarPath:   "/statics/images/avatar/admin_avatar.jpeg",
			JoinedAt:     &utc_now_time,
			IsAdmin:      true})
	// activity
	db.Create(
		&model.Activity{
			LastLoginAt: &utc_now_time, LastBuyAt: nil, LastChangePasswordAt: nil, UserID: 1, LoginsAt: fmt.Sprintf("%s|", utc_now_time)})
	// product
	db.Create(
		&model.Product{
			Name: "Soviet Union Computer", Description: "An old computer made in soviet union", Price: 450, UserID: 1,
			Categories: []*model.Product_Category{{Name: "Computer", Description: "Computers category", UserID: 1}},
			Tags:       []*model.Product_Tag{{Name: "Pc", Description: "PC tags"}},
			Inventory:  3})

	db.Create(
		&model.Product{
			Name: "Klider Book", Description: "Klider book by Mahmoud Dolat Abadi", Price: 89, UserID: 1,
			Categories: []*model.Product_Category{{Name: "Book", Description: "Books category", UserID: 1}},
			Tags:       []*model.Product_Tag{{Name: "Book", Description: "Book"}},
			Inventory:  74})

	db.Create(
		&model.Product{
			Name: "Helicopter", Description: "I have a helicopter but, I don't know model of it", Price: 9320, UserID: 1,
			Categories: []*model.Product_Category{{Name: "War", Description: "War category", UserID: 1}},
			Tags:       []*model.Product_Tag{{Name: "War", Description: "War things"}},
			Inventory:  1})

	db.Create(
		&model.Product{
			Name: "Uniqe Hammer", Description: "An old uniqe hammer made in Korea", Price: 100020, UserID: 1,
			Categories: []*model.Product_Category{{Name: "Hammer", Description: "", UserID: 1}},
			Tags:       []*model.Product_Tag{{Name: "Hammer", Description: "Hammers category"}},
			Inventory:  1})

	db.Create(
		&model.Product{
			Name: "IBM Watson Super-Computer", Description: "Watson Super-Computer", Price: 100020, UserID: 1,
			Categories: []*model.Product_Category{{BasicModel: model.BasicModel{ID: 1}}},
			Tags:       []*model.Product_Tag{{BasicModel: model.BasicModel{ID: 1}}},
			Inventory:  3})
	db.Create(
		&model.Product{
			Name: "F-14 Aircraft", Description: "Now you can buy f-14 :)", Price: 100020, UserID: 1,
			Categories: []*model.Product_Category{{BasicModel: model.BasicModel{ID: 3}}},
			Tags:       []*model.Product_Tag{{BasicModel: model.BasicModel{ID: 3}}},
			Inventory:  5})

	db.Create(&model.Product_Image{ProductID: 1, Path: "/statics/images/product/soviet_union_pc.jpg"})
	db.Create(&model.Product_Image{ProductID: 2, Path: "/statics/images/product/klider.jpg"})
	db.Create(&model.Product_Image{ProductID: 3, Path: "/statics/images/product/helicopter.jpg"})
	db.Create(&model.Product_Image{ProductID: 4, Path: "/statics/images/product/hammer.jpg"})
	db.Create(&model.Product_Image{ProductID: 5, Path: "/statics/images/product/watson_sc.jpg"})
	db.Create(&model.Product_Image{ProductID: 6, Path: "/statics/images/product/aircraft_f14.jpg"})

	db.Create(&model.Product_Comment{
		Text:      "This is a uniqe product",
		UserID:    1,
		ProductID: 1,
	})
	db.Create(&model.Product_Comment{
		Text:      "It should be one of Stalin super computer :)",
		UserID:    1,
		ProductID: 1,
	})
	// address
	db.Create(
		&model.Address{
			Name: "Joseph Stalin", Country: "Russia", City: "Stalingrad", Street: "Stalin",
			BuildingNumber: "1-Stalin", PostalCode: "1", Description: "Don't look at the package!",
			UserID: 1})

	db.Create(
		&model.Wishlist{
			UserID: 1, Products: []*model.Product{{BasicModel: model.BasicModel{ID: 1}}}})

	// wallet
	db.Create(&model.WalletInfo{Name: "None", Addr: "A-valid-addr", IsDefault: true, UserID: 1})

	// Discussion category
	db.Create(
		&model.Discussion_Category{
			Name: "General", Description: "General things", UserID: 1})
	db.Create(
		&model.Discussion_Category{
			Name: "Bitcoin", Description: "Bitcoin things", UserID: 1})
	db.Create(
		&model.Discussion_Category{
			Name: "Art", Description: "Art things", UserID: 1})

	// Discussion
	db.Create(
		&model.Discussion{
			Name: "Mining discussion", Description: "Mining discussion", UserID: 1, Categories: []*model.Discussion_Category{{BasicModel: model.BasicModel{ID: 2}}}})

	db.Create(
		&model.Discussion{
			Name: "Buying discussion", Description: "Buying discussion", UserID: 1, Categories: []*model.Discussion_Category{{BasicModel: model.BasicModel{ID: 2}}}})

	db.Create(
		&model.Discussion{
			Name: "Art discussions", Description: "Art discussions", UserID: 1, Categories: []*model.Discussion_Category{{BasicModel: model.BasicModel{ID: 3}}}})
	db.Create(
		&model.Discussion{
			Name: "Other", Description: "Other things", UserID: 1, Categories: []*model.Discussion_Category{{BasicModel: model.BasicModel{ID: 1}}}})

	// forum
	db.Create(
		&model.Forum{
			Name: "Mining in Russia", Description: "Mining in Russia forum", DiscussionID: 1, UserID: 1})
	db.Create(
		&model.Forum{
			Name: "Mining in US", Description: "Mining in US forum", DiscussionID: 1, UserID: 1})
	db.Create(
		&model.Forum{
			Name: "Buying in England", Description: "Buying in England forum", DiscussionID: 2, UserID: 1})
	db.Create(
		&model.Forum{
			Name: "Buying in Finland", Description: "Buying in Finland forum", DiscussionID: 2, UserID: 1,
		})
	db.Create(
		&model.Forum{
			Name: "Music arts", Description: "Music arts forum", DiscussionID: 3, UserID: 1,
		})
	db.Create(
		&model.Forum{
			Name: "Something other", Description: "Other forum", DiscussionID: 4, UserID: 1,
		})

	// topics
	db.Create(&model.Topic{Name: "Mining in Russia topic", Description: "How can i mine bitcoin in Russia?", ForumID: 1, UserID: 1})
	db.Create(&model.Topic{Name: "Mining in US topic", Description: "How can i mine bitcoin in US?", ForumID: 2, UserID: 1})
	db.Create(&model.Topic{Name: "Buying in England", Description: "How can i buy bitcoin in England", ForumID: 3, UserID: 1})
	db.Create(&model.Topic{Name: "Buying in Finland", Description: "How can i buy bitcoin in Finland", ForumID: 4, UserID: 1})
	db.Create(&model.Topic{Name: "Can AI make music", Description: "Is there any AI algorithm that can make musics", ForumID: 5, UserID: 1})
	db.Create(&model.Topic{Name: "I want to ask some questions out of from your forums", Description: "Who am i?\nWho is you?\nAm i a robot?", ForumID: 6, UserID: 1})

	// cart
	db.Create(
		&model.Cart{
			TotalPrice: 2322, IsOrdered: false, UserID: 1,
			CartItems: []*model.CartItem{{ProductID: 1, ProductQuantity: 4}, {ProductID: 3, ProductQuantity: 1}, {ProductID: 4, ProductQuantity: 1}, {ProductID: 5, ProductQuantity: 1}, {ProductID: 6, ProductQuantity: 1}},
		})

	// order status
	db.Create(&model.OrderStatus{Status: "Processing"})
	db.Create(&model.OrderStatus{Status: "Delivered"})
	// order
	db.Create(&model.Order{SenderWalletInfoID: 1, UserID: 1, OrderStatusID: 2, CartID: 1})
	// comment
	db.Create(&model.Topic_Comment{Text: "I can tell you, you can buy it very hard", UserID: 1, TopicID: 1})
	db.Create(&model.Topic_Comment{Text: "Easy as drinking water", UserID: 1, TopicID: 3})
	db.Create(&model.Topic_Comment{Text: "Easy as drinking water", UserID: 1, TopicID: 4})
	db.Create(&model.Topic_Comment{Text: "Yes, There is some algorithms that can do that, But not as human", UserID: 1, TopicID: 5})
}

// Incomplete
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
	// w = without
	w_user_pass_protocol_ip := dsn[strings.LastIndex(dsn, "/")+1:]
	return w_user_pass_protocol_ip[:strings.LastIndex(w_user_pass_protocol_ip, "?")]
}

func CreateDatabaseFromDSN(dsn string) error {
	// Create database
	dsn_without_database := dsn[:strings.LastIndex(dsn, "/")] + "/"
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
	var service_names = []string{"mysqld.service", "mysql.service"}
	for i := range service_names {
		command := fmt.Sprintf("systemctl start %s", service_names[i])
		_, err := exec.Command("bash", "-c", command).Output()
		if err == nil {
			return true
		}
	}
	return false
}
