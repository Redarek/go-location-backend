package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"location-backend/internal/domain/dto"
)

type matrixRepo struct {
	pool *pgxpool.Pool
}

func NewMatrixRepo(pool *pgxpool.Pool) *matrixRepo {
	return &matrixRepo{pool: pool}
}

func (r *matrixRepo) Create(ctx context.Context, createMatrixDTOs []*dto.CreateMatrixDTO) (err error) {
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
	pointsRows := make([][]interface{}, len(createMatrixDTOs))
	var matrixRows [][]interface{}
	for i, dto := range createMatrixDTOs {
		pointsRows[i] = []interface{}{dto.FloorID, dto.X, dto.Y}

		for _, matrixPoint := range dto.MatrixPoints {
			matrixRows = append(matrixRows, []interface{}{
				matrixPoint.SensorID, matrixPoint.Rssi24, matrixPoint.Rssi5, matrixPoint.Rssi6, matrixPoint.Distance,
			})
		}
	}

	// Вставка данных в таблицу points
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"points"},
		[]string{"floor_id", "x", "y"},
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
		[]string{"sensor_id", "rssi24", "rssi5", "rssi6", "distance"},
		pgx.CopyFromRows(matrixRows),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to bulk insert matrix points")
		return err
	}

	return
}
