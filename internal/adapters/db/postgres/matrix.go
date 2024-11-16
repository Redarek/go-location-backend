package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/entity"
)

// type MatrixFilter struct {
// 	ID      *int
// 	FloorID uuid.UUID
// 	// X            *float64
// 	// Y            *float64
// 	Band     string
// 	SensorID uuid.UUID

// 	// RSSI24   *float64
// 	// RSSI5    *float64
// 	// RSSI6    *float64
// 	// Distance *float64
// }

type matrixRepo struct {
	pool *pgxpool.Pool
}

func NewMatrixRepo(pool *pgxpool.Pool) *matrixRepo {
	return &matrixRepo{pool: pool}
}

func (r *matrixRepo) Create(ctx context.Context, points []*entity.Point, matrixPoints []*entity.MatrixPoint) (err error) {
	// Begin a transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				log.Error().Msg("failed to rollback transaction")
			}
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				log.Error().Msg("failed to commit transaction")
			}
		}
	}()

	// queryPoints := `INSERT INTO points (
	// 	floor_id,
	// 	x, y
	// )
	// VALUES ($1, $2, $3)`
	// queryMatrix := `INSERT INTO matrix (
	// 	sensor_id,
	// 	rssi24, rssi5, rssi6,
	// 	distance
	// )
	// VALUES ($1, $2, $3, $4, $5)`

	// for _, createMatrixDTO := range createMatrixDTOs {
	// 	_, err = tx.Exec(ctx, queryPoints,
	// 		createMatrixDTO.FloorID,
	// 		createMatrixDTO.X, createMatrixDTO.Y,
	// 	)
	// 	if err != nil {
	// 		log.Error().Err(err).Msgf("failed to insert point %v", createMatrixDTO)
	// 		return err
	// 	}
	// }

	// Подготовка данных для таблицы points и matrix
	pointsRows := [][]interface{}{}
	for _, point := range points {
		pointsRows = append(pointsRows, []interface{}{
			point.ID,
			point.FloorID,
			point.X, point.Y,
		})
	}

	matrixRows := [][]interface{}{}
	for _, matrixPoint := range matrixPoints {
		matrixRows = append(matrixRows, []interface{}{
			matrixPoint.PointID,
			matrixPoint.SensorID,
			matrixPoint.RSSI24, matrixPoint.RSSI5, matrixPoint.RSSI6,
			matrixPoint.Distance,
		})
	}

	// Вставка данных в таблицу points
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"points"},
		[]string{"id", "floor_id", "x", "y"},
		pgx.CopyFromRows(pointsRows),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to bulk insert points")
		return err
	}

	// Вставка данных в таблицу matrix
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"matrix"},
		[]string{"point_id", "sensor_id", "rssi24", "rssi5", "rssi6", "distance"},
		pgx.CopyFromRows(matrixRows),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to bulk insert matrix points")
		return err
	}

	log.Debug().Msg("matrix is saved to database")

	return
}

func (r *deviceRepo) SearchPoints(ctx context.Context, filter entity.SearchParameters) (points []*entity.Point, err error) {
	query := `SELECT 
		p.id,
		p.floor_id,
		p.x, p.y,

		-- m.rssi24, m.rssi5, m.rssi6,
		-- m.distance,
		COUNT(*) AS count
	FROM points p 
	JOIN matrix m ON p.id = m.point_id
	WHERE 
		p.floor_id = $1`
	args := []interface{}{}
	argIndex := 2

	if len(filter.SensorsBetween) > 0 {
		query += " AND (0"
		for sensor, between := range filter.SensorsBetween {
			query += fmt.Sprintf(" OR (sensor_id = $%d AND rssi%s BETWEEN $%d AND $%d)", argIndex, filter.Band, argIndex+1, argIndex+2)
			args = append(args, sensor, between.From, between.To)
			argIndex += 3
		}
		query += ")"
	}

	query += " GROUP BY point_id"

	// TODO убедиться, что не больше 3
	query += fmt.Sprintf(" HAVING count = $%d", argIndex)
	args = append(args, filter.SensorsBetween)
	// argIndex ++

	// //? индекс тут на единицу меньше, чем по факту
	// query += fmt.Sprintf("LIMIT NULLIF($%d, 0) OFFSET $%d", argIndex, argIndex + 1)
	// args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute query")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var point *entity.Point
		err = rows.Scan(
			&point.ID,
			&point.FloorID,
			&point.X, &point.Y,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan row")
			return
		}
		points = append(points, point)
	}

	if rows.Err() != nil {
		log.Error().Err(err).Msg("error iterating rows")
		return
	}

	return
}

func (r *matrixRepo) Delete(ctx context.Context, floorID uuid.UUID) (deletedCount int64, err error) {
	query := `DELETE FROM points 
	WHERE floor_id = $1`

	cmdTag, err := r.pool.Exec(ctx, query, floorID)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete points")
		return 0, err
	}

	deletedCount = cmdTag.RowsAffected()
	log.Debug().Msgf("%d rows affected", deletedCount)

	return deletedCount, nil
}
