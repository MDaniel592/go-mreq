# Build server

```
go build -tags=server -o bin/server src/main.go src/server.go
```

# Build client

```
go build -tags=client -o bin/client src/main.go src/client.go
```


1. `go mod init mreq`
2. `go get github.com/valyala/fasthttp`