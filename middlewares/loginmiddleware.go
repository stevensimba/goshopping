package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var enverr = godotenv.Load()
var store = sessions.NewCookieStore([]byte(os.Getenv("sessionkey")))

func Auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "mylogins")

		if _, ok := session.Values["username"]; !ok {
			http.Redirect(w, r, "/account/login", 302)
			return
		}
		HandlerFunc.ServeHTTP(w, r)

	}
}

func Hidefiles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if request endswith "/"
		if strings.HasSuffix(r.URL.Path, "/") {

			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
