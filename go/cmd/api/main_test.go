package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := newDatabase()
	if err != nil {
		t.Fatalf("newDatabase() error = %v", err)
	}

	return newRouter(db)
}

func TestPing(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if got, want := recorder.Body.String(), "{\"message\":\"pong\"}"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestGroguReturnsKnownQuote(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/grogu", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var body struct {
		Quote string `json:"quote"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if body.Quote == "" {
		t.Fatal("quote is empty")
	}

	allowedQuotes := make(map[string]struct{}, len(groguQuotes))
	for _, quote := range groguQuotes {
		allowedQuotes[quote] = struct{}{}
	}

	if _, ok := allowedQuotes[body.Quote]; !ok {
		t.Fatalf("quote = %q, want one of %v", body.Quote, groguQuotes)
	}
}

func TestGroguReturnsErrorWhenNoQuotesConfigured(t *testing.T) {
	originalQuotes := groguQuotes
	groguQuotes = nil
	t.Cleanup(func() {
		groguQuotes = originalQuotes
	})

	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/grogu", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}

	if got, want := recorder.Body.String(), "{\"error\":\"no grogu quotes configured\"}"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}
