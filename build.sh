#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main github.com/coldog/awslogs-syslog
docker build -t coldog/awslogs-syslog .
rm main
