package db

import (
	"github.com/google/uuid"
	"time"
)

type Service interface {
	GetUserByUsername(username string) (u User, err error)
	CreateUser(username, password string) (id uuid.UUID, err error)

	CreateSite(userUUID uuid.UUID, s *Site) (id uuid.UUID, err error)
	GetSite(siteUUID uuid.UUID) (s *Site, err error)
	GetSites(userUUID uuid.UUID) (s []*Site, err error)
	SoftDeleteSite(siteUUID uuid.UUID) (err error)
	IsSiteSoftDeleted(siteUUID uuid.UUID) (check bool, err error)
	RestoreSite(siteUUID uuid.UUID) (err error)
	PatchUpdateSite(site *Site) (err error)

	CreateBuilding(b *Building) (id int, err error)
	GetBuilding(buildingUUID uuid.UUID) (b *Building, err error)
	IsBuildingSoftDeleted(buildingUUID uuid.UUID) (isDeleted bool, err error)
	GetBuildings(siteUUID uuid.UUID) (bs []*Building, err error)
	SoftDeleteBuilding(buildingUUID uuid.UUID) (err error)
	RestoreBuilding(buildingUUID uuid.UUID) (err error)
	PatchUpdateBuilding(b *Building) (err error)

	CreateFloor(f *Floor) (id int, err error)
	GetFloor(floorUUID uuid.UUID) (f *Floor, err error)
	IsFloorSoftDeleted(floorUUID uuid.UUID) (isDeleted bool, err error)
	GetFloors(buildingUUID uuid.UUID) (fs []*Floor, err error)
	SoftDeleteFloor(floorUUID uuid.UUID) (err error)
	RestoreFloor(floorUUID uuid.UUID) (err error)
	PatchUpdateFloor(f *Floor) (err error)

	CreateWallType(wt *WallType) (id int, err error)
	GetWallType(wallTypeUUID uuid.UUID) (wt *WallType, err error)
	IsWallTypeSoftDeleted(wallTypeUUID uuid.UUID) (isDeleted bool, err error)
	GetWallTypes(siteUUID uuid.UUID) (wts []*WallType, err error)
	SoftDeleteWallType(wallTypeUUID uuid.UUID) (err error)
	RestoreWallType(wallTypeUUID uuid.UUID) (err error)
	PatchUpdateWallType(wt *WallType) (err error)

	CreateWall(w *Wall) (id int, err error)
	GetWall(wallUUID uuid.UUID) (w *Wall, err error)
	IsWallSoftDeleted(wallUUID uuid.UUID) (isDeleted bool, err error)
	GetWalls(floorUUID uuid.UUID) (ws []*Wall, err error)
	SoftDeleteWall(wallUUID uuid.UUID) (err error)
	RestoreWall(wallUUID uuid.UUID) (err error)
	PatchUpdateWall(w *Wall) (err error)

	CreateAccessPointType(apt *AccessPointType) (id int, err error)
	GetAccessPointType(accessPointTypeUUID uuid.UUID) (apt *AccessPointType, err error)
	IsAccessPointTypeSoftDeleted(accessPointTypeUUID uuid.UUID) (isDeleted bool, err error)
	GetAccessPointTypes(siteUUID uuid.UUID) (apts []*AccessPointType, err error)
	SoftDeleteAccessPointType(accessPointTypeUUID uuid.UUID) (err error)
	RestoreAccessPointType(accessPointTypeUUID uuid.UUID) (err error)

	CreateRadio(r *Radio) (id int, err error)
	GetRadio(radioUUID uuid.UUID) (r Radio, err error)
	IsRadioSoftDeleted(radioUUID uuid.UUID) (isDeleted bool, err error)
	GetRadios(accessPointTypeUUID uuid.UUID) (rs []*Radio, err error)
	SoftDeleteRadio(radioUUID uuid.UUID) (err error)
	RestoreRadio(radioUUID uuid.UUID) (err error)
	PatchUpdateRadio(r *Radio) (err error)

	CreateAccessPoint(ap *AccessPoint) (id int, err error)
	GetAccessPoint(accessPointUUID uuid.UUID) (ap *AccessPoint, err error)
	IsAccessPointSoftDeleted(accessPointUUID uuid.UUID) (isDeleted bool, err error)
	GetAccessPoints(floorUUID uuid.UUID) (aps []*AccessPoint, err error)
	SoftDeleteAccessPoint(accessPointUUID uuid.UUID) (err error)
	RestoreAccessPoint(accessPointUUID uuid.UUID) (err error)
	PatchUpdateAccessPoint(ap *AccessPoint) (err error)

	Health() map[string]string
}

type User struct {
	ID        uuid.UUID  `db:"id"`
	Username  string     `db:"username"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"` // Используйте указатель для NULL-значений
}

type Role struct {
	ID        uuid.UUID  `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type UserRole struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	RoleID    uuid.UUID  `db:"role_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type RefreshToken struct {
	ID     uuid.UUID  `db:"id"`
	Token  string     `db:"token"`
	Expiry *time.Time `db:"expiry"`
	UserID uuid.UUID  `db:"user_id"`
}

type Site struct {
	ID          uuid.UUID  `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	UserID      uuid.UUID  `db:"user_id"`
}

type Building struct {
	ID          uuid.UUID  `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	Country     string     `db:"country"`
	City        string     `db:"city"`
	Address     string     `db:"address"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	SiteID      uuid.UUID  `db:"site_id"`
}

type Floor struct {
	ID         uuid.UUID  `db:"id"`
	Name       *string    `db:"name"`
	Number     *int       `db:"number"`
	Image      *string    `db:"image"`
	Scale      *float64   `db:"scale"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
	BuildingID uuid.UUID  `db:"building_id"`
}

type AccessPoint struct {
	ID                uuid.UUID  `db:"id"`
	Name              string     `db:"name"`
	X                 *int       `db:"x"`
	Y                 *int       `db:"y"`
	Z                 *float64   `db:"z"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
	FloorID           uuid.UUID  `db:"floor_id"`
	AccessPointTypeID uuid.UUID  `db:"access_point_type_id"`
}

type AccessPointType struct {
	ID        uuid.UUID  `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	SiteID    uuid.UUID  `db:"site_id"`
}

type Radio struct {
	ID                uuid.UUID  `db:"id"`
	Number            string     `db:"number"`
	Channel           *int       `db:"channel"`
	WiFi              string     `db:"wifi"`
	Power             *int       `db:"power"`
	Bandwidth         string     `db:"bandwidth"`
	GuardInterval     *int       `db:"guard_interval"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
	AccessPointTypeID uuid.UUID  `db:"access_point_type_id"`
}

type Wall struct {
	ID         uuid.UUID  `db:"id"`
	X1         *int       `db:"x1"`
	Y1         *int       `db:"y1"`
	X2         *int       `db:"x2"`
	Y2         *int       `db:"y2"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
	FloorID    uuid.UUID  `db:"floor_id"`
	WallTypeID uuid.UUID  `db:"wall_type_id"`
}

type WallType struct {
	ID           uuid.UUID  `db:"id"`
	Name         string     `db:"name"`
	Color        string     `db:"color"`
	Attenuation1 *float64   `db:"attenuation1"`
	Attenuation2 *float64   `db:"attenuation2"`
	Attenuation3 *float64   `db:"attenuation3"`
	Thickness    *float64   `db:"thickness"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	SiteID       uuid.UUID  `db:"site_id"`
}
