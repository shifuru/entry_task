package user_apis

import (
	"entry_task/common/log"
	"entry_task/common/utils"
	"entry_task/logic/user_logic"
	"entry_task/model/database/db_user"
	"entry_task/view/http/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequest struct {
	Username string       `json:"username" binding:"required,max=64"`
	Password string       `json:"password" binding:"required,max=255"`
	Role     db_user.Role `json:"role"`
	TopicId  uint         `json:"topic_id"`
}

func UserSignUp(c *gin.Context) {
	requestId := utils.GetRequestId(c)
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(common.ParamErrorCodeWithMessage))
		return
	}
	log.Info(requestId, "get request: %+v", request)

	userId, err := user_logic.Register(requestId, &request.Username, &request.Password, &request.Role, &request.TopicId, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
		return
	}
	log.Info(requestId, "get response: %+v", gin.H{"id": userId})

	c.JSON(http.StatusOK, gin.H{"id": userId})
}
