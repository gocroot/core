package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/models"
	"github.com/whatsauth/watoken"
)

func do(t *testing.T, method, path, body string, header map[string]string) (int, models.Response) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	HandleRoutes(rec, req)
	var resp models.Response
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("%s %s: body bukan JSON valid: %v (%q)", method, path, err, rec.Body.String())
		}
	}
	if rec.Code != http.StatusNoContent && rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("%s %s: Content-Type = %q", method, path, rec.Header().Get("Content-Type"))
	}
	return rec.Code, resp
}

func TestRoutes(t *testing.T) {
	cases := []struct {
		method, path string
		wantCode     int
	}{
		{http.MethodGet, "/", http.StatusOK},
		{http.MethodGet, "/tidak-ada", http.StatusNotFound},
		{http.MethodPost, "/tidak-ada", http.StatusNotFound},
		{http.MethodPut, "/", http.StatusMethodNotAllowed},
		{http.MethodDelete, "/login", http.StatusMethodNotAllowed},
	}
	for _, c := range cases {
		if code, _ := do(t, c.method, c.path, "", nil); code != c.wantCode {
			t.Errorf("%s %s = %d, want %d", c.method, c.path, code, c.wantCode)
		}
	}
}

func TestLoginFlow(t *testing.T) {
	config.PrivateKey, config.PublicKey = watoken.GenerateKey()

	if code, _ := do(t, http.MethodPost, "/login", "bukan json", nil); code != http.StatusBadRequest {
		t.Errorf("body rusak = %d, want 400", code)
	}
	if code, _ := do(t, http.MethodPost, "/login", strings.Repeat(" ", 2<<20)+"{}", nil); code != http.StatusBadRequest {
		t.Errorf("body > 1 MiB = %d, want 400", code)
	}
	if code, _ := do(t, http.MethodPost, "/login", `{"username":"admin","password":"salah"}`, nil); code != http.StatusNotFound {
		t.Errorf("password salah = %d, want 404", code)
	}

	code, resp := do(t, http.MethodPost, "/login", `{"username":"admin","password":"admin123"}`, nil)
	if code != http.StatusOK || resp.Token == "" {
		t.Fatalf("login benar = %d, token=%q", code, resp.Token)
	}

	code, resp = do(t, http.MethodGet, "/login", "", map[string]string{"Token": resp.Token})
	if code != http.StatusOK || resp.Username != "admin" {
		t.Errorf("cek token = %d, username=%q", code, resp.Username)
	}
	if code, _ := do(t, http.MethodGet, "/login", "", map[string]string{"Token": "token-palsu"}); code != http.StatusBadRequest {
		t.Errorf("token palsu = %d, want 400", code)
	}
}

func TestPreflightViaRoutes(t *testing.T) {
	code, _ := do(t, http.MethodOptions, "/login", "", map[string]string{"Origin": "https://jscroot.if.co.id"})
	if code != http.StatusNoContent {
		t.Errorf("preflight = %d, want 204", code)
	}
	// OPTIONS tanpa origin yang diizinkan diperlakukan sebagai method biasa
	if code, _ := do(t, http.MethodOptions, "/login", "", nil); code != http.StatusMethodNotAllowed {
		t.Errorf("OPTIONS tanpa origin = %d, want 405", code)
	}
}
