package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

type RequestOption struct {
	Method      string
	URL         string
	Body        interface{}
	ContentType string
}

func Request(opts *RequestOption) (*http.Request, *httptest.ResponseRecorder) {
	var reqBody bytes.Buffer
	err := json.NewEncoder(&reqBody).Encode(opts.Body)
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(opts.Method, opts.URL, &reqBody)
	if opts.ContentType != "" {
		req.Header.Set(echo.HeaderContentType, opts.ContentType)
	} else {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	return req, rec
}
