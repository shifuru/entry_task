package middlewares

import (
	"entry_task/common/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		RequestId := strconv.FormatInt(time.Now().UnixMilli(), 10) + uuid.New().String()
		c.Set(constants.RequestId, &RequestId)
		c.Next()
	}
}
