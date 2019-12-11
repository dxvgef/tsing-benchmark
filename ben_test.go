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
	httprouterApp        http.Handler
	httprouterRecoverApp http.Handler
	tsingApp             http.Handler
	tsingRecoverApp      http.Handler
	echoApp              http.Handler
	echoRecoverApp       http.Handler
	ginApp               http.Handler
	ginRecoverApp        http.Handler
)

// nolint
func init() {
	path, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

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
	calcMem("httprouterRecover", func() {
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

	// --------------------- tsing -----------------------------
	var tsingHandler = func(ctx *tsing.Context) error {
		ctx.ResponseWriter.WriteHeader(204)
		return nil
	}
	calcMem("tsing", func() {
		config := tsing.Config{
			RootPath: path,
		}
		app := tsing.New(&config)
		for _, route := range githubAPI {
			app.Router.Handle(route.Method, route.Path, tsingHandler)
		}
		tsingApp = app
	})
	calcMem("tsingRecover", func() {
		config := tsing.Config{
			RootPath:     path,
			Recover:      true,
			EventHandler: func(event *tsing.Event) {},
		}
		app := tsing.New(&config)
		for _, route := range githubAPI {
			app.Router.Handle(route.Method, route.Path, tsingHandler)
		}
		tsingRecoverApp = app
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
	calcMem("echoRecover", func() {
		app := echo.New()
		// 启用echo的recover
		app.Use(middleware.Recover())
		for _, route := range githubAPI {
			app.Add(route.Method, route.Path, echoHandler)
		}
		echoRecoverApp = app
	})

	// --------------------- echo -----------------------------
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
	calcMem("ginRecover", func() {
		gin.SetMode(gin.ReleaseMode)
		app := gin.New()
		// 启用gin的recover
		app.Use(gin.Recovery())
		for _, route := range githubAPI {
			app.Handle(route.Method, route.Path, ginHandler)
		}
		ginRecoverApp = app
	})

}

func BenchmarkHttprouter(b *testing.B) {
	benchRoutes(b, httprouterApp, githubAPI)
}
func BenchmarkHttprouterRecover(b *testing.B) {
	benchRoutes(b, httprouterRecoverApp, githubAPI)
}

func BenchmarkTsing(b *testing.B) {
	benchRoutes(b, tsingApp, githubAPI)
}
func BenchmarkTsingRecover(b *testing.B) {
	benchRoutes(b, tsingRecoverApp, githubAPI)
}

func BenchmarkEcho(b *testing.B) {
	benchRoutes(b, echoApp, githubAPI)
}
func BenchmarkEchoRecover(b *testing.B) {
	benchRoutes(b, echoRecoverApp, githubAPI)
}

func BenchmarkGin(b *testing.B) {
	benchRoutes(b, ginApp, githubAPI)
}
func BenchmarkGinRecover(b *testing.B) {
	benchRoutes(b, ginRecoverApp, githubAPI)
}
