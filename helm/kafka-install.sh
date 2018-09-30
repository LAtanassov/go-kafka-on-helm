#!/bin/sh
helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
helm install --name go-kafka incubator/kafka