package main

import (
	"context"
	"crawler/internal/repository"
	"crawler/internal/service"
	"crawler/internal/sys"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := log.New(os.Stdout, "Spyder : ", log.LstdFlags|log.Lshortfile)
	cfg, err := sys.LoadConfig("conf.env")
	if err != nil {
		log.Fatalf("err : %v", err)
	}
	if err := run(cfg, log); err != nil {
		log.Fatalf("fatal error : %v \n", err)
	}
}

func run(cfg *sys.Config, log *log.Logger) error {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return err
	}

	if err = db.Ping(ctx); err != nil {
		return err
	}

	postRepo := repository.NewPostRepository(db, log)
	postService := service.NewPostService(postRepo, log)
	receiver := service.NewAPIResponseReceiver()
	worker := service.NewWorker(postService, receiver, log)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		worker.Start(ctx, time.Minute)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case <-shutdown:
		cancel()
		log.Println("Service stopped working")
	}

	return nil
}
