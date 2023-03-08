package api

import (
	"context"
	"fmt"

	pb "github.com/dlshle/authnz/proto"
	"google.golang.org/grpc"
)

func test() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	c := pb.NewAuthNZClient(conn)

	r, err := c.Authorize(context.Background(), &pb.AuthorizeRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", r)
}
