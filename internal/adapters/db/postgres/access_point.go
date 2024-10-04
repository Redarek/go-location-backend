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

type accessPointRepo struct {
	pool *pgxpool.Pool
}

func NewAccessPointRepo(pool *pgxpool.Pool) *accessPointRepo {
	return &accessPointRepo{pool: pool}
}

type tmpAccessPointRadio struct {
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
	AccessPointID *uuid.UUID          `db:"access_point_id"`
	CreatedAt     *pgtype.Timestamptz `db:"created_at"`
	UpdatedAt     *pgtype.Timestamptz `db:"updated_at"`
	DeletedAt     *pgtype.Timestamptz `db:"deleted_at"`
}

func (r *accessPointRepo) Create(ctx context.Context, createAccessPointDTO *dto.CreateAccessPointDTO) (accessPointID uuid.UUID, err error) {
	query := `INSERT INTO access_points (
			name, 
			color,
			x, 
			y,
			z,
			is_virtual,
			access_point_type_id,
			floor_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	row := r.pool.QueryRow(ctx, query,
		createAccessPointDTO.Name,
		createAccessPointDTO.Color,
		createAccessPointDTO.X,
		createAccessPointDTO.Y,
		createAccessPointDTO.Z,
		createAccessPointDTO.IsVirtual,
		createAccessPointDTO.AccessPointTypeID,
		createAccessPointDTO.FloorID,
	)
	err = row.Scan(&accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("failed to scan access point")
		return
	}

	return
}

func (r *accessPointRepo) GetOne(ctx context.Context, accessPointID uuid.UUID) (accessPoint *entity.AccessPoint, err error) {
	query := `SELECT 
			id, 
			name, 
			color,
			x, 
			y,
			z,
			is_virtual,
			access_point_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM access_points WHERE id = $1 AND deleted_at IS NULL`
	row := r.pool.QueryRow(ctx, query, accessPointID)

	accessPoint = &entity.AccessPoint{}
	err = row.Scan(
		&accessPoint.ID,
		&accessPoint.Name,
		&accessPoint.Color,
		&accessPoint.X,
		&accessPoint.Y,
		&accessPoint.Z,
		&accessPoint.IsVirtual,
		&accessPoint.AccessPointTypeID,
		&accessPoint.FloorID,
		&accessPoint.CreatedAt, &accessPoint.UpdatedAt, &accessPoint.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no access point found with ID %v", accessPointID)
			return nil, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve access point")
		return
	}
	log.Debug().Msgf("retrieved access point: %v", accessPoint)
	return
}

func (r *accessPointRepo) GetOneDetailed(ctx context.Context, accessPointID uuid.UUID) (apDetailed *entity.AccessPointDetailed, err error) {
	query := `SELECT
			ap.id,
			ap.name,
			ap.color,
			ap.x, ap.y, ap.z,
			ap.is_virtual,
			ap.access_point_type_id,
			ap.floor_id,
			ap.created_at, ap.updated_at, ap.deleted_at,

			apt.id,
			apt.name,
			apt.model,
			apt.color,
			apt.z,
			apt.is_virtual,
			apt.site_id,
			apt.created_at, apt.updated_at, apt.deleted_at,

			r.id,
			r.number,
			r.channel,
			r.channel2,
			r.channel_width,
			r.wifi,
			r.power,
			r.bandwidth,
			r.guard_interval,
			r.is_active,
			r.access_point_id,
			r.created_at, r.updated_at, r.deleted_at
		FROM access_points ap
		JOIN access_point_types apt ON ap.access_point_type_id = apt.id AND apt.deleted_at IS NULL
		LEFT JOIN access_point_radios r ON ap.id = r.access_point_id AND r.deleted_at IS NULL
		WHERE ap.id = $1 AND ap.deleted_at IS NULL`
	rows, err := r.pool.Query(ctx, query, accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve access point")
		return
	}
	defer rows.Close()

	apDetailed = &entity.AccessPointDetailed{}
	apt := entity.AccessPointType{}
	i := 0

	for rows.Next() {
		i++
		tmpRadio := tmpAccessPointRadio{}

		err = rows.Scan(
			&apDetailed.ID,
			&apDetailed.Name,
			&apDetailed.Color,
			&apDetailed.X, &apDetailed.Y, &apDetailed.Z,
			&apDetailed.IsVirtual,
			&apDetailed.AccessPointTypeID,
			&apDetailed.FloorID,
			&apDetailed.CreatedAt, &apDetailed.UpdatedAt, &apDetailed.DeletedAt,

			&apt.ID,
			&apt.Name,
			&apt.Model,
			&apt.Color,
			&apt.Z,
			&apt.IsVirtual,
			&apt.SiteID,
			&apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt,

			&tmpRadio.ID,
			&tmpRadio.Number,
			&tmpRadio.Channel,
			&tmpRadio.Channel2,
			&tmpRadio.ChannelWidth,
			&tmpRadio.WiFi,
			&tmpRadio.Power,
			&tmpRadio.Bandwidth,
			&tmpRadio.GuardInterval,
			&tmpRadio.IsActive,
			&tmpRadio.AccessPointID,
			&tmpRadio.CreatedAt, &tmpRadio.UpdatedAt, &tmpRadio.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan access point and related data")
			return
		}
		apDetailed.AccessPointType = apt

		if tmpRadio.ID != nil {
			radio := &entity.AccessPointRadio{
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
				AccessPointID: *tmpRadio.AccessPointID,
				CreatedAt:     *tmpRadio.CreatedAt,
				UpdatedAt:     *tmpRadio.UpdatedAt,
				DeletedAt:     tmpRadio.DeletedAt,
			}
			apDetailed.Radios = append(apDetailed.Radios, radio)
		}
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	if i == 0 {
		log.Debug().Msg("access point was not found")
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved access point with detailed info: %v", apDetailed)
	return
}

func (r *accessPointRepo) GetAll(ctx context.Context, floorID uuid.UUID, limit, offset int) (accessPoints []*entity.AccessPoint, err error) {
	query := `SELECT 
			id, 
			name, 
			color,
			x, 
			y,
			z,
			is_virtual,
			access_point_type_id,
			floor_id,
			created_at, updated_at, deleted_at
		FROM access_points 
		WHERE floor_id = $1 AND deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, floorID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve accessPoints")
		return
	}
	defer rows.Close()

	for rows.Next() {
		accessPoint := &entity.AccessPoint{}
		err = rows.Scan(
			&accessPoint.ID,
			&accessPoint.Name,
			&accessPoint.Color,
			&accessPoint.X,
			&accessPoint.Y,
			&accessPoint.Z,
			&accessPoint.IsVirtual,
			&accessPoint.AccessPointTypeID,
			&accessPoint.FloorID,
			&accessPoint.CreatedAt, &accessPoint.UpdatedAt, &accessPoint.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan access point")
			return
		}
		accessPoints = append(accessPoints, accessPoint)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	length := len(accessPoints)
	if length == 0 {
		log.Info().Msgf("access points for floor ID %v were not found", floorID)
		return nil, service.ErrNotFound
	}

	log.Debug().Msgf("retrieved %d access points", length)
	return
}

