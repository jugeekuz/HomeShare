package app_test

import (
	"time"
	"net/http"
	"net/http/httptest"
	"testing"

	"file-server/internal/app"
	"file-server/internal/job"
)

func TestPreflightRequest(t *testing.T) {
	job_timeout := 45 * time.Second
	jm := job.NewJobManager(job_timeout)

	server := app.SetupServer(jm)
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	t.Run("OPTIONS response check", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/upload", nil)
		req.Header.Set("Origin", "http://kuza.gr")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rr := httptest.NewRecorder()

		server.Handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", rr.Code)
		}
	})
}

func TestHeadersPost(t *testing.T) {
	expectedOrigin := "http://kuza.gr"

	job_timeout := 45 * time.Second
	jm := job.NewJobManager(job_timeout)

	server := app.SetupServer(jm)
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	t.Run("POST headers check", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/upload", nil)
		req.Header.Set("Origin", "http://kuza.gr")
		resp := httptest.NewRecorder()

		server.Handler.ServeHTTP(resp, req)

		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != expectedOrigin {
			t.Errorf("Expected Access-Control-Allow-Origin header to be %q, got %q", expectedOrigin, got)
		}
	})

}
func TestNonExistentEndpoints(t *testing.T) {
	job_timeout := 45 * time.Second
	jm := job.NewJobManager(job_timeout)

	server := app.SetupServer(jm)
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	t.Run("GET root", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("POST invalid path", func(t *testing.T) {
		resp, err := http.Post(ts.URL+"/invalid", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}
