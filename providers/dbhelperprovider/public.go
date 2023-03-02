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

func (dh *DBHelper) CreateNewUser(newUserRequest *models.CreateNewUserRequest) (*int, error) {
	var newUserID int
	var password null.String
	if newUserRequest.Password != "" {
		password = null.StringFrom(crypto.HashAndSalt(newUserRequest.Password))
	}

	txErr := dbutil.WithTransaction(dh.DB, func(tx *sqlx.Tx) error {

		SQL := `INSERT INTO users
				(name, email, phone, password, address, country_code, created_at)
				VALUES (trim($1), lower(trim($2)), $3, $4,$5,$6,$7)
				RETURNING id`

		args := []interface{}{
			newUserRequest.Name,
			newUserRequest.Email,
			newUserRequest.Phone,
			password,
			newUserRequest.Address,
			newUserRequest.CountryCode,
			time.Now().UTC(),
		}

		err := tx.Get(&newUserID, SQL, args...)
		if err != nil {
			logrus.Errorf("CreateNewUser: error creating user %v", err)
			return err
		}

		SQL = `INSERT INTO user_profiles
				(user_id, gender, date_of_birth)
				VALUES ($1, $2, $3)`

		_, err = tx.Exec(SQL, newUserID, newUserRequest.Gender, newUserRequest.DateOfBirth)
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
	SQL := `SELECT  users.id, users.name, email, phone, address,gender, date_of_birth
			FROM users 
			JOIN user_profiles up on up.user_id = users.id 
			WHERE users.id = $1`

	var fetchUserData models.FetchUserData
	err := dh.DB.Get(&fetchUserData, SQL, userID)
	if err != nil {
		logrus.Errorf("IsPhoneNumberAlreadyExist: error getting whether phone exist: %v", err)
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

func (dh *DBHelper) FetchUserSessionData(userID int, deviceId, modelName, osVersion, platform string) ([]models.FetchUserSessionsData, error) {
	//language=sql
	SQL := `SELECT  id, user_id,end_time,  token
			FROM sessions
			WHERE user_id = $1
			AND platform = $2
			AND model_name = $3
			AND os_version = $4
			AND device_id = $5`

	fetchUserSessionData := make([]models.FetchUserSessionsData, 0)
	err := dh.DB.Select(&fetchUserSessionData, SQL, userID, platform, modelName, osVersion, deviceId)
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
