#!/bin/bash
#
# This is an example of how to create hawloom-secret in your Kubernetes cluster

kubectl create secret generic hawloom-secret \
    "--from-literal=postgres-host=$POSTGRES_HOST" \
    "--from-literal=postgres-port=$POSTGRES_PORT" \
    "--from-literal=postgres-user=$POSTGRES_USER" \
    "--from-literal=postgres-password=$POSTGRES_PASSWORD" \
    --save-config --dry-run=client -o yaml | kubectl apply -f -

