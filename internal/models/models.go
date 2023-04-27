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
	Image          string    `json:"image"`
	CreatedAt      time.Time `json:"_"`
	UpdatedAt      time.Time `json:"_"`
}

// Order is the type for all orders
type Orders struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"_"`
	UpdatedAt     time.Time `json:"_"`
}

// Status is the type for orders statues
type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"_"`
	UpdatedAt time.Time `json:"_"`
}

// TransactionStatus is the type for transaction statues
type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"_"`
	UpdatedAt time.Time `json:"_"`
}

// Transaction is the type for transaction
type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	BankReturnCode      string    `json:"Bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"_"`
	UpdatedAt           time.Time `json:"_"`
}

// Transaction is the type for transaction
type User struct {
	ID        int       `json:"id"`
	FirstName int       `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"_"`
	UpdatedAt time.Time `json:"_"`
}

func (m *DBModel) GetWidget(id int) (Widgets, error) {
	// make context Background with timewout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widgets

	// query row
	// coalese is default value
	row := m.DB.QueryRowContext(ctx, `
		select
			id, name, description, inventory_level, price, coalesce(image, ''), created_at, updated_at
		from
			widgets
		where id = ?`, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)
	if err != nil {
		return widget, err
	}

	return widget, err
}

// InsertTransaction insert a new txn, and return its id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	// define contect
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into transaction
			(amount, currency, last_four, bank_return_coode, transaction_status_id, create_at, update_at)
		values (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.TransactionStatusID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertOrder insert a new order, and return its id
func (m *DBModel) Orders(ord Orders) (int, error) {
	// define contect
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into orders
			(widget_id, transaction_id, status_id, quantity, amount, create_at, update_at)
		values (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		ord.WidgetID,
		ord.TransactionID,
		ord.StatusID,
		ord.Quantity,
		ord.Amount,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
