package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SensorRadioTemplateMapper struct{}

func (*SensorRadioTemplateMapper) CreateHTTPtoDomain(httpDTO *http_dto.CreateSensorRadioTemplateDTO) (domainDTO *domain_dto.CreateSensorRadioTemplateDTO) {
	domainDTO = &domain_dto.CreateSensorRadioTemplateDTO{
		Number:        httpDTO.Number,
		Channel:       httpDTO.Channel,
		Channel2:      httpDTO.Channel2,
		ChannelWidth:  httpDTO.ChannelWidth,
		WiFi:          httpDTO.WiFi,
		Power:         httpDTO.Power,
		Bandwidth:     httpDTO.Bandwidth,
		GuardInterval: httpDTO.GuardInterval,
		SensorTypeID:  httpDTO.SensorTypeID,
	}

	return
}

func (*SensorRadioTemplateMapper) EntityDomainToHTTP(domainEntity *entity.SensorRadioTemplate) (httpDTO *http_dto.SensorRadioTemplateDTO) {
	httpDTO = &http_dto.SensorRadioTemplateDTO{
		ID:            domainEntity.ID,
		Number:        domainEntity.Number,
		Channel:       domainEntity.Channel,
		Channel2:      domainEntity.Channel2,
		ChannelWidth:  domainEntity.ChannelWidth,
		WiFi:          domainEntity.WiFi,
		Power:         domainEntity.Power,
		Bandwidth:     domainEntity.Bandwidth,
		GuardInterval: domainEntity.GuardInterval,
		SensorTypeID:  domainEntity.SensorTypeID,
		CreatedAt:     domainEntity.CreatedAt,
		UpdatedAt:     domainEntity.UpdatedAt,
		DeletedAt:     domainEntity.DeletedAt,
	}

	return
}

func (m *SensorRadioTemplateMapper) EntitiesDomainToHTTP(domainEntities []*entity.SensorRadioTemplate) (httpDTOs []*http_dto.SensorRadioTemplateDTO) {
	for _, domainDTO := range domainEntities {
		httpDTO := m.EntityDomainToHTTP(domainDTO)
		httpDTOs = append(httpDTOs, httpDTO)
	}

	return
}

func (*SensorRadioTemplateMapper) UpdateHTTPtoDomain(httpDTO *http_dto.PatchUpdateSensorRadioTemplateDTO) (domainDTO *domain_dto.PatchUpdateSensorRadioTemplateDTO) {
	domainDTO = &domain_dto.PatchUpdateSensorRadioTemplateDTO{
		ID:            httpDTO.ID,
		Number:        httpDTO.Number,
		Channel:       httpDTO.Channel,
		Channel2:      httpDTO.Channel2,
		ChannelWidth:  httpDTO.ChannelWidth,
		WiFi:          httpDTO.WiFi,
		Power:         httpDTO.Power,
		Bandwidth:     httpDTO.Bandwidth,
		GuardInterval: httpDTO.GuardInterval,
		SensorTypeID:  httpDTO.SensorTypeID,
	}

	return
}
