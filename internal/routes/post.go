package routes

import (
	"net/http"

	"urlshortener/internal/db"
	"urlshortener/internal/utils"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// CreateShortenedUrl takes url query param as a Original url.
// Then hashes it, checks in tarantool for existence. If exist and existed original url equals requested - returns it.
// Otherwise, save this pair in tarantool. If collision exists - raises 500 InternalServerError.
func CreateShortenedUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.URL.RawQuery == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !r.URL.Query().Has("source") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalUrl := r.URL.Query().Get("source")
	hashedUrl := utils.MakeHashFromAdress(originalUrl)

	existedUrl, err := db.GetUrl(hashedUrl)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if existedUrl == originalUrl {
		w.Write([]byte(hashedUrl))
		return
	}
	_, err = db.InsertUrl(hashedUrl, originalUrl)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(hashedUrl))
}
