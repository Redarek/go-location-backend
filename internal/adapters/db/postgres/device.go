package postgres

import (
	"context"

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
