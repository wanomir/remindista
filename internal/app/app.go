package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	httpv1 "github.com/vedomirr/remindista/internal/controller/http_v1"
	"github.com/vedomirr/remindista/internal/infrastructure/repository"
	tg "github.com/vedomirr/remindista/internal/service/telegram"
	"github.com/vedomirr/remindista/internal/service/updater"
	"github.com/vedomirr/remindista/internal/service/worker"

	"go.uber.org/zap"

	"github.com/vedomirr/remindista/pkg/psql"

	"github.com/vedomirr/d"
	"github.com/vedomirr/e"
	"github.com/vedomirr/l"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	exitStatusOk     = 0
	exitStatusFailed = 1
)

type service interface {
	Run(ctx context.Context)
	Stop()
}

type App struct {
	config *Config
	logger *zap.Logger

	ctx     context.Context
	errChan chan error

	server *http.Server
	http   *httpv1.HttpController

	updater service
	worker  service
}

func NewApp() (*App, error) {
	a := new(App)

	if err := a.init(); err != nil {
		return nil, e.Wrap("failed to init app", err)
	}

	return a, nil
}

func (a *App) Run() (exitCode int) {
	defer a.recoverFromPanic(&exitCode)
	var err error

	ctx, stop := signal.NotifyContext(a.ctx, os.Interrupt, os.Kill)
	defer stop()

	// run updater service
	go a.updater.Run(ctx)

	// run worker service
	go a.worker.Run(ctx)

	// listen for incoming api requests
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errChan <- err
		}
	}()

	go d.Run(a.config.Debug.ServerAddr)

	select {
	case err = <-a.errChan:
		a.logger.Error(e.Wrap("fatal error, service shutdown", err).Error())
		exitCode = exitStatusFailed
	case <-ctx.Done():
		a.logger.Info("service shutdown")
	}

	return exitStatusOk
}

func (a *App) init() (err error) {
	// config
	if err = a.readConfig(); err != nil {
		return e.Wrap("failed to read config", err)
	}

	a.ctx = context.Background()
	a.errChan = make(chan error)

	l.BuildLogger(a.config.Log.Level)
	a.logger = l.Logger()

	// database
	pool, err := psql.Connect(a.ctx,
		psql.WithHost(a.config.PG.Host),
		psql.WithDatabase(a.config.PG.Database),
		psql.WithUser(a.config.PG.User),
		psql.WithPassword(a.config.PG.Password),
		psql.WithUserAdmin(a.config.PG.UserAdmin),
		psql.WithPasswordAdmin(a.config.PG.PasswordAdmin),
		psql.WithMigrations(os.DirFS("db/migrations")),
		psql.WithLogger(a.logger),
	)
	if err != nil {
		return e.Wrap("failed to init db", err)
	}

	// telegram entity
	telegram, err := tg.NewTelegram(a.config.TG.Token)
	if err != nil {
		return e.Wrap("failed to init telegram", err)
	}
	repo := repository.NewPostgresDB(pool)

	a.updater = updater.NewUpdater(telegram, repo)
	a.worker = worker.NewWorker(telegram, repo, a.config.Worker.Interval)

	// http server
	a.server = &http.Server{
		Addr:         a.config.Target.Addr,
		Handler:      a.routes(),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return nil
}

func (a *App) readConfig() (err error) {
	a.config = new(Config)
	if err = cleanenv.ReadEnv(a.config); err != nil {
		return err
	}

	return nil
}

func (a *App) recoverFromPanic(exitCode *int) {
	if panicErr := recover(); panicErr != nil {
		a.logger.Error(fmt.Sprintf("recover from panic: %v, stacktrace: %s", panicErr, string(debug.Stack())))
		*exitCode = exitStatusFailed
	}
}
