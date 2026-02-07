package repositories

import (
	"database/sql"
	"fmt"
	"go-kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	trx, err := repo.db.Begin()

	if err != nil {
		return res, nil
	}

	// Main logic
	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productName string
		var productID, price, stock int

		err := trx.QueryRow("SELECT id, name, price, stock FROM products WHERE id=$1", item.ProductID).Scan(&productID, &productName, &price, &stock)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		subtotal := price * item.Quantity
		totalAmount += subtotal

		// Kurangi jumlah stock
		_, err = trx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)

		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert trx
	var transactionID int
	err = trx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)

	if err != nil {
		return nil, err
	}

	// insert trx detail
	for i, detail := range details {
		details[i].TransactionID = transactionID
		_, err = trx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)

		if err != nil {
			return nil, err
		}
	}

	if err := trx.Commit(); err != nil {
		return res, nil
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
