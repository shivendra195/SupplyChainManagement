package models

import (
	"errors"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/volatiletech/null"
	"time"
)

var ErrorPasswordNotMatched = errors.New("password not matched")

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

type AuthLoginRequest struct {
	Platform  string      `json:"platform"`
	ModelName null.String `json:"modelName"`
	OSVersion null.String `json:"osVersion"`
	DeviceID  null.String `json:"deviceId"`
	Token     string      `json:"token"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
}

type EmailAndPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FetchUserData struct {
	UserID      int         `json:"userId" db:"id"`
	Name        string      `json:"name" db:"name"`
	Address     string      `json:"address" db:"address"`
	Email       string      `json:"email" db:"email"`
	Phone       string      `json:"phone" db:"phone"`
	Gender      GenderType  `json:"gender" db:"gender"`
	DateOfBirth null.String `json:"dateOfBirth" db:"date_of_birth"`
}

type FetchUserSessionsData struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"userId" db:"user_id"`
	UUIDToken string    `json:"UUIDToken" db:"token"`
	EndTime   time.Time `json:"endTime" db:"end_time"`
}

type Response struct {
	Persons []Person `json:"persons"`
}

type Person struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserContextData struct {
	UserID      int         `json:"userId" db:"id"`
	SessionID   string      `json:"sessionID" db:"token"`
	Name        string      `json:"name" db:"name"`
	Address     string      `json:"address" db:"address"`
	Email       string      `json:"email" db:"email"`
	Phone       string      `json:"phone" db:"phone"`
	Gender      GenderType  `json:"gender" db:"gender"`
	DateOfBirth null.String `json:"dateOfBirth" db:"date_of_birth"`
}
