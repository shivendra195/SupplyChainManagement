package providers

import "github.com/shivendra195/supplyChainManagement/models"

type DBHelperProvider interface {
	CreateNewUser(newUserRequest *models.CreateNewUserRequest) (*int, error)
	IsUserAlreadyExists(emailID string) (isUserExist bool, user models.UserData, err error)
	IsPhoneNumberAlreadyExist(phone string) (bool, error)
	FetchUserData(userID int) (models.FetchUserData, error)
}
