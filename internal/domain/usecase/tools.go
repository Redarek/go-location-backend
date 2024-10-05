package usecase

import "github.com/rs/zerolog/log"

func LimitOffsetValidator(limit, offset, maxLimit int) (int, int) {
	if limit < 1 {
		log.Warn().Msgf("invalid limit: too low: %d, min is 1", limit)
		limit = 1
	}
	if limit > maxLimit {
		log.Warn().Msgf("invalid limit: too hight: %d, max is %d", limit, maxLimit)
		limit = maxLimit
	}
	if offset < 0 {
		log.Warn().Msgf("invalid offset: too low: %d", offset)
		offset = 0
	}

	return limit, offset
}
