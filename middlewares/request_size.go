package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MaxAllowedSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		//? Essentially writes max size to request.body when parsed by function BindJSON of gin
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}
