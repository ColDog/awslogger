FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY main /root/main
EXPOSE 514/tcp
ENTRYPOINT ["/root/main"]
