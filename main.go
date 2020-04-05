package main

import (
	"encoding/gob"
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
	router.GET("/login/", controller.LoginPage)
	router.GET("/admin/login/", controller.AdminLoginGET)
	router.GET("/admin/", controller.AdminPage)

	router.POST("/admin/login/", controller.AdminLoginPOST)

	router.ServeFiles("/assets/*filepath", http.Dir("./views/assets"))
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln(err)
	}
}
