package storage

import (
	"database/sql"
	"sertis_app/model"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var dataSourceName = "sertis:sertis@tcp(localhost:3306)/sertis"

//CreateDBConnection is a function to create db connection
func CreateDBConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName+"?parseTime=true&loc=Asia%2FBangkok")
	if err != nil {
		return nil, err
	}
	return db, nil
}

//CheckAlreadyHaveAccount is a function that check already have account
func CheckAlreadyHaveAccount(db *sql.DB, credentials model.Credentials) bool {
	var id int

	err := db.QueryRow("SELECT id FROM user WHERE username = ?", credentials.Username).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		zap.S().Info("CheckAlreadyHaveAccount ", err)

		return true
	}
	return false
}

//CreateAccount is a function that create account
func CreateAccount(db *sql.DB, credentials model.Credentials) error {
	stmt, err := db.Prepare("INSERT INTO `user` (`username`, `password`) VALUES (?, ?)")
	if err != nil {
		zap.S().Info("Prepare CreateAccount error ", err)
		return err
	}

	_, err = stmt.Exec(credentials.Username, credentials.Password)
	defer stmt.Close()
	if err != nil {
		zap.S().Info("CreateAccount insert error ", err)
		return err
	}
	return nil

}
