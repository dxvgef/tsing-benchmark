package tsing_benchmark

import (
	"log"
	"net/http"
	"testing"

	"github.com/dimfeld/httptreemux/v5"
	tsingV1 "github.com/dxvgef/tsing"
	tsingV2 "github.com/dxvgef/tsing/v2"
	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	tsingV1App           http.Handler
	tsingV1RecoverApp    http.Handler
	tsingV2App           http.Handler
	tsingV2RecoverApp    http.Handler
	httprouterApp        http.Handler
	httprouterRecoverApp http.Handler
	echoApp              http.Handler
	echoRecoverApp       http.Handler
	ginApp               http.Handler
	ginRecoverApp        http.Handler
	treemaxApp           http.Handler
)

// nolint
func init() {
	log.SetFlags(log.Lshortfile)
	// path, err := os.Getwd()
	// if err != nil {
	// 	os.Exit(1)
	// }

	// --------------------- tsing v2 -----------------------------
	var tsingV2Handler = func(ctx *tsingV2.Context) error {
		ctx.ResponseWriter.WriteHeader(204)
		return nil
	}
	calcMem("tsing v2", func() {
		app := tsingV2.New()
		for _, route := range githubAPI {
			app.Handle(route.Method, route.Path, tsingV2Handler)
		}
		tsingV2App = app
	})
	// --------------------- tsing v2 recover -----------------------------
	calcMem("tsing v2 recover", func() {
		app := tsingV2.New(tsingV2.Config{
			Recovery: true,
		})
		for _, route := range githubAPI {
			app.Handle(route.Method, route.Path, tsingV2Handler)
		}
		tsingV2RecoverApp = app
	})

	// --------------------- tsing v1 -----------------------------
	var tsingV1Handler = func(ctx *tsingV1.Context) error {
		ctx.ResponseWriter.WriteHeader(204)
		return nil
	}
	calcMem("tsing v1", func() {
		app := tsingV1.New(&tsingV1.Config{})
		for _, route := range githubAPI {
			app.Router.Handle(route.Method, route.Path, tsingV1Handler)
		}
		tsingV1App = app
	})
	calcMem("tsing v1 recover", func() {
		app := tsingV1.New(&tsingV1.Config{
			Recover: true,
		})
		for _, route := range githubAPI {
			app.Router.Handle(route.Method, route.Path, tsingV1Handler)
		}
		tsingV1RecoverApp = app
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
		router.PanicHandler = func(writer http.ResponseWriter, request *http.Request, i interface{}) {}
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
		app.Use(middleware.Recover())
		for _, route := range githubAPI {
			app.Add(route.Method, route.Path, echoHandler)
		}
		echoRecoverApp = app
	})

	// --------------------- httptreemux -----------------------------
	var treemuxHandler = func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.WriteHeader(204)
	}
	calcMem("httptreemux", func() {
		app := httptreemux.New()
		for _, route := range githubAPI {
			app.Handle(route.Method, route.Path, treemuxHandler)
		}
		treemaxApp = app
	})

}

func Benchmark_TsingV2(b *testing.B) {
	benchRoutes(b, tsingV2App, githubAPI)
}
func Benchmark_TsingV2_Recover(b *testing.B) {
	benchRoutes(b, tsingV2RecoverApp, githubAPI)
}

func Benchmark_TsingV1(b *testing.B) {
	benchRoutes(b, tsingV1App, githubAPI)
}

func Benchmark_TsingV1_Recover(b *testing.B) {
	benchRoutes(b, tsingV1RecoverApp, githubAPI)
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

func Benchmark_HTTPTreemux(b *testing.B) {
	benchRoutes(b, treemaxApp, githubAPI)
}
