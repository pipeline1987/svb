package models

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}
