package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"scoring"
)

func main() {
	// connect to the grpc server
	conn, err := grpc.Dial(":5051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// create a new client
	c := scoring.NewScoringClient(conn)

	ctx := context.Background()

	// write a value
	_, err = c.AddModelMap(ctx, &scoring.AddModelMapRequest{
		Name: "dense",
		KeyMapConfig: []*scoring.KeyMapConfig{
			&scoring.KeyMapConfig{
				DataType: scoring.KeyMapConfig_DataType_DT_INT32,
				Inkey:    "k",
				Outkey:   "keys",
				Shape:    []int64{3},
			},
			&scoring.KeyMapConfig{
				DataType: scoring.KeyMapConfig_DataType_DT_FLOAT,
				Inkey:    "f",
				Outkey:   "features",
				Shape:    []int64{3, 9},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Scoring client: Successfully configured model: dense")

	resp, err := c.Predict(ctx, &scoring.PredictRequest{
		ModelName: "dense",
		Feats: map[string]string{"k": "[1, 2, 3]",
			"f": "[1, 2, 3, 4, 5, 6, 7, 8, 9,1, 2, 3, 4, 5, 6, 7, 8, 9,1, 2, 3, 4, 5, 6, 7, 8, 9]"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Got output from model...")
	for k, v := range resp.ResponseMap {
		fmt.Println(k, v)
	}
	fmt.Printf("Scoring client: Successfully invoked model: dense")
}
