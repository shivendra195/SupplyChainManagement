package models

type GenderType string

type UserRoles string

type ENV string

type JWT string

type OrderStatus string

type Context string

const (
	JWTSecretKey JWT = ""

	PSQLDBURL ENV = "PSQL_DB_URL"

	Dealer     UserRoles = "dealer"
	Retailer   UserRoles = "retailer"
	Admin      UserRoles = "admin"
	SuperAdmin UserRoles = "super admin"
	All        UserRoles = "All"

	Male   GenderType = "male"
	Female GenderType = "female"
	Other  GenderType = "other"
	None   GenderType = "none"

	OpenOrderStatus       OrderStatus = "open"
	InStockOrderStatus    OrderStatus = "in stock"
	SoldOutOrderStatus    OrderStatus = "in transfer"
	InTransferOrderStatus OrderStatus = "sold out"
	OutOfStockOrderStatus OrderStatus = "out of stock"

	MinAgeLimit = 18

	DefaultLimit = 18

	UserContext Context = "userContext"
)
