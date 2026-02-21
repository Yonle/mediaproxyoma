## pleroma/akkoma alternative mediaproxy

**requirements**:
you need a working [go-bwhero](https://github.com/Yonle/go-bwhero) backend, and a golang compiler.

```
go build -o mediaproxyoma .
```


## running

set the two following variable names
- `BWHERO_HOST` for bandwidth hero server address (example: "http://localhost:8080/")
- `LISTEN` for listen address (syntax: "<listenaddr>:<port>")

running:
```
env BWHERO_HOST=http://localhost:8080/ LISTEN=0.0.0.0:8888 ./mediaproxyoma
```

optional to set: `USER_AGENT`

or, spin the entire thing alongside [go-bwhero](https://github.com/Yonle/go-bwhero) via docker compose:
```
docker compose up
```
it will be on localhost:8080.

then, configure your reverse proxy to forward any request going to /proxy/* to be forwarded to http://localhost:8080/ instead.

## optional, but IMPORTANT for production

by default, mediaproxyoma:
- has a default user agent linking to itself
- **didn't check for signature**

this is good for testing, but NOT for production server.

so, if you're planning to use this in production, you must configure the following environment variables:
- `USER_AGENT`
- `PLEROMA_SECRET_KEY_BASE`

you can obtain `PLEROMA_SECRET_KEY_BASE` in your `prod.secret.exs` config:
```
import Config

config :pleroma, Pleroma.Web.Endpoint,
   url: [host: "fedinet.waltuh.cyou", scheme: "https", port: 443],
   http: [ip: {127.0.0.1}, port: 4000],
   secret_key_base: "/xxxxxxxx/xxxxxxx",
   #                 ^ TAKE THIS
   live_view: [signing_salt: "XXXXXX"],
   signing_salt: "XXXXXX"
```

then set it in environment variable:
```
env \
  USER_AGENT="mediaproxyoma at fedinet.example.com; admin contact: mailto:admin@example.com" \
  PLEROMA_SECRET_KEY_BASE="/xxxxxxxx/xxxxxxx" \
  ./mediaproxyoma
```

notice: you **MUST KEEP THE SECRET KEY BASE SECRET**.

## misc if you wanna deal with old media

we have the following optional environment variables:
- `OLD_MEDIA_HOST`
- `OLD_MEDIA_PATHPREFIX`
- `NEW_MEDIA_HOST`
- `NEW_MEDIA_SCHEME`
- `NEW_MEDIA_PATHPREFIX`

say, you just have yourself migrated your media URL (and also it's files) from `https://eu2.somestorage.com/xxx:fedi/` to `https://media.waltuh.cyou/media/`, where that new URL is actually reverse proxying on a varnish backend at `http://127.0.0.1:6081` in the same machine/network as where your mediaproxyoma is running.

so, you need to run mediaproxyoma like this:
```
env \
  BWHERO_HOST=http://127.0.0.1:8111 \
  LISTEN=0.0.0.0:8080 \
  OLD_MEDIA_HOST=eu2.somestorage.com \
  OLD_MEDIA_PATHPREFIX=/xxx:fedi/ \
  NEW_MEDIA_HOST=127.0.0.1:6081 \
  NEW_MEDIA_SCHEME=http \
  NEW_MEDIA_PATHPREFIX=/media/ \
  ./mediaproxyoma
```

The following example will make mediaproxyoma proxy anything that goes to `https://eu2.somestorage.com/xxx:fedi/<filename>` going to `http://127.0.0.1:6081/media/<filename>` instead. It's useful for say, making full use of varnish cache.

## other probably useful environment variable

changing Allowed Origin for CORS: `ALLOW_ORIGIN` (default is `*`)