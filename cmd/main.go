package main

import (
	"entry_task/middlewares"
	"entry_task/view/http/apis"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(middlewares.RequestIDMiddleware())
	router.Use(middlewares.LogMiddleware())

	generalGroup := router.Group("/api/")

	apis.SetMessageApi(generalGroup)
	apis.SetUserApi(generalGroup)

	addressPort := "localhost:8080"
	_ = router.Run(addressPort)
}
