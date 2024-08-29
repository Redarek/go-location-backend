package user

import (
	"context"
)

type service struct {
	repository Repository
}

// func NewService(repository Repository) user.Service {
// 	return &service{repository: repository}
// }

func (s *service) GetUserByUsername(ctx context.Context, username string) (user *User, err error) {
	return s.repository.GetOne(username)
}

// // CreateUser creates a user
// func (s *service) CreateUser(username, password string) (id uuid.UUID, err error) {
// 	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
// 	row := s.Pool.QueryRow(context.Background(), query, username, password)
// 	err = row.Scan(&id)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Failed to create user")
// 	}
// 	return
// }

// // GetUserByUsername retrieves a user
// func (p *postgres) GetUserByUsername(username string) (user User, err error) {
// 	query := `SELECT id, username, password, created_at, updated_at, deleted_at FROM users WHERE username = $1 AND deleted_at IS NULL`
// 	row := p.Pool.QueryRow(context.Background(), query, username)
// 	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			log.Error().Err(err).Msgf("No user found with username %v", username)
// 			return
// 		}
// 		log.Error().Err(err).Msg("Failed to retrieve user")
// 		return
// 	}
// 	log.Debug().Msgf("Retrieved user: %v", user)
// 	return
// }
