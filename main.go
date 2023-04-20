package main

import (
	"fmt"
	"net/http"
	"regexp"
	"text/template"
	"time"
)

var templ *template.Template

type User struct {
	userName string
	email    string
	password string
	gustUser string
}

var UserDetails User

var emailinlogin string
var passwordinlogin string

func main() {
	UserDetails.gustUser = "gust user"
	templ = template.Must(template.ParseGlob("template/*.html"))
	http.HandleFunc("/", loginPage)
	http.HandleFunc("/sign", signupPage)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/logout", logoutPage)
	port := ":8080"
	http.ListenAndServe(port, nil)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		passwordinlogin = r.FormValue("password")
		emailinlogin = r.FormValue("email")

		if emailinlogin == "" || passwordinlogin == "" {
			fmt.Println("enter the login data")
			templ.ExecuteTemplate(w, "login.html", nil)
			return
		}

		if UserDetails.email == emailinlogin && UserDetails.password == passwordinlogin {
			fmt.Println("loged")
			exp := time.Now().Add(5 * time.Minute)
			cookie := http.Cookie{
				Name:    "user",
				Value:   UserDetails.email,
				Expires: exp,
			}
			http.SetCookie(w, &cookie)
		}
		templ.ExecuteTemplate(w, "homePage.html", UserDetails.userName)
	} else {
		cookie, err := r.Cookie("user")
		if err != nil {
			fmt.Println("erroor in cookeie")
			templ.ExecuteTemplate(w, "login.html", nil)
			return
		} else if cookie.Value == UserDetails.email && cookie.Expires.Before(time.Now()) {

			fmt.Println("home page sucseess")
			templ.ExecuteTemplate(w, "homePage.html", UserDetails.userName)
			return

		} else {
			templ.ExecuteTemplate(w, "login.html", nil)
			fmt.Println("please re logine")
		}
	}
}

func signupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("sighn post")
		userName := r.FormValue("username")
		password := r.FormValue("password")
		Cpassword := r.FormValue("confirm-password")
		email := r.FormValue("email")

		if userName == "" || password == "" || email == "" || Cpassword == "" {
			// 			// view engine
			templ.ExecuteTemplate(w, "sighnup.html", nil)
		}
		if password != Cpassword || !isValidEmail(email) {
			// view engine
			fmt.Println(" re entr data")
			return
		}

		UserDetails.userName = userName
		UserDetails.password = password
		UserDetails.email = email
		// view engine set
		templ.ExecuteTemplate(w, "login.html", "Acount created succesfuly pls login")
		fmt.Println("signup sucsess")
		fmt.Println(UserDetails)

	} else {
		fmt.Println("else in sighn")
		templ.ExecuteTemplate(w, "sighnup.html", nil)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("user")
	if cookie.Value == UserDetails.email && cookie.Expires.Before(time.Now()) {
		fmt.Println(UserDetails)
		fmt.Println("home page sucseess")
		fmt.Println(UserDetails)

		templ.ExecuteTemplate(w, "homePage.html", UserDetails.userName)
		//pas msg
		return
	} else {

		templ.ExecuteTemplate(w, "homePage.html", UserDetails.gustUser)

		//pas msg
	}
}

func logoutPage(w http.ResponseWriter, r *http.Request) {

	exp := time.Now().Add(-1)
	cookie := http.Cookie{
		Name:    "user",
		Value:   "",
		Expires: exp,
	} //dlt
	http.SetCookie(w, &cookie)
	//http.Redirect(w, r, "/", http.StatusOK)
	templ.ExecuteTemplate(w, "login.html", nil)
}

func isValidEmail(email string) bool {
	// Regular expression to validate email address
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile regular expression
	regex := regexp.MustCompile(pattern)

	// Match email address against regular expression
	return regex.MatchString(email)
}
