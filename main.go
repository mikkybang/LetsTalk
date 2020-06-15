package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"

	"github.com/metaclips/LetsTalk/controller"
	"github.com/metaclips/LetsTalk/model"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	gob.Register(time.Time{})
	model.InitDB()
	router := httprouter.New()
	go model.HubConstruct.Run()

	router.GET("/", controller.HomePage)
	router.GET("/ws", controller.ServeWs)
	router.GET("/login", controller.HomePageLoginGet)
	router.GET("/admin/login", controller.AdminLoginGET)
	router.GET("/admin", controller.AdminPage)
	router.GET("/search/:ID/:UUID/:Key", controller.SearchUser)

	router.POST("/login", controller.HomePageLoginPost)
	router.POST("/admin/login", controller.AdminLoginPOST)
	router.POST("/admin/upload", controller.UploadUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("HTTP_PLATFORM_PORT")
	}
	if port == "" {
		port = "8080"
	}

	router.ServeFiles("/assets/*filepath", http.Dir("./views/assets"))
	fmt.Println("Webserver UP")
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalln(err)
	}
}
