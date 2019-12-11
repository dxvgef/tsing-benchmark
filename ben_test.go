package tsing_benchmark

import (
	"net/http"
	"os"
	"testing"

	"github.com/dxvgef/tsing"
	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	httprouterApp http.Handler
	tsingApp      http.Handler
	echoApp       http.Handler
	ginApp        http.Handler
)

// nolint
func init() {
	calcMem("HttpRouter", func() {
		router := httprouter.New()
		// 启用recover
		// router.PanicHandler = func(writer http.ResponseWriter, request *http.Request, i interface{}) {}

		// router.NotFound = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {})
		// router.HandleMethodNotAllowed = true
		// router.MethodNotAllowed = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {})
		handler := func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI {
			router.Handle(route.Method, route.Path, handler)
		}
		httprouterApp = router
	})

	calcMem("Echo", func() {
		app := echo.New()
		// 启用echo的recover
		app.Use(middleware.Recover())
		handler := func(ctx echo.Context) error {
			ctx.Response().WriteHeader(204)
			return nil
		}
		for _, route := range githubAPI {
			app.Add(route.Method, route.Path, handler)
		}
		echoApp = app
	})

	calcMem("Gin", func() {
		gin.SetMode(gin.ReleaseMode)
		app := gin.New()
		// 启用gin的recover
		app.Use(gin.Recovery())
		handler := func(ctx *gin.Context) {
			ctx.Writer.WriteHeader(204)
		}
		for _, route := range githubAPI {
			app.Handle(route.Method, route.Path, handler)
		}
		ginApp = app
	})

	calcMem("Tsing", func() {
		path, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		// 通常用的配置
		regular := tsing.Config{
			RootPath:              path,
			RedirectTrailingSlash: true,
			HandleOPTIONS:         true,
			FixPath:               true,
			Recover:               false,
			EventHandler: func(event *tsing.Event) {

			},
			ErrorEvent:            true,
			NotFoundEvent:         false,
			MethodNotAllowedEvent: false,
			Trigger:               false,
			Trace:                 false,
			ShortPath:             true,
		}
		// 最高性能配置
		// performance := tsing.Config{
		// 	RootPath:              path,
		// 	RedirectTrailingSlash: false,
		// 	HandleOPTIONS:         false,
		// 	FixPath:               false,
		// 	Recover:               false,
		// 	EventHandler:          nil,
		// 	ErrorEvent:            false,
		// 	NotFoundEvent:         false,
		// 	MethodNotAllowedEvent: false,
		// 	Trigger:               false,
		// 	Trace:                 false,
		// 	ShortPath:             false,
		// }
		app := tsing.New(&regular)
		handler := func(ctx *tsing.Context) error {
			ctx.ResponseWriter.WriteHeader(204)
			return nil
		}
		for _, route := range githubAPI {
			app.Router.Handle(route.Method, route.Path, handler)
		}
		tsingApp = app
	})

}

func BenchmarkTsing(b *testing.B) {
	benchRoutes(b, tsingApp, githubAPI)
}

func BenchmarkHttpRouter(b *testing.B) {
	benchRoutes(b, httprouterApp, githubAPI)
}

func BenchmarkEcho(b *testing.B) {
	benchRoutes(b, echoApp, githubAPI)
}

func BenchmarkGin(b *testing.B) {
	benchRoutes(b, ginApp, githubAPI)
}
