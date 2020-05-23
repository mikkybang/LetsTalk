package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/metaclips/LetsTalk/controller"
	"github.com/metaclips/LetsTalk/model"
)

func main() {
	defer func() {
		fmt.Println("Webserver DOWN")
	}()

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

	router.ServeFiles("/assets/*filepath", http.Dir("./views/assets"))
	fmt.Println("Webserver UP")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln(err)
	}
}
