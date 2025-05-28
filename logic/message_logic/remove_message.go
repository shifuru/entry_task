package message_logic

import (
	"entry_task/model/database"
	"entry_task/model/database/db_message"
)

func RemovePost(requestId *string, messageId *uint, userId *uint) error {
	db := database.GetDB(requestId)

	err := db.Model(&db_message.Message{}).Where("id = ?", *messageId).UpdateColumn("deleted_id", *userId).Error
	if err != nil {
		return err
	}

	err = db.Model(&db_message.Message{}).Where("id = ?", *messageId).Delete(&db_message.Message{}).Error
	if err != nil {
		return err
	}

	return nil
}
