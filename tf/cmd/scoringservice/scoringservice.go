package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"encoding/json"

	"golang.org/x/net/context"

	framework "tensorflow/core/framework"
	pb "tensorflow_serving"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	"scoring"
	"util"
)

// Server is the main concept that will house our Model Configuration
type server struct {
	data   sync.Map
	client pb.PredictionServiceClient
}

func (s *server) AddModelMap(ctx context.Context, req *scoring.AddModelMapRequest) (*scoring.AddModelMapResponse, error) {
	name := req.GetName()
	keyMapConfig := req.GetKeyMapConfig()
	log.Printf("Adding model config of model : %q\n", name)
	s.data.Store(name, keyMapConfig)
	return &scoring.AddModelMapResponse{
		Status: true,
	}, status.New(codes.OK, "").Err()

}

func (s *server) Predict(ctx context.Context, req *scoring.PredictRequest) (pr *scoring.PredictResponse, err error) {
	name := req.GetModelName()
	feats := req.GetFeats()
	if v, ok := s.data.Load(name); ok {
		log.Printf("found Model Config %q\n", name)
		log.Printf("Predicting using model : %q\n", name)
		responseMap := map[string]string{}
		var pr *pb.PredictRequest
		var version int64
		version = 1
		p := v.([]*scoring.KeyMapConfig)
		pr, err = newDensePredictRequest(&name, &version, feats, p)

		resp, err := s.client.Predict(context.Background(), pr)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Println("Got output from model...")
		responseMap = getResponseFromOutputTensors(resp.Outputs)

		return &scoring.PredictResponse{
			Status:      true,
			ResponseMap: responseMap,
		}, status.New(codes.OK, "").Err()
	}
	return &scoring.PredictResponse{Status: false}, status.Errorf(codes.NotFound, "could not find model %s", name)
}

func getResponseFromOutputTensors(outputs map[string]*framework.TensorProto) map[string]string {
	responseMap := map[string]string{}
	for k, v := range outputs {
		responseMap[k] = convertToString(v)
	}
	return responseMap
}

func convertToString(tensor *framework.TensorProto) string {
	valuesText := []string{}
	switch tensor.GetDtype() {
	case framework.DataType_DT_FLOAT:
		values := tensor.GetFloatVal()
		for i := range values {
			number := values[i]
			text := strconv.FormatFloat(float64(number), 'E', -1, 32)
			valuesText = append(valuesText, text)
		}

		// Join our string slice.
		result := strings.Join(valuesText, ",")
		return result
	case framework.DataType_DT_DOUBLE:
		values := tensor.GetDoubleVal()
		for i := range values {
			number := values[i]
			text := strconv.FormatFloat(number, 'E', -1, 64)
			valuesText = append(valuesText, text)
		}

		// Join our string slice.
		result := strings.Join(valuesText, ",")
		return result
	case framework.DataType_DT_INT64:
		values := tensor.GetInt64Val()
		for i := range values {
			number := values[i]
			text := strconv.FormatInt(number, 10)
			valuesText = append(valuesText, text)
		}

		// Join our string slice.
		result := strings.Join(valuesText, ",")
		return result
	case framework.DataType_DT_INT32:
		values := tensor.GetIntVal()
		for i := range values {
			number := values[i]
			text := strconv.FormatInt(int64(number), 10)
			valuesText = append(valuesText, text)
		}

		// Join our string slice.
		result := strings.Join(valuesText, ",")
		return result
	default:
		return "Unknown data type, couldn't understand!"
	}
}

func newDensePredictRequest(modelName *string, modelVersion *int64, feats map[string]string, keyMapConfig []*scoring.KeyMapConfig) (pr *pb.PredictRequest, err error) {
	pr = util.NewPredictRequest(*modelName, *modelVersion)

	for _, element := range keyMapConfig {
		dataType := element.GetDataType()
		inkey := element.GetInkey()
		outkey := element.GetOutkey()
		shape := element.GetShape()
		var dataT framework.DataType

		switch dataType {
		case scoring.KeyMapConfig_DataType_DT_FLOAT:
			dataT = framework.DataType_DT_FLOAT
		case scoring.KeyMapConfig_DataType_DT_DOUBLE:
			dataT = framework.DataType_DT_DOUBLE
		case scoring.KeyMapConfig_DataType_DT_INT32:
			dataT = framework.DataType_DT_INT32
		default:
			err = errors.New("Unknown data type")
		}
		if err != nil {
			return nil, err
		}
		inval := feats[inkey]
		var prod int64 = 1
		for _, x := range shape {
			prod *= x
		}
		tt, err := getTensorInputArray(inval, dataT, prod)
		if err != nil {
			err = errors.New("Mapping failed")
			return nil, err
		}

		util.AddInput(pr, outkey, dataT, tt, shape, nil)
	}
	//util.AddInput(pr, "keys", framework.DataType_DT_INT32, []int32{1, 2, 3}, nil, nil)
	//util.AddInput(pr, "features", framework.DataType_DT_FLOAT, []float32{
	//	1, 2, 3, 4, 5, 6, 7, 8, 9,
	//	1, 2, 3, 4, 5, 6, 7, 8, 9,
	//	1, 2, 3, 4, 5, 6, 7, 8, 9,
	//}, []int64{3, 9}, nil)
	return pr, nil
}

func getTensorInputArray(featString string, dataT framework.DataType, prod int64) (tensor interface{}, err error) {
	switch dataT {
	case framework.DataType_DT_FLOAT:
		tt := make([]float32, prod)
		//todo Add conversion from featString to tt
		dec := json.NewDecoder(strings.NewReader(featString))
		err := dec.Decode(&tt)
		return tt, err
	case framework.DataType_DT_DOUBLE:
		tt := make([]float64, prod)
		dec := json.NewDecoder(strings.NewReader(featString))
		err := dec.Decode(&tt)
		return tt, err
	case framework.DataType_DT_INT32:
		tt := make([]int32, prod)
		dec := json.NewDecoder(strings.NewReader(featString))
		err := dec.Decode(&tt)
		return tt, err
	default:
		return nil, errors.New("Unknown data type in getTensorInputArray")
	}
}
func main() {
	// open a port to communicate on
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new grpc server
	s := grpc.NewServer()

	// create our server
	ser := &server{data: sync.Map{}}
	// register our service

	scoring.RegisterScoringServer(s, ser)

	// Let the world know we are starting and where we are listening
	log.Printf("starting gRPC service on %s\n", lis.Addr())

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	serverAddr := "127.0.0.1:8500"
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	ser.client = pb.NewPredictionServiceClient(conn)
	// start listening and responding
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
