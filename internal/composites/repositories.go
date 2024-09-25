package composites

import (
	repository "location-backend/internal/adapters/db/postgres"
)

type RepositoryComposite struct {
	healthRepo   repository.HealthRepo
	userRepo     repository.UserRepo
	roleRepo     repository.RoleRepo
	siteRepo     repository.SiteRepo
	buildingRepo repository.BuildingRepo
	floorRepo    repository.FloorRepo
	wallTypeRepo repository.WallTypeRepo
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
	}
}
