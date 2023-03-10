package dbhelperprovider

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shivendra195/supplyChainManagement/crypto"
	"github.com/shivendra195/supplyChainManagement/dbutil"
	"github.com/shivendra195/supplyChainManagement/models"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null"

	"time"
)

func (dh *DBHelper) CreateNewUser(newUserRequest *models.CreateNewUserRequest, userID int) (*int, error) {
	var newUserID int
	var password null.String
	if newUserRequest.Password != "" {
		password = null.StringFrom(crypto.HashAndSalt(newUserRequest.Password))
	}

	txErr := dbutil.WithTransaction(dh.DB, func(tx *sqlx.Tx) error {

		SQL := `INSERT INTO users
				(name, email, phone, password, address, country_code, created_at, created_by)
				VALUES (trim($1), lower(trim($2)), $3, $4,$5,$6,$7,$8)
				RETURNING id`

		args := []interface{}{
			newUserRequest.Name,
			newUserRequest.Email,
			newUserRequest.Phone,
			password,
			newUserRequest.Address,
			newUserRequest.CountryCode,
			time.Now().UTC(),
			userID,
		}

		err := tx.Get(&newUserID, SQL, args...)
		if err != nil {
			logrus.Errorf("CreateNewUser: error creating user %v", err)
			return err
		}

		SQL = `INSERT INTO user_profiles
				(user_id, gender, date_of_birth, profile_image_id, state, country)
				VALUES ($1, $2, $3, $4, $5, $6 )`

		_, err = tx.Exec(SQL, newUserID, newUserRequest.Gender, newUserRequest.DateOfBirth, newUserRequest.ProfileImageID, newUserRequest.State, newUserRequest.Country)
		if err != nil {
			logrus.Errorf("CreateNewUser: error creating user_profile %v", err)
			return err
		}

		SQL = `INSERT INTO user_roles
				 (user_id, role)
				 VALUES ($1,$2)`

		_, err = tx.Exec(SQL, newUserID, newUserRequest.Role)
		if err != nil {
			logrus.Errorf("CreateNewUser: error creating user roles err %v", err)
			return err
		}
		return nil
	})

	if txErr != nil {
		logrus.Errorf("CreateNewUser: error in creating user: %v", txErr)
		return nil, txErr
	}

	return &newUserID, nil
}

func (dh *DBHelper) IsUserAlreadyExists(emailID string) (isUserExist bool, user models.UserData, err error) {
	//	language=sql
	SQL := `SELECT id
			FROM users
			WHERE email = lower($1)
			  AND archived_at IS NULL
			  `

	err = dh.DB.Get(&user, SQL, emailID)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("isEmailAlreadyExist: unable to get user from email %v", err)
		return false, user, err
	}

	if err == sql.ErrNoRows {
		return false, user, nil
	}

	return true, user, nil
}

func (dh *DBHelper) IsPhoneNumberAlreadyExist(phone string) (bool, error) {
	// language=sql
	SQL := `SELECT count(*) > 0 
            FROM users
            WHERE archived_at IS NULL
            AND phone  = $1`

	var isPhoneAlreadyExist bool
	err := dh.DB.Get(&isPhoneAlreadyExist, SQL, phone)
	if err != nil {
		logrus.Errorf("IsPhoneNumberAlreadyExist: error getting whether phone exist: %v", err)
		return isPhoneAlreadyExist, err
	}

	return isPhoneAlreadyExist, nil
}

func (dh *DBHelper) FetchUserData(userID int) (models.FetchUserData, error) {
	//language=sql
	SQL := `SELECT  users.id, 
        			users.name, 
        			email, 
        			user_roles.role, 
        			phone,
        			address,
        			user_profiles.gender,
        			user_profiles.date_of_birth,
        			user_profiles.profile_image_id,
        			user_profiles.country,
					user_profiles.state,
					profile_image.url
			FROM users 
			    JOIN user_profiles on users.id = user_profiles.user_id
			    JOIN profile_image on user_profiles.profile_image_id = profile_image.id
			    JOIN user_roles on users.id = user_roles.user_id
			JOIN user_profiles up on up.user_id = users.id 
			WHERE users.id = $1`

	var fetchUserData models.FetchUserData
	err := dh.DB.Get(&fetchUserData, SQL, userID)
	if err != nil {
		logrus.Errorf("FetchUserData: error getting user data: %v", err)
		return fetchUserData, err
	}
	return fetchUserData, nil
}

