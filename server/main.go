package main

import (
	"context"
	"log"
	"net/http"

	pb "github.com/folucode/appointment-scheduler/proto"
	userconnect "github.com/folucode/appointment-scheduler/proto/protoconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/connect"
	"github.com/rs/cors"
)

type server struct{}

func (s *server) GetUser(ctx context.Context, req *connect.Request[pb.GetUserRequest]) (*connect.Response[pb.GetUserResponse], error) {
	log.Printf("Incoming Request to get user: %s", req.Msg.Id)

	res := connect.NewResponse(&pb.GetUserResponse{
		User: &pb.User{},
	})
	return res, nil
}

func main() {
	mux := http.NewServeMux()

	path, handler := userconnect.NewUserServiceHandler(&server{})
	mux.Handle(path, handler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"POST", "OPTIONS", "GET"},
		AllowedHeaders: []string{
			"Content-Type",
			"Connect-Protocol-Version",
			"Connect-Timeout-Ms",
		},
		Debug: true,
	})

	addr := ":8080"
	log.Printf("Server is listening on %s...", addr)

	err := http.ListenAndServe(addr, h2c.NewHandler(c.Handler(mux), &http2.Server{}))

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
