package models

import (
	"context"
	"fmt"
	"time"

	"github.com/stevensimba/shopcart/config"
	"github.com/stevensimba/shopcart/entities"
)

type ProductModel struct {
}

func (*ProductModel) Findall() ([]entities.Product, error) {
	db, err := config.DbConn()
	if err != nil {
		return nil, err
	} else {
		rows, err2 := db.Query("Select * from product")

		if err2 != nil {
			return nil, err2
		}

		var products []entities.Product
		for rows.Next() {
			var product entities.Product
			rows.Scan(&product.Id, &product.Name, &product.Price, &product.Quantity, &product.Photo)

			products = append(products, product)
		}
		return products, nil
	}
}

func (*ProductModel) Find(id int64) (entities.Product, error) {
	db, err := config.DbConn()
	if err != nil {
		return entities.Product{}, err
	} else {
		rows, err2 := db.Query("Select * from product where id = ?", id)

		if err2 != nil {
			return entities.Product{}, err2
		}
		var product entities.Product
		for rows.Next() {
			rows.Scan(&product.Id, &product.Name, &product.Price, &product.Quantity, &product.Photo)
		}
		return product, nil
	}
}

func InsertProduct(p entities.Product) {
	db, _ := config.DbConn()
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, "insert into product(name, price, quantity, photo) values(?, ?, ?, ?)", p.Name, p.Price, p.Quantity, p.Photo)
	if err != nil {
		fmt.Println("db error", err)
	}
	lastId, _ := res.LastInsertId()
	fmt.Println("product number", lastId)

}
