package mapper

import (
	"sync"

	domain_entity "location-backend/internal/domain/entity"
	location_entity "location-backend/plugins/location"
)

type deviceMapper struct{}

func (*deviceMapper) EntityDomainToLocation(domainEntity *domain_entity.DeviceDetailed) (locationEntity *location_entity.Device) {
	locationEntity = &location_entity.Device{
		ID:       domainEntity.ID,
		MAC:      domainEntity.MAC,
		FloorID:  domainEntity.FloorID,
		SensorID: *domainEntity.SensorID,
		RSSI:     domainEntity.RSSI,
		Band:     *domainEntity.Band,
	}

	return
}

func (m *deviceMapper) EntitiesDomainToLocation(domainEntities []*domain_entity.DeviceDetailed) (locationEntities []*location_entity.Device) {
	for _, domainEntity := range domainEntities {
		locationEntity := m.EntityDomainToLocation(domainEntity)
		locationEntities = append(locationEntities, locationEntity)
	}

	return
}

// Паттерн Singleton для DeviceMapper
var (
	deviceMapperInstance *deviceMapper
	deviceMapperOnce     sync.Once
)

// GetDeviceMapper возвращает единственный экземпляр DeviceMapper
func GetDeviceMapper() *deviceMapper {
	deviceMapperOnce.Do(func() {
		deviceMapperInstance = &deviceMapper{}
	})
	return deviceMapperInstance
}
