package user_logic

import (
	"entry_task/common/log"
	"entry_task/common/utils"
	"entry_task/config"
	"entry_task/model/cache/cache_authorize"
	"entry_task/model/database"
	"entry_task/model/database/db_user"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

func Register(requestId *string, username *string, password *string, role *db_user.Role, topicId *uint, c *gin.Context) (uint, error) {

	if *role != db_user.Admin {
		if *topicId == 0 {
			return 0, errors.New("not topic chosen")
		}
		*role = db_user.Normal
	} else {
		*topicId = 0
	}

	hash, err := utils.HashPwd(*password)
	if err != nil {
		return 0, err
	}

	user := db_user.User{
		Username:  *username,
		Password:  hash,
		Role:      *role,
		TopicId:   *topicId,
		CreatedAt: time.Now(),
	}

	db := database.GetDB(requestId)
	err = db.Create(&user).Error
	if err != nil {
		return 0, err
	}

	log.Info(utils.GetRequestId(c), "register successfully, user info: %+v", struct {
		UserId    uint   `json:"user_id"`
		Username  string `json:"username"`
		UserGroup uint   `json:"user_group"`
	}{user.Id, user.Username, user.TopicId},
	)

	return user.Id, nil
}

func Login(requestId *string, username *string, password *string, c *gin.Context) (string, uint, string, uint, error) {
	user := db_user.User{}

	db := database.GetDB(requestId)
	err := db.Where("username = ?", *username).First(&user).Error
	if err != nil {
		return "-1", 0, "", 0, errors.New("user not found")
	}

	if !utils.CheckPassword(password, &user.Password) {
		return "-2", 0, "", 0, errors.New("invalid password")
	}

	var tokenString string

	if config.ProjectConfig.Jwt.Mode == 0 {
		tokenString, err = utils.GenerateJWT(&user.Id, &user.Username, &user.TopicId)
	} else if config.ProjectConfig.Jwt.Mode == 1 {
		userCache := cache_authorize.User{
			Id:       user.Id,
			Username: *username,
			Group:    user.TopicId,
		}
		tokenString, err = userCache.Set(c.Request.Context())
	} else {
		return "", 0, "", 0, errors.New("unknow mode")
	}
	if err != nil {
		return "", 0, "", 0, err
	}

	log.Info(utils.GetRequestId(c), "login successfully, user info: %+v", struct {
		UserId    uint   `json:"user_id"`
		Username  string `json:"username"`
		UserGroup uint   `json:"user_group"`
	}{user.Id, user.Username, user.TopicId},
	)

	return tokenString, user.Id, user.Username, user.TopicId, nil
}
