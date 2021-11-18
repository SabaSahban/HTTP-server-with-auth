package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"http-server-with-auth/models"
	"io"
	"log"
	"net/http"
	"os"
)

var Upload = func(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	out, err := os.Create("public/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(out)
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	filepath := "http://localhost:8080/file/" + filename

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"filepath": filepath,
	})
	models.Validate(c, filename)
}

var Download = func(c *gin.Context) {
	//fmt.Println(c.Request.Context().Value("user").(uint))
	id := c.Value("user").(uint)
	file := &models.File{}
	err := c.ShouldBindJSON(file)
	if err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "invalid request body",
		})
		return
	}
	if !models.IsPermitted(file.FileName, id) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "you are not permitted",
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "you are allowed",
	})

	c.File("./public/" + file.FileName)
}
