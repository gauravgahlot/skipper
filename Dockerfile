FROM alpine:3.11

ENTRYPOINT ["skipper"]
EXPOSE 42113

COPY skipper /bin/
