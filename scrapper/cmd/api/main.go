package main

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"scrapper/cmd/api/handlers"
	"scrapper/internal/sys"
	"syscall"
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot create logger")
		return
	}

	conf, err := sys.LoadConfig("conf.env")
	if err != nil {
		log.Sugar().Errorf("error : %v", err)
		return
	}
	if err := run(conf, log.Sugar()); err != nil {
		log.Sugar().Errorf("error : %v", err)
	}
}

func run(conf *sys.Config, log *zap.SugaredLogger) error {
	ctx := context.Background()

	db, err := sys.Connect(ctx, conf)
	if err != nil {
		return err
	}

	router := handlers.API(db)

	srv := http.Server{
		Addr:    conf.ADDR,
		Handler: router,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT)

	serverErrors := make(chan error, 1)
	log.Infof("Starting server on [ADDR]: %s", conf.ADDR)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			serverErrors <- err
		}
	}()

	select {
	case <-shutdown:
		if err := srv.Shutdown(ctx); err != nil {
			log.Errorf("server shutdown error : %v", err)
			return err
		}
		log.Info("server was stopped gracefully")

	case err := <-serverErrors:
		log.Errorf("server error : %v", err)
		return err
	}

	return nil
}
