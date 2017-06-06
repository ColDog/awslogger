# AWSLogger

Reads from STDIN and pushes logs to Cloudwatch Logs.

## Example

```
./logger.sh | docker run -i coldog/awslogger -group=test -stream=test
```

## Usage
```
Usage of /root/main:
  -group string
    	AWS Cloudwatch group name, will be created if it doesn't exist
  -max-put-size int
    	Max put message size (default 10000)
  -message-buffer-size int
    	Message buffer size (default 20000)
  -read-buffer-size int
    	UDP read buffer size (default 1024)
  -region string
    	AWS region
  -stream string
    	AWS Cloudwatch stream name, will be created if it doesn't exist (default "b88774fdc25e")
  -tf string
    	GO time formatting string expects timestamp to be first element in log entries (default "2006-01-02T15:04:05-0700")
```
