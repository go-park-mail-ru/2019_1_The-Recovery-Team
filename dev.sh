#!/usr/bin/env bash

sudo lsof -i tcp:8080 | awk '{system("kill -9 " $2)}'
sudo lsof -i tcp:8081 | awk '{system("kill -9 " $2)}'
sudo lsof -i tcp:50051 | awk '{system("kill -9 " $2)}'
sudo lsof -i tcp:50052 | awk '{system("kill -9 " $2)}'
sudo lsof -i tcp:8082 | awk '{system("kill -9 " $2)}'

go run -race ./cmd/service/profile/main.go -local &
go run -race ./cmd/service/session/main.go -local &
go run -race ./cmd/server/main.go -local &
go run -race ./cmd/service/game/main.go -local &
go run -race ./cmd/service/chat/main.go -local &


