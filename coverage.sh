#!/usr/bin/env bash

go test -coverpkg=./... -coverprofile=coverage.out.tmp  ./...
cat coverage.out.tmp | grep -v -E "_easyjson.go|.pb.go|docs.go" > coverage.out
rm coverage.out.tmp
go tool cover -func=coverage.out
go tool cover -html=coverage.out
rm coverage.out