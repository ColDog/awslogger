# AWSLogger

Reads from STDIN and pushes logs to Cloudwatch Logs.

```
Usage: awslogger [options]

    Pipe logs from various sources into Cloudwatch Logs. All logs are piped through STDIN.

Example With journalctl:

    journalctl -o short-iso -f | docker run -i --name=logger.service coldog/awslogger

Options:
  -group string
    	AWS Cloudwatch group name, will be created if it doesn't exist.
  -max-put-size int
    	Max put message size. (default 5000)
  -message-buffer-size int
    	Message buffer size. (default 20000)
  -read-buffer-size int
    	Read buffer size. (default 1024)
  -stream string
    	AWS Cloudwatch stream name, will be created if it doesn't exist. Defaults to the hostname. (default "271806a8fda8")
  -time-fmt string
    	GO time formatting string. (default "2006-01-02T15:04:05-0700")
```
