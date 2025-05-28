package message_apis

import (
	"encoding/json"
	"entry_task/common/log"
	"entry_task/common/utils"
	"entry_task/logic/message_logic"
	"entry_task/view/http/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"net/http"
)

type NewMessageRequest struct {
	TopicId uint            `json:"topic_id" binding:"required,gte=0"`
	Content *string         `json:"content" binding:"required,max=2048"`
	Tags    *datatypes.JSON `json:"tags"`
}

type MessageResponse struct {
	Id uint `json:"id"`
}

func MessageNew(c *gin.Context) {
	requestId := utils.GetRequestId(c)
	group := c.GetUint("user_group")
	userId := c.GetUint("user_id")

	if group != 0 {
		c.JSON(http.StatusForbidden, common.CustomErrorResponse(-1, "permission denied"))
		return
	}

	var request NewMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(common.ParamErrorCodeWithMessage))
		return
	}
	log.Info(requestId, "get request: %+v", request)

	if request.Tags == nil {

		var empty []string

		raw, err := json.Marshal(empty)
		if nil != err {
			fmt.Println("error: ", err)
			return
		}

		nonTag := datatypes.JSON(raw)
		request.Tags = &nonTag
	}

	message, err := message_logic.NewMessage(requestId, &request.TopicId, request.Content, request.Tags, &userId)
	if nil != err {
		c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
		return
	}

	messageResponse := &MessageResponse{
		Id: message.ID,
	}
	log.Info(requestId, "get response: %+v", messageResponse)

	c.JSON(http.StatusOK, common.SuccessResponse(messageResponse))
	return
}
