package public_struct

import (
	"time"

	viewmodel "github.com/0ne-zero/f4h/public_struct/view_model"
)

type TopicForShowTopicViewModelWithUserID struct {
	UserID int
	viewmodel.TopicBriefViewModel
}
type LastPostViewModelWithUserID struct {
	UserID int
	viewmodel.LastPost
}
type UserBasicInformation struct {
	Username  string
	CreatedAt time.Time
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

type RequestBasicInformation struct {
	IP     string
	Path   string
	Method string
}
