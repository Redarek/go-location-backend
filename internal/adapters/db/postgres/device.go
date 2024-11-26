package postgres

import (
	"context"
	// "fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/service"
)

type deviceRepo struct {
	pool *pgxpool.Pool
}

func NewDeviceRepo(pool *pgxpool.Pool) *deviceRepo {
	return &deviceRepo{pool: pool}
}

// func (r *deviceRepo) Create(ctx context.Context, createDeviceDTO *dto.CreateDeviceDTO) (deviceID uuid.UUID, err error) {
// 	query := `INSERT INTO devices (
// 			name,
// 			description,
// 			country,
// 			city,
// 			address,
// 			site_id
// 		)
// 		VALUES ($1, $2, $3, $4, $5, $6)
// 		RETURNING id`
// 	row := r.pool.QueryRow(ctx, query,
// 		createDeviceDTO.Name,
// 		createDeviceDTO.Description,
// 		createDeviceDTO.Country,
// 		createDeviceDTO.City,
// 		createDeviceDTO.Address,
// 		createDeviceDTO.SiteID,
// 	)
// 	err = row.Scan(&deviceID)
// 	if err != nil {
// 		log.Error().Err(err).Msg("failed to scan device")
// 		return
// 	}

// 	return
// }

// TODO сделать нормальный search
func (r *deviceRepo) GetAll(ctx context.Context, mac string, floorID uuid.UUID, limit, offset int) (devices []*entity.Device, err error) {
	query := `SELECT 
			d.id, 
			d.mac,
			d.sensor_id,
			d.rssi,
			d.band,
			d.channel_width,
			d.last_contact_time 
		FROM devices d
		JOIN sensors s ON d.sensor_id = s.id AND s.deleted_at IS NULL
		WHERE d.mac = $1 AND s.floor_id = $2
		LIMIT NULLIF($3, 0) OFFSET $4`
	rows, err := r.pool.Query(ctx, query,
		mac,
		floorID,
		limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve devices")
		return
	}
	defer rows.Close()

	for rows.Next() {
		device := &entity.Device{}
		err = rows.Scan(
			&device.ID,
			&device.MAC,
			&device.SensorID,
			&device.RSSI,
			&device.Band,
			&device.ChannelWidth,
			&device.LastContactTime,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan device")
			return
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(devices)
	if length == 0 {
		log.Info().Msgf("devices with MAC %s were not found", mac)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d devices", length)
	return
}

func (r *deviceRepo) GetAllDetailedByMAC(ctx context.Context, mac string, limit, offset int) (devicesDetailed []*entity.DeviceDetailed, err error) {
	query := `SELECT 
			d.id, 
			d.mac,
			d.sensor_id,
			d.rssi,
			d.band,
			d.channel_width,
			d.last_contact_time,
			s.floor_id 
		FROM devices d
		JOIN sensors s ON d.sensor_id = s.id AND s.deleted_at IS NULL
		WHERE d.mac = $1
		LIMIT NULLIF($2, 0) OFFSET $3`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve devices")
		return
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		i++
		deviceDetailed := &entity.DeviceDetailed{}

		err = rows.Scan(
			&deviceDetailed.ID,
			&deviceDetailed.MAC,
			&deviceDetailed.SensorID,
			&deviceDetailed.RSSI,
			&deviceDetailed.Band,
			&deviceDetailed.ChannelWidth,
			&deviceDetailed.LastContactTime,
			&deviceDetailed.FloorID,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan device detailed and related data")
			return
		}

		devicesDetailed = append(devicesDetailed, deviceDetailed)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	if i == 0 {
		log.Debug().Msg("devices detailed were not found")
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved devices with detailed info: %v", devicesDetailed)
	return
}

// func (r *deviceRepo) Search(ctx context.Context, filter entity.SearchParameters, limit, offset int) (devices []*entity.Device, err error) {
// 	query := `SELECT
// 	d.id,
// 	d.mac,
// 	d.sensor_id,
// 	d.rssi,
// 	d.band,
// 	d.channel_width,
// 	d.last_contact_time
// FROM devices d
// JOIN sensors s ON d.sensor_id = s.id AND s.deleted_at IS NULL
// WHERE d.mac = $1 AND s.floor_id = $2
// LIMIT NULLIF($3, 0) OFFSET $4`
// 	args := []interface{}{}
// 	argIndex := 1 // Нумерация параметров начинается с 1 в pgx

// 	if filter.ID != nil {
// 		query += fmt.Sprintf(" AND id = $%d", argIndex)
// 		args = append(args, *filter.ID)
// 		argIndex++
// 	}
// 	if filter.Name != nil {
// 		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
// 		args = append(args, "%"+*filter.Name+"%")
// 		argIndex++
// 	}
// 	if filter.Email != nil {
// 		query += fmt.Sprintf(" AND email = $%d", argIndex)
// 		args = append(args, *filter.Email)
// 		argIndex++
// 	}
// 	if filter.Age != nil {
// 		query += fmt.Sprintf(" AND age = $%d", argIndex)
// 		args = append(args, *filter.Age)
// 		argIndex++
// 	}

// 	rows, err := r.pool.Query(ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute query: %w", err)
// 	}
// 	defer rows.Close()

// 	var users []User
// 	for rows.Next() {
// 		var user User
// 		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %w", err)
// 		}
// 		users = append(users, user)
// 	}

// 	if rows.Err() != nil {
// 		return nil, fmt.Errorf("error iterating rows: %w", rows.Err())
// 	}

// 	return users, nil
// }
