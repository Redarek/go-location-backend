package mapper

import (
	"sync"

	domain_entity "location-backend/internal/domain/entity"
	location_entity "location-backend/plugins/location"
)

type sensorMapper struct{}

func (*sensorMapper) EntityDomainToLocation(domainEntity *domain_entity.Sensor) (locationEntity *location_entity.Sensor) {
	locationEntity = &location_entity.Sensor{
		ID:                 domainEntity.ID,
		Name:               domainEntity.Name,
		X:                  domainEntity.X,
		Y:                  domainEntity.Y,
		Z:                  domainEntity.Z,
		MAC:                domainEntity.MAC,
		RxAntGain:          domainEntity.RxAntGain,
		HorRotationOffset:  domainEntity.HorRotationOffset,
		VertRotationOffset: domainEntity.VertRotationOffset,
		CorrectionFactor24: domainEntity.CorrectionFactor24,
		CorrectionFactor5:  domainEntity.CorrectionFactor5,
		CorrectionFactor6:  domainEntity.CorrectionFactor6,
		IsVirtual:          domainEntity.IsVirtual,
		Diagram:            domainEntity.Diagram,
		FloorID:            domainEntity.FloorID,
	}

	return
}

func (m *sensorMapper) EntitiesDomainToLocation(domainEntities []*domain_entity.Sensor) (locationEntities []*location_entity.Sensor) {
	for _, domainEntity := range domainEntities {
		locationEntity := m.EntityDomainToLocation(domainEntity)
		locationEntities = append(locationEntities, locationEntity)
	}

	return
}

// Паттерн Singleton для SensorMapper
var (
	sensorMapperInstance *sensorMapper
	sensorMapperOnce     sync.Once
)

// GetSensorMapper возвращает единственный экземпляр SensorMapper
func GetSensorMapper() *sensorMapper {
	sensorMapperOnce.Do(func() {
		sensorMapperInstance = &sensorMapper{}
	})
	return sensorMapperInstance
}
