package models

type Order struct {
	Quantity        string      `json:"quantity " db:"quantity"`
	ReferenceID     string      `json:"referenceID" db:""`
	Payload         interface{} `json:"payload"`
	OrderedBy       int         `json:"orderedBy" db:"ordered_by"`
	ShippingAddress string      `json:"shippingAddress" db:"shipping_address"`
}

type CreatedOrder struct {
	ID int `json:"id " db:"id"`
}

type OrderedItemID struct {
	ID int `json:"id " db:"id"`
}
