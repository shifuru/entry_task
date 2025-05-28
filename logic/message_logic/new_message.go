package message_logic

import (
	"encoding/json"
	"entry_task/model/database"
	"entry_task/model/database/db_message"
	"gorm.io/datatypes"
	"time"
)

func NewMessage(requestId *string, topicId *uint, content *string, tags *datatypes.JSON, createId *uint) (*db_message.Message, error) {
	nowTime := time.Now()
	message := db_message.Message{
		TopicId:   *topicId,
		Content:   *content,
		TagList:   *tags,
		CreatedId: *createId,
	}

	db := database.GetDB(requestId)
	err := db.Create(&message).Error
	if err != nil {
		return nil, err
	}

	var tagString []string
	err = json.Unmarshal(*tags, &tagString)
	if err != nil {
		return nil, err
	}

	nowTime = time.Now()

	for _, tagName := range tagString {
		tag := db_message.Tag{}

		err = db.FirstOrCreate(&tag, db_message.Tag{Tag: tagName}).Error
		if err != nil {
			return nil, err
		}

		conn := db_message.MessageTag{
			MessageId: message.ID,
			TagId:     tag.Id,
			TopicId:   message.TopicId,
			UpdatedAt: nowTime,
		}

		err = db.Create(&conn).Error
		if err != nil {
			return nil, err
		}

	}

	return &message, nil
}
