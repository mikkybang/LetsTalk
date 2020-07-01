package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/metaclips/LetsTalk/controller"
	"github.com/metaclips/LetsTalk/model"
)

func main() {
	file, err := os.Open("config.json") // For read access.
	if err != nil {
		log.Fatal("Error loading the config file")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}
	log.Println(Config.DB_HOST)

	gob.Register(time.Time{})
	model.InitDB()
	router := httprouter.New()
	go model.HubConstruct.Run()

	router.GET("/", controller.HomePage)
	router.GET("/ws", controller.ServeWs)
	router.GET("/login", controller.HomePageLoginGet)
	router.GET("/admin/login", controller.AdminLoginGET)
	router.GET("/admin", controller.AdminPage)

	router.POST("/login", controller.HomePageLoginPost)
	router.POST("/admin/login", controller.AdminLoginPOST)
	router.POST("/admin/upload", controller.UploadUser)

	port := Config.Port
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

func LoadConfiguraton() {

}
