package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("secretkey")))

// Ensure a user is loggedin before giving permission
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

// Hide static files (prevent direct navigation)
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
