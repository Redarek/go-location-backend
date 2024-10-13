package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type AccessPointTypeMapper struct{}

func (*AccessPointTypeMapper) CreateToDomain(httpDTO *http_dto.CreateAccessPointTypeDTO) (domainDTO *domain_dto.CreateAccessPointTypeDTO) {
	domainDTO = &domain_dto.CreateAccessPointTypeDTO{
		Name:      httpDTO.Name,
		Model:     httpDTO.Model,
		Color:     httpDTO.Color,
		Z:         httpDTO.Z,
		IsVirtual: httpDTO.IsVirtual,
		SiteID:    httpDTO.SiteID,
	}

	return
}

func (*AccessPointTypeMapper) ViewToHTTP(httpDTO *http_dto.AccessPointTypeDTO) (domainDTO *domain_dto.AccessPointTypeDTO) {
	domainDTO = &domain_dto.AccessPointTypeDTO{
		ID:        httpDTO.ID,
		Name:      httpDTO.Name,
		Model:     httpDTO.Model,
		Color:     httpDTO.Color,
		Z:         httpDTO.Z,
		IsVirtual: httpDTO.IsVirtual,
		SiteID:    httpDTO.SiteID,
		CreatedAt: httpDTO.CreatedAt,
		UpdatedAt: httpDTO.UpdatedAt,
		DeletedAt: httpDTO.DeletedAt,
	}

	return
}

// TODO fix to dto
func (*AccessPointTypeMapper) EntityDomainToHTTP(domainEntity *entity.AccessPointType) (httpDTO *http_dto.AccessPointTypeDTO) {
	httpDTO = &http_dto.AccessPointTypeDTO{
		ID:        domainEntity.ID,
		Name:      domainEntity.Name,
		Model:     domainEntity.Model,
		Color:     domainEntity.Color,
		Z:         domainEntity.Z,
		IsVirtual: domainEntity.IsVirtual,
		SiteID:    domainEntity.SiteID,
		CreatedAt: domainEntity.CreatedAt,
		UpdatedAt: domainEntity.UpdatedAt,
		DeletedAt: domainEntity.DeletedAt,
	}

	return
}

func (m *AccessPointTypeMapper) EntitiesDomainToHTTP(domainEntities []*entity.AccessPointType) (httpDTOs []*http_dto.AccessPointTypeDTO) {
	for _, domainDTO := range domainEntities {
		httpDTO := m.EntityDomainToHTTP(domainDTO)
		httpDTOs = append(httpDTOs, httpDTO)
	}

	return
}
