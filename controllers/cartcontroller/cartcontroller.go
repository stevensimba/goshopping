package cartcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/stevensimba/goshopping/entities"
	"github.com/stevensimba/goshopping/models"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("secretkey")))

// The cart page helps to display the cart object values
func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	strCart := session.Values["cart"].(string)

	var cart []entities.Item
	json.Unmarshal([]byte(strCart), &cart)

	// sub-total for each cart item ( price * no. of items )
	tmp, err := template.New("index.html").Funcs(template.FuncMap{
		"total": func(item entities.Item) float64 {
			return item.Product.Price * float64(item.Quantity)
		},
	}).ParseFiles("views/cart/index.html")

	if err != nil {
		fmt.Println(err)
	}

	data := map[string]interface{}{
		"cart": cart,
		// total(cart): sum of all items
		"totals": total(cart),
	}
	tmp.Execute(w, data)
}

// Buy creates a mysession cookie to save cart object data
// After a purchase redirect to the cart page
func Buy(w http.ResponseWriter, r *http.Request) {

	// use the id value in the http request to fetch the product from the db
	query := r.URL.Query()
	id, _ := strconv.ParseInt(query.Get("id"), 10, 64)

	var productModel models.ProductModel
	product, err := productModel.Find(id)
	if err != nil {
		fmt.Println(err)
	}

	session, _ := store.Get(r, "mysession")
	cart := session.Values["cart"]
	if cart == nil {
		var cart []entities.Item
		cart = append(cart, entities.Item{
			Product:  product,
			Quantity: 1,
		})
		bytesCart, _ := json.Marshal(cart)
		session.Values["cart"] = string(bytesCart)
		session.Save(r, w)
	} else {

		session, _ := store.Get(r, "mysession")
		strCart := session.Values["cart"].(string)

		// convert the json value into a struct
		var cart []entities.Item
		json.Unmarshal([]byte(strCart), &cart)
		index := exists(id, cart)

		if index == -1 {
			cart = append(cart, entities.Item{
				Product:  product,
				Quantity: 1,
			})
		} else {
			cart[index].Quantity++
		}
		bytesCart, _ := json.Marshal(cart)
		session.Values["cart"] = string(bytesCart)
	}

	session.Save(r, w)
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func Exitcart(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func exists(id int64, cart []entities.Item) int {
	for i := 0; i < len(cart); i++ {
		if cart[i].Product.Id == id {
			return i
		}
	}
	return -1
}

func total(cart []entities.Item) float64 {
	var sum float64 = 0

	for _, item := range cart {
		sum += item.Product.Price * float64(item.Quantity)

	}
	return sum
}

func Remove(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id, _ := strconv.ParseInt(query.Get("id"), 10, 64)

	session, _ := store.Get(r, "mysession")
	strCart := session.Values["cart"].(string)

	var cart []entities.Item
	json.Unmarshal([]byte(strCart), &cart)

	index := exists(id, cart)
	cart = append(cart[:index], cart[index+1:]...)

	bytesCart, _ := json.Marshal(cart)
	session.Values["cart"] = string(bytesCart)

	session.Save(r, w)
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}
