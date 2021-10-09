package productcontroller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"text/template"

	"github.com/stevensimba/goshopping/entities"
	"github.com/stevensimba/goshopping/models"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("secretkey")))

// the homepage handler uses FindAll() to fetch all product data & display
func Index(w http.ResponseWriter, r *http.Request) {
	var username string
	session, _ := store.Get(r, "mylogins")

	if session.Values["username"] != nil {
		username = session.Values["username"].(string)
	}

	var productModel models.ProductModel
	products, err := productModel.Findall()
	if err != nil {
		fmt.Println(err)
	}
	data := map[string]interface{}{
		"products": products,
		"username": username,
	}
	tmp, err := template.ParseFiles("views/product/index.html")

	if err != nil {
		fmt.Println(err)
	}
	tmp.Execute(w, data)
}

// The Product() function use id to fetch more details about a product
func Product(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id, _ := strconv.ParseInt(query.Get("id"), 10, 64)

	var productModel models.ProductModel
	product, err := productModel.Find(id)
	if err != nil {
		fmt.Println(err)
	}

	tmp, err := template.ParseFiles("views/product/product.html")

	if err != nil {
		fmt.Println(err)
	}
	tmp.Execute(w, product)
}

// A form used to submit new products
func AddProduct(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("views/product/addproduct.html")
	tmp.Execute(w, nil)
}

// The Process() function helps to save new products in the db;
// It validates the file size and saves a new photo in the static folder
func Process(w http.ResponseWriter, r *http.Request) {
	var product entities.Product
	product.Name = r.FormValue("productname")
	product.Price, _ = strconv.ParseFloat(r.FormValue("price"), 64)
	product.Quantity, _ = strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	ps := regexp.MustCompile(" ")
	photoname := ps.ReplaceAllString(product.Name, "-")
	product.Photo = photoname + ".jpg"
	models.InsertProduct(product)

	// 1024 * 1024 == 1MB
	if err := r.ParseMultipartForm(1024 * 1024 * 5); err != nil {
		http.Error(w, "the upload file is too big. Choose a smaller image", http.StatusInternalServerError)
	}
	//file, fileHeader, err
	file, _, _ := r.FormFile("photo")
	defer file.Close()

	dst, _ := os.Create(fmt.Sprintf("./static/images/%s", product.Photo))
	defer dst.Close()

	io.Copy(dst, file)
	http.Redirect(w, r, "/index", http.StatusSeeOther)

}
