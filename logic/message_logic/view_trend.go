package message_logic

import (
	"context"
	"encoding/json"
	"entry_task/model/cache"
	"entry_task/model/database"
	"entry_task/model/database/db_message"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type CountResult struct {
	TagId uint `json:"tag_id"`
	Count int  `json:"count"`
}

type NameResult struct {
	CountResult
	Tag string `json:"tag"`
}

type ReturnResult struct {
	Tags      []NameResult `json:"tags"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func GetTrend(requestId *string, topicId *uint, c context.Context) (*ReturnResult, error) {
	var ret ReturnResult

	redisKey := "trend"
	redisDB := cache.GetRedisClient()

	val, err := redisDB.Get(c, redisKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(val), &ret)
		if err != nil {
			return nil, err
		}
	} else if errors.Is(err, redis.Nil) {
		db := database.GetDB(requestId)
		var result []CountResult
		if *topicId != 0 {
			err = db.Model(&db_message.MessageTag{}).Select("tag_id, COUNT(*) as count").Where("topic_id = ?", topicId).Group("tag_id").Order("count DESC").Scan(&result).Error
			if err != nil {
				return nil, err
			}
		} else {
			err = db.Model(&db_message.MessageTag{}).Select("tag_id, COUNT(*) as count").Group("tag_id").Order("count DESC").Scan(&result).Error
			if err != nil {
				return nil, err
			}
		}
		if len(result) > 5 {
			result = result[:5]
		}
		ret.Tags = make([]NameResult, len(result))
		ret.UpdatedAt = time.Now()

		for i, v := range result {
			err = db.Raw("SELECT tag FROM tag_tab WHERE id = ?", v.TagId).Scan(&ret.Tags[i].Tag).Error
			if err != nil {
				return nil, err
			}
			ret.Tags[i].CountResult = v
		}

		trendBytes, err := json.Marshal(ret)
		if err == nil {
			_ = redisDB.Set(c, redisKey, string(trendBytes), 30*time.Minute).Err()
		}

	} else if err != nil {
		return nil, err
	}

	return &ret, nil

}
