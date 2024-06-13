package server

import (
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tranhuyducseven/goem-template/config"
	"github.com/tranhuyducseven/goem-template/pkg/openobserve"
	"github.com/tranhuyducseven/goem-template/pkg/utils"
)

type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	Middlewares []string
}

var RouteRecs map[string][]RouteInfo = make(map[string][]RouteInfo)

const trimModule = "github.com/tranhuyducseven/goem-template/internal/"

func setupAddHandlerEvent(e *echo.Echo) {
	e.OnAddRouteHandler = func(host string, route echo.Route, handler echo.HandlerFunc, middleware []echo.MiddlewareFunc) {
		routeName := route.Name
		if routeName == "" {
			return
		}

		if !strings.Contains(routeName, trimModule) {
			return
		}

		groupAndHandler := strings.TrimPrefix(routeName, trimModule)
		var g, h string
		if strings.Contains(groupAndHandler, "/") {
			groupAndHandlerSplit := strings.Split(groupAndHandler, "/")
			g, h = groupAndHandlerSplit[0], groupAndHandlerSplit[1]
		} else {
			h = groupAndHandler
			g = "_empty"
		}

		handlerName := strings.Split(h, ".")
		if len(handlerName) == 2 {
			h = handlerName[1]
		}

		var middlewareNames []string
		for _, m := range middleware {
			mdw := utils.GetFunctionName(m)
			middlewareNames = append(middlewareNames, strings.TrimPrefix(mdw, trimModule+"middleware."))
		}

		RouteRecs[g] = append(RouteRecs[g], RouteInfo{
			Method:      route.Method,
			Path:        route.Path,
			Handler:     h,
			Middlewares: middlewareNames,
		})
	}
}

func setupMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.Env.CORS_WHITE_LIST,
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(openobserve.Middleware())
	e.Use(utils.RequestLogMiddleware())
}

func setupErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = utils.HttpErrorHandler
}

func setupValidator(e *echo.Echo) {
	e.Validator = utils.NewValidator()
}

func (s *Server) printRoutes() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Module", "Method", "Path", "Handler", "Middlewares"})

	routes := make([]string, 0, len(RouteRecs))
	for module := range RouteRecs {
		routes = append(routes, module)
	}

	sort.Strings(routes)

	for _, module := range routes {
		routeInfos := RouteRecs[module]
		for i, routeInfo := range routeInfos {
			var m string = ""
			if i == 0 && module != "_empty" {
				m = module
			}

			coloredMethod := utils.ColorMethod(routeInfo.Method)
			t.AppendRow(table.Row{m, coloredMethod, routeInfo.Path, routeInfo.Handler, strings.Join(routeInfo.Middlewares, " -> ")})

		}
		t.AppendSeparator()
	}

	t.Render()
}
