package repositories

import (
	"database/sql"
	"go-kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetReport(startDate, endDate string) (*models.Report, error) {
	report := models.Report{}

	// Total Revenue
	if err := repo.db.QueryRow(
		"SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE created_at::date BETWEEN $1 AND $2",
		startDate, endDate,
	).Scan(&report.TotalRevenue); err != nil {
		return nil, err
	}

	// Total Transactions
	if err := repo.db.QueryRow(
		"SELECT COUNT(*) FROM transactions WHERE created_at::date BETWEEN $1 AND $2",
		startDate, endDate,
	).Scan(&report.TotalTransaction); err != nil {
		return nil, err
	}

	// Top Product by Quantity Sold
	var name sql.NullString
	var qty sql.NullInt64
	if err := repo.db.QueryRow(`
		SELECT p.name, SUM(td.quantity) AS qty
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at::date BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&name, &qty); err != nil {
		if err == sql.ErrNoRows {
			report.TopProduct = models.TopProduct{ProductName: "", Quantity: 0}
			return &report, nil
		}

		return nil, err
	}

	report.TopProduct = models.TopProduct{ProductName: name.String, Quantity: int(qty.Int64)}
	return &report, nil
}
