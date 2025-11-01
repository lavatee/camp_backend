package service

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lavatee/camp_backend/internal/model"
	"github.com/lavatee/camp_backend/internal/repository"
	"github.com/minio/minio-go/v7"
)

const (
	salt       = "ckd3sk0fcakh"
	accessTTL  = 15 * time.Minute
	refreshTTL = 20 * 24 * time.Hour
	tokenKey   = "mironpidoras"
)

type UsersService struct {
	repo   *repository.Repository
	s3     *minio.Client
	bucket string
}

func NewUsersService(repo *repository.Repository, s3 *minio.Client, bucket string) *UsersService {
	return &UsersService{
		repo:   repo,
		s3:     s3,
		bucket: bucket,
	}
}

func (s *UsersService) hashPassword(password string) string {
	sha := sha1.New()
	sha.Write([]byte(password))
	return fmt.Sprintf("%x", sha.Sum([]byte(salt)))
}

func (s *UsersService) SignUp(user model.User) (int, error) {
	user.PasswordHash = s.hashPassword(user.PasswordHash)
	user.PhotoURL = s.getMediaURL(fmt.Sprint(user.Id))
	return s.repo.Users.CreateUser(user)
}

func (s *UsersService) NewToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}
	return stringToken, nil
}

func (s *UsersService) SignIn(email string, password string) (string, string, error) {
	user, err := s.repo.Users.SignIn(email, s.hashPassword(password))
	if err != nil {
		return "", "", err
	}
	userId := user.Id
	accessClaims := jwt.MapClaims{
		"exp": time.Now().Add(accessTTL).Unix(),
		"id":  userId,
	}
	refreshClaims := jwt.MapClaims{
		"exp": time.Now().Add(refreshTTL).Unix(),
		"id":  userId,
	}
	access, err := s.NewToken(accessClaims)
	if err != nil {
		return "", "", err
	}
	refresh, err := s.NewToken(refreshClaims)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *UsersService) Refresh(refreshToken string) (string, string, error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token is invalid")
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		return "", "", err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		accessClaims := jwt.MapClaims{
			"exp": time.Now().Add(accessTTL).Unix(),
			"id":  claims["id"],
		}
		refreshClaims := jwt.MapClaims{
			"exp": time.Now().Add(refreshTTL).Unix(),
			"id":  claims["id"],
		}
		access, err := s.NewToken(accessClaims)
		if err != nil {
			return "", "", err
		}
		refresh, err := s.NewToken(refreshClaims)
		if err != nil {
			return "", "", err
		}
		return access, refresh, nil
	}
	return "", "", fmt.Errorf("token is invalid")
}

func (s *UsersService) ParseToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token is invalid")
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token is expired")
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token is expired")
}

func (s *UsersService) GetOneUser(userId int) (model.User, error) {
	return s.repo.Users.GetOneUser(userId)
}

func (s *UsersService) FindUserByTag(tag string) (model.User, error) {
	return s.repo.Users.FindUserByTag(tag)
}

func (s *UsersService) CheckTagUnique(tag string) bool {
	return s.repo.CheckTagUnique(tag)
}

func (s *UsersService) EditUserInfo(user model.User) error {
	return s.repo.Users.EditUserInfo(user)
}

func (s *UsersService) getMediaURL(key string) string {
	return fmt.Sprintf("https://5a1bc5f7-b5c2-4a61-969a-beacbd4d7999.selstorage.ru/%s", key)
}

func (s *UsersService) NewProfilePhoto(userId int, file multipart.File) (string, error) {
	key := fmt.Sprint(userId)
	url := s.getMediaURL(key)
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	_, err = s.s3.PutObject(context.Background(), s.bucket, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return url, nil
}
