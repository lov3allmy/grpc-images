package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/lov3allmy/tages/internal"
	pb "github.com/lov3allmy/tages/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	//dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
	//	"postgres",
	//	"postgrespw",
	//	"localhost",
	//	"5436",
	//	"postgres")

	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("connect to database: success")

	repo := internal.NewRepository(db)

	m, err := migrate.New("file://schema", dbSource)
	if err != nil {
		log.Fatalf("failed to create migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to up migration: %v", err)
	}
	defer func() {
		m.Close()
	}()
	log.Println("migration up: success")

	sm := &internal.ServerMiddleware{}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(sm.Interceptor))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		grpcServer.GracefulStop()
	}()

	pb.RegisterImageStorageServiceServer(grpcServer, internal.NewServer(repo))
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		log.Println("start GRPC server")
	}
}
