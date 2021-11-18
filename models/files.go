package models

import (
	"github.com/gin-gonic/gin"
)

type File struct {
	FileName string `json:"file_name"`
}

func Validate(c *gin.Context, fileName string){
	temp := &File{}
	GetDB().Table("files").Where("file_name =?", fileName).First(temp)
	if temp.FileName != "" {
		c.JSON(422, gin.H{
			"error": true,
			"message": "filename is already taken",
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"message": "file is created successfully",
	})
	Create(c, fileName)
}

func  Create(c *gin.Context, fileName string) {
	file := File{FileName: fileName}
	GetDB().Create(&file)
	c.JSON(200, gin.H{
		"error": false,
		"message": file.FileName,
	})
}