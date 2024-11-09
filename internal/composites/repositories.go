package composites

import (
	repository "location-backend/internal/adapters/db/postgres"
	"location-backend/internal/domain/service"
)

type RepositoryComposite struct {
	healthRepo   service.IHealthRepo
	userRepo     service.IUserRepo
	roleRepo     service.IRoleRepo
	siteRepo     service.ISiteRepo
	buildingRepo service.IBuildingRepo
	floorRepo    service.IFloorRepo
	wallTypeRepo service.IWallTypeRepo
	wallRepo     service.IWallRepo

	accessPointTypeRepo          service.IAccessPointTypeRepo
	accessPointRadioTemplateRepo service.IAccessPointRadioTemplateRepo
	accessPointRepo              service.IAccessPointRepo
	accessPointRadioRepo         service.IAccessPointRadioRepo

	sensorTypeRepo          service.ISensorTypeRepo
	sensorRadioTemplateRepo service.ISensorRadioTemplateRepo
	sensorRepo              service.ISensorRepo
	sensorRadioRepo         service.ISensorRadioRepo

	matrixRepo service.IMatrixRepo

	deviceRepo service.IDeviceRepo
}

func NewRepositoryComposite(composite *PostgresComposite) (repositoryComposite *RepositoryComposite) {
	return &RepositoryComposite{
		healthRepo:   repository.NewHealthRepo(composite.pool),
		userRepo:     repository.NewUserRepo(composite.pool),
		roleRepo:     repository.NewRoleRepo(composite.pool),
		siteRepo:     repository.NewSiteRepo(composite.pool),
		buildingRepo: repository.NewBuildingRepo(composite.pool),
		floorRepo:    repository.NewFloorRepo(composite.pool),
		wallTypeRepo: repository.NewWallTypeRepo(composite.pool),
		wallRepo:     repository.NewWallRepo(composite.pool),

		accessPointTypeRepo:          repository.NewAccessPointTypeRepo(composite.pool),
		accessPointRadioTemplateRepo: repository.NewAccessPointRadioTemplateRepo(composite.pool),
		accessPointRepo:              repository.NewAccessPointRepo(composite.pool),
		accessPointRadioRepo:         repository.NewAccessPointRadioRepo(composite.pool),

		sensorTypeRepo:          repository.NewSensorTypeRepo(composite.pool),
		sensorRadioTemplateRepo: repository.NewSensorRadioTemplateRepo(composite.pool),
		sensorRepo:              repository.NewSensorRepo(composite.pool),
		sensorRadioRepo:         repository.NewSensorRadioRepo(composite.pool),

		matrixRepo: repository.NewMatrixRepo(composite.pool),

		deviceRepo: repository.NewDeviceRepo(composite.pool),
	}
}
