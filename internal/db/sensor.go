package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	. "location-backend/internal/db/model"
)

type Diagram struct {
	Degree map[string]Degree `json:"degree"`
}

type Degree struct {
	HorGain  float64 `json:"hor_gain"`
	VertGain float64 `json:"vert_gain"`
}

// CreateSensor creates a sensor
func (p *postgres) CreateSensor(s *Sensor) (id uuid.UUID, err error) {
	query := `INSERT INTO sensors (
	name, 
	x, y, z, 
	mac, 
	ip, 
	alias, 
	interface_0, interface_1, interface_2, 
	rx_ant_gain, 
	hor_rotation_offset, vert_rotation_offset, 
	correction_factor_24, correction_factor_5, correction_factor_6, 
	is_virtual,
	diagram, 
	floor_id, 
	sensor_type_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query,
		s.Name,
		s.X, s.Y, s.Z,
		s.MAC,
		s.IP,
		s.Alias,
		s.Interface0, s.Interface1, s.Interface2,
		s.RxAntGain,
		s.HorRotationOffset, s.VertRotationOffset,
		s.CorrectionFactor24, s.CorrectionFactor5, s.CorrectionFactor6,
		s.IsVirtual,
		s.Diagram,
		s.FloorID,
		s.SensorTypeID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create sensor")
	}
	return
}

// GetSensor retrieves a sensor
func (p *postgres) GetSensor(sensorUUID uuid.UUID) (s *Sensor, err error) {
	query := `SELECT
			id,
			name,
			x, y, z,
			mac,
			ip,
			alias,
			interface_0, interface_1, interface_2,
			rx_ant_gain,
			hor_rotation_offset, vert_rotation_offset,
			correction_factor_24, correction_factor_5, correction_factor_6,
			is_virtual,
			diagram,
			sensor_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM sensors WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, sensorUUID)
	s = &Sensor{}
	err = row.Scan(&s.ID,
		&s.Name,
		&s.X, &s.Y, &s.Z,
		&s.MAC,
		&s.IP,
		&s.Alias,
		&s.Interface0, &s.Interface1, &s.Interface2,
		&s.RxAntGain,
		&s.HorRotationOffset, &s.VertRotationOffset,
		&s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6,
		&s.IsVirtual,
		&s.Diagram,
		&s.SensorTypeID,
		&s.FloorID,
		&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No sensor found with uuid %v", sensorUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve sensor")
		return
	}
	log.Debug().Msgf("Retrieved sensor: %v", s)
	return
}

