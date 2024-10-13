package mapper

import (
	http_dto "location-backend/internal/controller/http/dto"
	// domain_dto "location-backend/internal/domain/dto"
	"location-backend/internal/domain/entity"
)

type SensorMapper struct {
	sensorRadioMapper *SensorRadioMapper
}

// TODO убрать зависимость от entity
func (m *SensorMapper) DetailedToHTTP(domain *entity.SensorDetailed) (httpDTO *http_dto.SensorDetailedDTO) {
	httpDTO = &http_dto.SensorDetailedDTO{
		SensorDTO: http_dto.SensorDTO{
			ID:                 domain.ID,
			Name:               domain.Name,
			Color:              domain.Color,
			X:                  domain.X,
			Y:                  domain.Y,
			Z:                  domain.Z,
			MAC:                domain.MAC,
			IP:                 domain.IP,
			RxAntGain:          domain.RxAntGain,
			HorRotationOffset:  domain.HorRotationOffset,
			VertRotationOffset: domain.VertRotationOffset,
			CorrectionFactor24: domain.CorrectionFactor24,
			CorrectionFactor5:  domain.CorrectionFactor5,
			CorrectionFactor6:  domain.CorrectionFactor6,
			IsVirtual:          domain.IsVirtual,
			Diagram:            domain.Diagram,
			SensorTypeID:       domain.SensorTypeID,
			FloorID:            domain.FloorID,
			CreatedAt:          domain.CreatedAt,
			UpdatedAt:          domain.UpdatedAt,
			DeletedAt:          domain.DeletedAt,
		},
		SensorType: (http_dto.SensorTypeDTO)(domain.SensorType),
		Radios:     m.sensorRadioMapper.EntitiesDomainToHTTP(domain.Radios),
	}

	return
}

func (m *SensorMapper) DetailedToHTTPList(domains []*entity.SensorDetailed) (httpDTOs []*http_dto.SensorDetailedDTO) {
	for _, domain := range domains {
		httpDTOs = append(httpDTOs, m.DetailedToHTTP(domain))
	}

	return
}
