package miniblog

import (
	"database/sql"
	"errors"
	"sertis_app/model"
	storage "sertis_app/storage"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
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

//VerifyJWTToken is a function to verify token
func (b *Blog) VerifyJWTToken(token string) (*model.JWTClaims, error) {
	// Initialize a new instance of `Claims`
	claims := &model.JWTClaims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return b.jwtSecketKey, nil
	})
	if err != nil {
		zap.L().Debug("Token is in valid")
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("Token is in valid")
		}
		return nil, errors.New("Token is in valid")
	}
	if !tkn.Valid {
		zap.L().Debug("Token is in valid")
		return nil, errors.New("Token is in valid")
	}

	return claims, nil
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

//CreateNewCard is a function that create new card
func (b *Blog) CreateNewCard(card model.Card) error {
	err := storage.CreateCard(b.db, card)

	if err != nil {
		return err
	}
	return nil
}

//GetAllCard is a function that get cards
func (b *Blog) GetAllCard() []model.Card {
	cards, err := storage.GetAllCard(b.db)
	if err != nil {
		zap.S().Debug("GetAllCard ", err)
		return []model.Card{}
	}
	return cards
}

//UpdateCard is a function that update card
func (b *Blog) UpdateCard(card model.Card) error {
	err := storage.UpdateCard(b.db, card)
	if err != nil {
		return err
	}
	return nil
}

//DeleteCard is a function that delete card
func (b *Blog) DeleteCard(cardID int, userID int) error {
	err := storage.DeleteCard(b.db, cardID, userID)
	if err != nil {
		return err
	}
	return nil
}
