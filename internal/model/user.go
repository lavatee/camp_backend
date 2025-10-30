package model

type User struct {
	Id           int    `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	Name         string `json:"name" db:"name"`
	Tag          string `json:"tag" db:"tag"`
	About        string `json:"about" db:"about"`
	PhotoURL     string `json:"photo_url" db:"photo_url"`
	Language     string `json:"language" db:"language"`
}
