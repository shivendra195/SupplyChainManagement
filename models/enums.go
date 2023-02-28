package models

type GenderType string

type UserRoles string

type ENV string

const (
	PSQLDBURL ENV = "PSQL_DB_URL"

	Dealer   UserRoles = "dealer"
	Retailer UserRoles = "retailer"
	Admin    UserRoles = "adminprovider"

	Male   GenderType = "male"
	Female GenderType = "female"
	Other  GenderType = "other"
	None   GenderType = "none"
)
