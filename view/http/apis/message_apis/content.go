package message_apis

import (
	"entry_task/common/log"
	"entry_task/common/utils"
	"entry_task/logic/message_logic"
	"entry_task/model/database/db_message"
	"entry_task/view/http/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ContentResponse struct {
	Message  db_message.Message   `json:"message"`
	Comments []db_message.Comment `json:"comments"`
}

func ViewContent(c *gin.Context) {
	requestId := utils.GetRequestId(c)
	group := c.GetUint("user_group")

	messageId64, err := strconv.ParseUint(c.Param("messageId"), 10, 32)
	if err != nil || messageId64 == 0 {
		c.JSON(http.StatusNotFound, common.CustomErrorResponse(-1, "404 Not Found"))
		return
	}
	messageId := uint(messageId64)

	commentPage64, err := strconv.ParseInt(c.Param("commentPage"), 10, 32)
	if err != nil || commentPage64 < 1 {
		c.JSON(http.StatusNotFound, common.CustomErrorResponse(-1, "404 Not Found"))
		return
	}
	commentPage := int(commentPage64)

	message, comments, err := message_logic.GetContent(requestId, &messageId, &group, &commentPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
		return
	}

	contentResponse := ContentResponse{
		Message:  *message,
		Comments: *comments,
	}
	log.Info(requestId, "get response: %+v", contentResponse)

	c.JSON(http.StatusOK, common.SuccessResponse(contentResponse))
	return

}
