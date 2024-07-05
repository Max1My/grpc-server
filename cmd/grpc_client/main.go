package main

import (
	"context"
	"di_container/internal/model"
	descAccess "di_container/pkg/access_v1"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

var accessToken = flag.String("a", "", "access token")

const servicePort = 50051

func main() {
	flag.Parse()

	ctx := context.Background()
	md := metadata.New(map[string]string{"Authorization": "Bearer " + *accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", servicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial GRPC client: %s", err)
	}

	cl := descAccess.NewAccessV1Client(conn)

	_, err = cl.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: model.ExamplePath,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Access granted")
}
