package composites

import (
	handler "location-backend/internal/controller/http/v1"
)

type HandlerComposite struct {
	HealthHandler handler.Handler
	UserHandler   handler.Handler
	RoleHandler   handler.Handler
	SiteHandler   handler.Handler
}

func NewHandlerComposite(composite *UsecaseComposite) (serviceComposite *HandlerComposite) {
	return &HandlerComposite{
		HealthHandler: handler.NewHealthHandler(composite.healthUsecase),
		UserHandler:   handler.NewUserHandler(composite.userUsecase),
		RoleHandler:   handler.NewRoleHandler(composite.roleUsecase),
		SiteHandler:   handler.NewSiteHandler(composite.siteUsecase),
	}
}
