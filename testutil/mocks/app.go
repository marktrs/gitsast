package testutil

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/marktrs/gitsast/app"
)

type TestApp struct {
	*app.App
}

func StartTestApp(ctx context.Context) *TestApp {
	_, app, err := app.Start(ctx, "test", "test")
	if err != nil {
		panic(err)
	}

	mock := clock.NewMock()
	mock.Set(time.Date(2020, time.January, 1, 2, 3, 4, 5000, time.UTC))
	app.SetClock(mock)

	return &TestApp{
		App: app,
	}
}
