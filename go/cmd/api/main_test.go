package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubShipService struct {
	ships []string
}

func (s stubShipService) ListShips() []string {
	return append([]string(nil), s.ships...)
}

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := newDatabase()
	if err != nil {
		t.Fatalf("newDatabase() error = %v", err)
	}

	router := newRouter(db, stubShipService{ships: []string{"Millennium Falcon"}})

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

func TestShips(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := newDatabase()
	if err != nil {
		t.Fatalf("newDatabase() error = %v", err)
	}

	router := newRouter(db, stubShipService{
		ships: []string{"Millennium Falcon", "X-wing", "TIE Fighter"},
	})

	req := httptest.NewRequest(http.MethodGet, "/ships", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response shipsResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if response.Count != 3 {
		t.Fatalf("count = %d, want %d", response.Count, 3)
	}

	wantShips := []string{"Millennium Falcon", "X-wing", "TIE Fighter"}
	if len(response.Ships) != len(wantShips) {
		t.Fatalf("ships length = %d, want %d", len(response.Ships), len(wantShips))
	}

	for i, want := range wantShips {
		if got := response.Ships[i]; got != want {
			t.Fatalf("ships[%d] = %q, want %q", i, got, want)
		}
	}

	if response.Message != "" {
		t.Fatalf("message = %q, want empty", response.Message)
	}
}

func TestShipsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := newDatabase()
	if err != nil {
		t.Fatalf("newDatabase() error = %v", err)
	}

	router := newRouter(db, stubShipService{})

	req := httptest.NewRequest(http.MethodGet, "/ships", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response shipsResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if response.Count != 0 {
		t.Fatalf("count = %d, want %d", response.Count, 0)
	}

	if len(response.Ships) != 0 {
		t.Fatalf("ships length = %d, want %d", len(response.Ships), 0)
	}

	if got, want := response.Message, "No ships available."; got != want {
		t.Fatalf("message = %q, want %q", got, want)
	}
}
