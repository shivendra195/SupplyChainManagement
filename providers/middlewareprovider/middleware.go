package middlewareprovider

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/shivendra195/supplyChainManagement/models"
	"github.com/shivendra195/supplyChainManagement/providers"
	"github.com/shivendra195/supplyChainManagement/scmerrors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var secret = []byte("supersecretkey")

const (
	authorization = "Authorization"
	bearerScheme  = "bearer"
	space         = " "
	sessionHeader = "x-session-token"
	maxAge        = 300
	sessionClaims = "sessionToken"
	minimumTime   = 10
)

type middleware struct {
	DBHelper providers.DBHelperProvider
}

func NewMiddleware(dbhelper providers.DBHelperProvider) providers.MiddlewareProvider {
	return &middleware{
		DBHelper: dbhelper,
	}
}

func (AM *middleware) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string
			r.Header.Get("token")
			r.Header.Get("Connection")

			tokenParts := strings.Split(r.Header.Get(authorization), space)
			if len(tokenParts) != 2 {
				scmerrors.RespondClientErr(w, errors.New("token not Bearer"), http.StatusUnauthorized, "Invalid token", "Invalid token")
				return
			}

			if !strings.EqualFold(tokenParts[0], bearerScheme) {
				scmerrors.RespondClientErr(w, errors.New("token not Bearer"), http.StatusUnauthorized, "Invalid token", "Invalid token")
				return
			}
			token = tokenParts[1]
			claims, err := GetClaimsFromToken(token)
			if err != nil {
				scmerrors.RespondClientErr(w, err, http.StatusUnauthorized, "GetClaimsFromToken :Invalid token", "Invalid token")
				return
			}

			SessionId, verifiedClaims, UserData, err := AM.getUserDataFromClaims(claims)
			if err != nil {
				scmerrors.RespondClientErr(w, err, http.StatusUnauthorized, "getUserDataFromClaims: Invalid token", "Invalid token")
				return
			}

			if !verifiedClaims {
				scmerrors.RespondClientErr(w, errors.New("invalid token"), http.StatusUnauthorized, "Invalid token", "Invalid token")
				return
			}
			err = AM.DBHelper.UpdateSession(SessionId)
			if err != nil {
				scmerrors.RespondClientErr(w, err, http.StatusUnauthorized, "UpdateSession: error updating sessions ", "UpdateSession error updating sessions ")
				return
			}
			var userContextData models.UserContextData
			userContextData.UserID = UserData.UserID
			userContextData.Name = UserData.Name
			userContextData.Email = UserData.Email
			userContextData.Phone = UserData.Phone
			userContextData.Gender = UserData.Gender
			userContextData.DateOfBirth = UserData.DateOfBirth
			userContextData.Address = UserData.Address
			userContextData.SessionID = SessionId

			ctxWithUser := context.WithValue(r.Context(), "userContext", userContextData)
			rWithUser := r.WithContext(ctxWithUser)
			next.ServeHTTP(w, rWithUser)

		})
	}
}

func (AM *middleware) UserFromContext(ctx context.Context) models.UserContextData {
	return ctx.Value("userContext").(models.UserContextData)
}

func GetClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		return claims, nil
	}
	return jwt.MapClaims{}, err
}

func (AM *middleware) getUserDataFromClaims(claims jwt.MapClaims) (string, bool, models.FetchUserData, error) {
	var validToken bool
	var standardClaims jwt.StandardClaims
	var UserData models.FetchUserData
	data := make(map[string]interface{})
	data = claims["data"].(map[string]interface{})

	var issuer, name, sessionId string

	issuer = claims["iss"].(string)
	expirationTime := claims["exp"]
	name = data["name"].(string)
	email := data["email"]
	UUIDToken := data["uuidToken"].(string)
	deviceId := data["deviceId"].(string)
	modelName := data["modelName"].(string)
	osVersion := data["oSVersion"].(string)
	platform := data["platform"].(string)

	fmt.Println("\tdata := claims[\"data\"].(map[string]string)", data)
	fmt.Println("issuer := claims[\"iss\"].(string)", issuer)
	fmt.Println("\texpirationTime  := claims[\"exp\"].(string)", expirationTime)
	fmt.Println("fmt.Println(standardClaims.Issuer)", standardClaims.Issuer)

	UserIDInt, err := strconv.Atoi(issuer)
	if err != nil {
		logrus.Error("GetUserDataFromClaims: error converting userId string to integer ", err)
		fmt.Println("GetUserDataFromClaims: error converting userId string to integer ", err)
		return sessionId, validToken, UserData, errors.New(fmt.Sprintln("GetUserDataFromClaims: error converting userId string to integer & \n", err))
	}

	UserData, err = AM.DBHelper.FetchUserData(UserIDInt)
	if err != nil {
		logrus.Error("GetUserDataFromClaims: error fetching user Data from database ", err)
		return sessionId, validToken, UserData, errors.New(fmt.Sprintln("GetUserDataFromClaims: error fetching user Data from database  & \n", err))
	}

	UserSessionsData, err := AM.DBHelper.FetchUserSessionData(UserIDInt, deviceId, modelName, osVersion, platform)
	if err != nil {
		logrus.Error("GetUserDataFromClaims: error fetching user session  Data from database ", err)
		return sessionId, validToken, UserData, errors.New(fmt.Sprintln("GetUserDataFromClaims: error fetching user Data from database  & \n", err))
	}
	sessionId = UserSessionsData[0].UUIDToken
	fmt.Println("GetUserDataFromClaims: UserSessionsData[0].UUIDToken ", UserSessionsData[0].UUIDToken, UUIDToken, UserSessionsData[0].EndTime.Unix())
	if email == UserData.Email && name == UserData.Name && UserIDInt == UserData.UserID {
		if standardClaims.ExpiresAt < time.Now().Unix() {
			if UserSessionsData[0].UUIDToken == UUIDToken && UserSessionsData[0].EndTime.Unix() > time.Now().Unix() {
				return sessionId, true, UserData, nil
			} else {
				return sessionId, validToken, UserData, errors.New(fmt.Sprintln("invalid session id   & \n", err))
			}
		} else {
			return sessionId, validToken, UserData, errors.New(fmt.Sprintln("token is expired \n", err))
		}
	} else {
		return sessionId, validToken, UserData, errors.New(fmt.Sprintln("invalid token  \n", err))
	}
}
