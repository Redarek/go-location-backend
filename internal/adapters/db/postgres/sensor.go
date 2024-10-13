package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
	"location-backend/internal/domain/service"
)

type sensorRepo struct {
	pool *pgxpool.Pool
}

func NewSensorRepo(pool *pgxpool.Pool) *sensorRepo {
	return &sensorRepo{pool: pool}
}

type tmpSensorRadio struct {
	ID            *uuid.UUID          `db:"id"`
	Number        *int                `db:"number"`
	Channel       *int                `db:"channel"`
	Channel2      *int                `db:"channel2"`
	ChannelWidth  *string             `db:"channel_width"`
	WiFi          *string             `db:"wifi"`
	Power         *int                `db:"power"`
	Bandwidth     *string             `db:"bandwidth"`
	GuardInterval *int                `db:"guard_interval"`
	IsActive      *bool               `db:"is_active"`
	SensorID      *uuid.UUID          `db:"sensor_id"`
	CreatedAt     *pgtype.Timestamptz `db:"created_at"`
	UpdatedAt     *pgtype.Timestamptz `db:"updated_at"`
	DeletedAt     *pgtype.Timestamptz `db:"deleted_at"`
}

func (r *sensorRepo) Create(ctx context.Context, createSensorDTO *dto.CreateSensorDTO) (sensorID uuid.UUID, err error) {
	query := `INSERT INTO sensors (
			name, 
			color,
			x, y, z,
			mac,
			ip,
			rx_ant_gain,
			hor_rotation_offset, vert_rotation_offset,
			correction_factor_24, correction_factor_5, correction_factor_6,
			is_virtual,
			diagram,
			sensor_type_id,
			floor_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createSensorDTO.Name,
		createSensorDTO.Color,
		createSensorDTO.X, createSensorDTO.Y, createSensorDTO.Z,
		createSensorDTO.MAC,
		createSensorDTO.IP,
		createSensorDTO.RxAntGain,
		createSensorDTO.HorRotationOffset, createSensorDTO.VertRotationOffset,
		createSensorDTO.CorrectionFactor24, createSensorDTO.CorrectionFactor5, createSensorDTO.CorrectionFactor6,
		createSensorDTO.IsVirtual,
		createSensorDTO.Diagram,
		createSensorDTO.SensorTypeID,
		createSensorDTO.FloorID,
	)
	err = row.Scan(&sensorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan sensor")
		return
	}

	return
}

func (r *sensorRepo) GetOne(ctx context.Context, sensorID uuid.UUID) (sensor *entity.Sensor, err error) {
	query := `SELECT 
			id, 
			name, 
			color,
			x, y, z,
			mac,
			ip,
			rx_ant_gain,
			hor_rotation_offset, vert_rotation_offset,
			correction_factor_24, correction_factor_5, correction_factor_6,
			is_virtual,
			diagram,
			sensor_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM sensors WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, sensorID)

	sensor = &entity.Sensor{}
	err = row.Scan(
		&sensor.ID,
		&sensor.Name,
		&sensor.Color,
		&sensor.X, &sensor.Y, &sensor.Z,
		&sensor.MAC,
		&sensor.IP,
		&sensor.RxAntGain,
		&sensor.HorRotationOffset, &sensor.VertRotationOffset,
		&sensor.CorrectionFactor24, &sensor.CorrectionFactor5, &sensor.CorrectionFactor6,
		&sensor.IsVirtual,
		&sensor.Diagram,
		&sensor.SensorTypeID,
		&sensor.FloorID,
		&sensor.CreatedAt, &sensor.UpdatedAt, &sensor.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no sensor found with ID %v", sensorID)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve sensor")
		return
	}
	log.Debug().Msgf("retrieved sensor: %v", sensor)
	return
}

