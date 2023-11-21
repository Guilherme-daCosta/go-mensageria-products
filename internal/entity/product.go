package entity

import "github.com/google/uuid"

type ProductRepository interface {
	Create(product *Product) error
	FindAll() ([]*Product, error)
}

type Product struct {
	Id    string
	Name  string
	Price float64
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		Id:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}
