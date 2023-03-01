package models

import (
	"github.com/volatiletech/null"
	"time"
)

type CreateUserParams struct {
	Name        string    `json:"name" db:"name"`
	Age         int32     `json:"age" db:"age"`
	Password    string    `json:"password" db:"password"`
	Address     string    `json:"address" db:"address"`
	CountryCode string    `json:"countryCode" db:"country_code"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateSessionRequest struct {
	Platform  string      `json:"platform"`
	ModelName null.String `json:"modelName"`
	OSVersion null.String `json:"osVersion"`
	DeviceID  null.String `json:"deviceId"`
}
