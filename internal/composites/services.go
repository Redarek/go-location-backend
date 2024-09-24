package composites

import (
	"location-backend/internal/domain/service"
)

type ServiceComposite struct {
	healthService   service.HealthService
	userService     service.UserService
	roleService     service.RoleService
	siteService     service.SiteService
	buildingService service.BuildingService
	floorService    service.FloorService
}

func NewServiceComposite(composite *RepositoryComposite) (serviceComposite *ServiceComposite) {
	return &ServiceComposite{
		healthService:   service.NewHealthService(composite.healthRepo),
		userService:     service.NewUserService(composite.userRepo),
		roleService:     service.NewRoleService(composite.roleRepo),
		siteService:     service.NewSiteService(composite.siteRepo),
		buildingService: service.NewBuildingService(composite.buildingRepo),
		floorService:    service.NewFloorService(composite.floorRepo),
	}
}
