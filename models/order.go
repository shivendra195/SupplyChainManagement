package models

type Order struct {
	Quantity string `json:"quantity " db:"quantity"`
}

type CreatedOrder struct {
	ID int `json:"id " db:"id"`
}
