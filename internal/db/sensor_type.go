package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type SensorType struct {
	ID         uuid.UUID `json:"id" db:"id"`     // "id" INTEGER [pk, increment]
	Name       *string   `json:"name" db:"name"` //   "sensor_name" VARCHAR(45)
	Color      *string   `json:"color" db:"color"`
	Alias      *string   `json:"alias" db:"alias"`            //   "alias" VARCHAR(45)
	Interface0 *string   `json:"interface0" db:"interface_0"` //   "interface_0" VARCHAR(45) [not null]
	Interface1 *string   `json:"interface1" db:"interface_1"` //   "interface_1" VARCHAR(45)
	Interface2 *string   `json:"interface2" db:"interface_2"` //   "interface_2" VARCHAR(45)
	// TODO "state" sensors_state_enum [not null, default: "DOWN"]
	// TODO "state_change" DATETIME [not null]
	// TODO "packets_captured" INTEGER [not null, default: 0]
	// TODO  "uptime" TIME [not null]
	// TODO  "logs_path" VARCHAR(45)
	// TODO  "approved" TINYINT(1) [not null, default: FALSE]
	// TODO "mode" VARCHAR(45)
	// TODO "type" TINYINT(1)
	// TODO  "primary_channel_freq" FLOAT
	// TODO  "primary_channel_width" VARCHAR(45)
	// TODO  "primary_interval" FLOAT
	// TODO  "secondary_interval" FLOAT
	RxAntGain          *float64           `json:"rxAntGain" db:"rx_ant_gain"`                   //   "rx_ant_gain" FLOAT [not null, default: 0]
	HorRotationOffset  *int               `json:"horRotationOffset" db:"hor_rotation_offset"`   //   "hor_rotation_offset" INTEGER [not null, default: 0]
	VertRotationOffset *int               `json:"vertRotationOffset" db:"vert_rotation_offset"` //   "vert_rotation_offset" INTEGER [not null, default: 0]
	CorrectionFactor24 *float64           `json:"correctionFactor24" db:"correction_factor_24"` //   "correction_factor24" INTEGER [not null, default: 0]  -> FLOAT
	CorrectionFactor5  *float64           `json:"correctionFactor5" db:"correction_factor_5"`   //   "correction_factor5" INTEGER [not null, default: 0] -> FLOAT
	CorrectionFactor6  *float64           `json:"correctionFactor6" db:"correction_factor_6"`   //   "correction_factor6" INTEGER [not null, default: 0float64 -> FLOAT
	Diagram            *json.RawMessage   `json:"diagram" db:"diagram"`                         // Тип JSON
	CreatedAt          pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt          pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
	DeletedAt          pgtype.Timestamptz `json:"deletedAt" db:"deleted_at"`
	SiteID             uuid.UUID          `json:"siteId" db:"site_id"` //  "map_id" INTEGER
}

// CreateSensorType creates a sensor type
func (p *postgres) CreateSensorType(s *SensorType) (id uuid.UUID, err error) {
	query := `INSERT INTO sensor_types (name, color, alias, interface_0, interface_1, interface_2, rx_ant_gain, hor_rotation_offset, vert_rotation_offset, correction_factor_24, correction_factor_5, correction_factor_6, diagram, site_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, s.Name, s.Color, s.Alias, s.Interface0, s.Interface1, s.Interface2, s.RxAntGain, s.HorRotationOffset, s.VertRotationOffset, s.CorrectionFactor24, s.CorrectionFactor5, s.CorrectionFactor6, s.Diagram, s.SiteID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create sensor type")
	}
	return
}

// GetSensorType retrieves a sensor
func (p *postgres) GetSensorType(sensorTypeUUID uuid.UUID) (s *SensorType, err error) {
	query := `SELECT * FROM sensor_types WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, sensorTypeUUID)
	s = &SensorType{}
	err = row.Scan(&s.ID, &s.Name, &s.Color, &s.Alias, &s.Interface0, &s.Interface1, &s.Interface2, &s.RxAntGain, &s.HorRotationOffset, &s.VertRotationOffset, &s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6, &s.Diagram, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.SiteID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msgf("No sensor type found with uuid %v", sensorTypeUUID)
			return
		}
		log.Error().Err(err).Msg("Failed to retrieve sensor type")
		return
	}
	log.Debug().Msgf("Retrieved sensor type: %v", s)
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
func (p *postgres) GetSensorTypes(siteUUID uuid.UUID) (ss []*SensorType, err error) {
	query := `SELECT * FROM sensor_types WHERE site_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve sensor types")
		return
	}
	defer rows.Close()

	var s *SensorType
	for rows.Next() {
		s = new(SensorType)
		err = rows.Scan(&s.ID, &s.Name, &s.Color, &s.Alias, &s.Interface0, &s.Interface1, &s.Interface2, &s.RxAntGain, &s.HorRotationOffset, &s.VertRotationOffset, &s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6, &s.Diagram, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.SiteID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan sensor type")
			return
		}
		ss = append(ss, s)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Rows iteration error")
		return
	}

	log.Debug().Msgf("Retrieved %d sensor types", len(ss))
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
func (p *postgres) PatchUpdateSensorType(s *SensorType) (err error) {
	query := "UPDATE sensor_types SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if s.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, s.Name)
		paramID++
	}
	if s.Color != nil {
		updates = append(updates, fmt.Sprintf("color = $%d", paramID))
		params = append(params, s.Name)
		paramID++
	}
	if s.Alias != nil {
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
