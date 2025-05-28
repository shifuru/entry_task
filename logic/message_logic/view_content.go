package message_logic

import (
	"entry_task/model/database"
	"entry_task/model/database/db_message"
	"errors"
)

const commentPageSize = 200

func GetContent(requestId *string, messageId *uint, group *uint, commentPage *int) (*db_message.Message, *[]db_message.Comment, error) {
	db := database.GetDB(requestId)

	message := db_message.Message{}
	err := db.First(&message, *messageId).Error
	if err != nil {
		return nil, nil, err
	}

	if message.TopicId != *group && *group != 0 {
		return nil, nil, errors.New("permission denied")
	}

	var comments []db_message.Comment
	offset := (*commentPage - 1) * commentPageSize
	err = db.
		Where("message_id = ?", *messageId).
		Order("created_at DESC").
		Limit(commentPageSize).
		Offset(offset).
		Find(&comments).Error
	if err != nil {
		return nil, nil, err
	}

	return &message, &comments, nil
}
