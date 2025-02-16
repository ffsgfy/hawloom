#!/bin/bash
#
# This is an example of how to create hawloom-config in your Kubernetes cluster

kubectl create configmap hawloom-config \
    "--from-file=config.json=$(dirname $0)/../config/dev.json" \
    --save-config --dry-run=client -o yaml | kubectl apply -f -

