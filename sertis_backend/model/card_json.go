package model

// Card is a struct use for receive, response and query from database
type Card struct {
	ID       int    `json:"id,omitempty" db:"id"`
	UserID   int    `json:"userid,omitempty" db:"user_id"`
	Name     string `json:"name" db:"name"`
	Status   string `json:"status" db:"status"`
	Content  string `json:"content" db:"content"`
	Category string `json:"category" db:"category"`
	Author   string `json:"author,omitempty"`
}
