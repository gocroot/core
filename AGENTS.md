# AGENTS.md

Aturan untuk AI coding agent (Antigravity, Claude Code, Codex, Cursor, Copilot, dsb.) yang bekerja di proyek hasil template [GoCroot Core](https://github.com/gocroot/core). Dokumentasi lengkap: https://gocroot.if.co.id/ (versi untuk agent: https://gocroot.if.co.id/llms.txt).

## Autentikasi: PASETO, bukan JWT

- Token login memakai **PASETO v4.public** lewat `github.com/whatsauth/watoken`. **Jangan memakai JWT** (`golang-jwt/jwt`, `lestrrat-go/jwx`, dsb.) dan jangan menulis kriptografi sendiri.
- Terbitkan token: `token, err := watoken.Encode(userID, config.PrivateKey)`. Selalu cek `err`, karena token kosong dikembalikan saat key tidak valid. Untuk data tambahan pakai `watoken.EncodeWithStruct`, untuk durasi lain `watoken.EncodeWithStructDuration`.
- Verifikasi token: `userID, err := watoken.DecodeGetId(config.PublicKey, r.Header.Get("Token"))`. Token dikirim lewat header **`Token`**, bukan `Authorization: Bearer`.
- Token tidak valid atau kedaluwarsa dibalas **`401`**.
- Key dari env `PRIVATEKEY` dan `PUBLICKEY`. Pasangan key dibuat dengan `watoken.GenerateKey()`. Service yang hanya memeriksa token cukup memegang `PUBLICKEY`.
- Alasannya: PASETO tidak punya header `alg` yang bisa dimanipulasi (celah `alg: none` dan kebingungan RS256/HS256 di JWT), dan kriptografinya terkunci per versi (Ed25519).
- Pengecualian: memverifikasi token dari penyedia luar yang memang menerbitkan JWT (Google ID token/OIDC, Firebase) boleh memakai library verifikasi mereka. Token sesi yang diterbitkan proyek ini tetap PASETO.

## HTTP dan routing

- Hanya `net/http` standar. **Jangan menambah framework web atau router** (Gin, Fiber, Echo, Chi, gorilla/mux).
- Routing ada di `routes/routes.go`: CORS dulu, lalu `switch r.Method`, lalu `switch r.URL.Path`. Path tidak dikenal ke `controller.NotFound`, method tidak dikenal ke `controller.MethodNotAllowed`.
- Path parameter: `controller.PathParam(r.URL.Path, "/prefix/")`, ditaruh di `default` setelah path statis. Prefix yang sama dipakai di controller untuk membaca nilainya.
- Handler berbentuk `func(w http.ResponseWriter, r *http.Request)` di `controller/`.
- Respons: `controller.WriteJSON(w, status, models.Response{...})`. Body request: `controller.DecodeJSON(w, r, &dst)` (batas 1 MiB). Jangan membangun JSON dengan `fmt.Sprintf`.
- Header custom baru dari frontend ditambahkan ke `Access-Control-Allow-Headers` di `config/cors.go`.

## Konfigurasi dan keamanan

- Semua konfigurasi dari env var yang dibaca di `config.SetEnv()`: `PRIVATEKEY`, `PUBLICKEY`, `MONGOSTRING`, `PORT`, `IP`, `ALLOWED_ORIGINS`. Jangan hardcode secret, key, atau connection string.
- Jangan mencetak nilai secret ke log. Cetak status `terisi`/`KOSONG` saja.
- Jangan `panic`/`log.Fatal` di `SetEnv()` atau `init()` (Cloud Functions gagal start). Simpan error-nya, lalu balas `503` dari endpoint yang membutuhkannya.
- MongoDB memakai `github.com/gocroot/mgdb`, koneksi dibuat sekali di `SetEnv()` (contoh di `README.md`, bagian Database Connection).
- NoSQL injection: jangan memasukkan `bson.M`/`map[string]any` hasil decode body request langsung ke filter. Decode ke struct bertipe jelas, lalu susun filter sendiri.
- Login di `controller/auth.go` adalah demo (env `DEMO_USER`/`DEMO_PASSWORD`). Untuk produksi, ganti dengan pengecekan user di database; password disimpan sebagai hash (bcrypt/argon2), dibandingkan dengan fungsi hash tersebut.

## Berhenti dengan rapi (sinyal)

- `run/main.go` menangani `SIGINT`, `SIGTERM`, dan **`SIGHUP`** dengan `signal.Notify`, lalu `srv.Shutdown` dengan batas waktu `masaTutup` (8 detik), dan `srv.Close()` jika batas habis. **Jangan hapus penanganan sinyal ini, dan jangan menghapus `SIGHUP`**: fitur *Hot restart* hosting (misalnya alwaysdata) mengirim `SIGHUP`, yang tanpa penanganan membunuh proses Go seketika.
- Jangan mengganti dengan `http.ListenAndServe` langsung tanpa shutdown. Proses yang terbunuh di tengah request bisa meninggalkan data separuh tertulis, dan proses yang tidak pernah keluar membuat binary baru tidak pernah dipakai.
- Server atau worker baru yang berjalan lama (consumer antrean, job terjadwal) wajib memakai pola yang sama: tangkap ketiga sinyal, selesaikan pekerjaan yang sedang jalan dengan batas waktu, lalu pastikan proses keluar.
- `masaTutup` harus lebih pendek dari jeda sebelum `SIGKILL` (Docker dan Cloud Run 10 detik).

## Struktur dan entry point

- Dua entry point, logika inisialisasi tetap di `config`: `gcf.go` (Google Cloud Run functions, fungsi `WebHook`) dan `run/main.go` (server lokal, VPS, web hosting, Docker).
- Model di `models/` dengan tag `json` (dan `bson` untuk MongoDB).

## Verifikasi

```sh
gofmt -l .
go vet ./...
go test -race ./...
go build ./...
```
