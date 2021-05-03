package main

import (
	"context"
	"log"

	"github.com/nitinjangam/grpc-go/greetpb"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Client dialing on port 50055")
	cc, err := grpc.Dial("0.0.0.0:50055", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while Dial: %v", err)
	}
	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)
	clientStream(c)
	serverStream(c)
	bidiStream(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := greetpb.GreetRequest{
		ReqGreet: "Ping",
	}
	res, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet RPC: %v", res.GetResGreet())
}

func clientStream(c greetpb.GreetServiceClient) {
	l := []string{"My", "name", "is", "Nitin"}
	cl, err := c.GreetMany(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetMany RPC: %v", err)
	}
	for _, v := range l {
		req := greetpb.GreetRequest{
			ReqGreet: v,
		}
		if err = cl.Send(&req); err != nil {
			log.Fatalf("Error while sending message on GreetMany RPC: %v", err)
		}
	}
	resp, err := cl.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while closing and receiving on GreetMany RPC: %v", err)
	}
	log.Printf("Response from GreetMany RPC: %v", resp.GetResGreet())
}

func serverStream(c greetpb.GreetServiceClient) {
	req := greetpb.GreetRequest{
		ReqGreet: "Hello",
	}
	str, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Fatalf("Got error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		resp, err := str.Recv()
		if err != nil {
			break
		}
		log.Printf("Response from GreetManyTimes RPC: %v", resp)
	}
}

func bidiStream(c greetpb.GreetServiceClient) {
	req := greetpb.GreetRequest{
		ReqGreet: "Jangam",
	}
	str, err := c.GreetAll(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetAll RPC: %v", err)
	}

	for i := 0; i <= 5; i++ {
		if err := str.Send(&req); err != nil {
			log.Fatalf("Error while send to GreetAll Server: %v", err)
		}
	}
	for {
		res, err := str.Recv()
		if err != nil {
			log.Fatalf("Error while receive from GreetAll RPC: %v", err)
			break
		}
		log.Fatalf("Response from GreetAll RPC: %v", res.GetResGreet())
	}
}
