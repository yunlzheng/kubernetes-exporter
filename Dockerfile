FROM        quay.io/prometheus/busybox:latest
MAINTAINER  yunl.zheng <yunl.zheng@gmail.com>

COPY kubenates-exporter /bin/kubenates-exporter
ENTRYPOINT [ "/bin/kubenates-exporter" ]
