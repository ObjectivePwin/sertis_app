package storage

import (
	"sertis_app/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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
