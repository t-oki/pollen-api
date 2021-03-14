FROM golang:1.16.2-alpine3.13 as builder
ENV APPDIR $GOPATH/src/github.com/t-oki/pollen-api
ENV GO111MODULE on
RUN \
  apk add --update && \
  rm -rf /var/cache/apk/* && \
  mkdir -p $APPDIR \
ADD . $APPDIR/
WORKDIR $APPDIR
RUN go build -mod=vendor -ldflags "-s -w" -o pollen-api cmd/api/main.go

FROM alpine:3.13
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/t-oki/pollen-api/pollen-api ./
ENTRYPOINT ["./pollen-api"]