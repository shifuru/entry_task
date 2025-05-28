package message_apis

import (
	"encoding/json"
	"entry_task/common/log"
	"entry_task/common/utils"
	"entry_task/logic/message_logic"
	"entry_task/view/http/common"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"net/http"
	"strconv"
)

type UpdateRequest struct {
	TopicId uint            `json:"topic_id" binding:"required,gte=0"`
	Content *string         `json:"content" binding:"required,max=2048"`
	Tags    *datatypes.JSON `json:"tags"`
}

type UpdateResponse struct {
	Id uint `json:"id"`
}

func MessageUpdate(c *gin.Context) {
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

	var request UpdateRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ParamErrorCodeWithMessage))
		return
	}
	log.Info(requestId, "get request: %+v", request)

	if request.Tags == nil {

		var empty []string

		raw, err := json.Marshal(empty)
		if nil != err {
			c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
			return
		}

		nonTag := datatypes.JSON(raw)
		request.Tags = &nonTag
	}

	_, err = message_logic.UpdateMessage(requestId, &messageId, &request.TopicId, request.Content, request.Tags, &userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
		return
	}

	updateResponse := UpdateResponse{Id: messageId}
	log.Info(requestId, "get response: %+v", updateResponse)

	c.JSON(http.StatusOK, common.SuccessResponse(updateResponse))

}
