package controller

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/metaclips/FinalYearProject/model"
	"github.com/metaclips/FinalYearProject/values"
)

// ServeWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}
	if err := cookie.CheckCookie(r, w); err != nil {
		http.Error(w, "Unable to authenticate", 404)
		return
	}

	// todo: add authentications
	ws, err := model.Upgrader.Upgrade(w, r, nil)
	// Get user ID from cookies.
	log.Println("User", cookie.Email, "Connected")
	if err != nil {
		log.Println(err)
		return
	}
	c := &model.Connection{Send: make(chan []byte, 256), WS: ws}

	s := model.Subscription{c, cookie.Email}
	model.HubConstruct.Register <- s
	go s.WritePump()
	s.ReadPump(cookie.Email)
}
