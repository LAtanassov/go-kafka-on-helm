#!/bin/sh
helm install ./go-app-chart -f ./go-app-chart/values.yaml -f ./go-app-chart/values.consumer.yaml
helm install ./go-app-chart -f ./go-app-chart/values.yaml -f ./go-app-chart/values.producer.yaml