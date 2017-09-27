#!/bin/bash
version=v1.0.0-beta
docker build -t kubernetes-exporter .
docker tag kubernetes-exporter registry.cn-hangzhou.aliyuncs.com/wise2c/kubernetes-exporter:$version
docker push registry.cn-hangzhou.aliyuncs.com/wise2c/kubernetes-exporter:$version
