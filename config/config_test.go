package config

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestListenAddr(t *testing.T) {
	cases := []struct{ ip, port, want string }{
		{"", "", ":3000"},
		{"", "8080", ":8080"},
		{"", ":9000", ":9000"},
		{" ", " 8081 ", ":8081"},
		{"127.1.2.3", "8100", "127.1.2.3:8100"},
		{"::1", "8100", "[::1]:8100"},
		{"[::1]", "8100", "[::1]:8100"},
		{"2001:db8::5", "", "[2001:db8::5]:3000"},
	}
	for _, c := range cases {
		if got := listenAddr(c.ip, c.port); got != c.want {
			t.Errorf("listenAddr(%q, %q) = %q, want %q", c.ip, c.port, got, c.want)
		}
	}
}

func TestSetAllowedOrigins(t *testing.T) {
	t.Cleanup(func() { SetAllowedOrigins(defaultOrigins) })

	SetAllowedOrigins([]string{" https://a.example ", "", "https://b.example"})
	if !isAllowedOrigin("https://a.example") || !isAllowedOrigin("https://b.example") {
		t.Fatal("origin dari daftar baru harus diizinkan (spasi dibuang)")
	}
	if isAllowedOrigin("https://jscroot.if.co.id") {
		t.Fatal("origin default harus terganti oleh daftar baru")
	}
	if isAllowedOrigin("") {
		t.Fatal("string kosong tidak boleh menjadi origin")
	}
}

func TestCORSPreflightAndPlainRequest(t *testing.T) {
	const origin = "https://jscroot.if.co.id"

	// preflight dari origin yang diizinkan dijawab 204
	req := httptest.NewRequest(http.MethodOptions, "/login", nil)
	req.Header.Set("Origin", origin)
	rec := httptest.NewRecorder()
	if !SetAccessControlHeaders(rec, req) {
		t.Fatal("preflight dari origin yang diizinkan harus dijawab")
	}
	if rec.Code != http.StatusNoContent || rec.Header().Get("Access-Control-Allow-Origin") != origin {
		t.Fatalf("preflight: code=%d allow-origin=%q", rec.Code, rec.Header().Get("Access-Control-Allow-Origin"))
	}

	// origin yang tidak diizinkan: tidak ada header CORS, request diteruskan
	req = httptest.NewRequest(http.MethodOptions, "/login", nil)
	req.Header.Set("Origin", "https://evil.example")
	rec = httptest.NewRecorder()
	if SetAccessControlHeaders(rec, req) || rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("origin yang tidak diizinkan tidak boleh mendapat header CORS")
	}
	if rec.Header().Get("Vary") != "Origin" {
		t.Fatal("Vary: Origin tetap harus dikirim untuk origin yang tidak diizinkan")
	}
}

// Dijalankan dengan -race untuk memastikan perubahan origin aman saat request berjalan.
func TestOriginsConcurrentAccess(t *testing.T) {
	t.Cleanup(func() { SetAllowedOrigins(defaultOrigins) })
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); AddAllowedOrigin("https://x.example") }()
		go func() { defer wg.Done(); _ = isAllowedOrigin("https://x.example") }()
	}
	wg.Wait()
}
