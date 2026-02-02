# -- build

FROM docker.io/alpine:20260127 AS builder

RUN apk add --no-cache go

WORKDIR /src/
COPY . .

RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags="-s -w -buildid=" -buildvcs=false -o /out/exec.bin .

# -- after build

FROM docker.io/alpine:20260127

LABEL org.opencontainers.image.source="https://github.com/Yonle/mediaproxyoma" \
      org.opencontainers.image.description="pleroma/akkoma alternative mediaproxy backend" \
      org.opencontainers.image.licenses="BSD-3-Clause"

COPY --from=builder /out/exec.bin /bin/mediaproxyoma

ENV LISTEN=0.0.0.0:8080

CMD ["/bin/mediaproxyoma"]