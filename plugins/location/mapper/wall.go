package mapper

import (
	"sync"

	domain_entity "location-backend/internal/domain/entity"
	"location-backend/plugins/location"
)

type wallMapper struct{}

func (*wallMapper) EntityDomainToLocation(domainEntity *domain_entity.WallDetailed, scale float64) (locationEntity *location.Wall) {
	x1 := location.PixelsToMeters(float64(domainEntity.X1), scale)
	y1 := location.PixelsToMeters(float64(domainEntity.Y1), scale)
	x2 := location.PixelsToMeters(float64(domainEntity.X2), scale)
	y2 := location.PixelsToMeters(float64(domainEntity.Y2), scale)
	locationEntity = &location.Wall{
		ID:            domainEntity.ID,
		X1:            x1,
		Y1:            y1,
		X2:            x2,
		Y2:            y2,
		Attenuation24: domainEntity.WallType.Attenuation24,
		Attenuation5:  domainEntity.WallType.Attenuation5,
		Attenuation6:  domainEntity.WallType.Attenuation6,
		Thickness:     domainEntity.WallType.Thickness,
		FloorID:       domainEntity.FloorID,
	}

	return
}

func (m *wallMapper) EntitiesDomainToLocation(domainEntities []*domain_entity.WallDetailed, scale float64) (locationEntities []*location.Wall) {
	for _, domainEntity := range domainEntities {
		locationEntity := m.EntityDomainToLocation(domainEntity, scale)
		locationEntities = append(locationEntities, locationEntity)
	}

	return
}

// Паттерн Singleton для WallMapper
var (
	wallMapperInstance *wallMapper
	wallMapperOnce     sync.Once
)

// GetWallMapper возвращает единственный экземпляр WallMapper
func GetWallMapper() *wallMapper {
	wallMapperOnce.Do(func() {
		wallMapperInstance = &wallMapper{}
	})
	return wallMapperInstance
}
