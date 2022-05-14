package main

import (
	"github.com/asaskevich/EventBus"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/newpurr/stock/internal/conf"
	"github.com/newpurr/stock/internal/event"
	kratosutils2 "github.com/newpurr/stock/pkg/utils/kratosutils"
)

type appStarter struct {
	// Name is the name of the compiled software.
	name string
	// Version is the version of the compiled software.
	version string
	// flagconf is the config flag.
	flagconf string
	hostname string
}

func newAppStarter(name string, version string, flagconf string, hostname string) *appStarter {
	return &appStarter{name: name, version: version, flagconf: flagconf, hostname: hostname}
}

func (as appStarter) start() {
	var (
		app    *kratos.App
		bc     conf.Bootstrap
		logger log.Logger
	)
	/*--------------------------------------------------------------------------
	| read configuration
	|-------------------------------------------------------------------------- */
	loader, cleanup, err := kratosutils2.NewConfLoader(
		file.NewSource(as.flagconf),
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	if err = loader.Scan(&bc); err != nil {
		panic(err)
	}

	/*--------------------------------------------------------------------------
	| init global logger
	|-------------------------------------------------------------------------- */
	logger = log.With(
		kratosutils2.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", as.hostname,
		"service.name", as.name,
		"service.version", as.version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	/*--------------------------------------------------------------------------
	| initialize the application，and start and wait for stop signal
	|-------------------------------------------------------------------------- */
	app, cleanup, err = initApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, bus EventBus.Bus) *kratos.App {
	err := event.LoadListener(bus)
	if err != nil {
		panic(err)
	}

	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs, gs,
		),
	)
}
