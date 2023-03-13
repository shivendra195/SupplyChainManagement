package server

import (
	"encoding/json"
	"github.com/shivendra195/supplyChainManagement/crypto"
	"github.com/shivendra195/supplyChainManagement/models"
	"github.com/shivendra195/supplyChainManagement/providers/authProvider"
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
	uc := srv.MiddlewareProvider.UserFromContext(req.Context())

	err := json.NewDecoder(req.Body).Decode(&newUserReq)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "error creating user", "Error parsing request")
		return
	}

	if newUserReq.Email.String == "" {
		scmerrors.RespondClientErr(resp, errors.New("email cannot be empty"), http.StatusBadRequest, "Email  cannot be empty", "Email cannot be empty")
		return
	}

	name := strings.TrimSpace(newUserReq.Name)
	if name == "" {
		scmerrors.RespondClientErr(resp, errors.New("name cannot be empty"), http.StatusBadRequest, "Name cannot be empty", "Name cannot be empty")
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
		scmerrors.RespondClientErr(resp, errors.New("phone number cannot be empty"), http.StatusBadRequest, "Phone number cannot be empty", "Phone number cannot be empty")
		return
	}

	if newUserReq.Password == "" {
		scmerrors.RespondClientErr(resp, errors.New("password cannot be empty"), http.StatusBadRequest, "password cannot be empty", "password cannot be empty")
		return
	}

	if newUserReq.Gender == "" {
		scmerrors.RespondClientErr(resp, errors.New("gender cannot be empty"), http.StatusBadRequest, "gender cannot be empty", "gender cannot be empty")
		return
	}

	if newUserReq.DateOfBirth.String == "" {
		scmerrors.RespondClientErr(resp, errors.New("date of birth cannot be empty"), http.StatusBadRequest, "Date of birth cannot be empty", "Date of birth cannot be empty")
		return
	}

	dateOfBirth, err := time.Parse("2006-01-02", newUserReq.DateOfBirth.String)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "invalid date of birth", "error in parsing date of birth")
		return
	}

	age := time.Now().Year() - dateOfBirth.Year()

	if age < models.MinAgeLimit {
		scmerrors.RespondClientErr(resp, errors.New("user should be at least 18 years old"), http.StatusBadRequest, "user should be at least 18 years old", "error in parsing date of birth")
		return
	}

	uncleanPhoneNumber := newUserReq.Phone.String

	if strings.Count(uncleanPhoneNumber, "+") == 2 {
		uncleanPhoneNumber = uncleanPhoneNumber[strings.LastIndex(uncleanPhoneNumber, "+"):]
	}

	phone := strings.ReplaceAll(uncleanPhoneNumber, " ", "")

	num, err := libphonenumber.Parse(phone, "IN")
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
		scmerrors.RespondClientErr(resp, errors.New("phone number already exist"), http.StatusBadRequest, "this phone number is already linked with one of our account please use a different phone number", "unable to create a user")
		return
	}

	if newUserReq.Role == models.SuperAdmin {
		scmerrors.RespondClientErr(resp, errors.New("cannot create super admin "), http.StatusBadRequest, "cannot create super admin", "cannot create super admin")
		return
	}
	if newUserReq.State == "" {
		scmerrors.RespondClientErr(resp, errors.New("state cannot be empty "), http.StatusBadRequest, "state cannot be empty ", "state cannot be empty ")
		return
	}
	if newUserReq.Country == "" {
		scmerrors.RespondClientErr(resp, errors.New("country cannot be empty "), http.StatusBadRequest, "country cannot be empty ", "country cannot be empty")
		return
	}
	if uc.Role == string(models.Admin) && newUserReq.Role == models.Admin {
		scmerrors.RespondClientErr(resp, errors.New("admins cannot create admins"), http.StatusBadRequest, "admins cannot create admins. Only super admins can create admins", "country cannot be empty")
		return
	}

	// Creating user in the database
	userID, err := srv.DBHelper.CreateNewUser(&newUserReq, uc.UserID)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusInternalServerError, "Error registering new user", "")
		return
	}

	utils.EncodeJSONBody(resp, http.StatusCreated, map[string]interface{}{
		"message": "success",
		"userId":  userID,
	})
}

