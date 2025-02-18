package rest

import "github.com/gin-gonic/gin"

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, msg string) {
	c.AbortWithStatusJSON(
		statusCode,
		response{
			msg,
		})
}
