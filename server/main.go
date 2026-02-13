package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/folucode/appointment-scheduler/internal/db"
	pb "github.com/folucode/appointment-scheduler/proto"
	protoconnect "github.com/folucode/appointment-scheduler/proto/protoconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	user, err := s.Storage.FindUserByEmail(req.Msg.ContactInformation.Email)

	if err != nil {
		if !errors.Is(err, db.ErrUserNotFound) {
			return nil, err
		}

		createdUser, err := s.Storage.CreateUser(db.User{
			ID:    uuid.NewString(),
			Name:  req.Msg.ContactInformation.Name,
			Email: req.Msg.ContactInformation.Email,
		})
		user = &createdUser
		log.Printf("user data created: %+v", user)
		if err != nil {
			return nil, err
		}
	}

	newAppt := db.Appointment{
		ID:          uuid.NewString(),
		Description: req.Msg.Description,
		UserID:      user.ID,
		ContactInformation: db.Contact{
			Name:  req.Msg.ContactInformation.Name,
			Email: req.Msg.ContactInformation.Email,
		},
		StartTime: req.Msg.StartTime.AsTime(),
		EndTime:   req.Msg.EndTime.AsTime(),
	}

	err = s.Storage.CreateAppointment(newAppt)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&pb.Appointment{
		Id:          newAppt.ID,
		Description: newAppt.Description,
		UserId:      newAppt.UserID,
		ContactInformation: &pb.ContactInformation{
			Name:  newAppt.ContactInformation.Name,
			Email: newAppt.ContactInformation.Email,
		},
		StartTime: timestamppb.New(newAppt.StartTime),
		EndTime:   timestamppb.New(newAppt.EndTime),
	})

	return res, nil
}

func (s *AppointmentServer) GetUserAppointments(
	ctx context.Context,
	req *connect.Request[pb.GetUserAppointmentRequest],
) (*connect.Response[pb.GetUserAppointmentResponse], error) {
	log.Printf("Incoming Request to get user appointments: %+v", req.Msg)

	if req.Msg.UserId == "" {
		return connect.NewResponse(&pb.GetUserAppointmentResponse{
			Appointment: []*pb.Appointment{},
		}), nil
	}

	data, err := s.Storage.GetAppointments(req.Msg.UserId)
	if err != nil {
		return nil, err
	}

	result := []*pb.Appointment{}

	for _, appt := range data {
		result = append(result, &pb.Appointment{
			Id:          appt.ID,
			Description: appt.Description,
			ContactInformation: &pb.ContactInformation{
				Name:  appt.ContactInformation.Name,
				Email: appt.ContactInformation.Email,
			},
			StartTime: timestamppb.New(appt.StartTime),
			EndTime:   timestamppb.New(appt.EndTime),
		})
	}

	res := connect.NewResponse(&pb.GetUserAppointmentResponse{
		Appointment: result,
	})

	return res, nil
}

func main() {
	mux := http.NewServeMux()

	database := db.NewDatabase("data/db.json")

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

	err := http.ListenAndServe(addr, h2c.NewHandler(c.Handler(mux), &http2.Server{}))

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
