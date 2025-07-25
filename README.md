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