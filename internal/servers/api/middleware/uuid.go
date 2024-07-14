package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const XRequestID = "X-Request-Id"

func UUIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := uuid.New()
		c.Set(XRequestID, u)
		c.Writer.Header().Set(XRequestID, u.String())

		c.Next()
	}
}
