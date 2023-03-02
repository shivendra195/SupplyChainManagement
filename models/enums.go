package models

type GenderType string

type UserRoles string

type ENV string

type JWT string

const (
	JWTSecretKey JWT = ""

	PSQLDBURL ENV = "PSQL_DB_URL"

	Dealer   UserRoles = "dealer"
	Retailer UserRoles = "retailer"
	Admin    UserRoles = "admin"

	Male   GenderType = "male"
	Female GenderType = "female"
	Other  GenderType = "other"
	None   GenderType = "none"

	MinAgeLimit = 18
)
