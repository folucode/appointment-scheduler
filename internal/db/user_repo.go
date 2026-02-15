package db

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/folucode/appointment-scheduler/proto"
	"github.com/jackc/pgx/v5"
)

func (db *Database) FindUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	query := `SELECT id, name, email FROM users WHERE email=$1`

	var user pb.User

	err := db.Pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db *Database) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3) ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name RETURNING id, name, email`

	var createdUser pb.User

	err := db.Pool.QueryRow(ctx, query,
		user.Id,
		user.Name,
		user.Email,
	).Scan(
		&createdUser.Id,
		&createdUser.Name,
		&createdUser.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create and return user: %w", err)
	}

	return &createdUser, nil
}
