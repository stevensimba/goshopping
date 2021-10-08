package main

import (
	"net/http"

	"github.com/stevensimba/shopcart/controllers/accountcontroller"
	"github.com/stevensimba/shopcart/controllers/cartcontroller"
	"github.com/stevensimba/shopcart/controllers/productcontroller"
	"github.com/stevensimba/shopcart/middlewares"
)

func main() {
	// ending paths in slash /path/ accepts both
	// /path and /path/, it is better to end in "/"
	// with exception of get (/path?id=2) paths

	// to hide access to all files in /static/ path
	//either  put an index.html file in all subfolders
	// or use a middleware that captures requests ending
	// with "/", either redirect or file not found error
	pix := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", pix))

	// serve a pdf: <a href="/public/file.pdf">read</a>
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", middlewares.Hidefiles(http.StripPrefix("/public/", fs)))

	http.HandleFunc("/", productcontroller.Index)
	http.HandleFunc("/index/", productcontroller.Index)
	http.HandleFunc("/product/", productcontroller.Product)
	http.HandleFunc("/product/add", productcontroller.AddProduct)
	http.HandleFunc("/product/process", productcontroller.Process)

	http.HandleFunc("/cart/", middlewares.Auth(cartcontroller.Index))
	http.HandleFunc("/cart/index/", cartcontroller.Index)
	http.HandleFunc("/cart/buy", cartcontroller.Buy)
	http.HandleFunc("/exit", cartcontroller.Exitcart)
	http.HandleFunc("/Remove", cartcontroller.Remove)

	http.HandleFunc("/account/", accountcontroller.Register)
	http.HandleFunc("/account/register/", accountcontroller.Register)
	http.HandleFunc("/account/registerAuth", accountcontroller.RegisterAuth)
	http.HandleFunc("/login/", accountcontroller.Login)
	http.HandleFunc("/account/login/", accountcontroller.Login)
	http.HandleFunc("/account/loginAuth", accountcontroller.LoginAuth)
	http.HandleFunc("/account/logout", middlewares.Auth(accountcontroller.Logout))

	http.ListenAndServe(":3000", nil)

}