// GetSensorDetailed retrieves a sensor detailed
func (p *postgres) GetSensorDetailed(sensorUUID uuid.UUID) (s *SensorDetailed, err error) {
	query := `SELECT 
		s.id, 
		s.name, 

		st.color, 

		s.x, s.y, s.z, 
		s.mac, 
		s.ip, 
		s.alias, 
		s.interface_0, s.interface_1, s.interface_2, 
		s.rx_ant_gain, 
		s.hor_rotation_offset, s.vert_rotation_offset, 
		s.correction_factor_24, s.correction_factor_5, s.correction_factor_6,
		s.is_virtual,
		s.diagram, 
		s.created_at, s.updated_at, s.deleted_at, 
		s.floor_id, 
		s.sensor_type_id
	FROM sensors s
	LEFT JOIN sensor_types st ON s.sensor_type_id = st.id AND st.deleted_at IS NULL
	WHERE s.id = $1 AND s.deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, sensorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve sensor detailed information")
		return
	}
	defer rows.Close()

	s = new(SensorDetailed)
	st := new(SensorType)

	for rows.Next() {
		err = rows.Scan(
			&s.ID,
			&s.Name,

			&st.Color,

			&s.X, &s.Y, &s.Z,
			&s.MAC,
			&s.IP,
			&s.Alias,
			&s.Interface0, &s.Interface1, &s.Interface2,
			&s.RxAntGain,
			&s.HorRotationOffset, &s.VertRotationOffset,
			&s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6,
			&s.IsVirtual,
			&s.Diagram,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
			&s.FloorID,
			&s.SensorTypeID,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan sensor detailed information")
			return
		}
		s.Color = st.Color
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved sensor with detailed info: %v", s)
	return
}

// IsSensorSoftDeleted checks if the sensor has been soft deleted
func (p *postgres) IsSensorSoftDeleted(sensorUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM sensors WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, sensorUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No sensor found with uuid %v", sensorUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve sensor")
		return
	}
	isDeleted = deletedAt.Valid
	log.Debug().Msgf("Is sensor deleted: %v", isDeleted)
	return
}

// GetSensors retrieves sensors
func (p *postgres) GetSensors(floorUUID uuid.UUID) (ss []*Sensor, err error) {
	query := `SELECT
			id,
			name,
			x, y, z,
			mac,
			ip,
			alias,
			interface_0, interface_1, interface_2,
			rx_ant_gain,
			hor_rotation_offset, vert_rotation_offset,
			correction_factor_24, correction_factor_5, correction_factor_6,
			is_virtual,
			diagram,
			floor_id,
			sensor_type_id,
			created_at, updated_at, deleted_at
		FROM sensors WHERE floor_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve sensors")
		return
	}
	defer rows.Close()

	var s *Sensor
	for rows.Next() {
		s = new(Sensor)
		err = rows.Scan(
			&s.ID,
			&s.Name,
			&s.X, &s.Y, &s.Z,
			&s.MAC,
			&s.IP,
			&s.Alias,
			&s.Interface0, &s.Interface1, &s.Interface2,
			&s.RxAntGain,
			&s.HorRotationOffset, &s.VertRotationOffset,
			&s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6,
			&s.IsVirtual,
			&s.Diagram,
			&s.FloorID,
			&s.SensorTypeID,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan sensor")
			return
		}
		ss = append(ss, s)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d sensors", len(ss))
	return
}

// GetSensorsDetailed retrieves sensors detailed
func (p *postgres) GetSensorsDetailed(floorUUID uuid.UUID) (ss []*SensorDetailed, err error) {
	query := `SELECT 
		s.id, 
		s.name, 
		
		st.color, 
		
		s.x, s.y, s.z, 
		s.mac, 
		s.ip, 
		s.alias, 
		s.interface_0, s.interface_1, s.interface_2, 
		s.rx_ant_gain, 
		s.hor_rotation_offset, s.vert_rotation_offset, 
		s.correction_factor_24, s.correction_factor_5, s.correction_factor_6, 
		s.is_virtual,
		s.diagram, 
		s.created_at, s.updated_at, s.deleted_at, 
		s.floor_id, 
		s.sensor_type_id
	FROM sensors s
	LEFT JOIN sensor_types st ON s.sensor_type_id = st.id AND st.deleted_at IS NULL
	WHERE s.floor_id = $1 AND s.deleted_at IS NULL
	GROUP BY s.id, s.name, st.color, s.x, s.y, s.z, s.mac, s.ip, s.alias, s.interface_0, s.interface_1, s.interface_2, s.rx_ant_gain, s.hor_rotation_offset, s.vert_rotation_offset, s.correction_factor_24, s.correction_factor_5, s.correction_factor_6, s.diagram, s.created_at, s.updated_at, s.deleted_at, s.floor_id, s.sensor_type_id`

	rows, err := p.Pool.Query(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve access points")
		return
	}
	defer rows.Close()

	for rows.Next() {
		s := new(SensorDetailed)
		st := new(SensorType)
		err = rows.Scan(
			&s.ID,
			&s.Name,

			&st.Color,

			&s.X, &s.Y, &s.Z,
			&s.MAC,
			&s.IP,
			&s.Alias,
			&s.Interface0, &s.Interface1, &s.Interface2,
			&s.RxAntGain,
			&s.HorRotationOffset, &s.VertRotationOffset,
			&s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6,
			&s.IsVirtual,
			&s.Diagram,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
			&s.FloorID,
			&s.SensorTypeID,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan sensor detailed information")
			return
		}
		s.Color = st.Color

		ss = append(ss, s)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved sensors with detailed info: %+v", ss)
	return
}

// SoftDeleteSensor soft delete a sensor
func (p *postgres) SoftDeleteSensor(sensorUUID uuid.UUID) (err error) {
	query := `UPDATE sensors SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, sensorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete sensor")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No sensor found with the uuid: %v", sensorUUID)
		return
	}
	log.Debug().Msg("Sensor deleted_at timestamp updated successfully")
	return
}

// RestoreSensor restore a sensor
func (p *postgres) RestoreSensor(sensorUUID uuid.UUID) (err error) {
	query := `UPDATE sensors SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, sensorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore sensor")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No sensor found with the uuid: %v", sensorUUID)
		return
	}
	log.Debug().Msg("Sensor deleted_at timestamp set null successfully")
	return
}

