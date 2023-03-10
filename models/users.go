package models

import (
	"errors"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/lib/pq"
	"github.com/volatiletech/null"
	"time"
)

var ErrorPasswordNotMatched = errors.New("password not matched")

type CreateNewUserRequest struct {
	Name           string              `json:"name"  db:"name"`
	Email          null.String         `json:"email"  db:"email"`
	Role           UserRoles           `json:"role"  db:"role"`
	Password       string              `json:"password"  db:"password"`
	Address        string              `json:"address" db:"address"`
	Gender         GenderType          `json:"gender"  db:"gender"`
	CountryCode    null.String         `json:"countryCode" db:"country_code"`
	Phone          null.String         `json:"phone" db:"phone"`
	DateOfBirth    null.String         `json:"dateOfBirth" db:"date_of_birth"`
	CreatedAt      time.Time           `db:"created_at"`
	UpdatedAt      timestamp.Timestamp `db:"updated_at"`
	ProfileImageID int                 `json:"profileImageId" db:"profile_image_id"`
	State          string              `json:"state" db:"state"`
	Country        string              `json:"country" db:"country"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"  `
	NewPassword string `json:"newPassword"  `
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
	UserID          int         `json:"userId" db:"id"`
	Name            string      `json:"name" db:"name"`
	Role            string      `json:"role" db:"role"`
	Address         string      `json:"address" db:"address"`
	Email           string      `json:"email" db:"email"`
	Phone           string      `json:"phone" db:"phone"`
	Gender          GenderType  `json:"gender" db:"gender"`
	DateOfBirth     null.String `json:"dateOfBirth" db:"date_of_birth"`
	ProfileImageID  string      `json:"profileImageId" db:"profile_image_id"`
	ProfileImageURL string      `json:"profileImageURL" db:"url"`
	State           string      `json:"state" db:"state"`
	Country         string      `json:"country" db:"country"`
	CreatedAt       time.Time   `json:"createdAt" db:"created_at"`
}

type ListUsers struct {
	UserID          int         `json:"userId" db:"id"`
	Name            string      `json:"name" db:"name"`
	Role            string      `json:"role" db:"role"`
	Address         string      `json:"address" db:"address"`
	Email           string      `json:"email" db:"email"`
	Phone           string      `json:"phone" db:"phone"`
	Gender          GenderType  `json:"gender" db:"gender"`
	DateOfBirth     null.String `json:"dateOfBirth" db:"date_of_birth"`
	ProfileImageID  string      `json:"profileImageId" db:"profile_image_id"`
	ProfileImageURL string      `json:"profileImageURL" db:"url"`
	State           string      `json:"state" db:"state"`
	Country         string      `json:"country" db:"country"`
}

type AdminDashboardData struct {
	TotalUsers     int        `json:"totalUsers" db:"total_users"`
	TotalDealers   int        `json:"totalDealers" db:"total_dealers"`
	TotalRetailers int        `json:"totalRetailers" db:"total_retailers"`
	ListUserInfo   []UserInfo `json:"listUserInfo" db:"list_user_info"`
}

type UserInfo struct {
	UserID    int         `json:"userId" db:"id"`
	Name      string      `json:"name" db:"name"`
	Email     string      `json:"email" db:"email"`
	Role      string      `json:"role" db:"role"`
	CreatedAt null.String `json:"createdAt" db:"created_at"`
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
	Role        string      `json:"role" db:"role"`
	Address     string      `json:"address" db:"address"`
	Email       string      `json:"email" db:"email"`
	Phone       string      `json:"phone" db:"phone"`
	Gender      GenderType  `json:"gender" db:"gender"`
	DateOfBirth null.String `json:"dateOfBirth" db:"date_of_birth"`
}

type RecentOrders struct {
	UserID    int         `json:"userId" db:"user_id"`
	Name      string      `json:"name" db:"name"`
	OrderID   int         `json:"orderId" db:"order_id"`
	Quantity  string      `json:"quantity" db:"quantity"`
	Status    OrderStatus `json:"status" db:"order_status"`
	CreatedAt null.String `json:"createdAt" db:"created_at"`
}

type OrderSummary struct {
	TotalCreatedOrders int `json:"totalCreatedOrders" db:"total_created_orders"`
	OpenDeliveries     int `json:"openDeliveries" db:"open_deliveries"`
	InStock            int `json:"inStock" db:"in_stock"`
	InTransfer         int `json:"inTransfer" db:"in_transfer"`
	SoldOut            int `json:"soldOut" db:"sold_out"`
}

type GetUserDataByEmail struct {
	UserID      int         `json:"userId" db:"id"`
	Name        string      `json:"name" db:"name"`
	Role        UserRoles   `json:"role" db:"role"`
	Address     string      `json:"address" db:"address"`
	Email       string      `json:"email" db:"email"`
	Phone       string      `json:"phone" db:"phone"`
	Gender      GenderType  `json:"gender" db:"gender"`
	DateOfBirth null.String `json:"dateOfBirth" db:"date_of_birth"`
}

type EditProfile struct {
	Name           string      `json:"name"`
	Address        string      `json:"address" db:"address"`
	Email          string      `json:"email" db:"email"`
	Phone          string      `json:"phone" db:"phone"`
	CountryCode    null.String `json:"countryCode" db:"country_code"`
	DateOfBirth    null.String `json:"dateOfBirth" db:"date_of_birth"`
	ProfileImageID int         `json:"profileImageId" db:"profile_image_id"`
	State          string      `json:"state" db:"state"`
	Country        string      `json:"country" db:"country"`
}

type CountryAndState struct {
	CountryID   int            `json:"countryId" db:"id"`
	Country     string         `json:"country" db:"country"`
	CountryCode string         `json:"countryCode" db:"country_code"`
	StateID     pq.Int32Array  `json:"stateId" db:"state_id"`
	State       pq.StringArray `json:"state" db:"state"`
}
