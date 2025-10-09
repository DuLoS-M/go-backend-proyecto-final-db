package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExampleController struct{}

func NewExampleController() *ExampleController {
	return &ExampleController{}
}

func (ec *ExampleController) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello from ExampleController!"})
}
