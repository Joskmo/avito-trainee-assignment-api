package main

import (
	"log/slog"
	"net/http"
	"time"

	repo "github.com/Joskmo/avito-trainee-assignment-api/internal/storage/postgres/sqlc"
	"github.com/Joskmo/avito-trainee-assignment-api/internal/teams"
	"github.com/Joskmo/avito-trainee-assignment-api/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config config
	db     *pgxpool.Pool
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Add middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))

	// handlers
	// for healthcheck
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// for teams
	teamsService := teams.NewService(repo.New(app.db), app.db)
	teamsHandler := teams.NewHandler(teamsService)
	r.Get("/team/get", teamsHandler.GetTeamByName)
	r.Post("/team/add", teamsHandler.CreateTeam)

	// for users
	usersService := users.NewService(repo.New(app.db), app.db)
	usersHandler := users.NewHandler(usersService)
	r.Post("/users/setIsActive", usersHandler.SetUserActivity)

	// for PRs
	

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	slog.Info("server has started", "address", app.config.addr)

	return srv.ListenAndServe()
}
