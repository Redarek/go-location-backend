package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SensorRadioMapper struct{}

func (*SensorRadioMapper) CreateHTTPtoDomain(httpDTO *http_dto.CreateSensorRadioDTO) (domainDTO *domain_dto.CreateSensorRadioDTO) {
	domainDTO = &domain_dto.CreateSensorRadioDTO{
		Number:        httpDTO.Number,
		Channel:       httpDTO.Channel,
		Channel2:      httpDTO.Channel2,
		ChannelWidth:  httpDTO.ChannelWidth,
		WiFi:          httpDTO.WiFi,
		Power:         httpDTO.Power,
		Bandwidth:     httpDTO.Bandwidth,
		GuardInterval: httpDTO.GuardInterval,
		SensorID:      httpDTO.SensorID,
	}

	return
}

func (*SensorRadioMapper) EntityDomainToHTTP(domainEntity *entity.SensorRadio) (httpDTO *http_dto.SensorRadioDTO) {
	httpDTO = &http_dto.SensorRadioDTO{
		ID:            domainEntity.ID,
		Number:        domainEntity.Number,
		Channel:       domainEntity.Channel,
		Channel2:      domainEntity.Channel2,
		ChannelWidth:  domainEntity.ChannelWidth,
		WiFi:          domainEntity.WiFi,
		Power:         domainEntity.Power,
		Bandwidth:     domainEntity.Bandwidth,
		GuardInterval: domainEntity.GuardInterval,
		SensorID:      domainEntity.SensorID,
		CreatedAt:     domainEntity.CreatedAt,
		UpdatedAt:     domainEntity.UpdatedAt,
		DeletedAt:     domainEntity.DeletedAt,
	}

	return
}

func (m *SensorRadioMapper) EntitiesDomainToHTTP(domainEntities []*entity.SensorRadio) (httpDTOs []*http_dto.SensorRadioDTO) {
	for _, domainDTO := range domainEntities {
		httpDTO := m.EntityDomainToHTTP(domainDTO)
		httpDTOs = append(httpDTOs, httpDTO)
	}

	return
}

func (*SensorRadioMapper) UpdateHTTPtoDomain(httpDTO *http_dto.PatchUpdateSensorRadioDTO) (domainDTO *domain_dto.PatchUpdateSensorRadioDTO) {
	domainDTO = &domain_dto.PatchUpdateSensorRadioDTO{
		ID:            httpDTO.ID,
		Number:        httpDTO.Number,
		Channel:       httpDTO.Channel,
		Channel2:      httpDTO.Channel2,
		ChannelWidth:  httpDTO.ChannelWidth,
		WiFi:          httpDTO.WiFi,
		Power:         httpDTO.Power,
		Bandwidth:     httpDTO.Bandwidth,
		GuardInterval: httpDTO.GuardInterval,
		SensorID:      httpDTO.SensorID,
	}

	return
}
