package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/camp_backend/internal/model"
)

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{
		db: db,
	}
}

func (r *UsersPostgres) CreateUser(user model.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, password_hash, name, tag, about, photo_url) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Email, user.PasswordHash, user.Name, user.Tag, user.About, user.PhotoURL)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UsersPostgres) SignIn(email string, passwordHash string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 and password_hash = $2", usersTable)
	if err := r.db.Get(&user, query, email, passwordHash); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UsersPostgres) EditUserInfo(user model.User) error {
	query := fmt.Sprintf("UPDATE %s SET name = $1, about = $2, tag = $3 WHERE id = $4", usersTable)
	_, err := r.db.Exec(query, user.Name, user.About, user.Tag, user.Id)
	return err
}

func (r *UsersPostgres) CheckTagUnique(tag string) bool {
	var isUnique bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE tag = $1)", usersTable)
	row := r.db.QueryRow(query, tag)
	if err := row.Scan(&isUnique); err != nil {
		return false
	}
	return !isUnique
}

func (r *UsersPostgres) FindUserByTag(tag string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE tag = $1", usersTable)
	if err := r.db.Get(&user, query, tag); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UsersPostgres) GetOneUser(userId int) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", usersTable)
	if err := r.db.Get(&user, query, userId); err != nil {
		return model.User{}, err
	}
	return user, nil
}
