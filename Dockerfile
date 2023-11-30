FROM golang:alpine AS build
WORKDIR /build
COPY .git/ .git/
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY scripts/ scripts/
COPY go.mod .
COPY go.sum .
RUN ash scripts/version.sh
RUN go build -v -o bin/goat cmd/goat/main.go

FROM alpine
LABEL org.opencontainers.image.authors="B12-Touch GmbH (hello@b12-touch.de)" \
      org.opencontainers.image.url="https://github.com/studio-b12/goat" \
      org.opencontainers.image.source="https://github.com/studio-b12/goat" \
      org.opencontainers.image.licenses="BSD-Clause-3"
COPY --from=build /build/bin/goat /bin/goat
ENTRYPOINT ["/bin/goat"]