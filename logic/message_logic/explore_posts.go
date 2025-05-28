package message_logic

import (
	"entry_task/model/database"
	"entry_task/model/database/db_message"
	"sort"
)

const postPageSize = 2

func ExplorePosts(requestId *string, topicId *uint, tag *string, postPage *int) (*[]db_message.Message, error) {

	db := database.GetDB(requestId)
	var posts []db_message.Message
	offset := (*postPage - 1) * postPageSize
	if *topicId == 0 && tag == nil {
		err := db.
			Model(&db_message.Message{}).
			Order("updated_at DESC").
			Limit(postPageSize).
			Offset(offset).
			Find(&posts).
			Error
		if err != nil {
			return nil, err
		}
		return &posts, nil
	} else if *topicId == 0 && tag != nil {
		var tagId uint
		err := db.Raw("SELECT id FROM tag_tab WHERE tag = ?", *tag).Scan(&tagId).Error
		if err != nil {
			return nil, err
		}

		var targetId []uint
		err = db.
			Model(&db_message.MessageTag{}).
			Where("tag_id = ?", tagId).
			Order("updated_at DESC").
			Limit(postPageSize).
			Offset(offset).
			Pluck("message_id", &targetId).
			Error
		if err != nil {
			return nil, err
		}

		err = db.Model(&db_message.Message{}).Where("id IN ?", targetId).Find(&posts).Error
		if err != nil {
			return nil, err
		}
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].UpdatedAt.After(posts[j].UpdatedAt)
		})

		return &posts, nil

	} else if tag == nil {
		var targetId []uint
		err := db.
			Model(&db_message.MessageTag{}).
			Where("topic_id = ?", *topicId).
			Order("updated_at DESC").
			Limit(postPageSize).
			Offset(offset).
			Pluck("message_id", &targetId).
			Error
		if err != nil {
			return nil, err
		}

		err = db.Model(&db_message.Message{}).Where("id IN ?", targetId).Find(&posts).Error
		if err != nil {
			return nil, err
		}
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].UpdatedAt.After(posts[j].UpdatedAt)
		})

		return &posts, nil
	} else {
		var tagId uint
		err := db.Raw("SELECT id FROM tag_tab WHERE tag = ?", *tag).Scan(&tagId).Error
		if err != nil {
			return nil, err
		}

		var targetId []uint
		err = db.
			Model(&db_message.MessageTag{}).
			Where("topic_id = ? AND tag_id = ?", *topicId, tagId).
			Order("updated_at DESC").
			Limit(postPageSize).
			Offset(offset).
			Pluck("message_id", &targetId).
			Error
		if err != nil {
			return nil, err
		}

		err = db.Model(&db_message.Message{}).Where("id IN ?", targetId).Find(&posts).Error
		if err != nil {
			return nil, err
		}
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].UpdatedAt.After(posts[j].UpdatedAt)
		})

		return &posts, nil
	}

}
