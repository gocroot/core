# core
Core Backend Skeleton dari [GoCroot](https://gocroot.if.co.id/). Dokumentasi lengkap: **https://gocroot.if.co.id/**

Memakai AI coding agent? Aturan proyek (PASETO bukan JWT, `net/http` tanpa framework, dsb.) ada di [`AGENTS.md`](AGENTS.md), dibaca otomatis oleh Antigravity, Codex, Cursor, Copilot, dan Claude Code (lewat `CLAUDE.md`).

Repo ini adalah **template repository**. Jangan di-import sebagai library. Buat proyek sendiri darinya, lalu ganti nama modulnya.

## Mulai Proyek Baru

1. Klik **Use this template → Create a new repository** di GitHub, atau jalankan:
   ```sh
   gh repo create myapi --template gocroot/core --private --clone
   cd myapi
   ```
2. Ganti nama modul `github.com/gocroot/core` menjadi milikmu:
   ```sh
   go mod edit -module github.com/username/myapi
   grep -rl 'github.com/gocroot/core' --include='*.go' . | xargs sed -i 's#github.com/gocroot/core#github.com/username/myapi#g'
   go build ./...
   ```
   Windows (PowerShell):
   ```powershell
   go mod edit -module github.com/username/myapi
   Get-ChildItem -Recurse -Filter *.go | ForEach-Object {
       (Get-Content $_.FullName -Raw) -replace 'github.com/gocroot/core', 'github.com/username/myapi' | Set-Content $_.FullName -NoNewline
   }
   go build ./...
   ```
   Pemanggilan di kode tetap `controller.Login`, `models.LoginAuth`, dan seterusnya. Yang berubah hanya baris import.
3. Jalankan:
   ```sh
   export PRIVATEKEY=... PUBLICKEY=... MONGOSTRING=...
   export DEMO_USER=admin DEMO_PASSWORD=pilih-password   # opsional, login demo
   go run ./run
   ```

## Environment Variables

| Env | Wajib | Keterangan |
|---|---|---|
| `PRIVATEKEY`, `PUBLICKEY` | Ya | Pasangan key PASETO untuk token login |
| `MONGOSTRING` | Jika memakai MongoDB | Connection string MongoDB |
| `PORT` | Tidak | Port server lokal/VPS, default `3000`. Di Cloud Functions port diatur otomatis |
| `IP` | Tidak | Alamat IP untuk listen. Kosong = semua interface. Diisi otomatis oleh hosting seperti alwaysdata. IPv6 didukung |
| `ALLOWED_ORIGINS` | Tidak | Origin CORS dipisah koma. Default: daftar di `config/cors.go` |
| `DEMO_USER`, `DEMO_PASSWORD` | Tidak | Mengaktifkan login demo di `POST /login`. Jika salah satu kosong, login dibalas `503`. **Jangan diisi di produksi**: ganti `controller.Login` dengan pengecekan user di database |

Saat start, `go run ./run` hanya mencetak status env (`terisi`/`KOSONG`), tidak pernah nilainya.
   Windows: `go run .\run\main.go`

## Test Login

Isi `DEMO_USER` dan `DEMO_PASSWORD`, lalu:

```sh
curl -X POST http://localhost:3000/login -d '{"username":"admin","password":"pilih-password"}'
curl http://localhost:3000/login -H "Token: <token dari login>"   # 200, atau 401 jika token tidak valid/kedaluwarsa
```

<img width="1920" height="1041" alt="image" src="https://github.com/user-attachments/assets/8b41b088-e9cf-4340-b4e6-fb5a79aa82d9" />


## Publish with Local Port

Publikasikan dengang menggunakan Cloudflare Tunnel

```sh
winget install --id Cloudflare.cloudflared
```

```sh
cloudflared tunnel --url 127.0.0.1:3000
```

## Graceful Shutdown

`run/main.go` menangani `SIGINT` (Ctrl+C), `SIGTERM` (`systemctl`, `docker stop`, Cloud Run, Render, Fly.io), dan `SIGHUP` (*Hot restart* di panel hosting seperti alwaysdata). Request yang sedang berjalan diberi waktu selesai maksimal 8 detik (`masaTutup`), lalu proses dipastikan keluar supaya binary baru benar-benar dipakai. Jangan hapus penanganan sinyal ini. Panduan lengkap: [Berhenti dengan Rapi](https://gocroot.if.co.id/deploy/#berhenti-dengan-rapi-graceful-shutdown).

## Token PASETO

Login di core menerbitkan token [PASETO](https://paseto.io) v4.public (Ed25519) lewat [whatsauth/watoken](https://github.com/whatsauth/watoken), bukan JWT:

- Tidak ada header `alg` yang bisa diatur penyerang, jadi celah JWT seperti `alg: none` dan kebingungan RS256/HS256 tidak mungkin terjadi.
- Kriptografi terkunci per versi (v4.public selalu Ed25519), programmer tidak perlu dan tidak bisa salah memilih algoritma.
- Token ditandatangani dengan `PRIVATEKEY` dan diverifikasi dengan `PUBLICKEY`. Service lain cukup memegang `PUBLICKEY` untuk memeriksa token tanpa server auth pusat. Ini sifat tanda tangan asimetris (JWT RS256/EdDSA juga bisa), dan di PASETO menjadi satu-satunya pilihan.

Token tidak bisa dicabut sebelum kedaluwarsa, jadi masa berlakunya pendek (`watoken.Encode` = 2 jam). Butuh logout paksa? Tambahkan denylist, jangan memperpanjang masa token. Simpan `PRIVATEKEY` hanya di service yang menerbitkan token.

## Generate Token

to generate PRIVATEKEY and PUBLICKEY

Clone this repo:

https://github.com/whatsauth/watoken

run 
```sh
go test
```

## CI / CD to Google Cloud Function

> **Panduan saja.** Repo ini sengaja tidak menyertakan workflow aktif. Untuk mengaktifkan deploy otomatis, salin YAML di bawah ke `.github/workflows/main.yml` di repo proyekmu (hasil fork/clone dari core). Setelah itu, setiap push ke `main` akan men-deploy function.

### Prasyarat

1. Service account GCP dengan role **Cloud Functions Developer**, **Service Account User**, **Cloud Build Editor**, dan **Artifact Registry Administrator** (untuk langkah cleanup).
2. Buat key JSON service account tersebut, lalu tambahkan **GitHub Secrets** di repo (Settings → Secrets and variables → Actions):

| Secret | Isi |
|---|---|
| `GOOGLE_CREDENTIALS` | Isi file key JSON service account |
| `GCP_PROJECT_ID` | ID project GCP |
| `MONGOSTRING` | Connection string MongoDB |
| `PRIVATEKEY` | Private key PASETO (lihat [Generate Token](#generate-token)) |
| `PUBLICKEY` | Public key PASETO |
| `ALLOWED_ORIGINS` | *(opsional)* Daftar origin CORS dipisah koma, misalnya `https://app.example.com,https://admin.example.com`. Jika kosong, dipakai daftar default di `config/cors.go` |

3. Ganti `FUNCTION_NAME` dengan nama yang unik di project GCP-mu. Deploy ke nama yang sudah dipakai function lain akan **menimpa** function tersebut.

### Workflow

```yml
name: Google Cloud Function Deployment

on:
  push:
    branches:
      - main
  workflow_dispatch: # bisa dijalankan manual dari tab Actions

env:
  FUNCTION_NAME: gocroot-core   # ganti dengan nama function milikmu
  REGION: asia-southeast2
  RUNTIME: go127                # sesuaikan dengan direktif go di go.mod

jobs:
  deploy:
    name: Deploy Cloud Function
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v7

      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v3
        with:
          credentials_json: ${{ secrets.GOOGLE_CREDENTIALS }}
          project_id: ${{ secrets.GCP_PROJECT_ID }}

      - name: Set up Cloud SDK
        uses: google-github-actions/setup-gcloud@v3

      - name: Deploy Cloud Function
        run: |
          gcloud functions deploy ${{ env.FUNCTION_NAME }} \
            --gen2 \
            --region=${{ env.REGION }} \
            --runtime=${{ env.RUNTIME }} \
            --entry-point=WebHook \
            --trigger-http \
            --allow-unauthenticated \
            --timeout=540s \
            --memory=256Mi \
            --min-instances=0 \
            --max-instances=100 \
            --set-env-vars MONGOSTRING='${{ secrets.MONGOSTRING }}',PRIVATEKEY='${{ secrets.PRIVATEKEY }}',PUBLICKEY='${{ secrets.PUBLICKEY }}',ALLOWED_ORIGINS='${{ secrets.ALLOWED_ORIGINS }}' \
            --quiet

      - name: Verify deployment
        run: |
          URL=$(gcloud functions describe ${{ env.FUNCTION_NAME }} --region=${{ env.REGION }} --format="value(serviceConfig.uri)")
          echo "Function URL: $URL"
          # Endpoint "/" pada core mengembalikan {"status":"success","message":"Welcome"}
          curl -fsS "$URL/"

      - name: Cleaning Artifact Registry
        run: gcloud artifacts repositories delete gcf-artifacts --location=${{ env.REGION }} --quiet

      - name: Show recent logs (on failure)
        if: failure()
        run: gcloud functions logs read ${{ env.FUNCTION_NAME }} --region=${{ env.REGION }} --limit=50
```

### Catatan

- **Runtime:** runtime Go yang didukung Google punya tanggal deprecation. Cek [jadwal runtime](https://docs.cloud.google.com/functions/docs/runtime-support), lalu naikkan `RUNTIME` bersama direktif `go` di `go.mod`. `go127` didukung sampai deprecation sekitar Agu/Sep 2027.
- **Cleanup (hemat biaya):** setiap deploy Cloud Functions menghasilkan image buildpack berukuran ratusan MB di repository `gcf-artifacts`, padahal free tier Artifact Registry hanya 0,5 GB. Karena itu seluruh repository dihapus setelah deploy berhasil, dan storage tetap mendekati nol. Function yang sedang berjalan tidak terpengaruh, dan repository akan dibuat ulang otomatis pada deploy berikutnya. Konsekuensinya: rollback ke revisi lama tidak bisa dilakukan karena image-nya sudah terhapus. Kalau kamu butuh rollback, hapus langkah ini dan pantau penggunaan storage.
- **Secret:** `--set-env-vars` membuat nilai secret terlihat di konfigurasi function (Console GCP). Untuk produksi, simpan secret di Secret Manager dan pakai `--set-secrets PRIVATEKEY=privatekey:latest,...`. Service account runtime function perlu role **Secret Manager Secret Accessor**.
- **Uji di lokal** sebelum push: `go build ./... && go vet ./...`, lalu `go run ./run` dan buka `http://localhost:3000/`.

## Database Connection

Untuk koneksi database menggunakan mgdb:

https://github.com/gocroot/mgdb

```sh
go get github.com/gocroot/mgdb
```

Di `config/config.go`, tambahkan variabel untuk menampung koneksi database dan error-nya, lalu isi di dalam `SetEnv()`:

```go
var (
	// ... variabel yang sudah ada (Port, MongoString, PrivateKey, PublicKey)

	// Database variables
	MongoDB  *mongo.Database
	MongoErr error
)

func SetEnv() {
	once.Do(func() {
		// ... pembacaan env yang sudah ada
		MongoDB, MongoErr = mgdb.MongoConnect(mgdb.DBInfo{
			DBString: MongoString,
			DBName:   "mydatabase",
		})
		if MongoErr != nil {
			log.Println("MongoDB belum terhubung:", MongoErr) // jangan log.Fatal
		}
	})
}
```

**Jangan pakai `log.Fatal` atau `panic` di sini.** Di Cloud Functions, `SetEnv()` dipanggil saat container start. Kalau proses mati di tahap itu, deploy gagal dengan error *Container Healthcheck failed* dan endpoint lain yang tidak butuh database ikut mati. Simpan error-nya, lalu cek di controller yang butuh database:

```go
if config.MongoErr != nil {
	WriteJSON(w, http.StatusServiceUnavailable, models.Response{Status: "Error", Message: "Database tidak tersedia"})
	return
}
```

Pastikan juga `MONGOSTRING` berisi URI `mongodb://` atau `mongodb+srv://` yang valid. Nilai kosong atau rusak bisa membuat `mgdb` panic saat koneksi.
