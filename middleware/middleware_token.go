package middleware

import (
	gin_const "github.com/fellowme/gin_commom_library/const"
	gin_token "github.com/fellowme/gin_commom_library/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type errorResponse struct {
	Message string
	Status  int
}

func CheckUserToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := gin_token.VerifyUserId(c); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errorResponse{
				Message: gin_const.VerifyUserFail,
				Status:  gin_const.VerifyUserFailCode,
			})
		}
		c.Next()
	}
}
