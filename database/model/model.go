package model

import (
	"time"
)

// This a Base Model for other models. like gorm.Model but without DeleteAt field
type BasicModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//region User
type User struct {
	BasicModel
	Username     string `gorm:"NOT NULL;"`
	Email        string `gorm:"NOT NULL;"`
	PasswordHash string `gorm:"NOT NULL;"`
	IsSeller     bool   `gorm:"NOT NULL;"`
	// User has many Order
	Orders []*Order `gorm:"foreignkey:UserID;references:ID"`
	// User has many Cart
	Carts []*Cart
	// User has many Poll
	Polls []*Poll
	// User has many Product_Comment
	Comments []*Product_Comment
	// User has many Product_Comment
	Product_Categories []*Product_Category
	// User has many Role
	Roles []*Role `gorm:"many2many:user_roles_m2m;NOT NULL"`
	// User has one Activity
	Activity Activity
	// User has many WalletInfo
	WalletInfos []*WalletInfo
	// User has many Address
	Addresses []*Address
	// User has many Product
	Products []*Product
	// User has many Forum
	Forums []*Forum
	// User has many Forum_Category
	Forum_Categories []*Forum_Category
	// User has many Topic
	Topics []*Topic
	// User has many Topic_Category
	Topic_Categories []*Topic_Category
	// User has many Topic_Comment
	Topic_Comments []*Topic_Comment
	// User has many Poll_Vote
	Poll_Votes []*Poll_Vote
	// User has many Topic_Comment_Vote
	Topic_Comment_Votes []*Topic_Comment_Vote
	// User has many Product_Comment_Vote
	Product_Comment_Votes []*Product_Comment_Vote
	// User has many Topic_Vote
	Topic_Votes []*Topic_Vote
}

type Address struct {
	BasicModel
	Name           string `gorm:"NOT NULL;"`
	Country        string `gorm:"NOT NULL;"`
	Province       string
	City           string `gorm:"NOT NULL;"`
	Street         string `gorm:"NOT NULL;"`
	BuildingNumber string `gorm:"NOT NULL;"`
	PostalCode     string
	Description    string
	// Address has one User
	UserID uint `gorm:"NOT NULL;"`
}
type WalletInfo struct {
	BasicModel
	Name      string `gorm:"NOT NULL;"`
	Addr      string `gorm:"NOT NULL;"`
	IsDefault bool   `gorm:"NOT NULL;"`
	// WalletInfo has one User
	UserID uint `gorm:"NOT NULL;"`
	// WalletInfo has one Order for sender (money)
	OrderID uint
}
type Activity struct {
	BasicModel
	LastLoginAt          *time.Time `gorm:"NOT NULL;"`
	LastBuyAt            *time.Time
	LastChatAt           *time.Time
	LastChangePasswordAt *time.Time
	// Activity has one User
	UserID uint `gorm:"NOT NULL;"`
}
type Role struct {
	BasicModel
	Name        string `gorm:"NOT NULL;"`
	Description string `gorm:"NOT NULL;"`

	// Role has many User
	Users []*User `gorm:"many2many:user_roles_m2m"`
}

//endregion

//region Product
type Product struct {
	BasicModel
	Name        string  `gorm:"NOT NULL;"`
	Description string  `gorm:"NOT NULL;"`
	Price       float64 `gorm:"NOT NULL;"`
	Inventory   uint
	// Product has many ImagePath
	Images []*Product_Image `gorm:"NOT NULL;"`

	// Product has many Product_Category and conversely (many to many)
	Categories []*Product_Category `gorm:"many2many:product_category_m2m;NOT NULL;"`
	// Product has many Comment
	Comments []*Product_Comment
	// Product has many ProductTag
	Tags []*Product_Tag `gorm:"many2many:product_tag_m2m;NOT NULL"`
	// Product has one User
	UserID uint `gorm:"NOT NULL;"`
}
type Product_Category struct {
	BasicModel
	Name          string `gorm:"NOT NULL;"`
	Description   string `gorm:"NOT NULL;"`
	ParentID      *uint
	SubCategories []Product_Category `gorm:"foreignkey:ParentID"`
	// Category has many Product and conversely (many to many)
	Products []*Product `gorm:"many2many:product_category_m2m"`
	// Product_Category has one User
	UserID uint
}
type Product_Tag struct {
	BasicModel
	Name        string `gorm:"NOT NULL;"`
	Description string `gorm:"NOT NULL;"`

	// Tag has many Product and conversely (many to many)
	Products []*Product `gorm:"many2many:product_tag_m2m"`
}
type Product_Image struct {
	BasicModel
	Path      string `gorm:"NOT NULL;"`
	ProductID uint
}
type Product_Comment struct {
	BasicModel
	Text string `gorm:"NOT NULL;"`
	// Product_Comment has one Product_Comment_Vote
	Votes Product_Comment_Vote `gorm:"foreignkey:Product_CommentID;references:ID"`
	// Comment has one User
	UserID uint `gorm:"NOT NULL;"`
	// Commnet has one Product
	ProductID uint
}

//endregion

