package message_logic

import (
	"encoding/json"
	"entry_task/model/database"
	"entry_task/model/database/db_message"
	"errors"
	"fmt"
	"github.com/elliotchance/pie/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

func UpdateMessage(requestId *string, messageId *uint, topicId *uint, content *string, tags *datatypes.JSON, updateId *uint) (*db_message.Message, error) {
	message := db_message.Message{
		TopicId:   *topicId,
		Content:   *content,
		TagList:   *tags,
		UpdatedId: *updateId,
	}

	db := database.GetDB(requestId)
	err := db.Model(&message).Where("id = ?", *messageId).Updates(message).Error
	if err != nil {
		return nil, err
	}

	var tagString []string
	err = json.Unmarshal(*tags, &tagString)
	if err != nil {
		return nil, err
	}

	var oldTags []uint
	err = db.
		Model(&db_message.MessageTag{}).
		Where("message_id = ?", *messageId).
		Pluck("tag_id", &oldTags).
		Error
	if err != nil {
		return nil, err
	}

	nowTime := time.Now()

	newTags := make([]uint, len(tagString))
	for i, tagName := range tagString {
		tag := db_message.Tag{}
		err = db.
			FirstOrCreate(&tag, db_message.Tag{Tag: tagName}).
			Error
		if err != nil {
			return nil, err
		}
		newTags[i] = tag.Id

		conn := db_message.MessageTag{
			MessageId: *messageId,
			TagId:     tag.Id,
			TopicId:   *topicId,
			UpdatedAt: nowTime,
		}

		oldConn := db_message.MessageTag{}

		err = db.
			Model(&conn).
			Where("tag_id = ? AND message_id = ?", tag.Id, *messageId).
			First(&oldConn).
			Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Create(&conn).Error
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		err = db.
			Model(&conn).
			Where("tag_id = ? AND message_id = ?", tag.Id, *messageId).
			UpdateColumn("updated_at", nowTime).
			Error
		if err != nil {
			return nil, err
		}
	}

	_, removed := pie.Diff(oldTags, newTags)
	for _, removedTag := range removed {
		err = db.
			Model(&db_message.MessageTag{}).
			Where("tag_id = ? AND message_id = ?", removedTag, *messageId).
			Delete(&db_message.MessageTag{}).
			Error
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(message)

	return &message, nil
}
