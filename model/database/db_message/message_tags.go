package db_message

import "time"

type MessageTag struct {
	TagId     uint      `json:"tag_id"`
	MessageId uint      `json:"message_id"`
	TopicId   uint      `json:"topic_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MessageTag) TableName() string {
	return "message_tags"
}
