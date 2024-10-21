## Running in development

```sh
go run main.go --listen-address=:8080
```

## Developing the index.html template

```sh
npx http-server
```

## Building

```sh
docker build -t envirobly/error-service .

docker run -p 8080:63108 envirobly/error-service
```

### Running a published image from GitHub registry

```sh
docker run -p 8080:63108 ghcr.io/envirobly/envirobly/error-service
```
