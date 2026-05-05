package repository

import (
	"database/sql"
	"go-qr-order/internal/models"
	"strings"
	"time"
)

// SQLiteOrderRepo stores orders in SQLite (orders + order_items).
type SQLiteOrderRepo struct {
	db *sql.DB
}

func NewSQLiteOrderRepo(db *sql.DB) *SQLiteOrderRepo {
	return &SQLiteOrderRepo{db: db}
}

var _ OrderRepository = (*SQLiteOrderRepo)(nil)

func formatDBTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func parseDBTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05"} {
		if parsed, err := time.Parse(layout, s); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.TrimSuffix(strings.Repeat("?,", n), ",")
}

func (r *SQLiteOrderRepo) loadItemsForOrder(orderID int) ([]models.OrderItem, error) {
	rows, err := r.db.Query(`
		SELECT id, name, quantity, price
		FROM order_items WHERE order_id = ?
		ORDER BY id`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var it models.OrderItem
		if err := rows.Scan(&it.ID, &it.Name, &it.Quantity, &it.Price); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (r *SQLiteOrderRepo) Save(order models.Order) models.Order {
	tx, err := r.db.Begin()
	if err != nil {
		return order
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec(`
		INSERT INTO orders (table_id, session_id, total_price, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		order.TableID, order.SessionID, order.TotalPrice, order.Status,
		formatDBTime(order.CreatedAt), formatDBTime(order.UpdatedAt))
	if err != nil {
		return order
	}
	id, err := res.LastInsertId()
	if err != nil {
		return order
	}
	order.ID = int(id)

	newItems := make([]models.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		res, err := tx.Exec(`
			INSERT INTO order_items (order_id, name, quantity, price)
			VALUES (?, ?, ?, ?)`,
			order.ID, item.Name, item.Quantity, item.Price)
		if err != nil {
			return order
		}
		lid, err := res.LastInsertId()
		if err != nil {
			return order
		}
		item.ID = int(lid)
		newItems = append(newItems, item)
	}
	order.Items = newItems

	if err := tx.Commit(); err != nil {
		return order
	}
	return order
}

func (r *SQLiteOrderRepo) GetById(id int) (models.Order, bool) {
	row := r.db.QueryRow(`
		SELECT id, table_id, session_id, total_price, status, created_at, updated_at
		FROM orders WHERE id = ?`, id)

	var o models.Order
	var createdAt, updatedAt string
	err := row.Scan(&o.ID, &o.TableID, &o.SessionID, &o.TotalPrice, &o.Status, &createdAt, &updatedAt)
	if err != nil {
		return models.Order{}, false
	}
	o.CreatedAt = parseDBTime(createdAt)
	o.UpdatedAt = parseDBTime(updatedAt)

	items, err := r.loadItemsForOrder(id)
	if err != nil {
		o.Items = nil
		return o, true
	}
	o.Items = items
	return o, true
}

func (r *SQLiteOrderRepo) GetAll() []models.Order {
	rows, err := r.db.Query(`
		SELECT id, table_id, session_id, total_price, status, created_at, updated_at
		FROM orders ORDER BY id DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var ca, ua string
		if err := rows.Scan(&o.ID, &o.TableID, &o.SessionID, &o.TotalPrice, &o.Status, &ca, &ua); err != nil {
			continue
		}
		o.CreatedAt = parseDBTime(ca)
		o.UpdatedAt = parseDBTime(ua)
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	if len(orders) == 0 {
		return []models.Order{}
	}

	ids := make([]int, len(orders))
	for i := range orders {
		ids[i] = orders[i].ID
	}
	ph := placeholders(len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	q := `SELECT id, order_id, name, quantity, price FROM order_items WHERE order_id IN (` + ph + `) ORDER BY order_id, id`
	itemRows, err := r.db.Query(q, args...)
	if err != nil {
		return orders
	}
	defer itemRows.Close()

	byOrder := make(map[int][]models.OrderItem)
	for itemRows.Next() {
		var it models.OrderItem
		var oid int
		if err := itemRows.Scan(&it.ID, &oid, &it.Name, &it.Quantity, &it.Price); err != nil {
			continue
		}
		byOrder[oid] = append(byOrder[oid], it)
	}
	for i := range orders {
		orders[i].Items = byOrder[orders[i].ID]
	}
	return orders
}

func (r *SQLiteOrderRepo) UpdateOrderById(id int, status string) models.Order {
	now := time.Now()
	_, _ = r.db.Exec(`UPDATE orders SET status = ?, updated_at = ? WHERE id = ?`,
		status, formatDBTime(now), id)
	o, _ := r.GetById(id)
	return o
}

func (r *SQLiteOrderRepo) DeleteOrderById(id int) bool {
	res, err := r.db.Exec(`DELETE FROM orders WHERE id = ?`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}
