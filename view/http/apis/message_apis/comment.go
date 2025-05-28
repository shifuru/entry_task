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

type CommentRequest struct {
	Content *string `json:"content" binding:"required,max=2048"`
}

type CommentResponse struct {
	Id uint `json:"id"`
}

func MessageComment(c *gin.Context) {
	requestId := utils.GetRequestId(c)

	messageId64, err := strconv.ParseUint(c.Param("messageId"), 10, 32)
	if err != nil || messageId64 == 0 {
		c.JSON(http.StatusNotFound, common.CustomErrorResponse(-1, "404 Not Found"))
		return
	}
	messageId := uint(messageId64)

	group := c.GetUint("user_group")
	userId := c.GetUint("user_id")

	var request CommentRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(common.ParamErrorCodeWithMessage))
		return
	}
	log.Info(requestId, "get request: %+v", request)

	comment, err := message_logic.NewComment(requestId, &messageId, request.Content, &userId, &group)
	if nil != err {
		c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
		return
	}

	commentResponse := CommentResponse{Id: comment.Id}
	log.Info(requestId, "get response: %+v", commentResponse)

	c.JSON(http.StatusOK, common.SuccessResponse(commentResponse))
	return
}
