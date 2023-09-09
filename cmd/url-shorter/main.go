package main

import (
	"context"
	"fmt"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/image"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/url/save"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/user/create_user"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/users"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/users/create"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/handlers/slogpretty"
	"github.com/AlexandrLitkevich/qwery/internal/migration"
	"github.com/AlexandrLitkevich/qwery/internal/storage/store"
	"net/http"
	"time"

	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/AlexandrLitkevich/qwery/internal/config"
	mwLogger "github.com/AlexandrLitkevich/qwery/internal/http-server/middleware/logger"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/sl"
	"github.com/AlexandrLitkevich/qwery/internal/storage/sqlite"
	"github.com/joho/godotenv"
)

// TODO:sent from env
const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// NOTE:выгружаем в переменную среду приложения env
	envErr := godotenv.Load("../../.env")
	if envErr != nil {
		fmt.Println("error load env")
	}
	//Load config.yaml
	cfg := config.MustLoad()
	//implement logger
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))
	log.Info("Load env == true")

	ctx := context.Background()

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		os.Exit(1)
	}

	etcdStore, err := store.New(ctx, log, cfg)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	log.Info("Started migratioin")

	cmp, ops, err := migration.UpgradeTo0001()
	if err != nil {
		os.Exit(1)
	}
	if etcdStore.Cli != nil {
		defer etcdStore.Cli.Close()
	}

	log.Info("finished  migration")
	//На любой оперции он зависает
	txn := etcdStore.Cli.Txn(ctx)
	txn.If(cmp...)
	txn.Then(ops...)
	log.Info("migration precommt")

	resp, err := txn.Commit()
	log.Info("migration commt")

	log.Info("resp", resp)
	if err != nil {

		log.Error("fail to commit: %w", err)
	}
	if !resp.Succeeded {
		log.Error("fail to migrate schema: something changes in storage while migration")
	}
	log.Info("migration done")

	log.Info("storages is created succesfly")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer) // Поднимаем приложение при панике
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Post("/user", create_user.New(log, storage)) // Тут прям магия))))
	router.Post("/image", image.New(log, storage))
	router.Get("/user_etcd/{userId}", users.New(ctx, log, etcdStore))
	router.Post("/user_etcd", create.New(ctx, log, etcdStore))

	//TODO Get request home page

	log.Info("started server", slog.String("address", cfg.Address))
	//ACTION: init server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Timeout,
	}

	if err := srv.ListenAndServe(); err != nil { //Blocked function
		log.Error("failed to start server")

	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	// TODO: close storage

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	//NOTE: Для каждой площадки свои логи

	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
