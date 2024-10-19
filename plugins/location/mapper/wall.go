package mapper

import (
	"sync"

	domain_entity "location-backend/internal/domain/entity"
	location_entity "location-backend/plugins/location"
)

type wallMapper struct{}

func (*wallMapper) EntityDomainToLocation(domainEntity *domain_entity.WallDetailed) (locationEntity *location_entity.Wall) {
	locationEntity = &location_entity.Wall{
		ID:            domainEntity.ID,
		X1:            domainEntity.X1,
		Y1:            domainEntity.Y1,
		X2:            domainEntity.X2,
		Y2:            domainEntity.Y2,
		Attenuation24: domainEntity.WallType.Attenuation24,
		Attenuation5:  domainEntity.WallType.Attenuation5,
		Attenuation6:  domainEntity.WallType.Attenuation6,
		Thickness:     domainEntity.WallType.Thickness,
		FloorID:       domainEntity.FloorID,
	}

	return
}

func (m *wallMapper) EntitiesDomainToLocation(domainEntities []*domain_entity.WallDetailed) (locationEntities []*location_entity.Wall) {
	for _, domainEntity := range domainEntities {
		locationEntity := m.EntityDomainToLocation(domainEntity)
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
