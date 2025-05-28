package db_message

import (
	"time"
)

type Comment struct {
	Id        uint      `json:"id" gorm:"primary_key;column:id"`
	MessageId uint      `json:"message_id" gorm:"column:message_id"`
	Content   string    `json:"content" gorm:"size:2048"`
	CreatedAt time.Time `json:"created_at"`
	UserId    uint      `json:"user_id"`
}

func (Comment) TableName() string {
	return "comment_tab"
}
