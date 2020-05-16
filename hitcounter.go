package hitcounter

import (
	"fmt"
	"github.com/writeas/impart"
	"github.com/writeas/web-core/log"
	"image"
	"image/color"
	"image/gif"
	"net/http"
	"time"
)

var counts = map[string]uint64{}

func handlePixel(w http.ResponseWriter, r *http.Request) error {
	path := r.FormValue("p")
	if path == "" {
		return impart.HTTPError{http.StatusNotFound, ""}
	}

	// Count the hit
	counts[path]++

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
		return nil
	}
	return nil
}

func handleViewHits(w http.ResponseWriter, r *http.Request) error {
	path := r.FormValue("p")
	fmt.Fprintf(w, "%d", counts[path])
	return nil
}
