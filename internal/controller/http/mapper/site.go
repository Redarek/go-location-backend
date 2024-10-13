package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	// domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SiteMapper struct {
	buildingMapper        *BuildingMapper
	wallTypeMapper        *WallTypeMapper
	accessPointTypeMapper *AccessPointTypeMapper
	sensorTypeMapper      *SensorTypeMapper
}

// TODO убрать зависимость от entity
func (m *SiteMapper) DetailedToHTTP(domain *entity.SiteDetailed) (httpDTO *http_dto.SiteDetailedDTO) {
	httpDTO = &http_dto.SiteDetailedDTO{
		SiteDTO: http_dto.SiteDTO{
			ID:          domain.ID,
			Name:        domain.Name,
			Description: domain.Description,
			UserID:      domain.UserID,
			CreatedAt:   domain.CreatedAt,
			UpdatedAt:   domain.UpdatedAt,
			DeletedAt:   domain.DeletedAt,
		},

		Buildings:        m.buildingMapper.EntitiesDomainToHTTP(domain.Buildings),
		WallTypes:        m.wallTypeMapper.EntitiesDomainToHTTP(domain.WallTypes),
		AccessPointTypes: m.accessPointTypeMapper.EntitiesDomainToHTTP(domain.AccessPointTypes),
		SensorTypes:      m.sensorTypeMapper.EntitiesDomainToHTTP(domain.SensorTypes),
	}

	return
}

func (m *SiteMapper) DetailedToHTTPList(domains []*entity.SiteDetailed) (httpDTOs []*http_dto.SiteDetailedDTO) {
	for _, domain := range domains {
		httpDTOs = append(httpDTOs, m.DetailedToHTTP(domain))
	}

	return
}