func (dh *DBHelper) LogInUserUsingEmailAndRole(loginReq models.EmailAndPassword, role models.UserRoles) (userID int, message string, err error) {
	// language=SQL
	SQL := `
		SELECT id,   
			password
		FROM
			users
		WHERE
			email = $1
			AND archived_at IS NULL 
	`

	var user = struct {
		ID             int    `db:"id"`
		HashedPassword string `db:"password"`
	}{}

	if err = dh.DB.Get(&user, SQL, loginReq.Email); err != nil && err != sql.ErrNoRows {
		logrus.Errorf("LogInUserUsingEmailAndRole: error while getting user %v", err)
		return userID, "error getting user", err
	}

	isPasswordMatched := crypto.ComparePasswords(user.HashedPassword, loginReq.Password)

	if !isPasswordMatched {
		return userID, "Password Not Correct", models.ErrorPasswordNotMatched
	}

	var userRole models.UserRoles
	SQL = `
		SELECT
			role
		FROM user_roles
		WHERE user_id = $1
		  	  AND role = $2
			  AND archived_at IS NULL
	`

	err = dh.DB.Get(&userRole, SQL, user.ID, role)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("LogInUserUsingEmailAndRole: error while getting user role:  %v", err)
		return userID, "error getting user role", err
	}
	if err == sql.ErrNoRows {
		return userID, "user role not matched", errors.New("user does not have required access")
	}

	return user.ID, "", nil
}

func (dh *DBHelper) ChangePasswordByUserID(userID int, changePasswordRequest models.ChangePasswordRequest) (bool, error) {
	// language=SQL
	SQL := `
		SELECT id,   
			password
		FROM
			users
		WHERE
			id = $1
			AND archived_at IS NULL`

	var user = struct {
		ID             int    `db:"id"`
		HashedPassword string `db:"password"`
	}{}

	if err := dh.DB.Get(&user, SQL, userID); err != nil && err != sql.ErrNoRows {
		logrus.Errorf("GetPasswordByUserID: error while getting user %v", err)
		return false, err
	}

	isPasswordMatched := crypto.ComparePasswords(user.HashedPassword, changePasswordRequest.OldPassword)

	if !isPasswordMatched {
		return false, models.ErrorPasswordNotMatched
	}
	var password null.String
	if changePasswordRequest.NewPassword != "" {
		password = null.StringFrom(crypto.HashAndSalt(changePasswordRequest.NewPassword))
	}

	SQL = `UPDATE users
			SET   password = $2
			WHERE id = $1`

	_, err := dh.DB.Exec(SQL, user.ID, password)
	if err != nil {
		logrus.Errorf("LogInUserUsingEmailAndRole: error while getting user role:  %v", err)
		return false, err
	}
	return true, nil
}

func (dh *DBHelper) StartNewSession(userID int, request *models.CreateSessionRequest) (string, error) {

	// language=sql
	SQL := `INSERT INTO sessions 
			(user_id, start_time,end_time, platform, model_name, os_version, device_id, token) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)	RETURNING token, id`

	args := []interface{}{
		userID,
		time.Now(),
		time.Now().Add(1 * time.Hour),
		request.Platform,
		request.ModelName,
		request.OSVersion,
		request.DeviceID,
		uuid.New(),
	}

	type sessionDetails struct {
		Token     string `db:"token"`
		SessionID int64  `db:"id"`
	}
	var session sessionDetails
	err := dh.DB.Get(&session, SQL, args...)
	if err != nil {
		logrus.Errorf("StartNewSession: error while starting new session: %v\n", err)
		return session.Token, err
	}

	return session.Token, nil
}

func (dh *DBHelper) FetchUserSessionData(userID int) ([]models.FetchUserSessionsData, error) {
	//language=sql
	SQL := `SELECT  id, user_id,end_time,  token
			FROM sessions
			WHERE user_id = $1`

	fetchUserSessionData := make([]models.FetchUserSessionsData, 0)
	err := dh.DB.Select(&fetchUserSessionData, SQL, userID)
	if err != nil {
		logrus.Errorf("FetchUserSessionData: error getting user session data from database: %v", err)
		return fetchUserSessionData, err
	}
	return fetchUserSessionData, nil
}

func (dh *DBHelper) UpdateSession(sessionId string) error {
	//language=sql
	SQL := `UPDATE sessions
    		SET end_time = $2
			WHERE token = $1`

	_, err := dh.DB.Exec(SQL, sessionId, time.Now().Add(1*time.Hour))
	if err != nil {
		logrus.Errorf("FetchUserSessionData: error getting user session data from database: %v", err)
		return err
	}
	return nil
}

