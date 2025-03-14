package app

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"file-server/config"
	"file-server/internal/job"
)

var endpoints = []string{
	"/login", // POST
	"/refresh", // POST
	"/logout", // POST
	"/upload", // POST
	"/download", // GET
	"/share", // POST
	"/share-file", // POST
	"/share-files", // GET
}


func setupTestServer(t *testing.T) (*httptest.Server, string) {
	cfg := config.LoadConfig()

	jm := job.NewJobManager(10*time.Minute)

	dummyInitDatabase := func () (*sql.DB, error){
		return nil, nil
	}

	srv, err := SetupServer(jm, dummyInitDatabase)
	if err != nil {
		t.Fatalf("failed to setup server: %v", err)
	}

	ts := httptest.NewServer(srv.Handler)
	return ts, cfg.DomainOrigin
}

func TestRoutes(t *testing.T) {
	ts, _ := setupTestServer(t)
	defer ts.Close()

	// Non Existent Endpoint
	res, err := http.Get(ts.URL + "/non-existent")
	if err != nil {
		t.Fatalf("failed to GET non-existent route: %v", err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent route, got %d", res.StatusCode)
	}

	// Existing Endpoints
	for _, route := range endpoints {
		t.Run("ExistingRoute"+route, func(t *testing.T) {
			res, err := http.Get(ts.URL + route)
			if err != nil {
				t.Fatalf("failed to GET %s: %v", route, err)
			}
			if res.StatusCode == http.StatusNotFound {
				t.Errorf("expected route %s not to return 404", route)
			}
		})
	}
}


func TestCors(t *testing.T) {
	ts, allowedOrigin := setupTestServer(t)
	defer ts.Close()

	// CORS success
	for _, route := range endpoints {
		t.Run("CorsSuccess_"+route, func(t *testing.T) {

			req, err := http.NewRequest("GET", ts.URL+route, nil)
			if err != nil {
				t.Fatalf("failed to create request for %s: %v", route, err)
			}

			req.Header.Set("Origin", allowedOrigin)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to GET %s: %v", route, err)
			}

			if got := res.Header.Get("Access-Control-Allow-Origin"); got != allowedOrigin {
				t.Errorf("expected Access-Control-Allow-Origin to be %q, got %q", allowedOrigin, got)
			}

		})
	}

	// Preflight Error
	invalidOrigin := "http://notallowed.com"
	for _, route := range endpoints {
		t.Run("CorsPreflightError_"+route, func(t *testing.T) {
			req, err := http.NewRequest("OPTIONS", ts.URL+route, nil)
			if err != nil {
				t.Fatalf("failed to create OPTIONS request for %s: %v", route, err)
			}
			req.Header.Set("Origin", invalidOrigin)
			req.Header.Set("Access-Control-Request-Method", "GET")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to do OPTIONS for %s: %v", route, err)
			}

			if got := res.Header.Get("Access-Control-Allow-Origin"); got != "" {
				t.Errorf("expected no Access-Control-Allow-Origin header for invalid origin, got %q", got)
			}
		})
	}

	// Preflight Sucess
	for _, route := range endpoints {
		t.Run("CorsPreflightSuccess_"+route, func(t *testing.T) {
			req, err := http.NewRequest("OPTIONS", ts.URL+route, nil)
			if err != nil {
				t.Fatalf("failed to create OPTIONS request for %s: %v", route, err)
			}
			req.Header.Set("Origin", allowedOrigin)
			req.Header.Set("Access-Control-Request-Method", "GET")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to perform OPTIONS for %s: %v", route, err)
			}

			if got := res.Header.Get("Access-Control-Allow-Origin"); got != allowedOrigin {
				t.Errorf("expected Access-Control-Allow-Origin header %q, got %q", allowedOrigin, got)
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	ts, allowedOrigin := setupTestServer(t)
	defer ts.Close()

	for _, route := range endpoints {
		t.Run("Headers_"+route, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+route, nil)
			if err != nil {
				t.Fatalf("failed to create GET request for %s: %v", route, err)
			}
			req.Header.Set("Origin", allowedOrigin)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to GET %s: %v", route, err)
			}

			if got := res.Header.Get("Access-Control-Allow-Origin"); got != allowedOrigin {
				t.Errorf("expected Access-Control-Allow-Origin header %q, got %q", allowedOrigin, got)
			}

		})
	}
}