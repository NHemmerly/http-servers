// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const removeUsers = `-- name: RemoveUsers :exec
DELETE FROM users
`

func (q *Queries) RemoveUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, removeUsers)
	return err
}

const updateCreds = `-- name: UpdateCreds :one
UPDATE users
SET hashed_password = $1, email = $2, updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email, is_chirpy_red
`

type UpdateCredsParams struct {
	HashedPassword string
	Email          string
	ID             uuid.UUID
}

type UpdateCredsRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	IsChirpyRed bool
}

func (q *Queries) UpdateCreds(ctx context.Context, arg UpdateCredsParams) (UpdateCredsRow, error) {
	row := q.db.QueryRowContext(ctx, updateCreds, arg.HashedPassword, arg.Email, arg.ID)
	var i UpdateCredsRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}

const upgradeUser = `-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
`

func (q *Queries) UpgradeUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, upgradeUser, id)
	return err
}