func (dh *DBHelper) EndSession(sessionId string) error {
	//language=sql
	SQL := `UPDATE sessions
    		SET end_time = now()
			WHERE token = $1`

	_, err := dh.DB.Exec(SQL, sessionId)
	if err != nil {
		logrus.Errorf("FetchUserSessionData: error getting user session data from database: %v", err)
		return err
	}
	return nil
}

func (dh *DBHelper) Dashboard() (models.FetchUserData, error) {
	//language=sql
	SQL := `SELECT  COUNT(users.id) as total_users,
			        COUNT(case when up.role = 'dealer' then up.id end ) total_dealers ,
			        COUNT(case when up.role = 'retailer' then up.id end ) total_retailers
			FROM users
			         JOIN user_roles up on up.user_id = users.id
			WHERE   users.archived_at is null
			    AND up.archived_at is null`

	var fetchUserData models.FetchUserData
	err := dh.DB.Get(&fetchUserData, SQL)
	if err != nil {
		logrus.Errorf("IsPhoneNumberAlreadyExist: error getting whether phone exist: %v", err)
		return fetchUserData, err
	}
	return fetchUserData, nil
}

func (dh *DBHelper) RecentUsers(limit int) ([]models.UserInfo, error) {
	//language=sql
	SQL := `select 	new_users.id, 
					new_users.name, 
					new_users.email, 
					ur.role, 
					new_users.created_at
                    from users new_users
                             join user_roles ur on new_users.id = ur.user_id
                    order by created_at desc
                    limit $1`

	UserInfo := make([]models.UserInfo, 0)
	err := dh.DB.Select(&UserInfo, SQL, limit)
	if err != nil {
		logrus.Errorf("RecentUsers: error getting recent users list: %v", err)
		return UserInfo, err
	}
	return UserInfo, nil
}

func (dh *DBHelper) RecentOrders(limit int) ([]models.RecentOrders, error) {
	//language=sql
	SQL := `SELECT 	ur.id, 
       				ur.name, 
       				orders.id, 
       				orders.quantity, 
       				orders.order_status,
       				orders.created_at
                    FROM orders
                             JOIN users ur ON orders.ordered_by = ur.id
                    ORDER BY created_at DESC
                    LIMIT $1`

	UserInfo := make([]models.RecentOrders, 0)
	err := dh.DB.Get(&UserInfo, SQL, limit)
	if err != nil {
		logrus.Errorf("RecentOrders: error getting recent order list: %v", err)
		return UserInfo, err
	}
	return UserInfo, nil
}

func (dh *DBHelper) OrderSummary() (models.OrderSummary, error) {
	//language=sql
	SQL := `SELECT 	COUNT(orders.id) filter ( where  completed_at IS NULL ) as total_created_orders,
			        COUNT(orders.id) filter ( where  orders.order_status = 'open' ) as open_deliveries ,
			        COUNT(orders.id) filter ( where  orders.order_status = 'in stock' ) as in_stock,
        			COUNT(orders.id) filter ( where  orders.order_status = 'in transfer' ) as in_transfer ,
       				COUNT(orders.id) filter ( where  orders.order_status = 'sold out' )  as sold_out      				
                    FROM orders`

	var orderSummary models.OrderSummary
	err := dh.DB.Get(&orderSummary, SQL)
	if err != nil {
		logrus.Errorf("OrderSummary: error getting order summary data : %v", err)
		return orderSummary, err
	}
	return orderSummary, nil
}

func (dh *DBHelper) UsersAll(userID, limit, offset int, role models.UserRoles) ([]models.FetchUserData, error) {
	roleCheck := true
	if role == models.All {
		roleCheck = false
	}
	args := []interface{}{
		userID,
		!roleCheck,
		role,
		limit,
		offset,
	}
	fetchAllUserData := make([]models.FetchUserData, 0)

	//language=sql
	SQL := `SELECT  users.id,
        			users.name,
        			users.email,
        			ur.role,
        			users.created_at
			FROM users 
			JOIN user_roles ur ON users.id = ur.user_id
			WHERE users.id <> $1
			AND ($2 OR ur.role = $3)
			ORDER BY users.created_at DESC
			LIMIT $4
			OFFSET $5`

	err := dh.DB.Select(&fetchAllUserData, SQL, args...)
	if err != nil {
		logrus.Errorf("UsersAll: error getting all user info : %v", err)
		return fetchAllUserData, err
	}
	return fetchAllUserData, nil
}

