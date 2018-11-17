package main

import (
	"context"
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
		Name: "band",
	})
	if err != nil {
		log.Fatal(err)
	}

}