func (srv *Server) loginWithEmailPassword(resp http.ResponseWriter, req *http.Request) {
	var token string
	var authLoginRequest models.AuthLoginRequest
	err := json.NewDecoder(req.Body).Decode(&authLoginRequest)
	if err != nil {
		logrus.Error("NewFunction : unable to decode request body ", err)
	}

	if authLoginRequest.Password == "" {
		scmerrors.RespondClientErr(resp, errors.New("password can not be empty"), http.StatusBadRequest, "Empty password!", "password field can not be empty")
		return
	}

	if authLoginRequest.Email == "" {
		scmerrors.RespondClientErr(resp, errors.New("email can not be empty"), http.StatusBadRequest, "Please enter email to login", "email can not be empty")
		return
	}
	UserDataByEmail, err := srv.DBHelper.GetUserInfoByEmail(authLoginRequest.Email)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "error getting user info", "error getting user info")
		return
	}

	loginReq := models.EmailAndPassword{
		Email:    authLoginRequest.Email,
		Password: authLoginRequest.Password,
	}
	loginReq.Email = strings.ToLower(loginReq.Email)

	userID, errorMessage, err := srv.DBHelper.LogInUserUsingEmailAndRole(loginReq, UserDataByEmail.Role)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusInternalServerError, errorMessage, errorMessage)
		return
	}

	createUserSession := models.CreateSessionRequest{
		Platform:  authLoginRequest.Platform,
		ModelName: authLoginRequest.ModelName.String,
		OSVersion: authLoginRequest.OSVersion.String,
		DeviceID:  authLoginRequest.DeviceID.String,
	}

	UUIDToken, err := srv.DBHelper.StartNewSession(userID, &createUserSession)
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
	devClaims["UUIDToken"] = UUIDToken
	devClaims["userInfo"] = UserDataByEmail
	devClaims["UserSession"] = createUserSession

	token, err = authProvider.GenerateJWT(devClaims)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusInternalServerError, "error while login", "error while login")
		return
	}

	//devClaims := make(map[string]interface{})
	//devClaims["userInfo"] = userInfo
	//devClaims["UUIDToken"] = UUIDToken
	//devClaims["UserSession"] = createUserSession
	//
	//token, err = authProvider.GenerateJWT(devClaims)
	//
	//if err != nil {
	//	scmerrors.RespondClientErr(resp, err, http.StatusInternalServerError, "error while login", "error while login")
	//	return
	//}

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"userInfo": userInfo,
		"token":    token,
	})
}

func (srv *Server) logout(resp http.ResponseWriter, r *http.Request) {
	uc := srv.MiddlewareProvider.UserFromContext(r.Context())
	err := srv.DBHelper.EndSession(uc.SessionID)
	if err != nil {
		logrus.Error("error creating user in database", err)
	}
	utils.EncodeJSON200Body(resp, map[string]interface{}{
		"message": "successfully logout",
	})
}

func (srv *Server) fetchUser(resp http.ResponseWriter, r *http.Request) {
	uc := srv.MiddlewareProvider.UserFromContext(r.Context())
	fetchedAuthor, err := srv.DBHelper.FetchUserData(uc.UserID)
	if err != nil {
		logrus.Error("error creating user in database", err)
	}
	utils.EncodeJSON200Body(resp, fetchedAuthor)
}

func (srv *Server) changePassword(resp http.ResponseWriter, req *http.Request) {
	var changePasswordRequest models.ChangePasswordRequest
	uc := srv.MiddlewareProvider.UserFromContext(req.Context())

	err := json.NewDecoder(req.Body).Decode(&changePasswordRequest)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "error updating password", "Error parsing request")
		return
	}

	if changePasswordRequest.OldPassword == "" {
		scmerrors.RespondClientErr(resp, errors.New("old password cannot be empty"), http.StatusBadRequest, "old password cannot be empty", "password cannot be empty")
		return
	}
	if changePasswordRequest.NewPassword == "" {
		scmerrors.RespondClientErr(resp, errors.New("new password cannot be empty"), http.StatusBadRequest, "new password cannot be empty", "password cannot be empty")
		return
	}
	isPasswordUpdated, err := srv.DBHelper.ChangePasswordByUserID(uc.UserID, changePasswordRequest)
	if err != nil {
		logrus.Error("error updating password", err)
	}
	if !isPasswordUpdated {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "unable to update the password ", "unable to update the password ")
		return
	}

	utils.EncodeJSON200Body(resp, map[string]interface{}{
		"message": "successfully updated password",
	})
}

