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

type Sensor struct {
	ID         uuid.UUID `json:"id" db:"id"`                  // "id" INTEGER [pk, increment]
	Name       *string   `json:"name" db:"name"`              //   "sensor_name" VARCHAR(45)
	X          *int      `json:"x" db:"x"`                    //   "x" FLOAT
	Y          *int      `json:"y" db:"y"`                    //   "y" FLOAT
	Z          *float64  `json:"z" db:"z"`                    //   "z" FLOAT
	MAC        *string   `json:"mac" db:"mac"`                //   "sensor_mac" VARCHAR(17) [unique, not null]
	IP         *string   `json:"ip" db:"ip"`                  //   "sensor_ip" VARCHAR(64) [not null]
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
	FloorID            uuid.UUID          `json:"floorId" db:"floor_id"` //  "map_id" INTEGER
	SensorTypeID       uuid.UUID          `json:"sensorTypeId" db:"sensor_type_id"`
}

type SensorDetailed struct {
	ID         uuid.UUID `json:"id" db:"id"`     // "id" INTEGER [pk, increment]
	Name       *string   `json:"name" db:"name"` //   "sensor_name" VARCHAR(45)
	Color      *string   `json:"color" db:"color"`
	X          *int      `json:"x" db:"x"`                    //   "x" FLOAT
	Y          *int      `json:"y" db:"y"`                    //   "y" FLOAT
	Z          *float64  `json:"z" db:"z"`                    //   "z" FLOAT
	MAC        *string   `json:"mac" db:"mac"`                //   "sensor_mac" VARCHAR(17) [unique, not null]
	IP         *string   `json:"ip" db:"ip"`                  //   "sensor_ip" VARCHAR(64) [not null]
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
	FloorID            uuid.UUID          `json:"floorId" db:"floor_id"` //  "map_id" INTEGER
	SensorTypeID       uuid.UUID          `json:"sensorTypeId" db:"sensor_type_id"`
}

type Diagram struct {
	Degree map[string]Degree `json:"degree"`
}

type Degree struct {
	HorGain  float64 `json:"hor_gain"`
	VertGain float64 `json:"vert_gain"`
}

// CreateSensor creates a sensor
func (p *postgres) CreateSensor(s *Sensor) (id uuid.UUID, err error) {
	query := `INSERT INTO sensors (name, x, y, z, mac, ip, alias, interface_0, interface_1, interface_2, rx_ant_gain, hor_rotation_offset, vert_rotation_offset, correction_factor_24, correction_factor_5, correction_factor_6, diagram, floor_id, sensor_type_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
			RETURNING id`
	row := p.Pool.QueryRow(context.Background(), query, s.Name, s.X, s.Y, s.Z, s.MAC, s.IP, s.Alias, s.Interface0, s.Interface1, s.Interface2, s.RxAntGain, s.HorRotationOffset, s.VertRotationOffset, s.CorrectionFactor24, s.CorrectionFactor5, s.CorrectionFactor6, s.Diagram, s.FloorID, s.SensorTypeID)
	err = row.Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create sensor")
	}
	return
}

// GetSensor retrieves a sensor
func (p *postgres) GetSensor(sensorUUID uuid.UUID) (s *Sensor, err error) {
	query := `SELECT * FROM sensors WHERE id = $1 AND deleted_at IS NULL`
	row := p.Pool.QueryRow(context.Background(), query, sensorUUID)
	s = &Sensor{}
	err = row.Scan(&s.ID, &s.Name, &s.X, &s.Y, &s.Z, &s.MAC, &s.IP, &s.Alias, &s.Interface0, &s.Interface1, &s.Interface2, &s.RxAntGain, &s.HorRotationOffset, &s.VertRotationOffset, &s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6, &s.Diagram, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.FloorID, &s.SensorTypeID)
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
	query := `
	SELECT s.id, s.name, st.color, s.x, s.y, s.z, s.mac, s.ip, s.alias, s.interface_0, s.interface_1, s.interface_2, s.rx_ant_gain, s.hor_rotation_offset, s.vert_rotation_offset, s.correction_factor_24, s.correction_factor_5, s.correction_factor_6, s.diagram, s.created_at, s.updated_at, s.deleted_at, s.floor_id, s.sensor_type_id
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
			&s.ID, &s.Name, &st.Color, &s.X, &s.Y, &s.Z, &s.MAC, &s.IP, &s.Alias, &s.Interface0, &s.Interface1, &s.Interface2, &s.RxAntGain, &s.HorRotationOffset, &s.VertRotationOffset, &s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6, &s.Diagram, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.FloorID, &s.SensorTypeID,
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
	query := `SELECT * FROM sensors WHERE floor_id = $1 AND deleted_at IS NULL`
	rows, err := p.Pool.Query(context.Background(), query, floorUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve sensors")
		return
	}
	defer rows.Close()

	var s *Sensor
	for rows.Next() {
		s = new(Sensor)
		err = rows.Scan(&s.ID, &s.Name, &s.X, &s.Y, &s.Z, &s.MAC, &s.IP, &s.Alias, &s.Interface0, &s.Interface1, &s.Interface2, &s.RxAntGain, &s.HorRotationOffset, &s.VertRotationOffset, &s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6, &s.Diagram, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.FloorID, &s.SensorTypeID)
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
	query := `
	SELECT s.id, s.name, st.color, s.x, s.y, s.z, s.mac, s.ip, s.alias, s.interface_0, s.interface_1, s.interface_2, s.rx_ant_gain, s.hor_rotation_offset, s.vert_rotation_offset, s.correction_factor_24, s.correction_factor_5, s.correction_factor_6, s.diagram, s.created_at, s.updated_at, s.deleted_at, s.floor_id, s.sensor_type_id
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
			&s.ID, &s.Name, &st.Color, &s.X, &s.Y, &s.Z, &s.MAC, &s.IP, &s.Alias, &s.Interface0, &s.Interface1, &s.Interface2, &s.RxAntGain, &s.HorRotationOffset, &s.VertRotationOffset, &s.CorrectionFactor24, &s.CorrectionFactor5, &s.CorrectionFactor6, &s.Diagram, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.FloorID, &s.SensorTypeID,
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

	if s.Name != nil {
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
	if s.MAC != nil {
		updates = append(updates, fmt.Sprintf("mac = $%d", paramID))
		params = append(params, s.MAC)
		paramID++
	}
	if s.IP != nil {
		updates = append(updates, fmt.Sprintf("ip = $%d", paramID))
		params = append(params, s.IP)
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
