FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -mod vendor -a -installsuffix cgo -ldflags="-w -s" -o api

FROM alpine

WORKDIR /app
COPY --from=builder /src/api .
COPY migrations migrations
EXPOSE 8080
CMD ["./api"]
# CMD ["sh", "-c", "./api"]