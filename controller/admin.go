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
		"Admin":       true,
	}

	admin := model.Admin{StaffDetails: model.User{Email: email}}
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
			"Super": admin.Super,
			"Email": admin.StaffDetails.Email,
		},
	}.CreateCookie(w)

	http.Redirect(w, r, "/admin/", 302)
}

func AdminLoginGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.AdminCookieName, Collection: values.AdminCollectionName}
	if err := cookie.CheckCookie(r, w); err == nil {
		http.Redirect(w, r, "/admin/", 302)
		return
	}

	data := map[string]interface{}{
		"SigninError": false,
		"Login":       "/admin/login/",
		"Admin":       true,
	}

	if err := loginTmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		log.Println(err)
	}
}

func AdminPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.AdminCookieName, Collection: values.AdminCollectionName}
	if err := cookie.CheckCookie(r, w); err != nil {
		http.Redirect(w, r, "/admin/login/", 302)
		return
	}

	data := map[string]interface{}{
		"UploadSuccess": false,
	}

	tmpl, terr := template.New("admin.html").Delims("(%", "%)").ParseFiles("views/admin/admin.html", "views/admin/components/tabs.vue",
		"views/admin/components/adduser.vue", "views/admin/components/block.vue", "views/admin/components/messagescan.vue")
	if terr != nil {
		log.Println("could not load template in AdminPage function", terr)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("could not execute template in AdminPage function", err)
	}
}

func UploadUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.AdminCookieName, Collection: values.AdminCollectionName}
	err := cookie.CheckCookie(r, w)
	if err != nil {
		http.Redirect(w, r, "/adnin/login", 302)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	name := r.FormValue("name")
	dateOfBirth := r.FormValue("DOB")
	usersClass := r.FormValue("usersClass")
	faculty := r.FormValue("faculty")

	user := model.User{Email: email, Name: name, DOB: dateOfBirth, Class: usersClass, Faculty: faculty}
	err = model.UploadUser(user, r)

	data := map[string]interface{}{
		"UploadSuccess": true,
		"Error":         false,
	}
	if err != nil {
		data["UploadSuccess"] = false
		data["Error"] = true
	}

	tmpl, terr := template.New("admin.html").Delims("(%", "%)").ParseFiles("views/admin/admin.html", "views/admin/components/tabs.vue",
		"views/admin/components/adduser.vue", "views/admin/components/block.vue", "views/admin/components/messagescan.vue")
	if terr != nil {
		log.Println("could not load template in UploadUser function", terr)
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("could not execute template in UploadUser function", err)
	}
}