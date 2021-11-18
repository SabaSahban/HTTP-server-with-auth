package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Permission struct {
	gorm.Model
	FileName string
	UserName string
}

func (permission *Permission) Validate(c *gin.Context) {
	temp := &Permission{}
	err := GetDB().Table("permissions").Where("username = ?", permission.UserName).First(temp).Error
	if err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "connection error, please try again",
		})
		return
	}
	err = GetDB().Table("files").Where("permission = ?", permission.FileName).First(temp).Error
	if err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "connection error, please try again",
		})
		return
	}

	if temp.UserName == "" {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "this username is not registered",
		})
		return
	}

	if temp.FileName == "" {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "this file name is not registered",
		})
		return
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": "permission is given successfully",
	})
	permission.Create(c)
}

func (permission *Permission) Create(c *gin.Context) {
	GetDB().Create(permission)
	c.JSON(200, permission)
}

func IsPermitted(fileName string, id uint) bool {
	account := &Account{}
	userName := db.Table("accounts").Select("email").Where("id = ?", id).Find(account)
	fmt.Println(account)
	err := db.Where(
		db.Where("file_name = ?", fileName).Where(db.Where("user_name = ?", userName)),
	).Find(&Permission{}).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
