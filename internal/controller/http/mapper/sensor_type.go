package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
)

type SensorTypeMapper struct{}

func (*SensorTypeMapper) CreateToDomain(httpDTO *http_dto.CreateSensorTypeDTO) (domainDTO *domain_dto.CreateSensorTypeDTO) {
	domainDTO = &domain_dto.CreateSensorTypeDTO{
		Name:      httpDTO.Name,
		Model:     httpDTO.Model,
		Color:     httpDTO.Color,
		Z:         httpDTO.Z,
		IsVirtual: httpDTO.IsVirtual,
		SiteID:    httpDTO.SiteID,
	}

	return
}

func (*SensorTypeMapper) ViewToHTTP(httpDTO *http_dto.SensorTypeDTO) (domainDTO *domain_dto.SensorTypeDTO) {
	domainDTO = &domain_dto.SensorTypeDTO{
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
