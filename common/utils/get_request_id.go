package utils

import (
	"entry_task/common/constants"
	"github.com/gin-gonic/gin"
)

func GetRequestId(c *gin.Context) *string {
	requestId, _ := c.Get(constants.RequestId)
	empStr := ""
	if v, ok := requestId.(*string); ok {
		return v
	}
	return &empStr
}
