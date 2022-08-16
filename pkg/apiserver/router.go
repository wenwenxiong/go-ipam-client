package apiserver

import (
	"github.com/gin-gonic/gin"
	"wenwenxiong/go-ipam-client/pkg/client/goipam"
)

func RegisterRoutes(router *gin.Engine, client goipam.Client) {
	handler := newHandler(client)
	router.GET("/cidr", handler.GetPrefix)
	router.POST("/cidr", handler.CreatePrefix)
	router.DELETE("/cidr", handler.DeletePrefix)
	router.POST("/ip", handler.AcquireIP)
	router.DELETE("/ip", handler.ReleaseIP)
}
