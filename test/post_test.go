package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"urlshortener/internal/routes"
)

func TestCreateShortenedUrl(t *testing.T) {

	testCases := []struct {
		caseName   string
		path       string
		res        string
		statusCode int
	}{
		{
			caseName:   "Empty path",
			path:       "/",
			res:        "",
			statusCode: http.StatusMethodNotAllowed,
		},
		{
			caseName:   "Non-existent path",
			path:       "/?src=123",
			res:        "",
			statusCode: http.StatusBadRequest,
		},
		{
			caseName:   "Ok",
			path:       "/?source=https://vk.com/im",
			res:        "c423d",
			statusCode: http.StatusOK,
		},
	}
	for _, tcase := range testCases {
		t.Run(tcase.caseName, func(t *testing.T) {
			request := httptest.NewRequest("POST", tcase.path, nil)
			responseRecorder := httptest.NewRecorder()
			routes.CreateShortenedUrl(responseRecorder, request, nil)

			if responseRecorder.Result().StatusCode != tcase.statusCode {
				t.Errorf("Wanted status code '%d', but got '%d'", tcase.statusCode, responseRecorder.Result().StatusCode)
			}
			if responseRecorder.Body.String() != tcase.res {
				t.Errorf("Wanted response body '%s', but got '%s'", tcase.res, responseRecorder.Body.String())
			}
		})
	}
}
