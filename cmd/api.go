// Package main is the entry point of the application.
package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/pr"
	"github.com/Joskmo/avito-trainee-assignment-api/internal/stats"
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
	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	// for teams
	teamsService := teams.NewService(repo.New(app.db), app.db)
	teamsHandler := teams.NewHandler(teamsService)
	r.Get("/team/get", teamsHandler.GetTeamByName)
	r.Post("/team/add", teamsHandler.CreateTeam)
	r.Post("/team/deactivateUsers", teamsHandler.DeactivateUsers)

	// for users
	usersService := users.NewService(repo.New(app.db), app.db)
	usersHandler := users.NewHandler(usersService)
	r.Post("/users/setIsActive", usersHandler.SetUserActivity)

	// for PRs
	prService := pr.NewService(repo.New(app.db), app.db)
	prHandler := pr.NewHandler(prService)
	r.Post("/pullRequest/create", prHandler.CreatePR)
	r.Post("/pullRequest/merge", prHandler.MergePR)
	r.Post("/pullRequest/reassign", prHandler.ReassignReviewer)
	r.Get("/pullRequest/userReviews", prHandler.GetUserReviews)

	// for stats
	statsService := stats.NewService(repo.New(app.db), app.db)
	statsHandler := stats.NewHandler(statsService)
	r.Get("/stats", statsHandler.GetStats)

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
