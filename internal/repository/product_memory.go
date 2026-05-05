package repository

import "go-qr-order/internal/models"

type InMemoryProductRepo struct {
	products map[int]models.Product
	lastID   int
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	return &InMemoryProductRepo{
		products: make(map[int]models.Product),
	}
}

func (p *InMemoryProductRepo) Save(product models.Product) models.Product {
	p.lastID++

	product.ID = p.lastID
	p.products[product.ID] = product

	return product
}