func (r *accessPointRepo) GetAllDetailed(ctx context.Context, floorID uuid.UUID, limit, offset int) (accessPointsDetailed []*entity.AccessPointDetailed, err error) {
	query := `SELECT 
			ap.id, 
			ap.name, 
			ap.color,
			ap.x, ap.y, ap.z,
			ap.is_virtual,
			ap.access_point_type_id,
			ap.floor_id, 
			ap.created_at, ap.updated_at, ap.deleted_at,
			
			apt.id, 
			apt.name, 
			apt.model,
			apt.color, 
			apt.z,
			apt.is_virtual,
			apt.site_id, 
			apt.created_at, apt.updated_at, apt.deleted_at, 
			
			r.id, 
			r.number, 
			r.channel, 
			r.channel2,
			r.channel_width,
			r.wifi, 
			r.power, 
			r.bandwidth, 
			r.guard_interval, 
			r.is_active, 
			r.access_point_id,
			r.created_at, r.updated_at, r.deleted_at
		FROM access_points ap
		JOIN access_point_types apt ON ap.access_point_type_id = apt.id AND apt.deleted_at IS NULL
		LEFT JOIN access_point_radios r ON ap.id = r.access_point_id AND r.deleted_at IS NULL
		WHERE ap.floor_id = $1 AND ap.deleted_at IS NULL
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, floorID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve access point")
		return
	}
	defer rows.Close()

	apdMap := make(map[uuid.UUID]*entity.AccessPointDetailed)
	i := 0

	for rows.Next() {
		i++
		apd := &entity.AccessPointDetailed{}
		apt := entity.AccessPointType{}
		tmpRadio := tmpAccessPointRadio{}

		err = rows.Scan(
			&apd.ID,
			&apd.Name,
			&apd.Color,
			&apd.X, &apd.Y, &apd.Z,
			&apd.IsVirtual,
			&apd.AccessPointTypeID,
			&apd.FloorID,
			&apd.CreatedAt, &apd.UpdatedAt, &apd.DeletedAt,

			&apt.ID,
			&apt.Name,
			&apt.Model,
			&apt.Color,
			&apt.Z,
			&apt.IsVirtual,
			&apt.SiteID,
			&apt.CreatedAt, &apt.UpdatedAt, &apt.DeletedAt,

			&tmpRadio.ID,
			&tmpRadio.Number,
			&tmpRadio.Channel,
			&tmpRadio.Channel2,
			&tmpRadio.ChannelWidth,
			&tmpRadio.WiFi,
			&tmpRadio.Power,
			&tmpRadio.Bandwidth,
			&tmpRadio.GuardInterval,
			&tmpRadio.IsActive,
			&tmpRadio.AccessPointID,
			&tmpRadio.CreatedAt, &tmpRadio.UpdatedAt, &tmpRadio.DeletedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan access point and related data")
			return
		}

		var radio *entity.AccessPointRadio
		if tmpRadio.ID != nil {
			radio = &entity.AccessPointRadio{
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
				AccessPointID: *tmpRadio.AccessPointID,
				CreatedAt:     *tmpRadio.CreatedAt,
				UpdatedAt:     *tmpRadio.UpdatedAt,
				DeletedAt:     tmpRadio.DeletedAt,
			}
		}

		existingAP, exists := apdMap[apd.ID]
		if exists {
			// If access point is already in the map, append the new radio to its list
			if radio != nil {
				existingAP.Radios = append(existingAP.Radios, radio)
			}
		} else {
			// If it's a new access point, initialize and add to map
			apd.AccessPointType = apt
			if radio != nil {
				apd.Radios = append(apd.Radios, radio)
			}
			apdMap[apd.ID] = apd
		}
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("rows iteration error")
		return
	}

	if i == 0 {
		log.Debug().Msg("access point was not found")
		return nil, service.ErrNotFound
	}

	// Convert map to slice
	for _, apd := range apdMap {
		accessPointsDetailed = append(accessPointsDetailed, apd)
	}

	log.Debug().Msgf("retrieved access point with detailed info: %v", accessPointsDetailed)
	return
}

func (r *accessPointRepo) Update(ctx context.Context, updateAccessPointDTO *dto.PatchUpdateAccessPointDTO) (err error) {
	query := "UPDATE access_points SET updated_at = NOW(), "
	updates := []string{}
	params := []interface{}{}
	paramID := 1

	if updateAccessPointDTO.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", paramID))
		params = append(params, updateAccessPointDTO.Name)
		paramID++
	}
	if updateAccessPointDTO.Color != nil {
		updates = append(updates, fmt.Sprintf("color = $%d", paramID))
		params = append(params, updateAccessPointDTO.Color)
		paramID++
	}
	if updateAccessPointDTO.X != nil {
		updates = append(updates, fmt.Sprintf("x = $%d", paramID))
		params = append(params, updateAccessPointDTO.X)
		paramID++
	}
	if updateAccessPointDTO.Y != nil {
		updates = append(updates, fmt.Sprintf("y = $%d", paramID))
		params = append(params, updateAccessPointDTO.Y)
		paramID++
	}
	if updateAccessPointDTO.Z != nil {
		updates = append(updates, fmt.Sprintf("z = $%d", paramID))
		params = append(params, updateAccessPointDTO.Z)
		paramID++
	}
	if updateAccessPointDTO.IsVirtual != nil {
		updates = append(updates, fmt.Sprintf("is_virtual = $%d", paramID))
		params = append(params, updateAccessPointDTO.IsVirtual)
		paramID++
	}
	if updateAccessPointDTO.AccessPointTypeID != nil {
		updates = append(updates, fmt.Sprintf("access_point_type_id = $%d", paramID))
		params = append(params, updateAccessPointDTO.AccessPointTypeID)
		paramID++
	}

	if len(updates) == 0 {
		log.Info().Msg("no fields provided for update")
		return service.ErrNotUpdated
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramID)
	params = append(params, updateAccessPointDTO.ID)

	commandTag, err := r.pool.Exec(ctx, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute update")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point found with the ID: %v", updateAccessPointDTO.ID)
		return service.ErrNotFound
	}

	return
}

// Checks if the accessPoint has been soft deleted
func (r *accessPointRepo) IsAccessPointSoftDeleted(ctx context.Context, accessPointID uuid.UUID) (isDeleted bool, err error) {
	var deletedAt sql.NullTime // Use sql.NullTime to properly handle NULL values
	query := `SELECT deleted_at FROM access_points WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, accessPointID)
	err = row.Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Err(err).Msgf("no access point found with UUID %v", accessPointID)
			return false, service.ErrNotFound
		}
		log.Error().Err(err).Msg("failed to retrieve access point")
		return
	}

	isDeleted = deletedAt.Valid
	log.Debug().Msgf("is accessPoint deleted: %v", isDeleted)
	return
}

func (r *accessPointRepo) SoftDelete(ctx context.Context, accessPointID uuid.UUID) (err error) {
	query := `UPDATE access_points SET deleted_at = NOW() WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("failed to soft delete access point")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point found with the ID: %v", accessPointID)
		return service.ErrNotFound
	}

	log.Debug().Msg("access point deleted_at timestamp updated successfully")
	return
}

func (r *accessPointRepo) Restore(ctx context.Context, accessPointID uuid.UUID) (err error) {
	query := `UPDATE access_points SET deleted_at = NULL WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, accessPointID)
	if err != nil {
		log.Error().Err(err).Msg("failed to restore access point")
		return
	}
	if commandTag.RowsAffected() == 0 {
		log.Info().Msgf("no access point found with the UUID: %v", accessPointID)
		return service.ErrNotFound
	}

	log.Debug().Msg("access point deleted_at timestamp set NULL successfully")
	return
}
