package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/j4ng5y/contrived-example-public/pkg/storage"
	_ "github.com/mattn/go-sqlite3"
)

type (
	Option func(*App)

	Config struct {
		HTTP struct {
			Addr              string
			Port              uint
			ReadTimeout       time.Duration
			ReadHeaderTimeout time.Duration
			WriteTimeout      time.Duration
			IdleTimeout       time.Duration
		}
	}

	App struct {
		logger  *slog.Logger
		cfg     *Config
		storage storage.Repo
	}
)

func WithHTTPAddr(addr string) Option {
	return func(app *App) {
		app.cfg.HTTP.Addr = addr
	}
}

func WithHTTPPort(port uint64) Option {
	return func(app *App) {
		app.cfg.HTTP.Port = uint(port)
	}
}

func WithHTTPReadTimeout(d time.Duration) Option {
	return func(app *App) {
		app.cfg.HTTP.ReadTimeout = d
	}
}

func WithHTTPReadHeaderTimeout(d time.Duration) Option {
	return func(app *App) {
		app.cfg.HTTP.ReadHeaderTimeout = d
	}
}

func WithHTTPWriteTimeout(d time.Duration) Option {
	return func(app *App) {
		app.cfg.HTTP.WriteTimeout = d
	}
}

func WithHTTPIdleTimeout(d time.Duration) Option {
	return func(app *App) {
		app.cfg.HTTP.IdleTimeout = d
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(app *App) {
		app.logger = logger
	}
}

func WithMemoryRepo() Option {
	return func(app *App) {
		stg, err := storage.NewRepo(storage.RepoType_Memory)
		if err != nil {
			app.logger.Error("unable to initialize storage repository due to error", "error", err)
			os.Exit(1)
		}
		app.storage = stg
	}
}

func WithSqliteRepo() Option {
	return func(app *App) {
		stg, err := storage.NewRepo(storage.RepoType_Sqlite)
		if err != nil {
			app.logger.Error("unable to initialize storage repository due to error", "error", err)
		}
		app.storage = stg
	}
}

func New(opts ...Option) *App {
	app := &App{
		logger: nil,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (app *App) Run() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, os.Kill, syscall.SIGTSTP)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", app.cfg.HTTP.Addr, app.cfg.HTTP.Port),
		ReadTimeout:       app.cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: app.cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      app.cfg.HTTP.WriteTimeout,
		IdleTimeout:       app.cfg.HTTP.IdleTimeout,
		Handler:           app.Mux(),
	}

	app.logger.Info("Running HTTP Server", "addr", app.cfg.HTTP.Addr, "port", app.cfg.HTTP.Port)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				os.Exit(0)
			} else {
				app.logger.Error("error stopping server", "err", err)
				os.Exit(1)
			}
		}
	}()

	sig := <-done

	app.logger.Info("Stopping HTTP Server due to signal", "signal", sig.String())

	if err := httpServer.Shutdown(ctx); err != nil {
		app.logger.Error("Error gracefully stopping the HTTP server", "error", err)
		os.Exit(1)
	}

	cancel()
	return nil
}
