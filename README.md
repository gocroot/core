# core
Core Backend Skeleton

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
   go run ./run
   ```
   Windows: `go run .\run\main.go`

## Test Login

<img width="1920" height="1041" alt="image" src="https://github.com/user-attachments/assets/8b41b088-e9cf-4340-b4e6-fb5a79aa82d9" />


## Publish with Local Port

Publikasikan dengang menggunakan Cloudflare Tunnel

```sh
winget install --id Cloudflare.cloudflared
```

```sh
cloudflared tunnel --url 127.0.0.1:3000
```

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

1. Service account GCP dengan role **Cloud Functions Developer**, **Service Account User**, **Cloud Build Editor**, dan **Artifact Registry Administrator** (yang terakhir hanya diperlukan untuk langkah cleanup).
2. Buat key JSON service account tersebut, lalu tambahkan **GitHub Secrets** di repo (Settings → Secrets and variables → Actions):

| Secret | Isi |
|---|---|
| `GOOGLE_CREDENTIALS` | Isi file key JSON service account |
| `GCP_PROJECT_ID` | ID project GCP |
| `MONGOSTRING` | Connection string MongoDB |
| `PRIVATEKEY` | Private key PASETO (lihat [Generate Token](#generate-token)) |
| `PUBLICKEY` | Public key PASETO |

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
            --set-env-vars MONGOSTRING='${{ secrets.MONGOSTRING }}',PRIVATEKEY='${{ secrets.PRIVATEKEY }}',PUBLICKEY='${{ secrets.PUBLICKEY }}' \
            --quiet

      - name: Verify deployment
        run: |
          URL=$(gcloud functions describe ${{ env.FUNCTION_NAME }} --region=${{ env.REGION }} --format="value(serviceConfig.uri)")
          echo "Function URL: $URL"
          # Endpoint "/" pada core mengembalikan {"status":"success","message":"Welcome"}
          curl -fsS "$URL/"

      - name: Show recent logs (on failure)
        if: failure()
        run: gcloud functions logs read ${{ env.FUNCTION_NAME }} --region=${{ env.REGION }} --limit=50

  cleanup:
    name: Cleanup Old Images
    runs-on: ubuntu-latest
    needs: deploy
    if: success()

    steps:
      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v3
        with:
          credentials_json: ${{ secrets.GOOGLE_CREDENTIALS }}
          project_id: ${{ secrets.GCP_PROJECT_ID }}

      - name: Set up Cloud SDK
        uses: google-github-actions/setup-gcloud@v3

      - name: Delete images older than 7 days
        run: |
          REPO=${{ env.REGION }}-docker.pkg.dev/${{ secrets.GCP_PROJECT_ID }}/gcf-artifacts
          if gcloud artifacts repositories describe gcf-artifacts --location=${{ env.REGION }} &>/dev/null; then
            gcloud artifacts docker images list "$REPO" \
              --filter="createTime<-P7D" --format="value(IMAGE)" | \
              xargs -r -n1 gcloud artifacts docker images delete --delete-tags --quiet
          else
            echo "Repository gcf-artifacts tidak ada, cleanup dilewati"
          fi
```

### Catatan

- **Runtime:** runtime Go yang didukung Google punya tanggal deprecation. Cek [jadwal runtime](https://docs.cloud.google.com/functions/docs/runtime-support), lalu naikkan `RUNTIME` bersama direktif `go` di `go.mod`. `go127` didukung sampai deprecation sekitar Agu/Sep 2027.
- **Cleanup:** jangan menghapus seluruh repository `gcf-artifacts` setiap deploy. Repository itu dipakai bersama oleh semua function di region yang sama, jadi cukup hapus image lama seperti contoh di atas.
- **Secret:** `--set-env-vars` membuat nilai secret terlihat di konfigurasi function (Console GCP). Untuk produksi, simpan secret di Secret Manager dan pakai `--set-secrets PRIVATEKEY=privatekey:latest,...`. Service account runtime function perlu role **Secret Manager Secret Accessor**.
- **Uji di lokal** sebelum push: `go build ./... && go vet ./...`, lalu `go run ./run` dan buka `http://localhost:3000/`.

## Database Connection

Untuk koneksi database menggunakan mgdb:

https://github.com/gocroot/mgdb

```sh
go get github.com/gocroot/mgdb
```

dalam config tambahkan variabel untuk menampung koneksi database mongo dengan errornya

```go
var (
    once        sync.Once
    Port        string
    MongoString string
    PrivateKey  string
    PublicKey   string
    
    // Database variables
    MongoDB     *mongo.Database
    //opsional untuk pengecekan error koneksi
    MongoErr    error
)

// SetEnv dengan protection untuk multiple calls
func SetEnv() {
	once.Do(func() {
		// Load environment variables
		// Set default configurations
		// Initialize global settings
		Port = ":3000"
		MongoString = os.Getenv("MONGOSTRING")
		PrivateKey = os.Getenv("PRIVATEKEY")
		PublicKey = os.Getenv("PUBLICKEY")
        var mconn = mgdb.DBInfo{
            DBString: MongoString,
            DBName:   "mydatabase",
        }
        MongoDB, MongoErr = mgdb.MongoConnect(mconn)
        //opsional untuk pengecekan koneksi, bisa dilakukan di luar config
        if MongoErr != nil {
            log.Fatalf("Failed to connect to MongoDB: %v", MongoErr)
        }
	})
}

```
