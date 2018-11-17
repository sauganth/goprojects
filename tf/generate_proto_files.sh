#!/bin/bash

set -ex
# install gRPC and protoc plugin for Go, see http://www.grpc.io/docs/quickstart/go.html#generate-grpc-code
rm -rf vendor/tensorflow* vendor/scoring
mkdir vendor/tensorflow vendor/tensorflow_serving vendor/scoring
protoc -I generate_golang_files/ generate_golang_files/*.proto --go_out=plugins=grpc:vendor/tensorflow_serving
protoc -I generate_golang_files/ generate_golang_files/tensorflow/core/framework/* --go_out=plugins=grpc:vendor
protoc -I generate_grpc/ generate_grpc/* --go_out=plugins=grpc:vendor/scoring
