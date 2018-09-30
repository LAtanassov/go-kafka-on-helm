#!/bin/sh

go test ./...
go build -tags netgo ./...

docker build -t latanassov/go-producer:0.1.0 .
docker login
docker push latanassov/go-producer:0.1.0