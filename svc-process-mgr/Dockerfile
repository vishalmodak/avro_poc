FROM golang:1.10-stretch as builder
RUN apt-get install git
ENV appDir=/go/src/avro_poc/svc-process-mgr
RUN mkdir -p $appDir
WORKDIR $appDir
ADD . $appDir
RUN go get .
RUN go build .

FROM golang:1.10-stretch
ENV TZ Etc/GMT
RUN apt-get install tzdata
# RUN apt-get --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /go/src/avro_poc/svc-process-mgr/*.avsc /root/
COPY --from=builder /go/src/avro_poc/svc-process-mgr/svc-process-mgr .
CMD ["./svc-process-mgr"]
EXPOSE 9000
