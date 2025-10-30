package model

type Room struct {
	Id           int    `json:"id" db:"id"`
	UserId       int    `json:"user_id" db:"user_id"`
	UserPhotoURL string `json:"user_photo_url" db:"user_photo_url"`
	UserName     string `json:"user_name" db:"user_name"`
	UserLanguage string `json:"user_language" db:"language"`
	Translating  bool   `json:"translating" db:"translating"`
}
