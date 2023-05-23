package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"urlshortener/internal/db"
	"urlshortener/internal/routes"

	"github.com/julienschmidt/httprouter"
)

func TestGetSourceUrl(t *testing.T) {
	// seed with initial data
	db.GetConnection()
	db.InsertUrl("c423d", "https://vk.com/im")

	testCases := []struct {
		caseName   string
		value      string
		res        string
		statusCode int
	}{
		{
			caseName:   "Empty path",
			value:      "",
			res:        "",
			statusCode: http.StatusNotFound,
		},
		{
			caseName:   "Non-existent path",
			value:      "123456",
			res:        "",
			statusCode: http.StatusNotFound,
		},
		{
			caseName:   "Ok",
			value:      "c423d",
			res:        "",
			statusCode: http.StatusSeeOther,
		},
	}
	for _, tcase := range testCases {
		t.Run(tcase.caseName, func(t *testing.T) {
			var params httprouter.Params
			request := httptest.NewRequest("GET", "/", nil)
			responseRecorder := httptest.NewRecorder()
			params = []httprouter.Param{httprouter.Param{Key: "hash", Value: tcase.value}}
			routes.GetSourceUrl(responseRecorder, request, params)

			if responseRecorder.Result().StatusCode != tcase.statusCode {
				t.Error(request.URL.Path, request.Method)
				t.Errorf("Wanted status code '%d', but got '%d'", tcase.statusCode, responseRecorder.Result().StatusCode)
			}
			if responseRecorder.Body.String() != tcase.res {
				if tcase.caseName != "Ok" {
					// skips Redirect case, because of referrent body
					t.Errorf("Wanted response body '%s', but got '%s'", tcase.res, responseRecorder.Body.String())
				}
			}
		})
	}
}
