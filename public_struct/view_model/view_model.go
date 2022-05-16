package viewmodel

import (
	"html/template"
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
