package miniblog

import (
	"database/sql"
	"errors"
	"sertis_app/model"
	storage "sertis_app/storage"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Blog struct {
	db            *sql.DB
	activeSession map[int]string
	jwtSecketKey  []byte
}

//NewBlog is a function to create Blog object
func NewBlog(db *sql.DB) *Blog {
	blog := Blog{
		db:            db,
		activeSession: make(map[int]string),
		jwtSecketKey:  []byte("sertis_seckey_key"),
	}

	return &blog
}

//CreateAccount is a function to check already have account if not then create it
func (b *Blog) CreateAccount(creds model.Credentials) error {
	if storage.CheckAlreadyHaveAccount(b.db, creds) {
		return errors.New("Already Have Account")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	creds.Password = string(hashedPassword)

	err := storage.CreateAccount(b.db, creds)
	if err != nil {
		return err
	}
	return nil
}

//LoginAddCreateAccessToken is a function to check user and password
func (b *Blog) LoginAddCreateAccessToken(creds model.Credentials) (string, error) {
	id, err := storage.VerifyUserAndPassword(b.db, creds)
	if err != nil {
		return "", errors.New("your username or password is incorrect")
	}

	// Declare the expiration time of the token
	// here, we have kept it as 60 minutes
	expirationTime := time.Now().Add(60 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &model.JWTClaims{
		ID:       id,
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(b.jwtSecketKey)
	return tokenString, nil
}
