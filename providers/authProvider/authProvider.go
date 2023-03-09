package authProvider

import (
	"github.com/shivendra195/supplyChainManagement/models"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null"

	//"fmt"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

//func generateJWT(userToken string) (string, error) {
//	token := jwt.New(jwt.SigningMethodEdDSA)
//	claims := token.Claims.(jwt.MapClaims)
//	claims["exp"] = time.Now().Add(1 * time.Hour)
//	claims["authorized"] = true
//	claims["user"] = "username"
//
//	tokenString, err := token.SignedString(sampleSecretKey)
//	if err != nil {
//		fmt.Errorf("something went wrong: %s", err.Error())
//		return "", err
//	}
//	return tokenString, nil
//}

var secret = []byte("supersecretkey")

type JWTClaim struct {
	Platform  string      `json:"platform"`
	ModelName null.String `json:"modelName"`
	OSVersion null.String `json:"osVersion"`
	DeviceID  null.String `json:"deviceId"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	UUIDToken string      `json:"token"`
	jwt.StandardClaims
}

func GenerateJWT(devClaims map[string]interface{}) (tokenString string, err error) {
	// var userInfo models.GetUserDataByEmail
	var userSessionData models.CreateSessionRequest
	var ok bool
	var UUIDToken string
	userInfo, ok := devClaims["userInfo"].(models.GetUserDataByEmail)
	if !ok {
		logrus.Error("GenerateJWT:  error getting values out of the devClaims map 1")
	}
	UUIDToken, ok = devClaims["UUIDToken"].(string)
	if !ok {
		logrus.Error("GenerateJWT:  error getting values out of the devClaims map 2")
	}
	userSessionData, ok = devClaims["UserSession"].(models.CreateSessionRequest)
	if !ok {
		logrus.Error("GenerateJWT:  error getting values out of the devClaims map 3")
	}
	UserIDString := strconv.Itoa(userInfo.UserID)
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &jwt.MapClaims{
		"iss": UserIDString,
		"exp": time.Now().Add(time.Hour).Unix(),
		"data": map[string]string{
			"id":        UserIDString,
			"name":      userInfo.Name,
			"role":      string(userInfo.Role),
			"modelName": userSessionData.ModelName,
			"platform":  userSessionData.Platform,
			"oSVersion": userSessionData.OSVersion,
			"deviceId":  userSessionData.DeviceID,
			"email":     userInfo.Email,
			"username":  userInfo.Name,
			"uuidToken": UUIDToken,
			"expiresAt": expirationTime.String(),
			"issuer":    UserIDString,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secret)
	return
}
