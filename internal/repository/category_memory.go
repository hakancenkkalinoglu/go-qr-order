package repository

import (
	"go-qr-order/internal/models"
	"sort"
)

type InMemoryCategoryRepo struct {
	categories    map[int]models.Category
	lastID        int
	productLastID int
}

func NewInCategoryRepo() *InMemoryCategoryRepo {
	return &InMemoryCategoryRepo{
		categories: make(map[int]models.Category),
	}
}

func (c *InMemoryCategoryRepo) SaveCategory(category models.Category) models.Category {
	c.lastID++
	category.ID = c.lastID

	items := make([]models.Product, 0, len(category.CategoryItems))
	for _, p := range category.CategoryItems {
		c.productLastID++
		p.ID = c.productLastID
		p.CategoryID = category.ID
		items = append(items, p)
	}
	category.CategoryItems = items

	c.categories[category.ID] = category
	return category
}

func (c *InMemoryCategoryRepo) GetCategoryByID(id int) (models.Category, bool) {
	cat, ok := c.categories[id]
	return cat, ok
}

func (c *InMemoryCategoryRepo) GetAllCategories() []models.Category {
	keys := make([]int, 0, len(c.categories))
	for k := range c.categories {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	out := make([]models.Category, 0, len(keys))
	for _, k := range keys {
		out = append(out, c.categories[k])
	}
	return out
}

var _ CategoryRepository = (*InMemoryCategoryRepo)(nil)
