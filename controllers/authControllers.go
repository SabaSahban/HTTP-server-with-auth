package controllers

import (
	"github.com/gin-gonic/gin"
	"http-server-with-auth/models"
)

var Register = func (c *gin.Context) {
	account := &models.Account{}
	err := c.ShouldBindJSON(account)
	if err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "invalid request body",
		})
		return
	}
	account.Validate(c)
}


var Authenticate = func(c *gin.Context) {
	account := &models.Account{}
	err := c.ShouldBindJSON(account)
	if err != nil {
		c.JSON(404, "invalid request body")
		return
	}
	models.Login(account.Email, account.Password, c)
}
