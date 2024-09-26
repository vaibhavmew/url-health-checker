package controller

import (
	"example/common"
	"example/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	dialer *service.Dialer
}

func New(dialer *service.Dialer) *Handler {
	return &Handler{
		dialer: dialer,
	}
}

func (h *Handler) Get(c *gin.Context) {
	var req service.Req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.BadRequest(c, err)
		return
	}

	res, err := h.dialer.Middleware(req)
	if err != nil {
		common.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)

}

//for testing purposes
func (h *Handler) Modify(c *gin.Context) {
	var req service.Mock
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.BadRequest(c, err)
		return
	}

	if req.Action == "u" {
		h.dialer.Client.Up(req.URL)
	} else {
		h.dialer.Client.Down(req.URL)
	}

	c.JSON(http.StatusOK, nil)

}
