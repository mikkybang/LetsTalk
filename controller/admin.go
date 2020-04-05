package controller

import (
	"html/template"
	"log"
	"net/http"

	"github.com/metaclips/FinalYearProject/model"
	"github.com/metaclips/FinalYearProject/values"

	"github.com/julienschmidt/httprouter"
)

var loginTmpl *template.Template

func init() {
	var terr error
	loginTmpl, terr = template.New("login.html").Delims("(%", "%)").ParseFiles("views/loginpage/login.html")
	if terr != nil {
		log.Fatalln(terr)
	}
}

func AdminLoginPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	data := map[string]interface{}{
		"SigninError": false,
		"Login":       "/admin/login/",
	}

	admin := model.Admin{StaffDetails: model.Staff{Email: email}}
	if err := admin.CheckAdminDetails(password); err != nil {
		data["SigninError"] = true
		data["ErrorDetail"] = "Invalid signin details"

		if err := loginTmpl.ExecuteTemplate(w, "login.html", data); err != nil {
			log.Println(err)
		}
		return
	}

	model.CookieDetail{
		Email:      admin.StaffDetails.Email,
		Collection: values.AdminCollectionName,
		CookieName: values.AdminCookieName,
		Path:       "/admin",
		Data: map[string]interface{}{
			"UUID":  admin.StaffDetails.UUID,
			"Super": admin.Super,
			"Email": admin.StaffDetails.Email,
		},
	}.CreateCookie(w)

	http.Redirect(w, r, "/admin/", 302)
}

func AdminLoginGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.AdminCookieName, Collection: values.AdminCollectionName}
	if err := cookie.CheckCookie(r); err == nil {
		http.Redirect(w, r, "/admin/", 302)
		return
	}

	data := map[string]interface{}{
		"SigninError": false,
		"Login":       "/admin/login/",
	}

	if err := loginTmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		log.Println(err)
	}
}

func AdminPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.AdminCookieName, Collection: values.AdminCollectionName}
	if err := cookie.CheckCookie(r); err != nil {
		http.Redirect(w, r, "/admin/login/", 302)
		return
	}

	// compare database UUID with cookie UUID
	tmpl, terr := template.New("admin.html").Delims("(%", "%)").ParseFiles("views/admin/admin.html")
	if terr != nil {
		log.Fatalln(terr)
	}
	tmpl.Execute(w, nil)
}
