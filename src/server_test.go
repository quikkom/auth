package main

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	db "github.com/quikkom/auth/database"
	"github.com/quikkom/auth/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestAuthService_Login(t *testing.T) {
	log.SetOutput(io.Discard)

	// Initialize gRPC
	lis := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	proto.RegisterAuthServer(server, &AuthServer{})
	go func() {
		err := server.Serve(lis)
		if err != nil {
			log.Panicf("Couldn't start server: %s", err)
		}
	}()
	t.Cleanup(func() {
		server.GracefulStop()
	})

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return lis.DialContext(context.Background())
	}))

	if err != nil {
		log.Panicf("Couldn't connect to the server: %s", err)
	}

	client := proto.NewAuthClient(conn)
	t.Cleanup(func() {
		conn.Close()
	})

	db.CreateDBConnection()

	username := "test"
	password := "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3" // 123
	email := "test@test.com"
	db.DBConn.Exec(context.Background(), `
		INSERT INTO users (username, email, "password")
			VALUES ($1, $2, $3); 
	`, username, email, password)
	t.Cleanup(func() {
		db.DBConn.Exec(context.Background(), "DELETE FROM users WHERE username=$1", username)
		db.DBConn.Close(context.Background())
	})
	t.Run("Login user", func(t *testing.T) {
		res, err := client.Login(context.Background(), &proto.LoginRequest{
			Username: &username,
			Password: "123",
		})

		assert.NoError(t, err, "Unexpected error")
		assert.NotNil(t, res, "Response is nil")
		assert.Equal(t, username, res.Username, "Wrong user")
		assert.NotEmpty(t, res.Token, "Empty token")

	})

	t.Run("Wrong password", func(t *testing.T) {
		res, err := client.Login(context.Background(), &proto.LoginRequest{
			Username: &username,
			Password: "123!",
		})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "wrong password")
		assert.Nil(t, res)
	})

	t.Run("User not found", func(t *testing.T) {
		dummyUsername := "noone"
		res, err := client.Login(context.Background(), &proto.LoginRequest{
			Username: &dummyUsername,
			Password: "123!",
		})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "user not found")
		assert.Nil(t, res)
	})
}
