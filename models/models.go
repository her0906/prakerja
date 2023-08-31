package models

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Transaction struct {
	ID         int `json:"id"`
	ProductID  int `json:"product_id"`
	Quantity   int `json:"quantity"`
	TotalPrice int `json:"total_price"`
}
