Prerequisites
=============
1) Install protoc
2) go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
3) go get -u google.golang.org/grpc
4) export PATH=$PATH:$GOPATH/bin
5) Run generate_proto_files.sh
6) Start grpc tf serving by running start_grpc_tf_serving.sh to start tf serving on port 8500
7) Build "go build cmd/scoringClient/*"
8) Build "go build cmd/scoringservice/*"
9) Build "go build cmd/tf_client/*"
10) Run tf client "./grpc_tf_predict_client --server_addr 127.0.0.1:8500 --model_name dense --model_version 1" to connect to docker tf container directly
11) Start grpc scoring serving by running ./scoringservice to run scoring service on port 5051
12) Run scoring client to go through the scoring service to the tf serving to score "./scoringclient"
