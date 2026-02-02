## pleroma/akkoma alternative mediaproxy

all you need is a working [bandwidth hero](https://github.com/Yonle/go-bwhero) backend, and a golang compiler.

```
go build -o mediaproxy .
```

then, set the two following variable names
- `BWHERO_HOST` for bandwidth hero server address (example: "http://localhost:8080/")
- `LISTEN` for listen address (syntax: "<listenaddr>:<port>")

running:
```
env BWHERO_HOST=http://localhost:8080/ LISTEN=0.0.0.0:8888 ./mediaproxy
```

or, spin the entire thing alongside [go-bwhero](https://github.com/Yonle/go-bwhero) via docker compose:
```
docker compose up
```
it will be on localhost:8080.

then, configure your reverse proxy to forward any request going to /proxy/* to be forwarded to http://localhost:8888/ instead.