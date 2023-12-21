package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var Router router

type router struct{}

func (r *router) InitApiRouter(g *gin.Engine) {
	g.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}
