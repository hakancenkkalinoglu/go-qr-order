package repository

import (
	"database/sql"
	"go-qr-order/internal/models"
)

// SQLiteProductRepo stores catalog products in SQLite.
type SQLiteProductRepo struct {
	db *sql.DB
}

func NewSQLiteProductRepo(db *sql.DB) *SQLiteProductRepo {
	return &SQLiteProductRepo{db: db}
}

var _ ProductRepository = (*SQLiteProductRepo)(nil)

func (r *SQLiteProductRepo) Save(product models.Product) models.Product {
	res, err := r.db.Exec(`
		INSERT INTO products (category_id, name, quantity, price)
		VALUES (?, ?, ?, ?)`,
		product.CategoryID, product.Name, product.Quantity, product.Price)
	if err != nil {
		return product
	}
	id, err := res.LastInsertId()
	if err != nil {
		return product
	}
	product.ID = int(id)
	return product
}

func (r *SQLiteProductRepo) GetProductByID(id int) (models.Product, bool) {
	row := r.db.QueryRow(`
		SELECT id, category_id, name, quantity, price
		FROM products WHERE id = ?`, id)
	var p models.Product
	if err := row.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Quantity, &p.Price); err != nil {
		return models.Product{}, false
	}
	return p, true
}

func (r *SQLiteProductRepo) ListProductsByCategoryID(categoryID int) []models.Product {
	rows, err := r.db.Query(`
		SELECT id, category_id, name, quantity, price
		FROM products WHERE category_id = ?
		ORDER BY id`, categoryID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Quantity, &p.Price); err != nil {
			continue
		}
		list = append(list, p)
	}
	return list
}
