package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"myapp/internal/lib/validator"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func InitContext(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	e.Validator = validator.SetCustomValidator()

	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, target, body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

func EncodeReqBody(t *testing.T, body interface{}) io.Reader {
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}
	return bytes.NewBuffer(data)
}

func DecodeRecBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	var data T
	err := json.Unmarshal(rec.Body.Bytes(), &data)
	if err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}
	return data
}
