#FROM        quay.io/prometheus/busybox:latest
FROM golang:1.8
MAINTAINER  yunl.zheng <yunl.zheng@gmail.com>

COPY ./ /go/src/github.com/yunlzheng/kubernates-exporter/
WORKDIR /go/src/github.com/yunlzheng/kubernates-exporter/
RUN go build

RUN cp kubernates-exporter /bin/kubenates-exporter
ENTRYPOINT [ "/bin/kubenates-exporter" ]
