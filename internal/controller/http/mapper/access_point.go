package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	// domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type AccessPointMapper struct {
	aprMapper *AccessPointRadioMapper
}

// TODO убрать зависимость от entity
func (m *AccessPointMapper) DetailedToHTTP(domain *entity.AccessPointDetailed) (httpDTO *http_dto.AccessPointDetailedDTO) {
	httpDTO = &http_dto.AccessPointDetailedDTO{
		AccessPointDTO: http_dto.AccessPointDTO{
			ID:                domain.ID,
			Name:              domain.Name,
			Color:             domain.Color,
			X:                 domain.X,
			Y:                 domain.Y,
			Z:                 domain.Z,
			IsVirtual:         domain.IsVirtual,
			AccessPointTypeID: domain.AccessPointTypeID,
			FloorID:           domain.FloorID,
			CreatedAt:         domain.CreatedAt,
			UpdatedAt:         domain.UpdatedAt,
			DeletedAt:         domain.DeletedAt,
		},
		AccessPointType: (http_dto.AccessPointTypeDTO)(domain.AccessPointType),
		Radios:          m.aprMapper.EntitiesDomainToHTTP(domain.Radios),
	}

	return
}

func (m *AccessPointMapper) DetailedToHTTPList(domains []*entity.AccessPointDetailed) (httpDTOs []*http_dto.AccessPointDetailedDTO) {
	for _, domain := range domains {
		httpDTOs = append(httpDTOs, m.DetailedToHTTP(domain))
	}

	return
}
