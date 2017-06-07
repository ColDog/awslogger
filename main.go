package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"time"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var (
	groupName         string
	streamName        string
	region            string
	timeFormatString  = "2006-01-02T15:04:05-0700"
	readBufferSize    = 1024
	messageBufferSize = 20000
	maxPutSize        = 5000
)

func usage() {
	usage := `Usage: awslogger [options]

    Pipe logs from various sources into Cloudwatch Logs. All logs are piped through STDIN.

Example With journalctl:

    journalctl -o short-iso -f | docker run -i --name=logger.service coldog/awslogger

Options:
`
	fmt.Fprintf(os.Stderr, usage)
	flag.PrintDefaults()
}

func init() {
	hostname, _ := os.Hostname()
	flag.Usage = usage
	flag.StringVar(&groupName, "group", "", "AWS Cloudwatch group name, will be created if it doesn't exist.")
	flag.StringVar(&streamName, "stream", hostname, "AWS Cloudwatch stream name, will be created if it doesn't exist. Defaults to the hostname.")
	flag.StringVar(&region, "region", "", "AWS region")
	flag.StringVar(&timeFormatString, "time-fmt", timeFormatString, "GO time formatting string.")
	flag.IntVar(&readBufferSize, "read-buffer-size", readBufferSize, "Read buffer size.")
	flag.IntVar(&messageBufferSize, "message-buffer-size", messageBufferSize, "Message buffer size.")
	flag.IntVar(&maxPutSize, "max-put-size", maxPutSize, "Max put message size.")
	flag.Parse()
}

type logger struct {
	client   *cloudwatchlogs.CloudWatchLogs
	messages chan string
	seqToken *string
}

func (l *logger) setup() error {
	l.client.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{LogGroupName: &groupName})
	l.client.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &groupName,
		LogStreamName: &streamName,
	})

	streams, err := l.client.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        &groupName,
		LogStreamNamePrefix: &streamName,
	})
	if err != nil {
		return err
	}

	stream := streams.LogStreams[0]
	l.seqToken = stream.UploadSequenceToken
	return nil
}

func (l *logger) run() {
	buffer := []string{}
	for {
		select {
		case msg := <-l.messages:
			buffer = append(buffer, msg)
			if len(buffer) >= maxPutSize {
				l.flush(buffer)
				buffer = nil
			}
		case <-time.After(250 * time.Millisecond):
			if len(buffer) > 0 {
				l.flush(buffer)
				buffer = nil
			}
		}
	}
}

func (l *logger) flush(messages []string) {
	events := []*cloudwatchlogs.InputLogEvent{}

	for _, msg := range messages {
		t, err := time.Parse(timeFormatString, strings.Split(msg, " ")[0])
		if err != nil {
			continue
		}

		events = append(events, &cloudwatchlogs.InputLogEvent{
			Message:   &msg,
			Timestamp: aws.Int64(t.UnixNano() / int64(time.Millisecond)),
		})
	}

	out, err := l.client.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogEvents:     events,
		LogGroupName:  &groupName,
		LogStreamName: &streamName,
		SequenceToken: l.seqToken,
	})
	if err != nil {
		return
	}
	l.seqToken = out.NextSequenceToken
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:                        &region,
		CredentialsChainVerboseErrors: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("failed to get aws session: %v", err)
	}

	lg := &logger{
		client:   cloudwatchlogs.New(sess),
		messages: make(chan string, messageBufferSize),
	}

	err = lg.setup()
	if err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}

	go lg.run()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		lg.messages <- text
	}
}
