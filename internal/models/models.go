package models

import (
	"context"
	"database/sql"
	"time"
)

// DGModel is the type for database connections values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// NewModels return a model type with database conncetion pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Widgets is the type for all widgets
type Widgets struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"_"`
	UpdatedAt      time.Time `json:"_"`
}

type Orders struct {
}

func (m *DBModel) GetWidget(id int) (Widgets, error) {
	// make context Background with timewout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widgets

	// query row
	row := m.DB.QueryRowContext(ctx, "select id, name from widgets where id = ?", id)
	err := row.Scan(&widget.ID, &widget.Name)
	if err != nil {
		return widget, err
	}

	return widget, err
}
