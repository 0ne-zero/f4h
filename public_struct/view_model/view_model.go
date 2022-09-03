package viewmodel

import (
	"html/template"
	"time"

	"github.com/0ne-zero/f4h/database/model"
)

type ProductBasicViewModel struct {
	ID        int
	Name      string
	Price     float64
	ImagePath string
}
type DiscussionViewModel struct {
	model.Discussion
	ForumsName []string
	TopicCount int
	PostCount  int
}
type DiscussionCategoryViewModel struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
	// Forum_Category has many Forum
	Discussions []*DiscussionViewModel
	// Forum_Category has one User
	UserID uint
}
type ForumViewModel struct {
	ID          uint
	Name        string
	Description string
	TopicCount  int
	PostCount   int
	LastPost    LastPost `gorm:"-"`
}
type TopicBriefViewModel struct {
	ID             uint
	Name           string
	ReplyCount     int
	ViewCount      int
	CreatedAt      *time.Time
	AuthorUsername string
	LastPost       LastPost `gorm:"-"`
}
type SidebarCategoryViewModel struct {
	Name          string
	SubCategories []SidebarCategoryViewModel
}
type LastPost struct {
	AuthorUsername string
	CreatedAt      time.Time
}

type TopicUserViewModel struct {
	Username  string
	PostCount uint
	JoinedAt  *time.Time
}

type TopicForShowTopicViewModel struct {
	Title       string
	Description template.HTML
	CreatedAt   time.Time
	UserInfo    *TopicUserViewModel
	Tags        []TopicTagBasicInformation
}

type TopicCommentViewModel struct {
	Title     string
	Text      string
	CreatedAt time.Time
	UserInfo  *TopicUserViewModel
	Tags      []TopicTagBasicInformation
	Reply     *TopicCommentViewModel
}

// Basics
type TopicForShowTopicViewModelWithUserID struct {
	UserID int
	TopicBriefViewModel
}
type LastPostViewModelWithUserID struct {
	UserID int
	LastPost
}
type UserBasicInformation struct {
	Username  string
	CreatedAt time.Time
}
type TopicTagBasicInformation struct {
	Name string
}
type TopicBasicInformation struct {
	Name        string
	Description string
	CreatedAt   time.Time
	UserID      int
}
type TopicCommentBasicInformation struct {
	Text      string
	CreatedAt time.Time
	UserID    int
	ReplyID   int
}

type TopicForEditViewModel struct {
	Name        string
	Description string
	Tags        []TopicTagBasicInformation `gorm:"-"`
	ForumName   string                     `gorm:"-"`
}

// User Panel
type UserPanel_Overview_Front struct {
	JoinedAt              *time.Time
	LastLoginAt           *time.Time
	LastPasswordChangedAt *time.Time
	LastBuyAt             *time.Time
	TotalPosts            int
	TotalProducts         int
	TotalPolls            int
}
type UserPanel_Overview_Login struct {
	Login_At *time.Time
}
type UserPanel_Overview_LastBuy struct {
	OrderID     int
	Time        *time.Time
	TotalPrice  float64
	OrderStatus string
}
type UserPanel_Profile_EditAccount struct {
	Email string
}
type UserPanel_Profile_ManageAddress struct {
	Name           string
	Country        string
	Province       string
	City           string
	Street         string
	BuildingNumber string
	PostalCode     string
	Description    string
}
type UserPanel_Profile_ManageWallet struct {
	Name      string
	Addr      string
	IsDefault bool
}
type UserPanel_Profile_EditSignature struct {
	Signature string
}

type CartItem struct {
	// Cart item id
	ID int
	// Product id
	ProductID  int
	Name       string
	Price      float64
	Quantity   int
	ImagePath  string
	TotalPrice float64
}
type Cart struct {
	TotalPrice float64
	CartItems  []CartItem
}
type ImageViewData struct {
	Name string
	Path string
}
type ProductDetailsImagesViewData struct {
	MainImage      string
	NumberOfSlides int
	SubImages      []ImageViewData
}

type ProductDetailsDetail struct {
	ID        int
	Name      string
	Price     float64
	Inventory int
}
type ProductDetailsDescription struct {
	Username    string
	Time        time.Time
	Description string
}
type ProductDetailsUserProduct struct {
	ID        int
	ImagePath string
	Name      string
	Price     float64
}
type ProductDetailsComment struct {
	ID       int
	Username string
	Time     *time.Time
	Text     string
}
type ProductDetailsComments struct {
	ProductID int
	Comments  []ProductDetailsComment
}
type ProductDetailsTabs struct {
	NumberOfComments int
	DescriptionData  ProductDetailsDescription
	UserProductsData []ProductDetailsUserProduct
	CommentsData     ProductDetailsComments
}
type RecommendedItems struct {
	ID        int
	Name      string
	ImagePath string
	Price     float64
}
type RecommendedViewData struct {
	NumberOfSlides   int
	RecommendedItems []RecommendedItems
}
