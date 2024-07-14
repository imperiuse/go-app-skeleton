package api

import "github.com/gin-gonic/gin"

type Engine = gin.Engine

func NewEngine() *Engine {
	e := gin.New()
	e.ForwardedByClientIP = true
	e.RemoteIPHeaders = []string{"X-Real-IP", "X-Forwarded-For"}

	return e
}
