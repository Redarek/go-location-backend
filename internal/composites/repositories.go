package composites

import (
	repository "location-backend/internal/adapters/db/postgres"
	"location-backend/internal/domain/service"
)

type RepositoryComposite struct {
	healthRepo   service.HealthRepo
	userRepo     service.UserRepo
	roleRepo     service.RoleRepo
	siteRepo     service.SiteRepo
	buildingRepo service.BuildingRepo
	floorRepo    service.FloorRepo
	wallTypeRepo service.WallTypeRepo
	wallRepo     service.WallRepo

	accessPointTypeRepo          service.AccessPointTypeRepo
	accessPointRadioTemplateRepo service.AccessPointRadioTemplateRepo
	accessPointRepo              service.AccessPointRepo
	accessPointRadioRepo         service.AccessPointRadioRepo

	sensorTypeRepo          service.SensorTypeRepo
	sensorRadioTemplateRepo service.SensorRadioTemplateRepo
	sensorRepo              service.SensorRepo
	sensorRadioRepo         service.SensorRadioRepo
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
	}
}
