package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	// domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type WallTypeMapper struct{}

func (*WallTypeMapper) EntityDomainToHTTP(domainEntity *entity.WallType) (httpDTO *http_dto.WallTypeDTO) {
	httpDTO = &http_dto.WallTypeDTO{
		ID:            domainEntity.ID,
		Name:          domainEntity.Name,
		Color:         domainEntity.Color,
		Attenuation24: domainEntity.Attenuation24,
		Attenuation5:  domainEntity.Attenuation5,
		Attenuation6:  domainEntity.Attenuation6,
		Thickness:     domainEntity.Thickness,
		SiteID:        domainEntity.SiteID,
		CreatedAt:     domainEntity.CreatedAt,
		UpdatedAt:     domainEntity.UpdatedAt,
		DeletedAt:     domainEntity.DeletedAt,
	}

	return
}

func (m *WallTypeMapper) EntitiesDomainToHTTP(domainEntities []*entity.WallType) (httpDTOs []*http_dto.WallTypeDTO) {
	for _, domainDTO := range domainEntities {
		httpDTO := m.EntityDomainToHTTP(domainDTO)
		httpDTOs = append(httpDTOs, httpDTO)
	}

	return
}
