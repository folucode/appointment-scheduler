package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/folucode/appointment-scheduler/proto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (db *Database) CreateAppointment(ctx context.Context, appt *pb.Appointment) error {
	query := `
        INSERT INTO appointments (id, user_id, contact_name, contact_email, start_time, end_time, title, description, date)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := db.Pool.Exec(ctx, query,
		appt.Id,
		appt.UserId,
		appt.ContactInformation.Name,
		appt.ContactInformation.Email,
		appt.StartTime.AsTime(),
		appt.EndTime.AsTime(),
		appt.Title,
		appt.Description,
		appt.Date.AsTime(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23P01" {
				return errors.New("conflict: this time slot overlaps with an existing appointment")
			}
		}
		return err
	}

	return nil
}

func (db *Database) GetAppointment(ctx context.Context, id string) (*pb.Appointment, error) {
	query := `
		SELECT id, user_id, title, description, date, contact_name, contact_email, start_time, end_time 
		FROM appointments WHERE id = $1`

	var appt pb.Appointment
	var contact pb.ContactInformation
	var date, start, end time.Time

	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&appt.Id,
		&appt.UserId,
		&appt.Title,
		&appt.Description,
		&contact.Name,
		&contact.Email,
		&date,
		&start,
		&end,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}

	appt.ContactInformation = &contact
	appt.StartTime = timestamppb.New(start)
	appt.EndTime = timestamppb.New(end)
	appt.Date = timestamppb.New(date)

	return &appt, nil
}

func (db *Database) GetAppointments(ctx context.Context, userId string) ([]*pb.Appointment, error) {
	query := `
        SELECT id, user_id, contact_name, contact_email, start_time, end_time, date, title, description 
        FROM appointments WHERE user_id = $1`

	rows, err := db.Pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*pb.Appointment
	for rows.Next() {
		var a pb.Appointment
		var c pb.ContactInformation
		var start, end, date time.Time

		err := rows.Scan(
			&a.Id,
			&a.UserId,
			&c.Name,
			&c.Email,
			&start,
			&end,
			&date,
			&a.Description,
			&a.Title,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		a.ContactInformation = &c

		a.StartTime = timestamppb.New(start)
		a.EndTime = timestamppb.New(end)
		a.Date = timestamppb.New(date)

		result = append(result, &a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *Database) DeleteAppointment(ctx context.Context, id string) (bool, error) {
	query := `DELETE FROM appointments WHERE id = $1`

	commandTag, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}
