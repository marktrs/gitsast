package app

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/benbjohnson/clock"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"github.com/marktrs/gitsast/internal/queue"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bunrouter"
	"github.com/urfave/cli/v2"
)

type appCtxKey struct{}

type App struct {
	ctx context.Context
	cfg *AppConfig

	stopping uint32
	stopCh   chan struct{}

	onStop      appHooks
	onAfterStop appHooks

	clock clock.Clock

	queue queue.Handler

	router    *bunrouter.Router
	apiRouter *bunrouter.Group

	dbOnce sync.Once
	db     *bun.DB

	validator *validator.Validate
}

func AppFromContext(ctx context.Context) *App {
	return ctx.Value(appCtxKey{}).(*App)
}

func ContextWithApp(ctx context.Context, app *App) context.Context {
	ctx = context.WithValue(ctx, appCtxKey{}, app)
	return ctx
}

func New(ctx context.Context, cfg *AppConfig) *App {
	app := &App{
		cfg:       cfg,
		stopCh:    make(chan struct{}),
		clock:     clock.New(),
		validator: validator.New(),
	}

	app.ctx = ContextWithApp(ctx, app)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	app.initRouter()
	app.initQueue()

	return app
}

func StartFromCLI(c *cli.Context) (context.Context, *App, error) {
	return Start(c.Context, c.Command.Name, c.String("env"))
}

func (app *App) SetClock(clock clock.Clock) {
	app.clock = clock
}

func Start(ctx context.Context, service, envName string) (context.Context, *App, error) {
	cfg, err := LoadConfigFile(FS(), service, envName)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	return StartWithConfig(ctx, cfg)
}

func StartWithConfig(ctx context.Context, cfg *AppConfig) (context.Context, *App, error) {
	app := New(ctx, cfg)
	if err := onStart.Run(ctx, app); err != nil {
		return nil, nil, err
	}
	return app.Context(), app, nil
}

func WaitExitSignal() os.Signal {
	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return <-ch
}

func (app *App) DB() *bun.DB {
	app.dbOnce.Do(func() {
		db := bun.NewDB(sql.OpenDB(
			pgdriver.NewConnector(pgdriver.WithDSN(app.cfg.DB.DSN))),
			pgdialect.New(),
		)

		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithEnabled(app.IsDebug()),
			bundebug.FromEnv(""),
		))

		app.db = db
	})
	return app.db
}

func (app *App) Running() bool {
	return !app.Stopping()
}

func (app *App) Stopping() bool {
	return atomic.LoadUint32(&app.stopping) == 1
}

func (app *App) Stop() {
	_ = app.onStop.Run(app.ctx, app)
	_ = app.onAfterStop.Run(app.ctx, app)
}

func (app *App) OnStop(name string, fn HookFunc) {
	app.onStop.Add(newHook(name, fn))
}

func (app *App) OnAfterStop(name string, fn HookFunc) {
	app.onAfterStop.Add(newHook(name, fn))
}

func (app *App) initQueue() {
	app.queue = queue.NewHandler()
}

func (app *App) Context() context.Context {
	return app.ctx
}

func (app *App) Config() *AppConfig {
	return app.cfg
}

func (app *App) Clock() clock.Clock {
	return app.clock
}

func (app *App) Router() *bunrouter.Router {
	return app.router
}

func (app *App) Queue() queue.Handler {
	return app.queue
}

func (app *App) APIRouter() *bunrouter.Group {
	return app.apiRouter
}

func (app *App) IsDebug() bool {
	return app.cfg.Debug
}

func (app *App) Validator() *validator.Validate {
	return app.validator
}

func (app *App) SetQueue(q queue.Handler) {
	app.queue = q
}
