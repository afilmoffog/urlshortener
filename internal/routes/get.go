package routes

import (
	"net/http"
	
	"urlshortener/internal/db"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// GetSourceUrl takes url query "hash" param as hashed original url.
// Then, checks it in tarantool. If exists - redirects to original url. Else - raises StatusNotFound
func GetSourceUrl(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shortedUrl := ps.ByName("hash")
	originalUrl, err := db.GetUrl(shortedUrl)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if originalUrl == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusSeeOther)

}
