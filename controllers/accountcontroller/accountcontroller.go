package accountcontroller

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"unicode"

	"github.com/stevensimba/goshopping/config"
	"github.com/stevensimba/goshopping/entities"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("sessionkey")))

// Serve a registration form
func Register(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseGlob("views/accountcontroller/*.html")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

// Register a new user using the username, password, gender and age
func RegisterAuth(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	r.ParseForm()
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Age, _ = strconv.Atoi(r.FormValue("age"))

	// the default gender is female
	if r.FormValue("gender") == "true" {
		user.Male = true
	} else {
		user.Male = false
	}

	var (
		alphaNumeric = true
		pwdLength    = false
	)

	// validate the password is between 3 and 25 characters
	if 3 < len(user.Password) && len(user.Password) < 25 {
		pwdLength = true
	}

	// validate username: ensure every character is alphanumeric
	for _, char := range user.Username {
		//unicode.IsLower, IsUpper, IsSymbol(char), IsSpace(int(char))
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			fmt.Printf("not letter/number %c", char)
			alphaNumeric = false
		}
	}

	tpl, _ := template.ParseGlob("views/accountcontroller/*.html")
	db, _ := config.DbConn()
	var Uid string
	var hash []byte
	var insertStmt *sql.Stmt
	var result sql.Result

	// validate password and the username
	if !pwdLength || !alphaNumeric {
		tpl.ExecuteTemplate(w, "register.html", "Password: 3-5, username is letters & numbers")
		return
	} else {
		data := map[string]interface{}{
			"user": user,
		}

		// save data in the database, only if the username is new
		rows := db.QueryRow("select id from users where username = ? ", user.Username)
		err := rows.Scan(&Uid)

		if err != sql.ErrNoRows {
			tpl.ExecuteTemplate(w, "register.html", "username already exists")
		} else {
			// encrypt password
			hash, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

			// insert the registration data in the user table
			insertStmt, _ = db.Prepare("insert into users (username, password, age, male) values(?, ?, ?, ?);")
			defer insertStmt.Close()

			result, err = insertStmt.Exec(user.Username, string(hash), user.Age, user.Male)
			rowsAff, _ := result.RowsAffected()
			lastIns, _ := result.LastInsertId()
			fmt.Printf("rows affected: %d, lastid: %d \n", rowsAff, lastIns)
			if err != nil {
				fmt.Println("error inserting new user")
			}

			tmp, _ := template.ParseFiles(
				"views/templates/mytemplate.html",
				"views/accountcontroller/thanks.html",
			)
			tmp.ExecuteTemplate(w, "layout", data)
		}

	}

}

// serve the login form
func Login(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseGlob("views/accountcontroller/*.html")

	tpl.ExecuteTemplate(w, "login.html", nil)
	return
}

// Authenticate the user login information
var login entities.Login

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseGlob("views/accountcontroller/*.html")

	r.ParseForm()
	login.Username = r.FormValue("username")
	login.Password = r.FormValue("password")

	// if the user is already loggedin, redirect to the homepage
	session, _ := store.Get(r, "mylogins")
	if session.Values["username"] == login.Username {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	db, _ := config.DbConn()
	var Pwd string

	// match the username with password
	rows := db.QueryRow("select password from users where username = ? ", login.Username)

	// if the username does not exist send a mismatch error
	err := rows.Scan(&Pwd)
	if err == sql.ErrNoRows {
		tpl.ExecuteTemplate(w, "login.html", "username and password mismatch")
	} else {
		// if username exists but password does not match
		// send a check username / password error
		err = bcrypt.CompareHashAndPassword([]byte(Pwd), []byte(login.Password))
		if err != nil {
			tpl.ExecuteTemplate(w, "login.html", "check username and password")
		} else {
			// create a new session with the username logged in
			session, _ := store.Get(r, "mylogins")
			session.Values["username"] = login.Username
			session.Save(r, w)
			// fecth the base and the include html file and send a success message
			tmp, _ := template.ParseFiles(
				"views/templates/mytemplate.html",
				"views/accountcontroller/loggedin.html",
			)

			tmp.ExecuteTemplate(w, "layout", "You have successfully loggedin")
		}
	}
}

// To logout a user, se the session expiry in the past by a second
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mylogins")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
