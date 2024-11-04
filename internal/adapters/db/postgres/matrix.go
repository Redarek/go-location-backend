package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/entity"
)

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
