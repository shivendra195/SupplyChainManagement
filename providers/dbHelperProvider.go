package providers

import (
	"github.com/shivendra195/supplyChainManagement/models"
)

type DBHelperProvider interface {
	CreateNewUser(newUserRequest *models.CreateNewUserRequest, userID int) (*int, error)
	IsUserAlreadyExists(emailID string) (isUserExist bool, user models.UserData, err error)
	IsPhoneNumberAlreadyExist(phone string) (bool, error)
	FetchUserData(userID int) (models.FetchUserData, error)
	ChangePasswordByUserID(userID int, changePasswordRequest models.ChangePasswordRequest) (bool, error)
	UsersAll(userID, limit, offset int, role models.UserRoles) ([]models.FetchUserData, error)
	LogInUserUsingEmailAndRole(loginReq models.EmailAndPassword, role models.UserRoles) (userID int, message string, err error)
	StartNewSession(userID int, request *models.CreateSessionRequest) (string, error)
	FetchUserSessionData(userID int) ([]models.FetchUserSessionsData, error)
	UpdateSession(sessionId string) error
	EndSession(sessionId string) error
	Dashboard() (models.FetchUserData, error)
	RecentUsers(limit int) ([]models.UserInfo, error)
	RecentOrders(limit, offset int, isStatusCheck bool, orderStatus models.OrderStatus) ([]models.RecentOrders, error)
	OrderSummary() (models.OrderSummary, error)
	CreateOrder(order models.Order) (models.CreatedOrder, error)
	GetUserInfoByEmail(email string) (models.GetUserDataByEmail, error)
	EditProfile(userID int, editProfileRequest models.EditProfile) error
	GetCountryAndState() ([]models.CountryAndState, error)
	IsOrderAlreadyExists(referenceID string) (isOrderExist bool, orderID models.CreatedOrder, err error)
	Scan(order models.Order, orderID models.CreatedOrder) (int, error)
}
