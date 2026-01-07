# Šauktiniai

Wrote this to not enter asmens kodas.

Live at [nenoriu.fun](https://nenoriu.fun)

## Run locally

```bash
go run scripts/fetch.go
docker compose --profile dev up --build
```

## Run prod

```bash
wget https://raw.githubusercontent.com/tkozakas/sauktiniai/main/docker-compose.yml
wget https://raw.githubusercontent.com/tkozakas/sauktiniai/main/Caddyfile
wget https://raw.githubusercontent.com/tkozakas/sauktiniai/main/scripts/fetch.go
go run fetch.go
mv backend/data data
echo "TUNNEL_TOKEN=xxx" > .env
docker compose --profile prod up -d
```
