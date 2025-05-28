package user_apis

import (
	"entry_task/common/constants"
	"entry_task/common/log"
	"entry_task/common/utils"
	"entry_task/logic/user_logic"
	"entry_task/view/http/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,max=64"`
	Password string `json:"password" binding:"required,max=255"`
}

type LoginResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Group    uint   `json:"group"`
}

func UserLogIn(c *gin.Context) {
	requestId := utils.GetRequestId(c)
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(common.ParamErrorCodeWithMessage))
		return
	}
	log.Info(requestId, "get request: %+v", request)

	tokenString, userId, username, userGroup, err := user_logic.Login(requestId, &request.Username, &request.Password, c)
	if err != nil {
		if tokenString == "-1" {
			c.JSON(http.StatusBadRequest, common.CustomErrorResponse(-1, "User not found."))
		} else if tokenString == "-2" {
			c.JSON(http.StatusBadRequest, common.CustomErrorResponse(-1, "Invalid password."))
		} else {
			c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
		}
		return
	}

	loginResponse := LoginResponse{
		Id:       userId,
		Username: username,
		Group:    userGroup,
	}
	log.Info(requestId, "get response: %+v", loginResponse)

	c.SetCookie(constants.SessionCookieName, tokenString, 7*24*3600, "/", "localhost", false, false)
	c.JSON(http.StatusOK, common.SuccessResponse(loginResponse))
}
