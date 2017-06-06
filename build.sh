#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main github.com/coldog/awslogger
docker build -t coldog/awslogger .
rm main
docker push coldog/awslogger
