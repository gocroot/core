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
	if rec.Code != http.StatusNoContent {
		h := rec.Header()
		if h.Get("Content-Type") != "application/json" || h.Get("Cache-Control") != "no-store" || h.Get("X-Content-Type-Options") != "nosniff" {
			t.Errorf("%s %s: header keamanan kurang: %v", method, path, h)
		}
	}
	if !strings.Contains(strings.Join(rec.Header().Values("Vary"), ","), "Origin") {
		t.Errorf("%s %s: Vary: Origin harus selalu ada", method, path)
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
		{http.MethodGet, "/version", http.StatusOK},
		{http.MethodPost, "/version", http.StatusNotFound},
	}
	for _, c := range cases {
		if code, _ := do(t, c.method, c.path, "", nil); code != c.wantCode {
			t.Errorf("%s %s = %d, want %d", c.method, c.path, code, c.wantCode)
		}
	}
}

func TestVersion(t *testing.T) {
	old := config.Version
	t.Cleanup(func() { config.Version = old })
	config.Version = "abc123"
	code, resp := do(t, http.MethodGet, "/version", "", nil)
	if code != http.StatusOK || resp.Status != "success" || resp.Message != "abc123" {
		t.Errorf("GET /version = %d %+v, want 200 success abc123", code, resp)
	}
}

func TestLoginDisabledWithoutDemoEnv(t *testing.T) {
	config.DemoUser, config.DemoPassword = "", ""
	if code, _ := do(t, http.MethodPost, "/login", `{"username":"admin","password":"admin123"}`, nil); code != http.StatusServiceUnavailable {
		t.Errorf("tanpa DEMO_USER/DEMO_PASSWORD = %d, want 503", code)
	}
	config.DemoUser = "admin" // hanya satu yang terisi tetap mati
	if code, _ := do(t, http.MethodPost, "/login", `{"username":"admin","password":""}`, nil); code != http.StatusServiceUnavailable {
		t.Errorf("DEMO_PASSWORD kosong = %d, want 503", code)
	}
}

func TestLoginFlow(t *testing.T) {
	config.PrivateKey, config.PublicKey = watoken.GenerateKey()
	config.DemoUser, config.DemoPassword = "admin", "rahasia-demo"
	t.Cleanup(func() { config.DemoUser, config.DemoPassword = "", "" })

	if code, _ := do(t, http.MethodPost, "/login", "bukan json", nil); code != http.StatusBadRequest {
		t.Errorf("body rusak = %d, want 400", code)
	}
	if code, _ := do(t, http.MethodPost, "/login", strings.Repeat(" ", 2<<20)+"{}", nil); code != http.StatusBadRequest {
		t.Errorf("body > 1 MiB = %d, want 400", code)
	}
	for _, body := range []string{
		`{"username":"admin","password":"salah"}`,
		`{"username":"admin","password":"admin123"}`, // password lama tidak lagi berlaku
		`{"username":"bukan","password":"rahasia-demo"}`,
		`{"username":"admin","password":"rahasia-dem"}`, // prefix tidak cukup
	} {
		if code, _ := do(t, http.MethodPost, "/login", body, nil); code != http.StatusNotFound {
			t.Errorf("kredensial salah %s = %d, want 404", body, code)
		}
	}

	code, resp := do(t, http.MethodPost, "/login", `{"username":"admin","password":"rahasia-demo"}`, nil)
	if code != http.StatusOK || resp.Token == "" {
		t.Fatalf("login benar = %d, token=%q", code, resp.Token)
	}

	token := resp.Token
	code, resp = do(t, http.MethodGet, "/login", "", map[string]string{"Token": token})
	if code != http.StatusOK || resp.Username != "admin" {
		t.Errorf("cek token = %d, username=%q", code, resp.Username)
	}
	for _, tok := range []string{"token-palsu", "", token[:len(token)-4] + "AAAA"} {
		if code, _ := do(t, http.MethodGet, "/login", "", map[string]string{"Token": tok}); code != http.StatusUnauthorized {
			t.Errorf("token tidak valid %q = %d, want 401", tok, code)
		}
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