func (srv *Server) dashboard(resp http.ResponseWriter, r *http.Request) {
	dashboardData, err := srv.createDashboard()
	if err != nil {
		logrus.Error("dashboard: error creating dashboard data ", err)
	}
	utils.EncodeJSON200Body(resp, dashboardData)
}

func (srv *Server) createDashboard() (models.DashboardData, error) {
	var dashboardData models.DashboardData
	var err error
	dashboardData.RecentUsers, err = srv.DBHelper.RecentUsers(5)
	if err != nil {
		logrus.Error("createDashboard: error getting recent user data", err)
	}
	dashboardData.RecentOrders, err = srv.DBHelper.RecentOrders(5, 0, false, models.OpenOrderStatus)
	if err != nil {
		logrus.Error("createDashboard: error getting recent order data", err)
	}
	dashboardData.OrderSummary, err = srv.DBHelper.OrderSummary()
	if err != nil {
		logrus.Error("createDashboard: error getting order summary data", err)
	}
	return dashboardData, nil
}

func (srv *Server) Users(resp http.ResponseWriter, r *http.Request) {
	uc := srv.MiddlewareProvider.UserFromContext(r.Context())
	role := models.UserRoles(r.URL.Query().Get("role"))
	limit, offset, err := utils.GetLimitOffsetFromRequest(r, models.DefaultLimit)
	if err != nil {
		logrus.Error("Users :error getting limit offset from request", err)
	}
	UserList, err := srv.DBHelper.UsersAll(uc.UserID, limit, offset, role)
	if err != nil {
		logrus.Error("Users :error getting users list", err)
	}
	utils.EncodeJSON200Body(resp, UserList)
}

func (srv *Server) Order(resp http.ResponseWriter, r *http.Request) {
	var order models.Order
	var orderID models.CreatedOrder
	var isOrderExists bool
	orderedItemID := 0

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "error creating order", "Error parsing request")
		return
	}
	isOrderExists, orderID, err = srv.DBHelper.IsOrderAlreadyExists(order.ReferenceID)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "error checking order status", "error checking order status")
		return
	}
	if !isOrderExists {
		orderID, err = srv.DBHelper.CreateOrder(order)
		if err != nil {
			logrus.Error("error creating user in database", err)
			return
		}
	}

	orderedItemID, err = srv.DBHelper.Scan(order, orderID)
	if err != nil {
		logrus.Error("error creating user in database", err)
		return
	}

	utils.EncodeJSON200Body(resp, map[string]interface{}{
		"orderId":       orderID.ID,
		"orderedItemId": orderedItemID,
	})
}

