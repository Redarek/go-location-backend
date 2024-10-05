package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
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
