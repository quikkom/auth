package main

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	db "github.com/quikkom/auth/database"
	"github.com/quikkom/auth/env"
	"github.com/quikkom/auth/proto"
)

type AuthServer struct {
	proto.UnimplementedAuthServer
}

func sha256(data string) string {
	sha256 := crypto.SHA256.New()
	sha256.Write([]byte(data))
	return fmt.Sprintf("%x", sha256.Sum(nil))
}

func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(env.AUTH_SECRET))
	if err != nil {
		return "", fmt.Errorf("couldn't generate JWT: %s", err)
	}

	return tokenString, nil
}

func (a *AuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {

	if req.Email == nil && req.Username == nil {
		return nil, errors.New("missing user information")
	}

	identifier := ""
	identifierValue := ""

	if req.Email != nil {
		identifier = "email"
		identifierValue = *req.Email
	} else {
		identifier = "username"
		identifierValue = *req.Username
	}

	username := ""
	password := ""
	row := db.DBConn.QueryRow(context.Background(),
		fmt.Sprintf("SELECT username, password FROM users WHERE %s=$1", identifier),
		identifierValue)
	err := row.Scan(&username, &password)

	if err != nil {
		slog.Error(fmt.Sprintf("User not found: %v", err))
		return nil, errors.New("user not found")
	}

	if sha256(req.Password) != password {
		return nil, errors.New("wrong password")
	}

	token, err := generateToken(username)
	if err != nil {
		return nil, err
	}

	slog.Debug(fmt.Sprintf("User %s authenticated successfully", username))

	return &proto.LoginResponse{
		Token:    token,
		Username: username,
	}, nil
}
