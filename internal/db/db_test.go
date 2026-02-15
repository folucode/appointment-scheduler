package db

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	pb "github.com/folucode/appointment-scheduler/proto"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createTestDB(t *testing.T) *Database {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	absPath, err := filepath.Abs("../migrations")
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.New("file://"+absPath, connStr)
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}

	db, err := NewDatabase(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestFindUserByEmail(t *testing.T) {
	db := createTestDB(t)
	ctx := context.Background()

	t.Run("returns ErrUserNotFound when user doesn't exist", func(t *testing.T) {
		user, err := db.FindUserByEmail(ctx, "nonexistent@test.com")
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("successfully finds an existing user", func(t *testing.T) {
		_, err := db.Pool.Exec(ctx, "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
			uuid.NewString(), "Test User", "found@test.com")
		assert.NoError(t, err)

		user, err := db.FindUserByEmail(ctx, "found@test.com")
		assert.NoError(t, err)
		assert.Equal(t, "Test User", user.Name)
	})
}

func TestAppointmentOverlap(t *testing.T) {
	db := createTestDB(t)
	ctx := context.Background()

	userID := uuid.NewString()

	_, err := db.Pool.Exec(ctx, `
    INSERT INTO users (id, name, email)
    VALUES ($1, $2, $3)`, userID, "Test User", "test@example.com")

	require.NoError(t, err)

	t.Run("returns conflict error when appointments overlap", func(t *testing.T) {
		appt1 := &pb.Appointment{
			Id:          uuid.NewString(),
			UserId:      userID,
			Title:       "Test title",
			Description: "Test description",
			Date:        timestamppb.Now(),
			ContactInformation: &pb.ContactInformation{
				Name:  "Test",
				Email: "test@example.com",
			},
			StartTime: timestamppb.New(time.Now()),
			EndTime:   timestamppb.New(time.Now().Add(time.Hour)),
		}
		err := db.CreateAppointment(ctx, appt1)
		assert.NoError(t, err)

		appt2 := &pb.Appointment{
			Id:          uuid.NewString(),
			UserId:      userID,
			Title:       "Test title",
			Description: "Test description",
			Date:        timestamppb.Now(),
			ContactInformation: &pb.ContactInformation{
				Name:  "Test",
				Email: "test@example.com",
			},
			StartTime: timestamppb.New(time.Now().Add(30 * time.Minute)),
			EndTime:   timestamppb.New(time.Now().Add(90 * time.Minute)),
		}

		err = db.CreateAppointment(ctx, appt2)

		if assert.Error(t, err, "The DB should have rejected this overlap!") {
			assert.Contains(t, err.Error(), "conflict")
		} else {
			t.Log("Error was nil - the database allowed the overlapping appointment!")
		}
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflict")
	})
}
