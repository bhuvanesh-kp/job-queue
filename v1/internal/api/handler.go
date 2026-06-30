package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type controller struct{}

func NewController() *controller {
	return &controller{}
}

func (ctr *controller) AddNewJobToQueue(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if req.ReqType != "sleep" {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "Invalid operation",
		})
		return
	}

	// worker call

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
	})
}
