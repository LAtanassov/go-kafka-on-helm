#!/bin/sh

go test ./...
go build -tags netgo ./...

docker build -t latanassov/go-consumer:0.2.0 .
docker login
docker push latanassov/go-consumer:0.2.0