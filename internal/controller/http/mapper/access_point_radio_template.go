package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type AccessPointRadioTemplateMapper struct{}

func (*AccessPointRadioTemplateMapper) CreateHTTPtoDomain(httpDTO *http_dto.CreateAccessPointRadioTemplateDTO) (domainDTO *domain_dto.CreateAccessPointRadioTemplateDTO) {
	domainDTO = &domain_dto.CreateAccessPointRadioTemplateDTO{
		Number:            httpDTO.Number,
		Channel:           httpDTO.Channel,
		Channel2:          httpDTO.Channel2,
		ChannelWidth:      httpDTO.ChannelWidth,
		WiFi:              httpDTO.WiFi,
		Power:             httpDTO.Power,
		Bandwidth:         httpDTO.Bandwidth,
		GuardInterval:     httpDTO.GuardInterval,
		AccessPointTypeID: httpDTO.AccessPointTypeID,
	}

	return
}

func (*AccessPointRadioTemplateMapper) EntityDomainToHTTP(domainEntity *entity.AccessPointRadioTemplate) (httpDTO *http_dto.AccessPointRadioTemplateDTO) {
	httpDTO = &http_dto.AccessPointRadioTemplateDTO{
		ID:                domainEntity.ID,
		Number:            domainEntity.Number,
		Channel:           domainEntity.Channel,
		Channel2:          domainEntity.Channel2,
		ChannelWidth:      domainEntity.ChannelWidth,
		WiFi:              domainEntity.WiFi,
		Power:             domainEntity.Power,
		Bandwidth:         domainEntity.Bandwidth,
		GuardInterval:     domainEntity.GuardInterval,
		AccessPointTypeID: domainEntity.AccessPointTypeID,
		CreatedAt:         domainEntity.CreatedAt,
		UpdatedAt:         domainEntity.UpdatedAt,
		DeletedAt:         domainEntity.DeletedAt,
	}

	return
}

func (m *AccessPointRadioTemplateMapper) EntitiesDomainToHTTP(domainEntities []*entity.AccessPointRadioTemplate) (httpDTOs []*http_dto.AccessPointRadioTemplateDTO) {
	for _, domainDTO := range domainEntities {
		httpDTO := m.EntityDomainToHTTP(domainDTO)
		httpDTOs = append(httpDTOs, httpDTO)
	}

	return
}

func (*AccessPointRadioTemplateMapper) UpdateHTTPtoDomain(httpDTO *http_dto.PatchUpdateAccessPointRadioTemplateDTO) (domainDTO *domain_dto.PatchUpdateAccessPointRadioTemplateDTO) {
	domainDTO = &domain_dto.PatchUpdateAccessPointRadioTemplateDTO{
		ID:                httpDTO.ID,
		Number:            httpDTO.Number,
		Channel:           httpDTO.Channel,
		Channel2:          httpDTO.Channel2,
		ChannelWidth:      httpDTO.ChannelWidth,
		WiFi:              httpDTO.WiFi,
		Power:             httpDTO.Power,
		Bandwidth:         httpDTO.Bandwidth,
		GuardInterval:     httpDTO.GuardInterval,
		AccessPointTypeID: httpDTO.AccessPointTypeID,
	}

	return
}
