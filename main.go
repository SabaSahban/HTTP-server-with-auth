package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/mux"
	"http-server-with-auth/app"
	"http-server-with-auth/controllers"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()
	router.POST("/user/new", controllers.Register)
	router.POST("/user/login", app.AuthorizeJWT, controllers.Authenticate)
	router.POST("/user/upload", controllers.Upload)
	router.GET("/user/download", app.AuthorizeJWT, controllers.Download)
	router.POST("/user/give/permission", controllers.GivePermission)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8040" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
