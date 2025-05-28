package message_logic

import (
	"entry_task/model/database"
	"entry_task/model/database/db_message"
	"errors"
	"time"
)

func NewComment(requestId *string, messageId *uint, content *string, userId *uint, group *uint) (*db_message.Comment, error) {
	db := database.GetDB(requestId)

	var topicId uint
	err := db.Raw("SELECT topic_id FROM message_tab WHERE id = ?", messageId).Scan(&topicId).Error
	if err != nil {
		return nil, err
	} else if topicId == 0 {
		return nil, errors.New("message not found")
	}
	if topicId != *group && *group != 0 {
		return nil, errors.New("permission denied")
	}

	nowTime := time.Now()
	comment := db_message.Comment{
		MessageId: *messageId,
		Content:   *content,
		UserId:    *userId,
		CreatedAt: nowTime,
	}

	err = db.Create(&comment).Error
	if err != nil {
		return nil, err
	}

	return &comment, nil
}
