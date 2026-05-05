package repository

import (
	"go-qr-order/internal/models"
	"sort"
)

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

func (p *InMemoryProductRepo) GetProductByID(id int) (models.Product, bool) {
	pr, ok := p.products[id]
	return pr, ok
}

func (p *InMemoryProductRepo) ListProductsByCategoryID(categoryID int) []models.Product {
	var out []models.Product
	for _, pr := range p.products {
		if pr.CategoryID == categoryID {
			out = append(out, pr)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

var _ ProductRepository = (*InMemoryProductRepo)(nil)
