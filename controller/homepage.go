package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/metaclips/FinalYearProject/model"
	"github.com/metaclips/FinalYearProject/values"
)

// ToDo: initially parse html template one time only
func HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}
	if err := cookie.CheckCookie(r, w); err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	uuid, ok := cookie.Data["UUID"].(string)
	if !ok {
		http.Error(w, "Could not retrieve UUID", 404)
		log.Println("Could not retrieve UUID in homepage")
		return
	}

	userInfo, err := model.GetAllUserRooms(cookie.Email)
	if err != nil {
		log.Println("Could not fetch users room", cookie.Email)
	}

	data := map[string]interface{}{
		"Email":      cookie.Email,
		"UUID":       uuid,
		"UsersRooms": userInfo.RoomsJoined,
		"Requests":   userInfo.JoinRequest,
		"Name":       values.Users[cookie.Email],
	}

	// Use (%%) instead of {{}} for templates.
	tmpl := template.Must(template.New("home.html").Delims("(%", "%)").ParseFiles("views/homepage/home.html",
		"views/homepage/components/SideBar.vue", "views/homepage/components/ChattingComponent.vue"))

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func HomePageLoginGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := map[string]interface{}{
		"SigninError": false,
		"Login":       "/login",
		"Admin":       false,
	}

	tmpl, terr := template.New("login.html").Delims("(%", "%)").ParseFiles("views/loginpage/login.html")
	if terr != nil {
		log.Fatalln(terr)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func HomePageLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := model.User{
		Email: email,
	}.CreateUserLogin(password, w)

	if err != nil {
		data := map[string]interface{}{
			"SigninError": true,
			"Login":       "/login",
			"Admin":       false,
		}

		tmpl, terr := template.New("login.html").Delims("(%", "%)").ParseFiles("views/loginpage/login.html")
		if terr != nil {
			log.Fatalln(terr)
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println(err)
		}
		return
	}

	http.Redirect(w, r, "/", 302)
}

// todo: Use API instead..
func SearchUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := params.ByName("ID")
	uniqueID := params.ByName("UUID")
	key := params.ByName("Key")
	if id == "" {
		log.Println("No id was specified while searching for user")
		http.Error(w, "Not found", 404)
		return
	}

	// todo: do we need to still validate???
	// Are matric number meant to be confidential???
	err := model.User{}.ValidateUser(id, uniqueID)
	if err != nil {
		log.Println("No id was specified while searching for user")
		http.Error(w, "Not found", 404)
		return
	}

	users := model.GetUser(key, id)
	data := map[string]interface{}{
		"UsersFound": users,
	}
	bytes, err := json.MarshalIndent(&data, "", "\t")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, "error sending information", 400)
		return
	}
}
