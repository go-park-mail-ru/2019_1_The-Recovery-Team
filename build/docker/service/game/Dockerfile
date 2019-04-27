FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -mod vendor -a -installsuffix cgo -ldflags="-w -s" -o game cmd/service/game/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /src/game .
COPY build/config/ build/config/

EXPOSE 8081