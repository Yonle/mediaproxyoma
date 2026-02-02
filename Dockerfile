FROM docker.io/alpine:edge

LABEL org.opencontainers.image.source="https://github.com/Yonle/mediaproxyoma" \
      org.opencontainers.image.description="pleroma/akkoma alternative mediaproxy backend" \
      org.opencontainers.image.licenses="BSD-3-Clause"

WORKDIR /a
COPY . .

RUN apk add --no-cache go \
    && env CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -buildvcs=false -o /a/mp . \
    && apk del go 

ENV LISTEN=0.0.0.0:8080

CMD ["/a/mp"]