//region Payment
type Order struct {
	BasicModel

	// Order belongs to WalletInfo for sender
	SenderWalletInfo   *WalletInfo `gorm:"foreignkey:SenderWalletInfoID;references:ID"`
	SenderWalletInfoID uint
	// Order has many WalletInfo for recivers
	ReciversWalletInfo []*WalletInfo

	Cart Cart
	// Order has one Cart
	CartID uint
	// Order has one User
	UserID uint
	// Order has one OrderStatus
	OrderStatusID uint `gorm:"NOT NULL;"`
}
type OrderItem struct {
	BasicModel
	// OrderItem has one Product
	ProductID uint    `gorm:"NOT NULL;"`
	Product   Product `gorm:"NOT NULL;"`
	// OrderItem has one Cart
	CartID uint
}
type OrderStatus struct {
	BasicModel
	Status string `gorm:"NOT NULL;"`
	// OrderStatus has many Order
	Orders []*Order
}
type Cart struct {
	BasicModel

	TotalPrice        float64 `gorm:"NOT NULL;"`
	IsOrdered         bool    `gorm:"NOT NULL;"`
	OrderItemQuantity uint    `gorm:"NOT NULL;"`
	// Cart has one User
	UserID uint `gorm:"NOT NULL;"`

	// Cart has many OrderItem
	OrderItems []*OrderItem `gorm:"NOT NULL;"`
}

//endregion

//region Votes
type Product_Comment_Vote struct {
	Positive uint
	Negative uint
	// Product_Comment_Vote has one Product_Comment
	Product_CommentID uint
	// Product_Comment_Vote has one User
	UserID uint
}
type Topic_Comment_Vote struct {
	Positive uint
	Negative uint
	// Topic_Comment_Vote has one Topic_Comment
	Topic_CommentID uint
	// Topic_Comment_Vote has one User
	UserID uint
}
type Poll_Vote struct {
	Positive uint
	Negative uint
	// Poll_Vote has one Poll
	PollID uint
	// Poll_Vote has one User
	UserID uint
}
type Topic_Vote struct {
	Positive uint
	Negative uint
	// Topic_Vote has one Topic
	TopicID uint
	// Topic_Vote has one User
	UserID uint
}

//endregion

//region Forum

type Forum struct {
	BasicModel
	Name        string
	Description string
	// Forum has many Topic
	Topics []*Topic
	// Forum has many Forum_Category
	Categories []*Forum_Category `gorm:"many2many:forum_category_m2m"`
	// Forum has a User
	UserID uint
}
type Topic struct {
	BasicModel
	Name        string
	Description string
	// Topic_Comment has one Topic_Comment_Vote
	Votes Topic_Vote `gorm:"foreignkey:TopicID;references:ID"`
	// Topic has one User
	UserID uint
	// Topic has one Forum
	ForumID uint
	// Topic has many Category
	Categories []*Topic_Category `gorm:"many2many:topic_category_m2m"`
	// Topic has many Topic_Comment
	Comments []*Topic_Comment
	// Topic has many Topic_Tag
	Tags []*Topic_Tag `gorm:"many2many:topic_tag_m2m"`
}
type Topic_Tag struct {
	BasicModel
	Name        string
	Description string
	// Topic_Tag has many Topic
	Topics []*Topic `gorm:"many2many:topic_tag_m2m"`
}
type Forum_Category struct {
	BasicModel
	Name             string `gorm:"NOT NULL;"`
	Description      string `gorm:"NOT NULL;"`
	ParentCategoryID uint
	// Forum_Category has many Forum
	Forums []*Forum `gorm:"many2many:forum_category_m2m"`
	// Forum_Category has one User
	UserID uint
}
type Topic_Comment struct {
	BasicModel
	Text string `gorm:"NOT NULL;"`
	// Topic_Comment has one User
	UserID uint `gorm:"NOT NULL;"`
	// Topic_Comment has one Topic
	TopicID uint
	// Topic_Comment has one Topic_Comment_Vote
	Votes Topic_Comment_Vote `gorm:"foreignkey:Topic_CommentID;references:ID"`
	// Topic_Comment has one Topic_Comment
	// ReplyID is a Topic_Comment ID
	ReplyID uint
}
type Topic_Category struct {
	BasicModel
	Name             string `gorm:"NOT NULL;"`
	Description      string `gorm:"NOT NULL;"`
	ParentCategoryID uint

	// Topic_Category has many Topic
	Topics []*Topic `gorm:"many2many:topic_category_m2m"`
	// Topic_Category has one User
	UserID uint
}

//endregion

//region Website
type Request struct {
	ID     uint      `gorm:"primarykey"`
	IP     string    `gorm:"NOT NULL"`
	Url    string    `gorm:"NOT NULL"`
	Method string    `gorm:"NOT NULL"`
	Time   time.Time `gorm:"NOT NULL"`
}

//endregion

// Poll
type Poll struct {
	BasicModel
	Name        string `gorm:"NOT NULL;"`
	Description string `gorm:"NOT NULL;"`
	// Poll has one Poll_Vote
	Votes Poll_Vote `gorm:"foreignkey:PollID;references:ID"`
	// Poll has one User
	UserID uint `gorm:"NOT NULL;"`
}