func (srv *Server) editProfile(resp http.ResponseWriter, r *http.Request) {
	var editProfileRequest models.EditProfile
	uc := srv.MiddlewareProvider.UserFromContext(r.Context())
	err := json.NewDecoder(r.Body).Decode(&editProfileRequest)
	if err != nil {
		logrus.Error("editProfile : unable to decode request body ", err)
	}
	name := strings.TrimSpace(editProfileRequest.Name)
	if name == "" {
		scmerrors.RespondClientErr(resp, errors.New("name cannot be empty"), http.StatusBadRequest, "Name cannot be empty", "Name cannot be empty")
		return
	}

	if uc.Role == string(models.Dealer) || uc.Role == string(models.Retailer) && uc.Email != editProfileRequest.Email {
		scmerrors.RespondClientErr(resp, errors.New("dealers and retailers are not authorized to update their emails"), http.StatusBadRequest, "cannot update email", "dealers and retailers are not authorized to update their emails")
		return
	}

	// checking if the user is already exist
	isUserExist, _, err := srv.DBHelper.IsUserAlreadyExists(editProfileRequest.Email)
	if err != nil {
		scmerrors.RespondGenericServerErr(resp, err, "error in processing request")
		return
	}

	if isUserExist {
		scmerrors.RespondClientErr(resp, errors.New("error updating user profile"), http.StatusBadRequest, "this email is already linked with one of our account please use a different email address", "unable to create a user with duplicate email address")
		return
	}

	editProfileRequest.Email = strings.ToLower(editProfileRequest.Email)
	if editProfileRequest.Name == "" {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "name cannot be empty", "name cannot be empty")
		return
	}
	if editProfileRequest.Email == "" {
		scmerrors.RespondClientErr(resp, errors.New("email cannot be empty"), http.StatusBadRequest, "email cannot be empty", "email cannot be empty")
		return
	}
	if editProfileRequest.Address == "" {
		scmerrors.RespondClientErr(resp, errors.New("address cannot be empty"), http.StatusBadRequest, "address cannot be empty", "address cannot be empty")
		return
	}
	if editProfileRequest.Country == "" {
		scmerrors.RespondClientErr(resp, errors.New("country cannot be empty"), http.StatusBadRequest, "country cannot be empty", "name cannot be empty")
		return
	}
	if editProfileRequest.State == "" {
		scmerrors.RespondClientErr(resp, errors.New("state cannot be empty"), http.StatusBadRequest, "state cannot be empty", "state cannot be empty")
		return
	}
	if editProfileRequest.Phone == "" {
		scmerrors.RespondClientErr(resp, errors.New("phone cannot be empty"), http.StatusBadRequest, "phone cannot be empty", "phone  cannot be empty")
		return
	}
	uncleanPhoneNumber := editProfileRequest.Phone

	if strings.Count(uncleanPhoneNumber, "+") == 2 {
		uncleanPhoneNumber = uncleanPhoneNumber[strings.LastIndex(uncleanPhoneNumber, "+"):]
	}

	phone := strings.ReplaceAll(uncleanPhoneNumber, " ", "")

	num, err := libphonenumber.Parse(phone, "IN")
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
	editProfileRequest.Phone = phoneNumber

	err = srv.DBHelper.EditProfile(uc.UserID, editProfileRequest)
	if err != nil {
		logrus.Error("error updating user profile in database", err)
	}
	utils.EncodeJSON200Body(resp, map[string]interface{}{
		"message": "profile updated successfully",
	})
}

func (srv *Server) getCountryAndState(resp http.ResponseWriter, r *http.Request) {
	countryAndState, err := srv.DBHelper.GetCountryAndState()
	if err != nil {
		logrus.Error("getCountryAndState: error getting Country and State ", err)
	}
	utils.EncodeJSON200Body(resp, countryAndState)
}

func (srv *Server) scan(resp http.ResponseWriter, r *http.Request) {
	var order models.Order
	var orderID models.CreatedOrder
	var isOrderExists bool
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "error creating order", "Error parsing request")
		return
	}
	isOrderExists, orderID, err = srv.DBHelper.IsOrderAlreadyExists(order.ReferenceID)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "error checking order status", "error checking order status")
		return
	}
	if !isOrderExists {
		orderID, err = srv.DBHelper.CreateOrder(order)
		if err != nil {
			logrus.Error("error creating user in database", err)
			return
		}
	}
	utils.EncodeJSON200Body(resp, map[string]interface{}{
		"orderId": orderID.ID,
	})
}

func (srv *Server) FetchOrder(resp http.ResponseWriter, r *http.Request) {
	orderStatus := models.OrderStatus(r.URL.Query().Get("orderStatus"))
	limit, offset, err := utils.GetLimitOffsetFromRequest(r, models.DefaultLimit)
	if err != nil {
		logrus.Error("FetchOrder :error getting limit offset from request", err)
	}
	recentOrders, err := srv.DBHelper.RecentOrders(limit, offset, true, orderStatus)
	if err != nil {
		logrus.Error("createDashboard: error getting recent order data", err)
	}
	utils.EncodeJSON200Body(resp, recentOrders)
}

//func (srv *Server) orderDetails(resp http.ResponseWriter, r *http.Request) {
//	orderStatus := models.OrderStatus(r.URL.Query().Get("orderStatus"))
//
//	recentOrders, err := srv.DBHelper.RecentOrders(limit, offset, true, orderStatus)
//	if err != nil {
//		logrus.Error("createDashboard: error getting recent order data", err)
//	}
//	utils.EncodeJSON200Body(resp, recentOrders)
//}
