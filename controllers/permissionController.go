package controllers

import (
	"github.com/gin-gonic/gin"
	"http-server-with-auth/models"
)

var GivePermission = func(c* gin.Context) {
	permission := &models.Permission{}
	err := c.ShouldBindJSON(permission)
	if err != nil{
		c.JSON(422, gin.H{
			"error": true,
			"message": "invalid request body",
		})
		return
	}
	permission.Create(c)
}
