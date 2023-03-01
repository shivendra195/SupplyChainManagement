package dbhelperprovider

import (
	"database/sql"
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
	SQL := `SELECT  name, email, phone, address,gender, date_of_birth
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
