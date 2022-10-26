package hitcounter

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/writeas/impart"
	"github.com/writeas/web-core/log"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"
)

const (
	countsFilename = "counts.json"
)

// Initialize and start web server
func Serve() error {
	// Set up routes
	r := mux.NewRouter()
	initRoutes(r)

	// Load counts
	data, err := ioutil.ReadFile(countsFilename)
	if err != nil {
		log.Info("failed to ReadFile: %s", err)
		log.Info("All counts starting from 0")
	} else {
		err = json.Unmarshal(data, &counts)
		if err != nil {
			log.Error("failed to Unmarshal %s: %s", countsFilename, err)
			log.Error("All counts starting from 0")
		}
	}

	// Start server
	serve := fmt.Sprintf("%s:%d", "", 6767)
	log.Info("Serving on %s", serve)
	err = http.ListenAndServe(serve, r)
	return err
}

func SaveCounts() error {
	b, err := json.Marshal(counts)
	if err != nil {
		return fmt.Errorf("unable to marshal counts: %s", err)
	}
	err = ioutil.WriteFile(countsFilename, b, 0666)
	if err != nil {
		return fmt.Errorf("counts not saved: %s", err)
	}
	return nil
}

func initRoutes(r *mux.Router) {
	r.HandleFunc("/hits", handleHandler(handleViewHits))
	r.HandleFunc("/hit", handleHandler(handleView))
	r.HandleFunc("/hit.gif", handleHandler(handlePixel))
	r.HandleFunc("/", handleHandler(handleHome))
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
				if status == http.StatusNotFound && err.Message != "" {
					fmt.Fprintf(w, err.Message)
				}
			} else {
				status = http.StatusInternalServerError
			}
		} else {
			status = http.StatusOK
		}
	}
}
