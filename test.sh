#!/bin/bash

./logger.sh | docker run -i --env-file=.env coldog/awslogs -group=test -stream=test -region=us-west-2
