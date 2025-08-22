# core
Core Backend Skeleton


```sh
go run .\run\main.go
```

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

https://github.com/whatsauth/watoken

## CI / CD to Google Clound Function

save as file into folder

.github/workflows/main.yml

```yml
name: Google Cloud Function Deployment
on:
  push:
    branches:
      - main
jobs:
    Deploy:
      name: Deploy
      runs-on: ubuntu-latest
      steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: GCP Authentication
        id: 'auth'
        uses: 'google-github-actions/auth@v2'
        with:
          credentials_json: '${{ secrets.GOOGLE_CREDENTIALS }}'
      - name: Debug GCP credentials
        env:
          GOOGLE_APPLICATION_CREDENTIALS: ${{ secrets.GOOGLE_CREDENTIALS }}
        run: |
          echo "$GOOGLE_APPLICATION_CREDENTIALS" > credentials.json
      - name: 'Set up Cloud SDK'
        uses: 'google-github-actions/setup-gcloud@v2'
      - name: 'Use gcloud CLI'
        run: 'gcloud info'
      - name: 'Deploy a gen 2 cloud function'
        run: |
          gcloud functions deploy gocroot \
            --region=asia-southeast2 \
            --allow-unauthenticated \
            --entry-point=WebHook \
            --gen2 \
            --runtime=go125 \
            --trigger-http \
            --timeout=540s \
            --set-env-vars MONGOSTRING='${{ secrets.MONGOSTRING }}',PRIVATEKEY='${{ secrets.PRIVATEKEY }}',PUBLICKEY='${{ secrets.PUBLICKEY }}'
      - name: 'Cek eksistensi fungsi'
        run: 'gcloud functions describe gocroot --region=asia-southeast2'
      - name: 'Cek log debugging'
        run: 'gcloud functions logs read gocroot --region=asia-southeast2'
      - name: 'Cleaning Artifact Registry'
        run: 'gcloud artifacts repositories delete gcf-artifacts --location=asia-southeast2 --quiet'
```