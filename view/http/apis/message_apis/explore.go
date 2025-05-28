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

type ExploreRequest struct {
	TopicId uint    `json:"topic_id"`
	Tag     *string `json:"tag"`
}

type ExploreResponse struct {
	Messages *[]db_message.Message `json:"messages"`
}

func MessageExplore(c *gin.Context) {
	requestId := utils.GetRequestId(c)
	group := c.GetUint("user_group")

	postPage64, err := strconv.ParseInt(c.Param("postPage"), 10, 32)
	if err != nil || postPage64 < 1 {
		c.JSON(http.StatusNotFound, common.CustomErrorResponse(-1, "404 Not Found"))
		return
	}
	postPage := int(postPage64)

	var request ExploreRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(common.ParamErrorCodeWithMessage))
		return
	}
	log.Info(requestId, "get request: %+v", request)

	if request.TopicId == 0 && group != 0 {
		request.TopicId = group
	} else if request.TopicId != group && group != 0 {
		c.JSON(http.StatusForbidden, common.CustomErrorResponse(-1, "permission denied"))
		return
	}

	posts, err := message_logic.ExplorePosts(requestId, &request.TopicId, request.Tag, &postPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
		return
	}

	exploreResponse := ExploreResponse{
		Messages: posts,
	}
	log.Info(requestId, "get response: %+v", exploreResponse)

	c.JSON(http.StatusOK, common.SuccessResponse(exploreResponse))
	return
}
