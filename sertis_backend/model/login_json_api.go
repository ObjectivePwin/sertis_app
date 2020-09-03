package model

// Credentials is a struct use for signin or sign up
type Credentials struct {
	ID       int    `db:"id"`
	Password string `json:"password,omitempty" db:"password"`
	Username string `json:"username" db:"username"`
}
