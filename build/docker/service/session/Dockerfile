FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -mod vendor -a -installsuffix cgo -ldflags="-w -s" -o session cmd/service/session/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /src/session .
COPY wait-for-it.sh .

RUN ["chmod", "+x", "./wait-for-it.sh"]

EXPOSE 50052