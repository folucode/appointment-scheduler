package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/folucode/appointment-scheduler/internal/db"
	pb "github.com/folucode/appointment-scheduler/proto"
	protoconnect "github.com/folucode/appointment-scheduler/proto/protoconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/rs/cors"
)

type AppointmentServer struct {
	protoconnect.UnimplementedAppointmentServiceHandler
	Storage *db.Database
}

type UserServer struct {
	protoconnect.UnimplementedUserServiceHandler
	Storage *db.Database
}

func (s *AppointmentServer) CreateAppointment(
	ctx context.Context,
	req *connect.Request[pb.CreateAppointmentRequest],
) (*connect.Response[pb.Appointment], error) {
	log.Printf("Incoming Request to create an appointment: %+v", req.Msg)

	user, err := s.Storage.FindUserByEmail(ctx, req.Msg.ContactInformation.Email)

	if err != nil {
		if !errors.Is(err, db.ErrUserNotFound) {
			return nil, err
		}

		createdUser, err := s.Storage.CreateUser(ctx, &pb.User{
			Id:    uuid.NewString(),
			Name:  req.Msg.ContactInformation.Name,
			Email: req.Msg.ContactInformation.Email,
		})

		if err != nil {
			return nil, err
		}

		user = createdUser
	}

	newAppt := &pb.Appointment{
		Id:          uuid.NewString(),
		Description: req.Msg.Description,
		UserId:      user.Id,
		ContactInformation: &pb.ContactInformation{
			Name:  req.Msg.ContactInformation.Name,
			Email: req.Msg.ContactInformation.Email,
		},
		StartTime: req.Msg.StartTime,
		EndTime:   req.Msg.EndTime,
	}

	err = s.Storage.CreateAppointment(ctx, newAppt)

	if err != nil {
		log.Printf("Error saving to database: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to save appointment"))
	}

	return connect.NewResponse(newAppt), nil
}

func (s *AppointmentServer) GetUserAppointments(
	ctx context.Context,
	req *connect.Request[pb.GetUserAppointmentRequest],
) (*connect.Response[pb.GetUserAppointmentResponse], error) {
	log.Printf("Incoming Request to get user appointments: %+v", req.Msg)

	if req.Msg.UserId == "" {
		return connect.NewResponse(&pb.GetUserAppointmentResponse{
			Appointments: []*pb.Appointment{},
		}), errors.New("user ID not supplied")
	}

	data, err := s.Storage.GetAppointments(ctx, req.Msg.UserId)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&pb.GetUserAppointmentResponse{
		Appointments: data,
	})

	return res, nil
}

// func (s *AppointmentServer) DeleteAppointment(ctx context.Context, req *connect.Request[pb.DeleteAppointmentRequest]) (*connect.Response[pb.DeleteAppointmentResponse], error) {
// 	log.Printf("Incoming Request to get user appointments: %+v", req.Msg)

// 	if req.Msg.Id == "" {
// 		return &connect.Response[pb.DeleteAppointmentResponse]{}, nil
// 	}

// 	success, err := s.Storage.DeleteAppointment(req.Msg.Id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&pb.DeleteAppointmentResponse{
// 		Success: success,
// 	}), nil
// }

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL not set")
	}

	log.Println("Running database migrations...")

	err := db.RunMigrations(connString)
	if err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}
	mux := http.NewServeMux()

	database, err := db.NewDatabase(context.Background(), connString)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	apptPath, apptHandler := protoconnect.NewAppointmentServiceHandler(&AppointmentServer{Storage: database})
	userPath, userHandler := protoconnect.NewUserServiceHandler(&UserServer{Storage: database})

	mux.Handle(apptPath, apptHandler)
	mux.Handle(userPath, userHandler)

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

	httpErr := http.ListenAndServe(addr, h2c.NewHandler(c.Handler(mux), &http2.Server{}))

	if httpErr != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
