package storage

import (
	"sertis_app/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	creds := model.Credentials{Username: "jojo", Password: "121912"}
	query := "INSERT INTO `user`"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("jojo", "121912").WillReturnResult(sqlmock.NewResult(1, 1))

	errCreateAccount := CreateAccount(db, creds)
	assert.NoError(t, errCreateAccount)
}

func TestVerifyUserAndPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	creds := model.Credentials{Username: "jojo", Password: "121912"}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	query := "SELECT id, password FROM user WHERE username = ?"

	rows := sqlmock.NewRows([]string{"id", "password"}).AddRow(1, hashedPassword)

	mock.ExpectQuery(query).WithArgs("jojo").WillReturnRows(rows)

	id, err := VerifyUserAndPassword(db, creds)
	assert.NotZero(t, id)
	assert.NoError(t, err)
}

func TestCreateCard(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	card := model.Card{UserID: 1, Name: "test", Status: "n/a", Category: "cat", Content: "content"}
	query := "INSERT INTO `card`"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(1, "test", "n/a", "content", "cat").WillReturnResult(sqlmock.NewResult(1, 1))

	errCreateCard := CreateCard(db, card)
	assert.NoError(t, errCreateCard)
}
