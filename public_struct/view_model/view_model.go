package viewmodel

import (
	"time"

	"github.com/0ne-zero/f4h/database/model"
)

type ProductViewModel struct {
	ID    int
	Name  string
	Price float64
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

type TopicViewModel struct {
	Title       string
	Description string
	CreatedAt   string
	UserInfo    *TopicUserViewModel
}

type TopicCommentViewModel struct {
	Title       string
	Description string
	CreatedAt   string
	UserInfo    *TopicUserViewModel
}

type ShowTopicViewModel struct {
	Topic    TopicViewModel
	Comments []TopicCommentViewModel
}
