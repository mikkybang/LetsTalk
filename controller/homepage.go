package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/metaclips/LetsTalk/model"
	"github.com/metaclips/LetsTalk/values"
)

// TODO: initially parse html template one time only
func HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}
	if err := cookie.CheckCookie(r, w); err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	data := struct {
		Email string
		UUID  string
		Name  string
	}{
		cookie.Email,
		cookie.Data.UUID,
		values.MapEmailToName[cookie.Email],
	}

	// Use (%%) instead of {{}} for templates.
	tmpl := template.Must(template.New("home.html").Delims("(%", "%)").ParseFiles(
		"views/homepage/home.html",
		"views/homepage/components/SideBar.vue", "views/homepage/components/ChattingComponent.vue"))

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func HomePageLoginGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := setLoginDetails(false, false, "", "/login")

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
	email = strings.ToLower(email)
	password := r.FormValue("password")

	err := model.User{
		Email: email,
	}.CreateUserLogin(password, w)

	if err != nil {
		data := setLoginDetails(true, false, "", "/login")

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

// TODO: Use as API instead..
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

	err := model.User{Email: id}.ValidateUser(uniqueID)
	if err != nil {
		log.Println("No id was specified while searching for user")
		http.Error(w, "Not found", 404)
		return
	}

	data := struct {
		UsersFound []string
	}{
		model.GetUser(key, id),
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		http.Error(w, values.ErrMarshal.Error(), 400)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, values.ErrWrite.Error(), 400)
		return
	}
}

func setLoginDetails(errors, isAdmin bool, errorDetail, link string) struct {
	SigninError bool
	Admin       bool
	Login       string
	ErrorDetail string
} {

	return struct {
		SigninError bool
		Admin       bool
		Login       string
		ErrorDetail string
	}{
		errors,
		isAdmin,
		link,
		errorDetail,
	}
}
