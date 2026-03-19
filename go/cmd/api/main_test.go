package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := newDatabase()
	if err != nil {
		t.Fatalf("newDatabase() error = %v", err)
	}

	router := newRouter(db)

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
