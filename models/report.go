package models

type Report struct {
	TotalRevenue     int        `json:"total_revenue"`
	TotalTransaction int        `json:"total_transaksi"`
	TopProduct       TopProduct `json:"top_product"`
}

type TopProduct struct {
	ProductName string `json:"name"`
	Quantity    int    `json:"qty_terjual"`
}
