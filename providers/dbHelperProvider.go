package providers

import (
	"github.com/shivendra195/supplyChainManagement/models"
)

type DBHelperProvider interface {
	CreateNewUser(newUserRequest *models.CreateNewUserRequest) (*int, error)
	IsUserAlreadyExists(emailID string) (isUserExist bool, user models.UserData, err error)
	IsPhoneNumberAlreadyExist(phone string) (bool, error)
	FetchUserData(userID int) (models.FetchUserData, error)
	LogInUserUsingEmailAndRole(loginReq models.EmailAndPassword, role models.UserRoles) (userID int, message string, err error)
	StartNewSession(userID int, request *models.CreateSessionRequest) (string, error)
	FetchUserSessionData(userID int, deviceId, modelName, osVersion, platform string) ([]models.FetchUserSessionsData, error)
	UpdateSession(sessionId string) error
	EndSession(sessionId string) error
}
