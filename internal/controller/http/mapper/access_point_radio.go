package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type AccessPointRadioMapper struct{}

func (*AccessPointRadioMapper) CreateHTTPtoDomain(httpDTO *http_dto.CreateAccessPointRadioDTO) (domainDTO *domain_dto.CreateAccessPointRadioDTO) {
	domainDTO = &domain_dto.CreateAccessPointRadioDTO{
		Number:        httpDTO.Number,
		Channel:       httpDTO.Channel,
		Channel2:      httpDTO.Channel2,
		ChannelWidth:  httpDTO.ChannelWidth,
		WiFi:          httpDTO.WiFi,
		Power:         httpDTO.Power,
		Bandwidth:     httpDTO.Bandwidth,
		GuardInterval: httpDTO.GuardInterval,
		AccessPointID: httpDTO.AccessPointID,
	}

	return
}

func (*AccessPointRadioMapper) EntityDomainToHTTP(domainEntity *entity.AccessPointRadio) (httpDTO *http_dto.AccessPointRadioDTO) {
	httpDTO = &http_dto.AccessPointRadioDTO{
		ID:            domainEntity.ID,
		Number:        domainEntity.Number,
		Channel:       domainEntity.Channel,
		Channel2:      domainEntity.Channel2,
		ChannelWidth:  domainEntity.ChannelWidth,
		WiFi:          domainEntity.WiFi,
		Power:         domainEntity.Power,
		Bandwidth:     domainEntity.Bandwidth,
		GuardInterval: domainEntity.GuardInterval,
		AccessPointID: domainEntity.AccessPointID,
		CreatedAt:     domainEntity.CreatedAt,
		UpdatedAt:     domainEntity.UpdatedAt,
		DeletedAt:     domainEntity.DeletedAt,
	}

	return
}

func (m *AccessPointRadioMapper) EntitiesDomainToHTTP(domainEntities []*entity.AccessPointRadio) (httpDTOs []*http_dto.AccessPointRadioDTO) {
	for _, domainDTO := range domainEntities {
		httpDTO := m.EntityDomainToHTTP(domainDTO)
		httpDTOs = append(httpDTOs, httpDTO)
	}

	return
}

func (*AccessPointRadioMapper) UpdateHTTPtoDomain(httpDTO *http_dto.PatchUpdateAccessPointRadioDTO) (domainDTO *domain_dto.PatchUpdateAccessPointRadioDTO) {
	domainDTO = &domain_dto.PatchUpdateAccessPointRadioDTO{
		ID:            httpDTO.ID,
		Number:        httpDTO.Number,
		Channel:       httpDTO.Channel,
		Channel2:      httpDTO.Channel2,
		ChannelWidth:  httpDTO.ChannelWidth,
		WiFi:          httpDTO.WiFi,
		Power:         httpDTO.Power,
		Bandwidth:     httpDTO.Bandwidth,
		GuardInterval: httpDTO.GuardInterval,
		AccessPointID: httpDTO.AccessPointID,
	}

	return
}
