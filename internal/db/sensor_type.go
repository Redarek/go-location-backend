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

// CreateSensorType creates a sensor type
func (p *postgres) CreateSensorType(st *SensorType) (id uuid.UUID, err error) {
	query := `INSERT INTO sensor_types (
			name, 
			color, 
			alias, 
			interface_0, interface_1, interface_2, 
			rx_ant_gain, 
			hor_rotation_offset, vert_rotation_offset, 
			correction_factor_24, correction_factor_5, correction_factor_6, 
			diagram, 
			is_virtual,
			site_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query,
		st.Name,
		st.Color,
		st.Alias,
		st.Interface0, st.Interface1, st.Interface2,
		st.RxAntGain,
		st.HorRotationOffset, st.VertRotationOffset,
		st.CorrectionFactor24, st.CorrectionFactor5, st.CorrectionFactor6,
		st.Diagram,
		st.IsVirtual,
		st.SiteID,
	)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create sensor type")
	}
	return
}

// GetSensorType retrieves a sensor
func (p *postgres) GetSensorType(sensorTypeUUID uuid.UUID) (st *SensorType, err error) {
	query := `SELECT
			id,
			name,
			color,
			alias,
			interface_0, interface_1, interface_2,
			rx_ant_gain,
			hor_rotation_offset, vert_rotation_offset,
			correction_factor_24, correction_factor_5, correction_factor_6,
			diagram,
			is_virtual,
			site_id,
			created_at, updated_at, deleted_at
		FROM sensor_types WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, sensorTypeUUID)
	st = &SensorType{}
	err = row.Scan(
		&st.ID,
		&st.Name,
		&st.Color,
		&st.Alias,
		&st.Interface0, &st.Interface1, &st.Interface2,
		&st.RxAntGain,
		&st.HorRotationOffset, &st.VertRotationOffset,
		&st.CorrectionFactor24, &st.CorrectionFactor5, &st.CorrectionFactor6,
		&st.Diagram,
		&st.IsVirtual,
		&st.SiteID,
		&st.CreatedAt, &st.UpdatedAt, &st.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No sensor type found with uuid %v", sensorTypeUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve sensor type")
		return
	}
	log.Debug().Msgf("Retrieved sensor type: %v", st)
	return
}

// IsSensorTypeSoftDeleted checks if the sensor type has been soft deleted
func (p *postgres) IsSensorTypeSoftDeleted(sensorTypeUUID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM sensor_types WHERE id = $1`
	row := p.Pool.QueryRow(context.Background(), query, sensorTypeUUID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No sensor type found with uuid %v", sensorTypeUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve sensor type")
		return
	}
	isDeleted = deletedAt.Valid
	log.Debug().Msgf("Is sensor type deleted: %v", isDeleted)
	return
}

// GetSensorTypes retrieves sensor types
func (p *postgres) GetSensorTypes(siteUUID uuid.UUID) (sts []*SensorType, err error) {
	query := `SELECT
			id,
			name,
			color,
			alias,
			interface_0, interface_1, interface_2,
			rx_ant_gain,
			hor_rotation_offset, vert_rotation_offset,
			correction_factor_24, correction_factor_5, correction_factor_6,
			diagram,
			is_virtual,
			site_id,
			created_at, updated_at, deleted_at
		FROM sensor_types WHERE site_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve sensor types")
		return
	}
	defer rows.Close()

	var st *SensorType
	for rows.Next() {
		st = new(SensorType)
		err = rows.Scan(
			&st.ID,
			&st.Name,
			&st.Color,
			&st.Alias,
			&st.Interface0, &st.Interface1, &st.Interface2,
			&st.RxAntGain,
			&st.HorRotationOffset, &st.VertRotationOffset,
			&st.CorrectionFactor24, &st.CorrectionFactor5, &st.CorrectionFactor6,
			&st.Diagram,
			&st.IsVirtual,
			&st.SiteID,
			&st.CreatedAt, &st.UpdatedAt, &st.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan sensor type")
			return
		}
		sts = append(sts, st)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d sensor types", len(sts))
	return
}

