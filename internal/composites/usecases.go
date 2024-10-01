package composites

import (
	"location-backend/internal/domain/usecase"
)

type UsecaseComposite struct {
	healthUsecase                   *usecase.HealthUsecase
	userUsecase                     *usecase.UserUsecase
	roleUsecase                     *usecase.RoleUsecase
	siteUsecase                     *usecase.SiteUsecase
	buildingUsecase                 *usecase.BuildingUsecase
	floorUsecase                    *usecase.FloorUsecase
	wallTypeUsecase                 *usecase.WallTypeUsecase
	wallUsecase                     *usecase.WallUsecase
	accessPointTypeUsecase          *usecase.AccessPointTypeUsecase
	accessPointRadioTemplateUsecase *usecase.AccessPointRadioTemplateUsecase
	accessPointUsecase              *usecase.AccessPointUsecase
	accessPointRadioUsecase         *usecase.AccessPointRadioUsecase
}

func NewUsecaseComposite(composite *ServiceComposite) (serviceComposite *UsecaseComposite) {
	return &UsecaseComposite{
		healthUsecase:                   usecase.NewHealthUsecase(composite.healthService),
		userUsecase:                     usecase.NewUserUsecase(composite.userService),
		roleUsecase:                     usecase.NewRoleUsecase(composite.roleService),
		siteUsecase:                     usecase.NewSiteUsecase(composite.siteService),
		buildingUsecase:                 usecase.NewBuildingUsecase(composite.buildingService, composite.siteService),
		floorUsecase:                    usecase.NewFloorUsecase(composite.floorService, composite.buildingService),
		wallTypeUsecase:                 usecase.NewWallTypeUsecase(composite.wallTypeService, composite.siteService),
		wallUsecase:                     usecase.NewWallUsecase(composite.wallService, composite.floorService),
		accessPointTypeUsecase:          usecase.NewAccessPointTypeUsecase(composite.accessPointTypeService, composite.siteService),
		accessPointRadioTemplateUsecase: usecase.NewAccessPointRadioTemplateUsecase(composite.accessPointRadioTemplateService, composite.accessPointTypeService),
		accessPointUsecase:              usecase.NewAccessPointUsecase(composite.accessPointService, composite.accessPointTypeService, composite.floorService),
		accessPointRadioUsecase:         usecase.NewAccessPointRadioUsecase(composite.accessPointRadioService, composite.accessPointService),
	}
}
