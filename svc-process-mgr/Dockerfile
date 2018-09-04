FROM golang:1.9-alpine as builder
RUN apk add --update glide git
ENV appDir=/go/src/lss_poc/svc-process-mgr
RUN mkdir -p $appDir
WORKDIR $appDir
ADD . $appDir
RUN go get .
RUN go build .
CMD ["/go/src/lss_poc/svc-process-mgr/svc-process-mgr"]

FROM alpine:latest
ENV TZ Etc/GMT
RUN apk add --update tzdata
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /go/src/lss_poc/svc-process-mgr/svc-process-mgr .
CMD ["./svc-process-mgr"]
EXPOSE 9000
