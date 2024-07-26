package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	db "github.com/quikkom/auth/database"
	"github.com/quikkom/auth/env"
	"github.com/quikkom/auth/proto"
	"google.golang.org/grpc"
)

// Creates a new listener for all interfaces and listens the given port.
func newListener(port int) (net.Listener, string) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	listen, err := net.Listen("tcp", addr)

	if err != nil {
		log.Panicf("Couldn't listen: %v", err)
	}

	return listen, addr
}

// Finds the port number which server runs on.
func getPort() int {
	env := env.FindEnv("PORT", "5000")
	port, err := strconv.Atoi(env)

	if err != nil {
		log.Panicf("Invalid port number: %s", env)
	}

	return port
}

// Registers the gRPC server
func registerServer() *grpc.Server {
	server := grpc.NewServer()
	proto.RegisterAuthServer(server, &AuthServer{})

	return server
}

// Creates a handler for Ctrl+C, SIGTERM signal
func setupHandleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		slog.Debug("Closing database connection")

		if db.DBConn != nil {
			db.DBConn.Close(context.Background())
		}
		os.Exit(0)
	}()
}

func main() {
	env.Fill()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	setupHandleSignal()

	listener, addr := newListener(getPort())
	db.CreateDBConnection()
	server := registerServer()

	slog.Info(fmt.Sprintf("Starting server at %s", addr))
	err := server.Serve(listener)
	if err != nil {
		log.Fatalf("Error when starting server: %v", err)
	}
}
