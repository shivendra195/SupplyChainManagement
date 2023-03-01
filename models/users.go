package models

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/volatiletech/null"
	"time"
)

type CreateNewUserRequest struct {
	Name        string              `json:"name"  db:"name"`
	Email       null.String         `json:"email"  db:"email"`
	Role        UserRoles           `json:"role"  db:"role"`
	Password    string              `json:"password"  db:"password"`
	Address     string              `json:"address" db:"address"`
	Gender      GenderType          `json:"gender"  db:"gender"`
	CountryCode null.String         `json:"countryCode" db:"country_code"`
	Phone       null.String         `json:"phone" db:"phone"`
	DateOfBirth null.String         `json:"dateOfBirth" db:"date_of_birth"`
	CreatedAt   time.Time           `db:"created_at"`
	UpdatedAt   timestamp.Timestamp `db:"updated_at"`
}

type UserData struct {
	UserID int `json:"userId" db:"id"`
}

type FetchUserData struct {
	Name        string      `json:"name" db:"name"`
	Address     string      `json:"address" db:"address"`
	Email       string      `json:"email" db:"email"`
	Phone       string      `json:"phone" db:"phone"`
	Gender      GenderType  `json:"gender" db:"gender"`
	DateOfBirth null.String `json:"dateOfBirth" db:"date_of_birth"`
}

type Response struct {
	Persons []Person `json:"persons"`
}

type Person struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