func (r *sensorRepo) GetOneByMAC(ctx context.Context, mac string) (sensor *entity.Sensor, err error) {
	query := `SELECT 
			id, 
			name, 
			color,
			x, y, z,
			mac,
			ip,
			rx_ant_gain,
			hor_rotation_offset, vert_rotation_offset,
			correction_factor_24, correction_factor_5, correction_factor_6,
			is_virtual,
			diagram,
			sensor_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM sensors WHERE mac = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, mac)

	sensor = &entity.Sensor{}
	err = row.Scan(
		&sensor.ID,
		&sensor.Name,
		&sensor.Color,
		&sensor.X, &sensor.Y, &sensor.Z,
		&sensor.MAC,
		&sensor.IP,
		&sensor.RxAntGain,
		&sensor.HorRotationOffset, &sensor.VertRotationOffset,
		&sensor.CorrectionFactor24, &sensor.CorrectionFactor5, &sensor.CorrectionFactor6,
		&sensor.IsVirtual,
		&sensor.Diagram,
		&sensor.SensorTypeID,
		&sensor.FloorID,
		&sensor.CreatedAt, &sensor.UpdatedAt, &sensor.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no sensor found with MAC %v", mac)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve sensor")
		return
	}
	log.Debug().Msgf("retrieved sensor: %v", sensor)
	return
}

func (r *sensorRepo) GetOneDetailed(ctx context.Context, sensorID uuid.UUID) (sensorDetailed *entity.SensorDetailed, err error) {
	query := `SELECT
			s.id,
			s.name, 
			s.color,
			s.x, s.y, s.z,
			s.mac,
			s.ip,
			s.rx_ant_gain, 
			s.hor_rotation_offset, s.vert_rotation_offset,
			s.correction_factor_24, s.correction_factor_5, s.correction_factor_6,
			s.is_virtual,
			s.diagram,
			s.sensor_type_id,
			s.floor_id,
			s.created_at, s.updated_at, s.deleted_at,

			st.id,
			st.name,
			st.model,
			st.color,
			st.z,
			st.is_virtual,
			st.site_id,
			st.created_at, st.updated_at, st.deleted_at,

			r.id,
			r.number,
			r.channel, r.channel2,
			r.channel_width,
			r.wifi,
			r.power,
			r.bandwidth,
			r.guard_interval,
			r.is_active,
			r.sensor_id,
			r.created_at, r.updated_at, r.deleted_at
		FROM sensors s
		JOIN sensor_types st ON s.sensor_type_id = st.id AND st.deleted_at IS NULL
		LEFT JOIN sensor_radios r ON s.id = r.sensor_id AND r.deleted_at IS NULL
		WHERE s.id = $1 AND s.deleted_at IS NULL`
	rows, err := r.pool.Query(ctx, query, sensorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve sensor")
		return
	}
	defer rows.Close()

	sensorDetailed = &entity.SensorDetailed{}
	sensorType := entity.SensorType{}
	i := 0

	for rows.Next() {
		i++
		tmpRadio := tmpSensorRadio{}

		err = rows.Scan(
			&sensorDetailed.ID,
			&sensorDetailed.Name,
			&sensorDetailed.Color,
			&sensorDetailed.X, &sensorDetailed.Y, &sensorDetailed.Z,
			&sensorDetailed.MAC,
			&sensorDetailed.IP,
			&sensorDetailed.RxAntGain,
			&sensorDetailed.HorRotationOffset, &sensorDetailed.VertRotationOffset,
			&sensorDetailed.CorrectionFactor24, &sensorDetailed.CorrectionFactor5, &sensorDetailed.CorrectionFactor6,
			&sensorDetailed.IsVirtual,
			&sensorDetailed.Diagram,
			&sensorDetailed.SensorTypeID,
			&sensorDetailed.FloorID,
			&sensorDetailed.CreatedAt, &sensorDetailed.UpdatedAt, &sensorDetailed.DeletedAt,

			&sensorType.ID,
			&sensorType.Name,
			&sensorType.Model,
			&sensorType.Color,
			&sensorType.Z,
			&sensorType.IsVirtual,
			&sensorType.SiteID,
			&sensorType.CreatedAt, &sensorType.UpdatedAt, &sensorType.DeletedAt,

			&tmpRadio.ID,
			&tmpRadio.Number,
			&tmpRadio.Channel, &tmpRadio.Channel2,
			&tmpRadio.ChannelWidth,
			&tmpRadio.WiFi,
			&tmpRadio.Power,
			&tmpRadio.Bandwidth,
			&tmpRadio.GuardInterval,
			&tmpRadio.IsActive,
			&tmpRadio.SensorID,
			&tmpRadio.CreatedAt, &tmpRadio.UpdatedAt, &tmpRadio.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan sensor and related data")
			return
		}
		sensorDetailed.SensorType = sensorType

		if tmpRadio.ID != nil {
			radio := &entity.SensorRadio{
				ID:            *tmpRadio.ID,
				Number:        *tmpRadio.Number,
				Channel:       *tmpRadio.Channel,
				Channel2:      tmpRadio.Channel2,
				ChannelWidth:  *tmpRadio.ChannelWidth,
				WiFi:          *tmpRadio.WiFi,
				Power:         *tmpRadio.Power,
				Bandwidth:     *tmpRadio.Bandwidth,
				GuardInterval: *tmpRadio.GuardInterval,
				IsActive:      *tmpRadio.IsActive,
				SensorID:      *tmpRadio.SensorID,
				CreatedAt:     *tmpRadio.CreatedAt,
				UpdatedAt:     *tmpRadio.UpdatedAt,
				DeletedAt:     tmpRadio.DeletedAt,
			}
			sensorDetailed.Radios = append(sensorDetailed.Radios, radio)
		}
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	if i == 0 {
		log.Debug().Msg("sensor was not found")
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved sensor with detailed info: %v", sensorDetailed)
	return
}

func (r *sensorRepo) GetAll(ctx context.Context, floorID uuid.UUID, limit, offset int) (sensors []*entity.Sensor, err error) {
	query := `SELECT 
			id, 
			name, 
			color,
			x, y, z,
			mac,
			ip,
			rx_ant_gain,
			hor_rotation_offset, vert_rotation_offset,
			correction_factor_24, correction_factor_5, correction_factor_6,
			is_virtual,
			diagram,
			sensor_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM sensors 
		WHERE floor_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, floorID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve sensors")
		return
	}
	defer rows.Close()

	for rows.Next() {
		sensor := &entity.Sensor{}
		err = rows.Scan(
			&sensor.ID,
			&sensor.Name,
			&sensor.Color,
			&sensor.X, &sensor.Y, &sensor.Z,
			&sensor.MAC,
			&sensor.IP,
			&sensor.RxAntGain,
			&sensor.HorRotationOffset, &sensor.VertRotationOffset,
			&sensor.CorrectionFactor24, &sensor.CorrectionFactor5, &sensor.CorrectionFactor6,
			&sensor.IsVirtual,
			&sensor.Diagram,
			&sensor.SensorTypeID,
			&sensor.FloorID,
			&sensor.CreatedAt, &sensor.UpdatedAt, &sensor.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan sensor")
			return
		}
		sensors = append(sensors, sensor)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(sensors)
	if length == 0 {
		log.Info().Msgf("sensors for floor ID %v were not found", floorID)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d sensors", length)
	return
}

func (r *sensorRepo) GetAllDetailed(ctx context.Context, floorID uuid.UUID, limit, offset int) (sensorsDetailed []*entity.SensorDetailed, err error) {
	query := `SELECT 
			s.id, 
			s.name, 
			s.color,
			s.x, s.y, s.z,
			s.mac,
			s.ip,
			s.rx_ant_gain,
			s.hor_rotation_offset, s.vert_rotation_offset,
			s.correction_factor_24, s.correction_factor_5, s.correction_factor_6,
			s.is_virtual,
			s.diagram,
			s.sensor_type_id,
			s.floor_id,
			s.created_at, s.updated_at, s.deleted_at,
			
			st.id, 
			st.name, 
			st.model,
			st.color, 
			st.z,
			st.is_virtual,
			st.site_id, 
			st.created_at, st.updated_at, st.deleted_at, 
			
			r.id, 
			r.number, 
			r.channel, r.channel2,
			r.channel_width,
			r.wifi, 
			r.power, 
			r.bandwidth, 
			r.guard_interval, 
			r.is_active, 
			r.sensor_id,
			r.created_at, r.updated_at, r.deleted_at
		FROM sensors s
		JOIN sensor_types st ON s.sensor_type_id = st.id AND st.deleted_at IS NULL
		LEFT JOIN sensor_radios r ON s.id = r.sensor_id AND r.deleted_at IS NULL
		WHERE s.floor_id = $1 AND s.deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, floorID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve sensor")
		return
	}
	defer rows.Close()

	apdMap := make(map[uuid.UUID]*entity.SensorDetailed)
	i := 0

	for rows.Next() {
		i++
		sensorDetailed := &entity.SensorDetailed{}
		sensorType := entity.SensorType{}
		tmpRadio := tmpSensorRadio{}

		err = rows.Scan(
			&sensorDetailed.ID,
			&sensorDetailed.Name,
			&sensorDetailed.Color,
			&sensorDetailed.X, &sensorDetailed.Y, &sensorDetailed.Z,
			&sensorDetailed.MAC,
			&sensorDetailed.IP,
			&sensorDetailed.RxAntGain,
			&sensorDetailed.HorRotationOffset, &sensorDetailed.VertRotationOffset,
			&sensorDetailed.CorrectionFactor24, &sensorDetailed.CorrectionFactor5, &sensorDetailed.CorrectionFactor6,
			&sensorDetailed.IsVirtual,
			&sensorDetailed.Diagram,
			&sensorDetailed.SensorTypeID,
			&sensorDetailed.FloorID,
			&sensorDetailed.CreatedAt, &sensorDetailed.UpdatedAt, &sensorDetailed.DeletedAt,

			&sensorType.ID,
			&sensorType.Name,
			&sensorType.Model,
			&sensorType.Color,
			&sensorType.Z,
			&sensorType.IsVirtual,
			&sensorType.SiteID,
			&sensorType.CreatedAt, &sensorType.UpdatedAt, &sensorType.DeletedAt,

			&tmpRadio.ID,
			&tmpRadio.Number,
			&tmpRadio.Channel, &tmpRadio.Channel2,
			&tmpRadio.ChannelWidth,
			&tmpRadio.WiFi,
			&tmpRadio.Power,
			&tmpRadio.Bandwidth,
			&tmpRadio.GuardInterval,
			&tmpRadio.IsActive,
			&tmpRadio.SensorID,
			&tmpRadio.CreatedAt, &tmpRadio.UpdatedAt, &tmpRadio.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan sensor detailed and related data")
			return
		}

		var radio *entity.SensorRadio
		if tmpRadio.ID != nil {
			radio = &entity.SensorRadio{
				ID:            *tmpRadio.ID,
				Number:        *tmpRadio.Number,
				Channel:       *tmpRadio.Channel,
				Channel2:      tmpRadio.Channel2,
				ChannelWidth:  *tmpRadio.ChannelWidth,
				WiFi:          *tmpRadio.WiFi,
				Power:         *tmpRadio.Power,
				Bandwidth:     *tmpRadio.Bandwidth,
				GuardInterval: *tmpRadio.GuardInterval,
				IsActive:      *tmpRadio.IsActive,
				SensorID:      *tmpRadio.SensorID,
				CreatedAt:     *tmpRadio.CreatedAt,
				UpdatedAt:     *tmpRadio.UpdatedAt,
				DeletedAt:     tmpRadio.DeletedAt,
			}
		}

		existingAP, exists := apdMap[sensorDetailed.ID]
		if exists {
			// If sensor is already in the map, append the new radio to its list
			if radio != nil {
				existingAP.Radios = append(existingAP.Radios, radio)
			}
		} else {
			// If it's a new sensor, initialize and add to map
			sensorDetailed.SensorType = sensorType
			if radio != nil {
				sensorDetailed.Radios = append(sensorDetailed.Radios, radio)
			}
			apdMap[sensorDetailed.ID] = sensorDetailed
		}
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	if i == 0 {
		log.Debug().Msg("sensors detailed were not found")
		return nil, service.ErrNotFound
	}

	// Convert map to slice
	for _, apd := range apdMap {
		sensorsDetailed = append(sensorsDetailed, apd)
	}

	log.Debug().Msgf("retrieved sensors with detailed info: %v", sensorsDetailed)
	return
}

func (r *sensorRepo) Update(ctx context.Context, updateSensorDTO *dto.PatchUpdateSensorDTO) (err error) {
	query := "UPDATE sensors SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateSensorDTO.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, updateSensorDTO.Name)
		paramID++
	}
	if updateSensorDTO.Color != nil {
		updates = append(updates, fmt.Sprintf("color = $%d", paramID))
		params = append(params, updateSensorDTO.Color)
		paramID++
	}
	if updateSensorDTO.X != nil {
		updates = append(updates, fmt.Sprintf("x = $%d", paramID))
		params = append(params, updateSensorDTO.X)
		paramID++
	}
	if updateSensorDTO.Y != nil {
		updates = append(updates, fmt.Sprintf("y = $%d", paramID))
		params = append(params, updateSensorDTO.Y)
		paramID++
	}
	if updateSensorDTO.Z != nil {
		updates = append(updates, fmt.Sprintf("z = $%d", paramID))
		params = append(params, updateSensorDTO.Z)
		paramID++
	}
	if updateSensorDTO.MAC != nil {
		updates = append(updates, fmt.Sprintf("mac = $%d", paramID))
		params = append(params, updateSensorDTO.MAC)
		paramID++
	}
	if updateSensorDTO.IP != nil {
		updates = append(updates, fmt.Sprintf("ip = $%d", paramID))
		params = append(params, updateSensorDTO.IP)
		paramID++
	}
	if updateSensorDTO.RxAntGain != nil {
		updates = append(updates, fmt.Sprintf("rx_ant_gain = $%d", paramID))
		params = append(params, updateSensorDTO.RxAntGain)
		paramID++
	}
	if updateSensorDTO.HorRotationOffset != nil {
		updates = append(updates, fmt.Sprintf("hor_rotation_offset = $%d", paramID))
		params = append(params, updateSensorDTO.HorRotationOffset)
		paramID++
	}
	if updateSensorDTO.VertRotationOffset != nil {
		updates = append(updates, fmt.Sprintf("vert_rotation_offset = $%d", paramID))
		params = append(params, updateSensorDTO.VertRotationOffset)
		paramID++
	}
	if updateSensorDTO.CorrectionFactor24 != nil {
		updates = append(updates, fmt.Sprintf("correction_factor_24 = $%d", paramID))
		params = append(params, updateSensorDTO.CorrectionFactor24)
		paramID++
	}
	if updateSensorDTO.CorrectionFactor5 != nil {
		updates = append(updates, fmt.Sprintf("correction_factor_5 = $%d", paramID))
		params = append(params, updateSensorDTO.CorrectionFactor5)
		paramID++
	}
	if updateSensorDTO.CorrectionFactor6 != nil {
		updates = append(updates, fmt.Sprintf("correction_factor_6 = $%d", paramID))
		params = append(params, updateSensorDTO.CorrectionFactor6)
		paramID++
	}
	if updateSensorDTO.Diagram != nil {
		updates = append(updates, fmt.Sprintf("diagram = $%d", paramID))
		params = append(params, updateSensorDTO.Diagram)
		paramID++
	}

	if updateSensorDTO.IsVirtual != nil {
		updates = append(updates, fmt.Sprintf("is_virtual = $%d", paramID))
		params = append(params, updateSensorDTO.IsVirtual)
		paramID++
	}
	if updateSensorDTO.SensorTypeID != nil {
		updates = append(updates, fmt.Sprintf("sensor_type_id = $%d", paramID))
		params = append(params, updateSensorDTO.SensorTypeID)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return service.ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateSensorDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor found with the ID: %v", updateSensorDTO.ID)
		return service.ErrNotFound
	}

	return
}

// Checks if the sensor has been soft deleted
func (r *sensorRepo) IsSensorSoftDeleted(ctx context.Context, sensorID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM sensors WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, sensorID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no sensor found with UUID %v", sensorID)
			return false, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve sensor")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is sensor deleted: %v", isDeleted)
	return
}

func (r *sensorRepo) SoftDelete(ctx context.Context, sensorID uuid.UUID) (err error) {
	query := `UPDATE sensors SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, sensorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete sensor")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor found with the ID: %v", sensorID)
		return service.ErrNotFound
	}

	log.Debug().Msg("sensor deleted_at timestamp updated successfully")
	return
}

func (r *sensorRepo) Restore(ctx context.Context, sensorID uuid.UUID) (err error) {
	query := `UPDATE sensors SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, sensorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore sensor")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no sensor found with the UUID: %v", sensorID)
		return service.ErrNotFound
	}

	log.Debug().Msg("sensor deleted_at timestamp set NULL successfully")
	return
}
