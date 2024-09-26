package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func InternalServerError(c *gin.Context, err error) {
	response := Response{
		Status:     "error",
		StatusCode: http.StatusInternalServerError,
		Data:       nil,
		Message:    err.Error(),
	}
	c.JSON(http.StatusInternalServerError, response)
}

func BadRequest(c *gin.Context, err error) {
	response := Response{
		Status:     "error",
		StatusCode: http.StatusBadRequest,
		Data:       nil,
		Message:    err.Error(),
	}
	c.JSON(http.StatusBadRequest, response)
}
