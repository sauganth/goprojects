Prerequisites
=============
1) Install protoc
2) go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
3) go get -u google.golang.org/grpc
4) export PATH=$PATH:$GOPATH/bin
5) Run generate_proto_files.sh
6) Start grpc tf serving by running start_grpc_tf_serving.sh
7) Build "go build -x"
7) Run client "./tf --server_addr 127.0.0.1:8500 --model_name dense --model_version 1"
