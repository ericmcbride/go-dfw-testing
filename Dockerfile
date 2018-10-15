FROM alpine:3.7

RUN apk -U add ca-certificates
COPY service service

CMD ["./service"]

