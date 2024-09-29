package composites

import (
	handler "location-backend/internal/controller/http/v1"
)

type HandlerComposite struct {
	HealthHandler                   handler.Handler
	UserHandler                     handler.Handler
	RoleHandler                     handler.Handler
	SiteHandler                     handler.Handler
	BuildingHandler                 handler.Handler
	FloorHandler                    handler.Handler
	WallTypeHandler                 handler.Handler
	WallHandler                     handler.Handler
	AccessPointTypeHandler          handler.Handler
	AccessPointRadioTemplateHandler handler.Handler
}

func NewHandlerComposite(composite *UsecaseComposite) (serviceComposite *HandlerComposite) {
	return &HandlerComposite{
		HealthHandler:                   handler.NewHealthHandler(composite.healthUsecase),
		UserHandler:                     handler.NewUserHandler(composite.userUsecase),
		RoleHandler:                     handler.NewRoleHandler(composite.roleUsecase),
		SiteHandler:                     handler.NewSiteHandler(composite.siteUsecase),
		BuildingHandler:                 handler.NewBuildingHandler(composite.buildingUsecase),
		FloorHandler:                    handler.NewFloorHandler(composite.floorUsecase),
		WallTypeHandler:                 handler.NewWallTypeHandler(composite.wallTypeUsecase),
		WallHandler:                     handler.NewWallHandler(composite.wallUsecase),
		AccessPointTypeHandler:          handler.NewAccessPointTypeHandler(composite.accessPointTypeUsecase),
		AccessPointRadioTemplateHandler: handler.NewAccessPointRadioTemplateHandler(composite.accessPointRadioTemplateUsecase),
	}
}
