package accountcontroller

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"unicode"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/stevensimba/shopcart/config"
	"github.com/stevensimba/shopcart/entities"
	"golang.org/x/crypto/bcrypt"
)

var enverr = godotenv.Load()

var store = sessions.NewCookieStore([]byte(os.Getenv("sessionkey")))

func Register(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseGlob("views/accountcontroller/*.html")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

func RegisterAuth(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	r.ParseForm()
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Age, _ = strconv.Atoi(r.FormValue("age"))

	if r.FormValue("gender") == "true" {
		user.Male = true
	} else {
		user.Male = false
	}
	var (
		alphaNumeric = true
		pwdLength    = false
	)

	// validation
	if 3 < len(user.Password) && len(user.Password) < 10 {
		pwdLength = true
	}

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

	if !pwdLength || !alphaNumeric {
		tpl.ExecuteTemplate(w, "register.html", "Password: 3-5, username is letters & numbers")
		return

	} else {

		data := map[string]interface{}{
			"user": user,
		}
		rows := db.QueryRow("select id from users where username = ? ", user.Username)
		err := rows.Scan(&Uid)

		if err != sql.ErrNoRows {

			tpl.ExecuteTemplate(w, "register.html", "username already exists")
		} else {
			hash, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

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
func Login(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseGlob("views/accountcontroller/*.html")

	tpl.ExecuteTemplate(w, "login.html", nil)
	return
}

var login entities.Login

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseGlob("views/accountcontroller/*.html")

	r.ParseForm()
	login.Username = r.FormValue("username")
	login.Password = r.FormValue("password")

	session, _ := store.Get(r, "mylogins")
	if session.Values["username"] == login.Username {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	db, _ := config.DbConn()
	var Pwd string

	rows := db.QueryRow("select password from users where username = ? ", login.Username)

	err := rows.Scan(&Pwd)
	if err == sql.ErrNoRows {
		tpl.ExecuteTemplate(w, "login.html", "username and password mismatch")
	} else {

		err = bcrypt.CompareHashAndPassword([]byte(Pwd), []byte(login.Password))
		if err != nil {
			tpl.ExecuteTemplate(w, "login.html", "check username and password")
		} else {

			// if mylogins did not exists it will be created
			// _, ok := session.Values["username"]
			// if !ok { http.Redirect(w, r, "account/login", 302)}

			session, _ := store.Get(r, "mylogins")
			session.Values["username"] = login.Username
			session.Save(r, w)

			tmp, _ := template.ParseFiles(
				"views/templates/mytemplate.html",
				"views/accountcontroller/loggedin.html",
			)

			tmp.ExecuteTemplate(w, "layout", "You have successfully loggedin")
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mylogins")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

