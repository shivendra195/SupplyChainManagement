package server

import (
	"encoding/json"
	"github.com/shivendra195/supplyChainManagement/crypto"
	"github.com/shivendra195/supplyChainManagement/models"
	"github.com/shivendra195/supplyChainManagement/scmerrors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
	//"log"
	"errors"
	"github.com/ttacon/libphonenumber"
	"net/http"
	"strings"
	"time"

	"github.com/shivendra195/supplyChainManagement/utils"
)

func (srv *Server) register(resp http.ResponseWriter, req *http.Request) {
	var newUserReq models.CreateNewUserRequest

	err := json.NewDecoder(req.Body).Decode(&newUserReq)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "error creating user", "Error parsing request")
		return
	}

	if newUserReq.Email.String == "" {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "Email  cannot be empty", "Email cannot be empty")
		return
	}

	name := strings.TrimSpace(newUserReq.Name)
	if name == "" {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "Name cannot be empty", "Name cannot be empty")
		return
	}
	// checking if the user is already exist
	isUserExist, _, err := srv.DBHelper.IsUserAlreadyExists(newUserReq.Email.String)
	if err != nil {
		scmerrors.RespondGenericServerErr(resp, err, "error in processing request")
		return
	}

	if isUserExist {
		scmerrors.RespondClientErr(resp, errors.New("error creating user"), http.StatusBadRequest, "this email is already linked with one of our account please use a different email address", "unable to create a user with duplicate email address")
		return
	}

	newUserReq.Email.String = strings.ToLower(newUserReq.Email.String)
	if newUserReq.Name == "" {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "name cannot be empty", "name cannot be empty")
		return
	}

	if newUserReq.Phone.String == "" {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "Phone number cannot be empty", "Phone number cannot be empty")
		return
	}

	if newUserReq.Password == "" {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "password cannot be empty", "password cannot be empty")
		return
	}

	if newUserReq.Gender == "" {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "gender cannot be empty", "gender cannot be empty")
		return
	}

	if newUserReq.DateOfBirth.String == "" {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "Date of birth cannot be empty", "Date of birth cannot be empty")
		return
	}

	dateOfBirth, err := time.Parse("2006-01-02", newUserReq.DateOfBirth.String)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "invalid date of birth", "error in parsing date of birth")
		return
	}

	age := time.Now().Year() - dateOfBirth.Year()

	if age < models.MinAgeLimit {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "user should be at least 18 years old", "error in parsing date of birth")
		return
	}

	uncleanPhoneNumber := newUserReq.Phone.String

	if strings.Count(uncleanPhoneNumber, "+") == 2 {
		uncleanPhoneNumber = uncleanPhoneNumber[strings.LastIndex(uncleanPhoneNumber, "+"):]
	}

	phone := strings.ReplaceAll(uncleanPhoneNumber, " ", "")

	num, err := libphonenumber.Parse(phone, "US")
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "Phone Number not in a correct format", "invalid format for phone number")
		return
	}

	isValidNumber := libphonenumber.IsValidNumber(num)

	if !isValidNumber {
		scmerrors.RespondClientErr(resp, errors.New("invalid phone number"), http.StatusBadRequest, "invalid phone number", "invalid phone number")
		return
	}

	phoneNumber := libphonenumber.Format(num, libphonenumber.E164)
	newUserReq.Phone = null.StringFrom(phoneNumber)

	if !crypto.IsGoodPassword(newUserReq.Password) {
		scmerrors.RespondClientErr(resp, errors.New("password length should be at least 6"), http.StatusBadRequest, "password length should be at least 6", "password length should be at least 6")
		return
	}

	isMobileAlreadyExist, err := srv.DBHelper.IsPhoneNumberAlreadyExist(phoneNumber)
	if err != nil {
		scmerrors.RespondGenericServerErr(resp, err, "unable to create user")
		return
	}

	if isMobileAlreadyExist {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "this phone number is already linked with one of our account please use a different phone number", "unable to create a user")
		return
	}

	// Creating user in the database
	userID, err := srv.DBHelper.CreateNewUser(&newUserReq)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusInternalServerError, "Error registering new user aaaaaaaaaaa", "")
		return
	}

	utils.EncodeJSONBody(resp, http.StatusCreated, map[string]interface{}{
		"message": "success",
		"userId":  userID,
	})
}

func (srv *Server) fetchUser(resp http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	var userData models.UserData

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		logrus.Error("NewFunction : unable to decode request body ", err)
	}

	fetchedAuthor, err := srv.DBHelper.FetchUserData(userData.UserID)
	if err != nil {
		logrus.Error("error creating user in database", err)
	}

	utils.EncodeJSON200Body(resp, fetchedAuthor)
}

func (srv *Server) loginWithEmailPassword(resp http.ResponseWriter, req *http.Request) {
	var token string
	var authLoginRequest models.AuthLoginRequest

	if authLoginRequest.Password == "" {
		scmerrors.RespondClientErr(resp, errors.New("password can not be empty"), http.StatusBadRequest, "Empty password!", "password field can not be empty")
		return
	}

	if authLoginRequest.Email == "" {
		scmerrors.RespondClientErr(resp, errors.New("email can not be empty"), http.StatusBadRequest, "Please enter email to login", "email can not be empty")
		return
	}

	loginReq := models.EmailAndPassword{
		Email:    authLoginRequest.Email,
		Password: authLoginRequest.Password,
	}
	loginReq.Email = strings.ToLower(loginReq.Email)
	userID, errorMessage, err := srv.DBHelper.LogInUserUsingEmailAndRole(loginReq, models.Admin)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusInternalServerError, errorMessage, errorMessage)
		return
	}

	createUserSession := models.CreateSessionRequest{
		Platform:  authLoginRequest.Platform,
		ModelName: authLoginRequest.ModelName,
		OSVersion: authLoginRequest.OSVersion,
		DeviceID:  authLoginRequest.DeviceID,
	}

	newSessionToken, err := srv.DBHelper.StartNewSession(userID, &createUserSession)
	if err != nil {
		scmerrors.RespondGenericServerErr(resp, err, "error in creating session")
		return
	}

	userInfo, err := srv.DBHelper.FetchUserData(userID)
	if err != nil {
		scmerrors.RespondGenericServerErr(resp, err, "error in getting user info")
		return
	}

	devClaims := make(map[string]interface{})
	devClaims["token"] = newSessionToken
	devClaims["userInfo"] = userInfo

	token, err = srv.MiddleProvider.CustomTokenAuthWithClaims(devClaims)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusInternalServerError, "error while login", "error while login")
		return
	}

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{})

}