func (dh *DBHelper) CreateOrder(userID int, address string, order models.Order) (models.CreatedOrder, error) {
	Args := []interface{}{
		userID,
		order.Quantity,
		uuid.New(),
		address,
		time.Now().UTC(),
		time.Now().UTC(),
	}

	//language=sql
	SQL := `INSERT INTO orders 
    		(ordered_by, quantity, reference_no, shipping_address, created_at, updated_at) 
			VALUES ($1,$2,$3,$4,$5,$6) returning id`

	var orderData models.CreatedOrder
	err := dh.DB.Get(&orderData, SQL, Args...)
	if err != nil {
		logrus.Errorf("CreateOrder: error creating order : %v", err)
		return orderData, err
	}
	return orderData, nil
}

func (dh *DBHelper) Upload(userID int) (models.FetchUserData, error) {
	//language=sql
	SQL := `SELECT  users.id, users.name, email, phone, address,gender, date_of_birth
			FROM users 
			JOIN user_profiles up on up.user_id = users.id 
			WHERE users.id = $1`

	var fetchUserData models.FetchUserData
	err := dh.DB.Get(&fetchUserData, SQL, userID)
	if err != nil {
		logrus.Errorf("FetchUserData: error getting user data: %v", err)
		return fetchUserData, err
	}
	return fetchUserData, nil
}

func (dh *DBHelper) GetUserInfoByEmail(email string) (models.GetUserDataByEmail, error) {
	//language=sql
	SQL := `SELECT  users.id, users.name, user_roles.role,email, phone, address,gender, date_of_birth
			FROM users 
			JOIN user_roles ON users.id = user_roles.user_id
			JOIN user_profiles up on up.user_id = users.id 
			WHERE users.email = $1`

	var getUserDataByEmail models.GetUserDataByEmail
	err := dh.DB.Get(&getUserDataByEmail, SQL, email)
	if err != nil {
		logrus.Errorf("FetchUserData: error getting user data: %v", err)
		return getUserDataByEmail, err
	}
	return getUserDataByEmail, nil
}

func (dh *DBHelper) EditProfile(userID int, editProfileRequest models.EditProfile) error {

	txErr := dbutil.WithTransaction(dh.DB, func(tx *sqlx.Tx) error {
		//language=sql
		SQL := `UPDATE users
    		SET  name = $1,
    		     email = $2, 
    		     phone = $3, 
    		     address = $4, 
    		     country_code = $5, 
    		     updated_at = now()
    		WHERE users.id = $6`

		Args := []interface{}{
			editProfileRequest.Name,
			editProfileRequest.Email,
			editProfileRequest.Phone,
			editProfileRequest.Address,
			editProfileRequest.CountryCode.String,
			userID,
		}

		_, err := tx.Exec(SQL, Args...)
		if err != nil {
			logrus.Errorf("EditProfile: error getting user data: %v", err)
			return err
		}

		//language=sql
		SQL = `UPDATE user_profiles
    		SET  country = $2,
    		     state = $3, 
    		     profile_image_id = $4,
    		     updated_at = now()
    		WHERE user_id = $1`

		Args = []interface{}{
			userID,
			editProfileRequest.Country,
			editProfileRequest.State,
			editProfileRequest.ProfileImageID,
		}

		_, err = tx.Exec(SQL, Args...)
		if err != nil {
			logrus.Errorf("EditProfile: error getting user data: %v", err)
			return err
		}
		return nil
	})
	if txErr != nil {
		logrus.Errorf("EditProfile: error in updating user profile: %v", txErr)
		return txErr
	}
	return nil
}

func (dh *DBHelper) GetCountryAndState() ([]models.CountryAndState, error) {
	//language=sql
	SQL := `SELECT  country.id,
         			country.country,
         			country.country_code,
         			array_remove(array_agg(s.id), NULL) state_id,
         			array_remove(array_agg(s.state), NULL) state
			FROM country
         	LEFT JOIN state s on country.id = s.country_id
			WHERE country.archived_at IS NULL
  			AND s.archived_at IS NULL
			GROUP BY country.id, country.country, country.country_code`

	countryAndState := make([]models.CountryAndState, 0)
	err := dh.DB.Select(&countryAndState, SQL)
	if err != nil {
		logrus.Errorf("GetCountryAndState: error getting order summary data : %v", err)
		return countryAndState, err
	}
	return countryAndState, nil
}
