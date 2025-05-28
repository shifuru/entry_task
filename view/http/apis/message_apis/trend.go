package message_apis

import (
	"entry_task/common/log"
	"entry_task/common/utils"
	"entry_task/logic/message_logic"
	"entry_task/view/http/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TrendRequest struct {
	TopicId *uint `json:"topic_id"`
}

type TrendResponse struct {
	Trend *message_logic.ReturnResult `json:"trend"`
}

func MessageTrend(c *gin.Context) {
	requestId := utils.GetRequestId(c)

	group := c.GetUint("user_group")

	if group == 0 {
		var request TrendRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorResponse(common.ParamErrorCodeWithMessage))
			return
		}
		log.Info(requestId, "get request: %+v", request)

		if request.TopicId != nil {
			group = *request.TopicId
		}
	}

	trend, err := message_logic.GetTrend(requestId, &group, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, err.Error()))
		return
	}

	trendResponse := TrendResponse{Trend: trend}
	log.Info(requestId, "get response: %+v", trendResponse)

	c.JSON(http.StatusOK, common.SuccessResponse(trendResponse))
}
