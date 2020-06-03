package tsing_benchmark

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/dxvgef/tsing"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	httprouterApp        http.Handler
	httprouterRecoverApp http.Handler
	tsingApp             http.Handler
	tsingRecoverApp      http.Handler
	echoApp              http.Handler
	echoRecoverApp       http.Handler
	ginApp               http.Handler
	ginRecoverApp        http.Handler
	chiApp               http.Handler
	chiRecoverApp        http.Handler
)

// nolint
func init() {
	log.SetFlags(log.Lshortfile)
	path, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	// --------------------- tsing -----------------------------
	var tsingHandler = func(ctx *tsing.Context) error {
		ctx.ResponseWriter.WriteHeader(204)
		return nil
	}
	calcMem("tsing", func() {
		app := tsing.New(&tsing.Config{})
		for _, route := range githubAPI {
			app.Router.Handle(route.Method, route.Path, tsingHandler)
		}
		tsingApp = app
	})
	calcMem("tsing recover", func() {
		config := tsing.Config{
			RootPath:           path,
			UnescapePathValues: true,
			EventTrace:         true,
			EventSource:        true,
			EventShortPath:     true,
			Recover:            true,
			EventHandler:       func(e *tsing.Event) {},
		}
		app := tsing.New(&config)
		for _, route := range githubAPI {
			app.Router.Handle(route.Method, route.Path, tsingHandler)
		}
		tsingRecoverApp = app
	})

	// --------------------- httprouter -----------------------------
	var httprouterHandler = func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.WriteHeader(204)
	}
	calcMem("httprouter", func() {
		router := httprouter.New()
		for _, route := range githubAPI {
			router.Handle(route.Method, route.Path, httprouterHandler)
		}
		httprouterApp = router
	})
	calcMem("httprouter recover", func() {
		router := httprouter.New()
		// 启用recover
		router.PanicHandler = func(writer http.ResponseWriter, request *http.Request, i interface{}) {}
		// router.NotFound = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {})
		// router.HandleMethodNotAllowed = true
		// router.MethodNotAllowed = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {})
		for _, route := range githubAPI {
			router.Handle(route.Method, route.Path, httprouterHandler)
		}
		httprouterRecoverApp = router
	})

	// --------------------- gin -----------------------------
	var ginHandler = func(ctx *gin.Context) {
		ctx.Writer.WriteHeader(204)
	}
	calcMem("gin", func() {
		gin.SetMode(gin.ReleaseMode)
		app := gin.New()
		for _, route := range githubAPI {
			app.Handle(route.Method, route.Path, ginHandler)
		}
		ginApp = app
	})
	calcMem("gin recover", func() {
		gin.SetMode(gin.ReleaseMode)
		app := gin.New()
		// 启用gin的recover
		app.Use(gin.Recovery())
		for _, route := range githubAPI {
			app.Handle(route.Method, route.Path, ginHandler)
		}
		ginRecoverApp = app
	})

	// --------------------- echo -----------------------------
	var echoHandler = func(ctx echo.Context) error {
		ctx.Response().WriteHeader(204)
		return nil
	}
	calcMem("echo", func() {
		app := echo.New()
		for _, route := range githubAPI {
			app.Add(route.Method, route.Path, echoHandler)
		}
		echoApp = app
	})
	calcMem("echo recover", func() {
		app := echo.New()
		// 启用echo的recover
		app.Use(middleware.Recover())
		for _, route := range githubAPI {
			app.Add(route.Method, route.Path, echoHandler)
		}
		echoRecoverApp = app
	})

	// --------------------- chi -----------------------------
	var chiHandler = func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(204)
	}
	calcMem("chi", func() {
		app := chi.NewRouter()
		for _, route := range githubAPI {
			app.MethodFunc(route.Method, route.Path, chiHandler)
		}
		chiApp = app
	})
	calcMem("chi recover", func() {
		app := chi.NewRouter()
		app.Use(chiMiddleware.Recoverer)
		for _, route := range githubAPI {
			app.MethodFunc(route.Method, route.Path, chiHandler)
		}
		chiRecoverApp = app
	})
}

func Benchmark_Tsing(b *testing.B) {
	benchRoutes(b, tsingApp, githubAPI)
}

func Benchmark_Tsing_Recover(b *testing.B) {
	benchRoutes(b, tsingRecoverApp, githubAPI)
}

func Benchmark_Httprouter(b *testing.B) {
	benchRoutes(b, httprouterApp, githubAPI)
}

func Benchmark_Httprouter_Recover(b *testing.B) {
	benchRoutes(b, httprouterRecoverApp, githubAPI)
}

func Benchmark_Gin(b *testing.B) {
	benchRoutes(b, ginApp, githubAPI)
}

func Benchmark_Gin_Recover(b *testing.B) {
	benchRoutes(b, ginRecoverApp, githubAPI)
}

func Benchmark_Echo(b *testing.B) {
	benchRoutes(b, echoApp, githubAPI)
}
func Benchmark_Echo_Recover(b *testing.B) {
	benchRoutes(b, echoRecoverApp, githubAPI)
}

func Benchmark_chi(b *testing.B) {
	benchRoutes(b, chiApp, githubAPI)
}
func Benchmark_chi_Recover(b *testing.B) {
	benchRoutes(b, chiRecoverApp, githubAPI)
}
