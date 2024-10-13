package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
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

// TODO fix to dto
func (*SensorTypeMapper) EntityDomainToHTTP(domainEntity *entity.SensorType) (httpDTO *http_dto.SensorTypeDTO) {
	httpDTO = &http_dto.SensorTypeDTO{
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

func (m *SensorTypeMapper) EntitiesDomainToHTTP(domainEntities []*entity.SensorType) (httpDTOs []*http_dto.SensorTypeDTO) {
	for _, domainDTO := range domainEntities {
		httpDTO := m.EntityDomainToHTTP(domainDTO)
		httpDTOs = append(httpDTOs, httpDTO)
	}

	return
}
