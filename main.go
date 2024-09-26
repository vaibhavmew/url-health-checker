package main

import (
	"example/controller"
	"example/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//dialer = service.New(service.NewClient())
	dialer := service.New(service.NewMockClient())

	go dialer.StartHealth()

	handler := controller.New(dialer)

	r.GET("/", handler.Get)
	r.GET("/modify", handler.Modify)
	r.Run(":8080")

}
