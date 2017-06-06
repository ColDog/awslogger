#!/bin/bash

./logger.sh | docker run -i --env-file=.env coldog/awslogger -group=test -stream=test -region=us-west-2
