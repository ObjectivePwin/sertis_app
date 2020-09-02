package miniblog

import (
	"database/sql"
	"errors"
	"sertis_app/model"
	storage "sertis_app/storage"
)

type Blog struct {
	db            *sql.DB
	activeSession map[int]string
}

//NewBlog is a function to create Blog object
func NewBlog(db *sql.DB) *Blog {
	blog := Blog{
		db:            db,
		activeSession: make(map[int]string),
	}

	return &blog
}

func (b *Blog) CreateAccount(creds model.Credentials) error {
	if storage.CheckAlreadyHaveAccount(b.db, creds) {
		return errors.New("Already Have Account")
	}

	err := storage.CreateAccount(b.db, creds)
	if err != nil {
		return err
	}
	return nil
}
