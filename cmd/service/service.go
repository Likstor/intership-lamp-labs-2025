package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"service/internal/delivery/http/v1/handlers"
	"service/internal/delivery/http/v1/middlewares"
	"service/internal/infrastructure/db/postgres"
	"service/internal/pkg/config"
	"service/internal/pkg/logs"
	"service/internal/pkg/pgclient"
	"service/internal/usecase"
)

const (
	op = "cmd.service"
)

func main() {
	file, w, err := logs.Setup(context.Background(), "./files/logs/service/")
	if err != nil {
		return
	}
	defer file.Close()
	defer w.Flush()

	logs.Info(
		context.Background(),
		"App starting",
		op,
	)

	cfg, err := config.GetConfig(context.Background())
	if err != nil {
		logs.Error(
			context.Background(),
			err.Error(),
			op,
		)

		return
	}

	dbClientCfg := pgclient.Config{
		Username: cfg.Storage.User,
		Password: cfg.Storage.Password,
		Host:     "postgres",
		Port:     "5432",
		Database: cfg.Storage.Name,
	}

	client, err := pgclient.NewClient(context.Background(), 5, dbClientCfg)
	if err != nil {
		logs.Error(
			context.Background(),
			err.Error(),
			op,
		)

		return
	}
	defer client.Close()

	mainMux := http.NewServeMux()

	notesRepository := postgres.NotesRepository{
		Pool: client,
	}

	notesUseCase := usecase.NotesService{
		Repository: notesRepository,
	}

	notesMux := handlers.NotesHandler{
		UseCase: notesUseCase,
	}.GetMux()

	mainMux.Handle("/", notesMux)

	wrappedMainMux := middleware.Cors(mainMux)                     // 4
	wrappedMainMux = middleware.Logger(wrappedMainMux)             // 3
	wrappedMainMux = middleware.SetupTrace(wrappedMainMux)         // 2
	wrappedMainMux = middleware.SetupContextValues(wrappedMainMux) // 1

	server := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           wrappedMainMux,
		ReadTimeout:       cfg.Server.ReadTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	}

	run(context.Background(), server)

	logs.Info(
		context.Background(),
		"App closed",
		op,
	)
}

func run(ctx context.Context, server *http.Server) {
	var wg sync.WaitGroup
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	wg.Add(1)
	go func() {
		logs.Info(
			ctx,
			"HTTP server started",
			op,
		)
		err := server.ListenAndServe()
		if err == http.ErrServerClosed {
			logs.Info(
				ctx,
				"HTTP server closed",
				op,
			)
			wg.Done()
			return
		}

		logs.Error(
			ctx,
			"HTTP server closed with error: "+err.Error(),
			op,
		)
	}()

	wg.Add(1)
	go func() {
		<-ctx.Done()

		err := server.Shutdown(ctx)
		if err != nil {
			logs.Error(
				ctx,
				"HTTP server shutdown error: "+err.Error(),
				op,
			)
		}

		wg.Done()
	}()

	wg.Wait()
}
