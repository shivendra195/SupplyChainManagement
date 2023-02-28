package models

import (
	"github.com/volatiletech/null"
)

type CreateNewUserRequest struct {
	Name        string      `json:"name"`
	Email       null.String `json:"email"`
	Phone       null.String `json:"phone"`
	Password    string      `json:"password"`
	Gender      GenderType  `json:"gender"`
	DateOfBirth null.String `json:"dateOfBirth"`
}
