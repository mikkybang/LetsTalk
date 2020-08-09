package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/metaclips/LetsTalk/controller"
	"github.com/metaclips/LetsTalk/model"
	"github.com/metaclips/LetsTalk/values"
)

func main() {
	if err := values.LoadConfiguration("./config.json"); err != nil {
		log.Fatalln("unable to load config", err)
	}

	gob.Register(time.Time{})
	model.InitDB()
	go model.HubConstruct.Run()

	router := httprouter.New()

	router.GET("/", controller.HomePage)
	router.GET("/ws", controller.ServeWs)
	router.GET("/login", controller.HomePageLoginGet)
	router.GET("/admin/login", controller.AdminLoginGET)
	router.GET("/admin", controller.AdminPage)

	router.POST("/login", controller.HomePageLoginPost)
	router.POST("/admin/login", controller.AdminLoginPOST)
	router.POST("/admin/upload", controller.UploadUser)

	port := values.Config.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	router.ServeFiles("/assets/*filepath", http.Dir("./views/assets"))
	log.Println("Webserver UP")

	// Optional use of TLS due to Heroku serving TLS at low level.
	if values.Config.TLS.CertPath != "" && values.Config.TLS.KeyPath != "" {
		if err := http.ListenAndServeTLS(":"+port, values.Config.TLS.CertPath, values.Config.TLS.KeyPath, router); err != nil {
			log.Fatalln(err)
		}

		return
	}

	// Note: without HTTPS users wont be able to login as SetCookie uses Secure flag.
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalln(err)
	}
}
