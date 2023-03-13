package models

import "time"

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

type FetchOrderDetails struct {
	ID              string      `json:"orderId " db:"id"`
	OrderedBy       int         `json:"orderedBy" db:"ordered_by"`
	Role            UserRoles   `json:"role" db:"role"`
	ShippingAddress string      `json:"shippingAddress" db:"shipping_address"`
	Quantity        int         `json:"quantity" db:"quantity"`
	Status          OrderStatus `json:"status" db:"status"`
	OrderedID       int         `json:"orderedId" db:"order_id"`
	CreatedAt       time.Time   `json:"createdAt" db:"created_at"`
	Items           []Item      `json:"items" db:"items"`
}

type Item struct {
	Id       int         `json:"itemId" db:"item_id"`
	ItemQRID int         `json:"itemQRID" db:"qr_id"`
	Data     interface{} `json:"data" db:"data"`
}
