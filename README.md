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

to generate PRIVATEKEY and PUBLICKEY

Clone this repo:

https://github.com/whatsauth/watoken

run 
```sh
go test
```

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


hasil improve (eksperimental)

```yml
name: Google Cloud Function Deployment

on:
  push:
    branches:
      - main
  workflow_dispatch: # Allow manual triggering

env:
  PROJECT_ID: ${{ secrets.GCP_PROJECT_ID }}
  FUNCTION_NAME: gocroot
  REGION: asia-southeast2

jobs:
  deploy:
    name: Deploy Cloud Function
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Authenticate to Google Cloud
      id: 'auth'
      uses: 'google-github-actions/auth@v2'
      with:
        credentials_json: '${{ secrets.GOOGLE_CREDENTIALS }}'
        project_id: ${{ secrets.GCP_PROJECT_ID }}
    
    - name: Set up Cloud SDK
      uses: 'google-github-actions/setup-gcloud@v2'
    
    - name: Configure gcloud project
      run: |
        gcloud config set project ${{ env.PROJECT_ID }}
        gcloud config list
    
    - name: Deploy Cloud Function
      run: |
        gcloud functions deploy ${{ env.FUNCTION_NAME }} \
          --region=${{ env.REGION }} \
          --allow-unauthenticated \
          --entry-point=WebHook \
          --gen2 \
          --runtime=go121 \
          --trigger-http \
          --timeout=540s \
          --memory=256Mi \
          --min-instances=0 \
          --max-instances=100 \
          --set-env-vars MONGOSTRING='${{ secrets.MONGOSTRING }}',PRIVATEKEY='${{ secrets.PRIVATEKEY }}',PUBLICKEY='${{ secrets.PUBLICKEY }}' \
          --quiet
    
    - name: Verify deployment
      run: |
        echo "Checking function status..."
        gcloud functions describe ${{ env.FUNCTION_NAME }} --region=${{ env.REGION }} --format="value(status)"
        
        echo "Function URL:"
        gcloud functions describe ${{ env.FUNCTION_NAME }} --region=${{ env.REGION }} --format="value(serviceConfig.uri)"
    
    - name: Run health check
      run: |
        FUNCTION_URL=$(gcloud functions describe ${{ env.FUNCTION_NAME }} --region=${{ env.REGION }} --format="value(serviceConfig.uri)")
        echo "Testing function at: $FUNCTION_URL"
        # Add your health check endpoint here
        # curl -f "$FUNCTION_URL/health" || exit 1
    
    - name: Check recent logs (on failure)
      if: failure()
      run: |
        echo "Checking recent logs for debugging..."
        gcloud functions logs read ${{ env.FUNCTION_NAME }} --region=${{ env.REGION }} --limit=50

  cleanup:
    name: Cleanup Old Artifacts
    runs-on: ubuntu-latest
    needs: deploy
    if: success()
    
    steps:
    - name: Authenticate to Google Cloud
      uses: 'google-github-actions/auth@v2'
      with:
        credentials_json: '${{ secrets.GOOGLE_CREDENTIALS }}'
        project_id: ${{ secrets.GCP_PROJECT_ID }}
    
    - name: Set up Cloud SDK
      uses: 'google-github-actions/setup-gcloud@v2'
    
    - name: Clean up old artifacts
      run: |
        # List repositories first
        gcloud artifacts repositories list --location=${{ env.REGION }}
        
        # Only delete if gcf-artifacts exists and has old images
        if gcloud artifacts repositories describe gcf-artifacts --location=${{ env.REGION }} &>/dev/null; then
          echo "Cleaning up old container images..."
          # Delete images older than 7 days
          gcloud artifacts docker images list ${{ env.REGION }}-docker.pkg.dev/${{ env.PROJECT_ID }}/gcf-artifacts \
            --filter="createTime<'-P7D'" --format="value(IMAGE)" | \
            xargs -r gcloud artifacts docker images delete --quiet
        else
          echo "No gcf-artifacts repository found, skipping cleanup"
        fi
```