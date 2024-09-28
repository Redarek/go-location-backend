package composites

import (
	"location-backend/internal/domain/service"
	"location-backend/internal/domain/usecase"
)

type ServiceComposite struct {
	healthService   usecase.HealthService
	userService     usecase.UserService
	roleService     usecase.RoleService
	siteService     usecase.SiteService
	buildingService usecase.BuildingService
	floorService    usecase.FloorService
	wallTypeService usecase.WallTypeService
	wallService     usecase.WallService
}

func NewServiceComposite(composite *RepositoryComposite) (serviceComposite *ServiceComposite) {
	return &ServiceComposite{
		healthService:   service.NewHealthService(composite.healthRepo),
		userService:     service.NewUserService(composite.userRepo),
		roleService:     service.NewRoleService(composite.roleRepo),
		siteService:     service.NewSiteService(composite.siteRepo),
		buildingService: service.NewBuildingService(composite.buildingRepo),
		floorService:    service.NewFloorService(composite.floorRepo),
		wallTypeService: service.NewWallTypeService(composite.wallTypeRepo),
		wallService:     service.NewWallService(composite.wallRepo, composite.wallTypeRepo),
	}
}
