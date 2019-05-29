FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -mod vendor -a -installsuffix cgo -ldflags="-w -s" -o api cmd/server/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /src/api .
COPY build build
COPY wait-for-it.sh .

VOLUME ["/app/upload/img"]

RUN ["chmod", "+x", "./wait-for-it.sh"]