// SoftDeleteSensorType soft delete a sensor type
func (p *postgres) SoftDeleteSensorType(sensorTypeUUID uuid.UUID) (err error) {
	query := `UPDATE sensor_types SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, sensorTypeUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to soft delete sensor type")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No sensor type found with the uuid: %v", sensorTypeUUID)
		return
	}
	log.Debug().Msg("Sensor type deleted_at timestamp updated successfully")
	return
}

// RestoreSensorType restore a sensor type
func (p *postgres) RestoreSensorType(sensorTypeUUID uuid.UUID) (err error) {
	query := `UPDATE sensor_types SET deleted_at = NULL WHERE id = $1`
	commandTag, err := p.Pool.Exec(context.Background(), query, sensorTypeUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to restore sensor type")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Error().Msgf("No sensor type found with the uuid: %v", sensorTypeUUID)
		return
	}
	log.Debug().Msg("Sensor type deleted_at timestamp set null successfully")
	return
}

// PatchUpdateSensorType updates only the specified fields of a sensor type
func (p *postgres) PatchUpdateSensorType(st *SensorType) (err error) {
	query := "UPDATE sensor_types SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if st.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, st.Name)
		paramID++
	}
	if st.Color != nil {
		updates = append(updates, fmt.Sprintf("color = $%d", paramID))
		params = append(params, st.Color)
		paramID++
	}
	if st.Alias != nil {
		updates = append(updates, fmt.Sprintf("alias = $%d", paramID))
		params = append(params, st.Alias)
		paramID++
	}
	if st.Interface0 != nil {
		updates = append(updates, fmt.Sprintf("interface_0 = $%d", paramID))
		params = append(params, st.Interface0)
		paramID++
	}
	if st.Interface1 != nil {
		updates = append(updates, fmt.Sprintf("interface_1 = $%d", paramID))
		params = append(params, st.Interface1)
		paramID++
	}
	if st.Interface2 != nil {
		updates = append(updates, fmt.Sprintf("interface_2 = $%d", paramID))
		params = append(params, st.Interface2)
		paramID++
	}
	if st.RxAntGain != nil {
		updates = append(updates, fmt.Sprintf("rx_ant_gain = $%d", paramID))
		params = append(params, st.RxAntGain)
		paramID++
	}
	if st.HorRotationOffset != nil {
		updates = append(updates, fmt.Sprintf("hor_rotation_offset = $%d", paramID))
		params = append(params, st.HorRotationOffset)
		paramID++
	}
	if st.VertRotationOffset != nil {
		updates = append(updates, fmt.Sprintf("vert_rotation_offset = $%d", paramID))
		params = append(params, st.VertRotationOffset)
		paramID++
	}
	if st.CorrectionFactor24 != nil {
		updates = append(updates, fmt.Sprintf("correction_factor_24 = $%d", paramID))
		params = append(params, st.CorrectionFactor24)
		paramID++
	}
	if st.CorrectionFactor5 != nil {
		updates = append(updates, fmt.Sprintf("correction_factor_5 = $%d", paramID))
		params = append(params, st.CorrectionFactor5)
		paramID++
	}
	if st.CorrectionFactor6 != nil {
		updates = append(updates, fmt.Sprintf("correction_factor_6 = $%d", paramID))
		params = append(params, st.CorrectionFactor6)
		paramID++
	}
	if st.Diagram != nil {
		updates = append(updates, fmt.Sprintf("diagram = $%d", paramID))
		params = append(params, st.Diagram)
		paramID++
	}

	updates = append(updates, fmt.Sprintf("is_virtual = $%d", paramID))
	params = append(params, st.IsVirtual)
	paramID++

	if len(updates) == 0 {
		log.Error().Msg("No fields provided for update")
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, st.ID)

	_, err = p.Pool.Exec(context.Background(), query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute update")
		return
	}

	return
}
