package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/KirillLich/todoapi/internal/config"
	"github.com/KirillLich/todoapi/internal/todo"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	//TODO: fix issue with relativity of the path
	//TODO: default from file and custom from env variables and flags
	cfg := config.MustLoad("../configs/local.yaml")
	log := setupLogger(cfg)

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DB.User, cfg.DB.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.DB.Host, cfg.DB.Port),
		Path:   cfg.DB.Name,
	}
	q := u.Query()
	q.Set("sslmode", cfg.DB.SSLMode)
	u.RawQuery = q.Encode()

	dsn := u.String()

	pool, err := sql.Open("pgx", dsn)
	defer pool.Close()
	if err != nil {
		log.Error("error while opening db", slog.String("Error", err.Error()))
		os.Exit(1)
	}
	err = pool.Ping()
	if err != nil {
		log.Error("error while pinging db", slog.String("Error", err.Error()))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	h := todo.NewHandler(todo.NewService(todo.NewRepository(pool)), log)
	//how to insert logger into handler
	mux.HandleFunc("/todos", h.Todos)
	mux.HandleFunc("/todos/", h.TodosId)

	log.Debug(fmt.Sprintf("server starting on %d", cfg.Server.Port))
	err = http.ListenAndServe(cfg.Server.Address+":"+strconv.Itoa(cfg.Server.Port), mux)
	if err != nil {
		log.Error("error while starting the server", slog.String("Error", err.Error()))
		os.Exit(1)
	}
}

func setupLogger(cfg config.Config) *slog.Logger {
	var log *slog.Logger

	//TODO: make more cases
	switch cfg.Env {
	case "local":
		log = setupConsoleLogger(&slog.HandlerOptions{Level: slog.LevelDebug})
	}

	return log
}

func setupConsoleLogger(opts *slog.HandlerOptions) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, opts)

	return slog.New(handler)
}
