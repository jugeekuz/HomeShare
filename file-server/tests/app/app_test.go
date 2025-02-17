package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"file-server/internal/app"
)


func TestPreflightRequest(t *testing.T) {
	server := app.SetupServer()
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	t.Run("OPTIONS response check", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/uploads", nil)
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

	server := app.SetupServer()
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	t.Run("POST headers check", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/uploads", nil)
		req.Header.Set("Origin", "http://kuza.gr")
		resp := httptest.NewRecorder()

		server.Handler.ServeHTTP(resp, req)

		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != expectedOrigin {
			t.Errorf("Expected Access-Control-Allow-Origin header to be %q, got %q", expectedOrigin, got)
		}
	})

}
func TestNonExistentEndpoints(t *testing.T) {
	server := app.SetupServer()
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