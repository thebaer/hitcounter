package hitcounter

import (
	"fmt"
	"github.com/writeas/impart"
	"github.com/writeas/web-core/bots"
	"github.com/writeas/web-core/log"
	"image"
	"image/color"
	"image/gif"
	"io/ioutil"
	"net/http"
	"time"
)

var counts = map[string]uint64{}

// CountView counts the view if it meets certain criteria, returning whether or not it was actually counted.
func CountView(r *http.Request) bool {
	if r.Method == "HEAD" || bots.IsBot(r.UserAgent()) {
		return false
	}

	path := r.FormValue("p")
	if path == "" {
		return false
	}

	// Count the hit
	counts[path]++
	SaveCounts()
	return true
}

func handlePixel(w http.ResponseWriter, r *http.Request) error {
	CountView(r)

	m := image.NewRGBA(image.Rect(0, 0, 1, 1))
	m.Set(0, 0, color.RGBA{255, 255, 255, 0})

	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Expires", "Wed, 11 Nov 2002 11:11:11 GMT")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")

	err := gif.Encode(w, m, nil)
	if err != nil {
		log.Error("Unable to encode gif: %s", err)
		return err
	}
	return nil
}

func handleView(w http.ResponseWriter, r *http.Request) error {
	CountView(r)

	return handleViewHits(w, r)
}

func handleViewHits(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	path := r.FormValue("p")
	fmt.Fprintf(w, "%d", counts[path])
	return nil
}

func handleHome(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		return impart.HTTPError{http.StatusNotFound, "Not found"}
	}
	data, err := ioutil.ReadFile("index.html")
	if err != nil {
		return impart.HTTPError{http.StatusNotFound, "Not found"}
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(data)
	return err
}
