package app

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	http.Handler
	log Logger
}

type Config struct {
	Logger Logger
}

func NewHandler(config *Config) http.Handler {
	app := &App{}

	if config.Logger != nil {
		app.log = config.Logger
	} else {
		app.log = noopLogger{}
	}

	router := mux.NewRouter()

	box := packr.New("app", "./build")
	router.Use(app.createLoggingMiddleware(app.log.Debugf))
	router.PathPrefix("/").Handler(app.handleStatic(box)).Methods(http.MethodGet)

	app.Handler = router

	return app
}

func (a *App) createLoggingMiddleware(log func(string, ...interface{})) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log("accessing %v", r.RequestURI)
			next.ServeHTTP(w, r)
		})
	}
}

func (a *App) handleStatic(box *packr.Box) http.Handler {
	return http.FileServer(box)
}
