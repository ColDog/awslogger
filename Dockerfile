FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY main /root/main
ENTRYPOINT ["/root/main"]
