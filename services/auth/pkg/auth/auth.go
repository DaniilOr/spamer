package auth

import (
	"context"
	"errors"
	"github.com/DaniilOr/spamer/services/auth/pkg/jwt/symmetric"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPass = errors.New("invalid password")

type Service struct {
	pool *pgxpool.Pool
}
type UserDetails struct {
	ID    int64
	Password []byte
	Login string
	Roles []string
}
type Data struct {
	UserID int64    `json:"userId"`
	Roles  []string `json:"roles"`
	Issued int64    `json:"iat"`
	Expire int64    `json:"exp"`
}
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) Login(ctx context.Context, login string, password string) (string, int64, error) {

	var personal UserDetails
	//log.Println("Quering db")
	err := s.pool.QueryRow(ctx, `
		SELECT id, login, password, roles FROM users WHERE login = $1
	`, login).Scan(&personal.ID, &personal.Login, &personal.Password, &personal.Roles)
	if err != nil {
		if err != pgx.ErrNoRows {
			return "", 0, ErrUserNotFound
		}
		return "", 0, err
	}
	//log.Printf("Got password %v and have %v\n", personal.Password, []byte(password))
	//log.Printf("%v", personal)
	err = bcrypt.CompareHashAndPassword(personal.Password, []byte(password))
	if err != nil {
		log.Println("Error with comparison")
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", 0, ErrInvalidPass
		}
		return "", 0, err
	}
	data := &Data{
		UserID: personal.ID,
		Roles:  personal.Roles,
		Issued: time.Now().Unix(),
		Expire: time.Now().Add(time.Minute * 10).Unix(),
	}
	key := []byte("some secter key goes here")
	log.Printf("Encrypting %v", data)
	token, err := symmetric.Encode(data, key)

	if err != nil{
		return "", 0, err
	}
	log.Println("Inserting into tokens")
	_, err = s.pool.Exec(ctx, `INSERT INTO tokens (token, userid) VALUES ($1, $2)`, token, data.UserID)
	if err != nil {
		return "", 0, err
	}
	log.Printf("token: %v %d", token, data.Expire)
	return token, data.Expire, nil
}

func (s *Service) UserID(ctx context.Context, token string) (userID int64, err error) {
	err = s.pool.QueryRow(ctx, `
		SELECT userid FROM tokens WHERE token = $1
	`, token).Scan(&userID)
	if err != nil {
		if err != pgx.ErrNoRows {
			return 0, ErrUserNotFound
		}
		return 0, err
	}

	return userID, nil
}