// PatchUpdateSensor updates only the specified fields of a sensor
func (p *postgres) PatchUpdateSensor(s *Sensor) (err error) {
	query := "UPDATE sensors SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if s.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, s.Name)
		paramID++
	}
	if s.X != nil {
		updates = append(updates, fmt.Sprintf("x = $%d", paramID))
		params = append(params, s.X)
		paramID++
	}
	if s.Y != nil {
		updates = append(updates, fmt.Sprintf("y = $%d", paramID))
		params = append(params, s.Y)
		paramID++
	}
	if s.Z != nil {
		updates = append(updates, fmt.Sprintf("z = $%d", paramID))
		params = append(params, s.Z)
		paramID++
	}
	if s.MAC != "" {
		updates = append(updates, fmt.Sprintf("mac = $%d", paramID))
		params = append(params, s.MAC)
		paramID++
	}
	if s.IP != "" {
		updates = append(updates, fmt.Sprintf("ip = $%d", paramID))
		params = append(params, s.IP)
		paramID++
	}
	if s.Alias != "" {
		updates = append(updates, fmt.Sprintf("alias = $%d", paramID))
		params = append(params, s.Alias)
		paramID++
	}
	if s.Interface0 != nil {
		updates = append(updates, fmt.Sprintf("interface_0 = $%d", paramID))
		params = append(params, s.Interface0)
		paramID++
	}
	if s.Interface1 != nil {
		updates = append(updates, fmt.Sprintf("interface_1 = $%d", paramID))
		params = append(params, s.Interface1)
		paramID++
	}
	if s.Interface2 != nil {
		updates = append(updates, fmt.Sprintf("interface_2 = $%d", paramID))
		params = append(params, s.Interface2)
		paramID++
	}
	if s.RxAntGain != nil {
		updates = append(updates, fmt.Sprintf("rx_ant_gain = $%d", paramID))
		params = append(params, s.RxAntGain)
		paramID++
	}
	if s.HorRotationOffset != nil {
		updates = append(updates, fmt.Sprintf("hor_rotation_offset = $%d", paramID))
		params = append(params, s.HorRotationOffset)
		paramID++
	}
	if s.VertRotationOffset != nil {
		updates = append(updates, fmt.Sprintf("vert_rotation_offset = $%d", paramID))
		params = append(params, s.VertRotationOffset)
		paramID++
	}
	if s.CorrectionFactor24 != nil {
		updates = append(updates, fmt.Sprintf("correction_factor_24 = $%d", paramID))
		params = append(params, s.CorrectionFactor24)
		paramID++
	}
	if s.CorrectionFactor5 != nil {
		updates = append(updates, fmt.Sprintf("correction_factor_5 = $%d", paramID))
		params = append(params, s.CorrectionFactor5)
		paramID++
	}
	if s.CorrectionFactor6 != nil {
		updates = append(updates, fmt.Sprintf("correction_factor_6 = $%d", paramID))
		params = append(params, s.CorrectionFactor6)
		paramID++
	}
	if s.Diagram != nil {
		updates = append(updates, fmt.Sprintf("is_virtual = $%d", paramID))
		params = append(params, s.IsVirtual)
		paramID++
	}
	if s.Diagram != nil {
		updates = append(updates, fmt.Sprintf("diagram = $%d", paramID))
		params = append(params, s.Diagram)
		paramID++
	}

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, s.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

	return
}
