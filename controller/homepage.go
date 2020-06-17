package controller

import (
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
		Email, UUID, Name string
	}{
		cookie.Email, cookie.Data.UUID,
		values.MapEmailToName[cookie.Email],
	}

	// Use (%%) instead of {{}} for templates.
	tmpl := template.Must(template.New("home.html").Delims("(%", "%)").ParseFiles(
		"views/homepage/home.html",
		"views/homepage/components/SideBar.vue", "views/homepage/components/ChattingComponent.vue", "views/homepage/components/CallUI.vue"))

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

func setLoginDetails(errors, isAdmin bool, errorDetail, link string) struct {
	SigninError, Admin bool
	Login, ErrorDetail string
} {

	return struct {
		SigninError, Admin bool
		Login, ErrorDetail string
	}{
		errors,
		isAdmin,
		link,
		errorDetail,
	}
}
