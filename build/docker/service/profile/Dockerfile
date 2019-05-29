FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -mod vendor -a -installsuffix cgo -ldflags="-w -s" -o profile cmd/service/profile/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /src/profile .
COPY build/schema/ build/schema/
COPY build/config/ build/config/
COPY wait-for-it.sh .


RUN ["chmod", "+x", "./wait-for-it.sh"]

EXPOSE 50051