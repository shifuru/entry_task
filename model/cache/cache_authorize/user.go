package cache_authorize

import (
	"context"
	"encoding/json"
	"entry_task/common/errorcode"
	"entry_task/common/utils"
	"entry_task/model/cache"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Group    uint   `json:"group"`
}

func (user *User) getKey(uuid string) string {
	return fmt.Sprintf("user_key_%s", uuid)
}

func (user *User) generateKey() string {
	return utils.UuidHex()
}

func (user *User) Set(ctx context.Context) (string, error) {
	userStr, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	sessionKey := user.generateKey()

	err = cache.GetRedisClient().Set(ctx, user.getKey(sessionKey), userStr, 7*24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return sessionKey, nil
}

func (user *User) Get(ctx context.Context, uuid string) error {
	if uuid == "" {
		return errorcode.GetSessionFailed
	}
	value, err := cache.GetRedisClient().Get(ctx, user.getKey(uuid)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errorcode.GetSessionFailed
		}
		return err
	}
	err = json.Unmarshal([]byte(value), user)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) Delete(ctx context.Context, uuid string) error {
	err := cache.GetRedisClient().Del(ctx, user.getKey(uuid)).Err()
	if err != nil {
		return err
	}
	return nil
}
