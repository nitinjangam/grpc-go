package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/nitinjangam/grpc-go/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, in *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Greet function called with: %v", in)
	req := in.GetReqGreet()
	log.Printf("Go req from client: %v", req)
	resp := greetpb.GreetResponse{
		ResGreet: "Pong",
	}
	return &resp, nil
}

func (*server) GreetMany(str greetpb.GreetService_GreetManyServer) error {
	log.Println("GreetMany function called")
	resp := ""
	for {
		req, err := str.Recv()
		if err != nil {
			break
		}
		resp += req.GetReqGreet() + " "
	}
	if err := str.SendAndClose(&greetpb.GreetResponse{
		ResGreet: resp,
	}); err != nil {
		log.Fatalf("Error while Send and close in GreetMany RPC: %v", err)
	}
	return nil
}

func (*server) GreetManyTimes(in *greetpb.GreetRequest, str greetpb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes function called with: %v", in)

	for i := 0; i < 5; i++ {
		resp := &greetpb.GreetResponse{
			ResGreet: "Nitin",
		}
		if err := str.Send(resp); err != nil {
			log.Fatalf("Got error while sending data to GreetManyTimes RPC: %v", err)
		}
	}
	return nil
}

func (*server) GreetAll(str greetpb.GreetService_GreetAllServer) error {
	log.Println("GreetAll function called")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	c := make(chan *greetpb.GreetRequest, 10)
	q := make(chan bool)
	go func() {
		defer wg.Done()
		for {
			req, err := str.Recv()
			if err != nil {
				q <- true
				break
			} else {
				c <- req
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		res := greetpb.GreetResponse{
			ResGreet: "Nitin",
		}
		for {
			select {
			case <-c:
				if err := str.Send(&res); err != nil {
					log.Fatalf("Error while send to GreetAll RPC: %v", err)
					break
				}
			case <-q:
				break
			}
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Started listening on port 50055")
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
