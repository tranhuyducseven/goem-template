package server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/tranhuyducseven/goem-template/config"
	"github.com/tranhuyducseven/goem-template/pkg/db"
	"github.com/tranhuyducseven/goem-template/pkg/openobserve"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	Raw           *echo.Echo
	traceProvider *sdktrace.TracerProvider
}

func New() *Server {
	config.LoadEnv()

	appName := config.Env.APP_NAME
	if config.Env.IsDev {
		appName = appName + "-dev"
	}
	openobserve.Init(openobserve.OpenObserveConfig{
		Endpoint:    config.Env.OPENOBSERVE_ENDPOINT,
		Credential:  config.Env.OPENOBSERVE_CREDENTIAL,
		ServiceName: appName,
		Env:         config.Env.ENV,
	})

	config.InitLogger()

	e := echo.New()
	e.HideBanner = true
	tp := openobserve.SetupTraceHTTP()

	setupAddHandlerEvent(e)
	setupMiddleware(e)
	setupErrorHandler(e)
	setupRoute(e)
	setupValidator(e)

	return &Server{e, tp}
}

func (s *Server) Start(addr string) error {
	loadSvcs()
	s.printRoutes()
	return s.Raw.Start(addr)
}

func (s *Server) Close() {
	closeSvcs()
	s.Raw.Close()
	err := s.traceProvider.Shutdown(context.Background())
	if err != nil {
		log.Err(err).Msg("Error shutting down trace provider")
	}
}

func loadSvcs() {
	db.Init()
}

func closeSvcs() {
	db.Close()
}
