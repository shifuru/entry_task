package db_message

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content   string         `json:"content" gorm:"size:2048"`
	TopicId   uint           `json:"topic_id"`
	TagList   datatypes.JSON `json:"tag_list"`
	CreatedId uint           `json:"created_id"`
	UpdatedId uint           `json:"updated_id"`
	DeletedId uint           `json:"deleted_id"`
}

func (Message) TableName() string {
	return "message_tab"
}
