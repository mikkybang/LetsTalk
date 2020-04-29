package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/metaclips/FinalYearProject/controller"
	"github.com/metaclips/FinalYearProject/model"
)

func main() {
	gob.Register(time.Time{})
	model.InitDB()
	router := httprouter.New()

	router.GET("/", controller.HomePage)
	router.GET("/login", controller.HomePageLoginGet)
	router.GET("/admin/login", controller.AdminLoginGET)
	router.GET("/admin", controller.AdminPage)

	router.POST("/login", controller.HomePageLoginPost)
	router.POST("/admin/login", controller.AdminLoginPOST)
	router.POST("/admin/upload", controller.UploadUser)

	router.ServeFiles("/assets/*filepath", http.Dir("./views/assets"))
	fmt.Println("Webserver UP")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Webserver DOWN")
}
