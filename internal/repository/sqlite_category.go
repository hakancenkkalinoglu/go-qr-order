package repository

import (
	"database/sql"
	"go-qr-order/internal/models"
)

// SQLiteCategoryRepo stores categories and nested products in SQLite.
type SQLiteCategoryRepo struct {
	db *sql.DB
}

func NewSQLiteCategoryRepo(db *sql.DB) *SQLiteCategoryRepo {
	return &SQLiteCategoryRepo{db: db}
}

var _ CategoryRepository = (*SQLiteCategoryRepo)(nil)

func (r *SQLiteCategoryRepo) loadProductsForCategory(categoryID int) ([]models.Product, error) {
	rows, err := r.db.Query(`
		SELECT id, category_id, name, quantity, price
		FROM products WHERE category_id = ?
		ORDER BY id`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Quantity, &p.Price); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *SQLiteCategoryRepo) SaveCategory(category models.Category) models.Category {
	tx, err := r.db.Begin()
	if err != nil {
		return category
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec(`INSERT INTO categories (name) VALUES (?)`, category.Name)
	if err != nil {
		return category
	}
	cid, err := res.LastInsertId()
	if err != nil {
		return category
	}
	category.ID = int(cid)

	items := make([]models.Product, 0, len(category.CategoryItems))
	for _, p := range category.CategoryItems {
		res, err := tx.Exec(`
			INSERT INTO products (category_id, name, quantity, price)
			VALUES (?, ?, ?, ?)`,
			category.ID, p.Name, p.Quantity, p.Price)
		if err != nil {
			return category
		}
		pid, err := res.LastInsertId()
		if err != nil {
			return category
		}
		p.ID = int(pid)
		p.CategoryID = category.ID
		items = append(items, p)
	}
	category.CategoryItems = items

	if err := tx.Commit(); err != nil {
		return category
	}
	return category
}

func (r *SQLiteCategoryRepo) GetCategoryByID(id int) (models.Category, bool) {
	row := r.db.QueryRow(`SELECT id, name FROM categories WHERE id = ?`, id)
	var c models.Category
	if err := row.Scan(&c.ID, &c.Name); err != nil {
		return models.Category{}, false
	}
	items, err := r.loadProductsForCategory(id)
	if err != nil {
		c.CategoryItems = nil
		return c, true
	}
	c.CategoryItems = items
	return c, true
}

func (r *SQLiteCategoryRepo) GetAllCategories() []models.Category {
	rows, err := r.db.Query(`SELECT id, name FROM categories ORDER BY id`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var cats []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			continue
		}
		cats = append(cats, c)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	if len(cats) == 0 {
		return []models.Category{}
	}

	ids := make([]int, len(cats))
	for i := range cats {
		ids[i] = cats[i].ID
	}
	ph := placeholders(len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	q := `SELECT id, category_id, name, quantity, price FROM products WHERE category_id IN (` + ph + `) ORDER BY category_id, id`
	itemRows, err := r.db.Query(q, args...)
	if err != nil {
		return cats
	}
	defer itemRows.Close()

	byCat := make(map[int][]models.Product)
	for itemRows.Next() {
		var p models.Product
		if err := itemRows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Quantity, &p.Price); err != nil {
			continue
		}
		byCat[p.CategoryID] = append(byCat[p.CategoryID], p)
	}
	for i := range cats {
		cats[i].CategoryItems = byCat[cats[i].ID]
	}
	return cats
}
