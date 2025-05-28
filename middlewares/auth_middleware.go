package middlewares

import (
	"entry_task/common/constants"
	"entry_task/common/utils"
	"entry_task/config"
	"entry_task/model/cache/cache_authorize"
	"entry_task/view/http/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(constants.SessionCookieName)
		var token string
		if cookie != nil {
			token = cookie.Value
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, common.CustomErrorResponse(-1, "No token provided"))
			c.Abort()
			return
		}

		var user cache_authorize.User

		if config.ProjectConfig.Jwt.Mode == 0 {
			user.Id, user.Username, user.Group, err = utils.ParseJWT(token)
		} else if config.ProjectConfig.Jwt.Mode == 1 {
			err = user.Get(c.Request.Context(), token)
		} else {
			c.JSON(http.StatusInternalServerError, common.CustomErrorResponse(-1, "unknown mode"))
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, common.CustomErrorResponse(-1, err.Error()))
			c.Abort()
			return
		}

		c.Set("user_id", user.Id)
		c.Set("username", user.Username)
		c.Set("user_group", user.Group)
		c.Next()

	}
}
