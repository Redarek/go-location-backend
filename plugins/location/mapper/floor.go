package mapper

import (
	"sync"

	domain_entity "location-backend/internal/domain/entity"
	location_entity "location-backend/plugins/location"
)

type floorMapper struct{}

func (*floorMapper) EntityDomainToLocation(domainEntity *domain_entity.Floor) (locationEntity *location_entity.Floor) {
	locationEntity = &location_entity.Floor{
		ID:                   domainEntity.ID,
		Name:                 domainEntity.Name,
		Number:               domainEntity.Number,
		Image:                domainEntity.Image,
		Heatmap:              domainEntity.Heatmap,
		WidthInPixels:        domainEntity.WidthInPixels,
		HeightInPixels:       domainEntity.HeightInPixels,
		Scale:                domainEntity.Scale,
		CellSizeMeter:        domainEntity.CellSizeMeter,
		NorthAreaIndentMeter: domainEntity.NorthAreaIndentMeter,
		SouthAreaIndentMeter: domainEntity.SouthAreaIndentMeter,
		WestAreaIndentMeter:  domainEntity.WestAreaIndentMeter,
		EastAreaIndentMeter:  domainEntity.EastAreaIndentMeter,
	}

	return
}

func (m *floorMapper) EntitiesDomainToLocation(domainEntities []*domain_entity.Floor) (locationEntities []*location_entity.Floor) {
	for _, domainEntity := range domainEntities {
		locationEntity := m.EntityDomainToLocation(domainEntity)
		locationEntities = append(locationEntities, locationEntity)
	}

	return
}

// Паттерн Singleton для FloorMapper
var (
	floorMapperInstance *floorMapper
	floorMapperOnce     sync.Once
)

// GetFloorMapper возвращает единственный экземпляр FloorMapper
func GetFloorMapper() *floorMapper {
	floorMapperOnce.Do(func() {
		floorMapperInstance = &floorMapper{}
	})
	return floorMapperInstance
}
