package hitcounter

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/writeas/impart"
	"github.com/writeas/web-core/log"
	"net/http"
	"runtime/debug"
	"time"
)

// Initialize and start web server
func Serve() error {
	r := mux.NewRouter()
	initRoutes(r)

	serve := fmt.Sprintf("%s:%d", "", 6767)
	log.Info("Serving on %s", serve)
	err := http.ListenAndServe(serve, r)
	return err
}

func initRoutes(r *mux.Router) {
	r.HandleFunc("/hits", handleHandler(handleViewHits))
	r.HandleFunc("/hit.gif", handleHandler(handlePixel))
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func handleHandler(f handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var status int
		start := time.Now()

		defer func() {
			if e := recover(); e != nil {
				log.Error("%s: %s", e, debug.Stack())
				status = http.StatusInternalServerError
			}

			log.Info("\"%s %s\" %d %s \"%s\"", r.Method, r.RequestURI, status, time.Since(start), r.UserAgent())
		}()

		err := f(w, r)
		if err != nil {
			if err, ok := err.(impart.HTTPError); ok {
				status = err.Status
				if status >= 300 && status < 400 {
					w.Header().Set("Location", err.Message)
				}
				w.WriteHeader(status)
			} else {
				status = http.StatusInternalServerError
			}
		}
		status = http.StatusOK
	}
}
