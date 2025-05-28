package message_apis

import (
	"entry_task/common/log"
	"entry_task/common/utils"
	"entry_task/logic/message_logic"
	"entry_task/view/http/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RemoveResponse struct {
	Id uint `json:"id"`
}

func RemoveMessage(c *gin.Context) {
	requestId := utils.GetRequestId(c)

	messageId64, err := strconv.ParseUint(c.Param("messageId"), 10, 32)
	if err != nil || messageId64 == 0 {
		c.JSON(http.StatusNotFound, common.CustomErrorResponse(-1, "404 Not Found"))
		return
	}
	messageId := uint(messageId64)

	group := c.GetUint("user_group")
	userId := c.GetUint("user_id")

	if group != 0 {
		c.JSON(http.StatusForbidden, common.CustomErrorResponse(-1, "permission denied"))
		return
	}

	err = message_logic.RemovePost(requestId, &messageId, &userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
		return
	}

	removeResponse := RemoveResponse{
		Id: messageId,
	}
	log.Info(requestId, "get response: %+v", removeResponse)

	c.JSON(http.StatusOK, common.SuccessResponse(removeResponse))

}
