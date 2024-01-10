package db

import (
	"context"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type User struct {
	ID                int64
	Name              string
	PhoneNumber       string
	OTP               string
	OTPExpirationTime time.Time
}

type DB struct {
    *pgxpool.Pool
}


func NewDB(pool *pgxpool.Pool) *DB {
    return &DB{Pool: pool}
}

func (db *DB) CreateUser(ctx context.Context, name, phoneNumber string) (int64, error) {
	var userID int64
	err := db.QueryRow(ctx, "INSERT INTO users (name, phone_number) VALUES ($1, $2) RETURNING id", name, phoneNumber).Scan(&userID)
	return userID, err
}

func (db *DB) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error) {
	var user User
	err := db.QueryRow(ctx, "SELECT * FROM users WHERE phone_number = $1", phoneNumber).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.OTP, &user.OTPExpirationTime)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (db *DB) GenerateOTP(ctx context.Context, otp string, expirationTime time.Time, phoneNumber string) (*User, error) {
	var user User
	query := "UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE phone_number = $3 RETURNING id, name, phone_number, otp, otp_expiration_time"
	err := db.QueryRow(ctx, query, otp, expirationTime, phoneNumber).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.OTP, &user.OTPExpirationTime)
	return &user, err
}

func (db *DB) VerifyOTP(ctx context.Context, phoneNumber string) (*User, error) {
	var user User
	query := "SELECT otp, otp_expiration_time FROM users WHERE phone_number = $1"
	err := db.QueryRow(ctx, query, phoneNumber).Scan(&user.OTP, &user.OTPExpirationTime)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &user, err
}