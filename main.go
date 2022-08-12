package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/lov3allmy/tages/cmd"
	pb "github.com/lov3allmy/tages/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	dbSource := "postgresql://postgres:postgrespw@postgres:5432/postgres?sslmode=disable"
	repo, err := cmd.NewRepository(dbSource)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.DB.Close()

	m, err := migrate.New("file://schema", dbSource)
	if err != nil {
		log.Fatalf("failed to create migration: %v", err)
	}

	err = m.Force(1)
	if err != nil {
		log.Fatalf("failed to force migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to up migration: %v", err)
	}
	defer func() {
		m.Down()
		m.Close()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	s := grpc.NewServer()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		<-sigCh
		s.GracefulStop()
		wg.Done()
	}()

	pb.RegisterImageStorageServiceServer(s, cmd.NewServer(repo))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	wg.Wait()
}
