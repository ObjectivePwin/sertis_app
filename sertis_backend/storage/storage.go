package storage

import (
	"database/sql"
	"errors"
	"sertis_app/model"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var dataSourceName = "sertis:sertis@tcp(mariadb-master:3306)/sertis"

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

	err := db.QueryRow("SELECT id FROM user WHERE username = ? limit 1", credentials.Username).Scan(&id)
	if err != nil {
		zap.S().Info("CheckAlreadyHaveAccount ", err)

		return false
	}

	return true
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

//VerifyUserAndPassword is a function that verify account is have return id
func VerifyUserAndPassword(db *sql.DB, credentials model.Credentials) (int, error) {
	var id int
	var password string

	err := db.QueryRow("SELECT id, password FROM user WHERE username = ?", credentials.Username).Scan(&id, &password)
	if err != nil {
		zap.S().Info("VerifyUserAndPassword query error ", err)

		return 0, errors.New("Not Found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(credentials.Password))
	if err != nil {
		zap.S().Info("VerifyUserAndPassword password not match ", err)
		return 0, errors.New("Password Not Match")

	}
	return id, nil
}

//CreateCard is a function that insert card to database
func CreateCard(db *sql.DB, card model.Card) error {
	stmt, err := db.Prepare("INSERT INTO `card` (`user_id`, `name`, `status`, `content`, `category`) VALUES (?, ?, ?, ?, ?)")
	_, err = stmt.Exec(card.UserID, card.Name, card.Status, card.Content, card.Category)

	defer stmt.Close()
	if err != nil {
		zap.S().Info("CreateCard insert error ", err)
		return err
	}
	return nil
}

//GetAllCard is a function that query all card from database
func GetAllCard(db *sql.DB) ([]model.Card, error) {
	rows, err := db.Query("SELECT c.`id`, `username`, `name`, `status`, `content`, `category` FROM `card` c, user u WHERE c.user_id = u.id")
	if err != nil {
		zap.S().Info("CreateCard insert error ", err)
		return nil, err
	}

	cards := []model.Card{}
	for rows.Next() {
		card := model.Card{}
		rowError := rows.Scan(&card.ID, &card.Author, &card.Name, &card.Status, &card.Content, &card.Category)
		if rowError != nil {
			return nil, rowError
		}
		cards = append(cards, card)
	}

	return cards, nil
}

//UpdateCard is a function that update card in database
func UpdateCard(db *sql.DB, card model.Card) error {
	stmt, err := db.Prepare("UPDATE card SET name = ?, status = ?, content = ?, category = ? WHERE id = ? and user_id = ?")
	_, err = stmt.Exec(card.Name, card.Status, card.Content, card.Category, card.ID, card.UserID)

	defer stmt.Close()
	if err != nil {
		zap.S().Info("UpdateCard insert error ", err)
		return err
	}
	return nil
}

//DeleteCard is a function that delete card in database
func DeleteCard(db *sql.DB, cardID int, userID int) error {
	stmt, err := db.Prepare("DELETE FROM `card` WHERE id = ? and user_id = ?")
	_, err = stmt.Exec(cardID, userID)

	defer stmt.Close()
	if err != nil {
		zap.S().Info("DeleteCard insert error ", err)
		return err
	}
	return nil
}
