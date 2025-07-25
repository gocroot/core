package models

// Model struct
type LoginAuth struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
