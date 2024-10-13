package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	// domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type BuildingMapper struct{}

// func (*SensorRadioTemplateMapper) CreateHTTPtoDomain(httpDTO *http_dto.CreateSensorRadioTemplateDTO) (domainDTO *domain_dto.CreateSensorRadioTemplateDTO) {
// 	domainDTO = &domain_dto.CreateSensorRadioTemplateDTO{
// 		Number:        httpDTO.Number,
// 		Channel:       httpDTO.Channel,
// 		Channel2:      httpDTO.Channel2,
// 		ChannelWidth:  httpDTO.ChannelWidth,
// 		WiFi:          httpDTO.WiFi,
// 		Power:         httpDTO.Power,
// 		Bandwidth:     httpDTO.Bandwidth,
// 		GuardInterval: httpDTO.GuardInterval,
// 		SensorTypeID:  httpDTO.SensorTypeID,
// 	}

// 	return
// }

func (*BuildingMapper) EntityDomainToHTTP(domainEntity *entity.Building) (httpDTO *http_dto.BuildingDTO) {
	httpDTO = &http_dto.BuildingDTO{
		ID:          domainEntity.ID,
		Name:        domainEntity.Name,
		Description: domainEntity.Description,
		Country:     domainEntity.Country,
		City:        domainEntity.City,
		Address:     domainEntity.Address,
		SiteID:      domainEntity.SiteID,
		CreatedAt:   domainEntity.CreatedAt,
		UpdatedAt:   domainEntity.UpdatedAt,
		DeletedAt:   domainEntity.DeletedAt,
	}

	return
}

func (m *BuildingMapper) EntitiesDomainToHTTP(domainEntities []*entity.Building) (httpDTOs []*http_dto.BuildingDTO) {
	for _, domainDTO := range domainEntities {
		httpDTO := m.EntityDomainToHTTP(domainDTO)
		httpDTOs = append(httpDTOs, httpDTO)
	}

	return
}

// func (*SensorRadioTemplateMapper) UpdateHTTPtoDomain(httpDTO *http_dto.PatchUpdateSensorRadioTemplateDTO) (domainDTO *domain_dto.PatchUpdateSensorRadioTemplateDTO) {
// 	domainDTO = &domain_dto.PatchUpdateSensorRadioTemplateDTO{
// 		ID:            httpDTO.ID,
// 		Number:        httpDTO.Number,
// 		Channel:       httpDTO.Channel,
// 		Channel2:      httpDTO.Channel2,
// 		ChannelWidth:  httpDTO.ChannelWidth,
// 		WiFi:          httpDTO.WiFi,
// 		Power:         httpDTO.Power,
// 		Bandwidth:     httpDTO.Bandwidth,
// 		GuardInterval: httpDTO.GuardInterval,
// 		SensorTypeID:  httpDTO.SensorTypeID,
// 	}

// 	return
// }
