package composites

import (
	"location-backend/internal/domain/service"
	"location-backend/internal/domain/usecase"
)

type ServiceComposite struct {
	healthService   usecase.IHealthService
	userService     usecase.IUserService
	roleService     usecase.IRoleService
	siteService     usecase.ISiteService
	buildingService usecase.IBuildingService
	floorService    usecase.IFloorService
	wallTypeService usecase.IWallTypeService
	wallService     usecase.IWallService

	accessPointTypeService          usecase.IAccessPointTypeService
	accessPointRadioTemplateService usecase.IAccessPointRadioTemplateService
	accessPointService              usecase.IAccessPointService
	accessPointRadioService         usecase.IAccessPointRadioService

	sensorTypeService          usecase.ISensorTypeService
	sensorRadioTemplateService usecase.ISensorRadioTemplateService
	sensorService              usecase.ISensorService
	sensorRadioService         usecase.ISensorRadioService

	matrixService usecase.IMatrixService

	deviceService usecase.IDeviceService
}

func NewServiceComposite(composite *RepositoryComposite) (serviceComposite *ServiceComposite) {
	return &ServiceComposite{
		healthService:   service.NewHealthService(composite.healthRepo),
		userService:     service.NewUserService(composite.userRepo),
		roleService:     service.NewRoleService(composite.roleRepo),
		siteService:     service.NewSiteService(composite.siteRepo, composite.buildingRepo, composite.wallTypeRepo, composite.accessPointTypeRepo, composite.sensorTypeRepo),
		buildingService: service.NewBuildingService(composite.buildingRepo),
		floorService:    service.NewFloorService(composite.floorRepo),
		wallTypeService: service.NewWallTypeService(composite.wallTypeRepo),
		wallService:     service.NewWallService(composite.wallRepo, composite.wallTypeRepo),

		accessPointTypeService:          service.NewAccessPointTypeService(composite.accessPointTypeRepo, composite.accessPointRadioTemplateRepo),
		accessPointRadioTemplateService: service.NewAccessPointRadioTemplateService(composite.accessPointRadioTemplateRepo),
		accessPointService:              service.NewAccessPointService(composite.accessPointRepo, composite.accessPointTypeRepo, composite.accessPointRadioRepo),
		accessPointRadioService:         service.NewAccessPointRadioService(composite.accessPointRadioRepo),

		sensorTypeService:          service.NewSensorTypeService(composite.sensorTypeRepo),
		sensorRadioTemplateService: service.NewSensorRadioTemplateService(composite.sensorRadioTemplateRepo),
		sensorService:              service.NewSensorService(composite.sensorRepo, composite.sensorTypeRepo),
		sensorRadioService:         service.NewSensorRadioService(composite.sensorRadioRepo),

		matrixService: service.NewMatrixService(composite.matrixRepo),

		deviceService: service.NewDeviceService(composite.deviceRepo),
	}
}
