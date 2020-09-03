package model

//ResponseAPI is a struct for Response Json API
type ResponseAPI struct {
	Success     bool   `json:"success"`
	Error       string `json:"error_message"`
	AccessToken string `json:"access_token,omitempty"`
	Cards       []Card `json:"cards,omitempty"`
}
