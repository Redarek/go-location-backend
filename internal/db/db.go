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

	CreateBuilding(b *Building) (id uuid.UUID, err error)
	GetBuilding(buildingUUID uuid.UUID) (b *Building, err error)
	IsBuildingSoftDeleted(buildingUUID uuid.UUID) (isDeleted bool, err error)
	GetBuildings(siteUUID uuid.UUID) (bs []*Building, err error)
	SoftDeleteBuilding(buildingUUID uuid.UUID) (err error)
	RestoreBuilding(buildingUUID uuid.UUID) (err error)
	PatchUpdateBuilding(b *Building) (err error)

	CreateFloor(f *Floor) (id uuid.UUID, err error)
	GetFloor(floorUUID uuid.UUID) (f *Floor, err error)
	IsFloorSoftDeleted(floorUUID uuid.UUID) (isDeleted bool, err error)
	GetFloors(buildingUUID uuid.UUID) (fs []*Floor, err error)
	SoftDeleteFloor(floorUUID uuid.UUID) (err error)
	RestoreFloor(floorUUID uuid.UUID) (err error)
	PatchUpdateFloor(f *Floor) (err error)

	CreateWallType(wt *WallType) (id uuid.UUID, err error)
	GetWallType(wallTypeUUID uuid.UUID) (wt *WallType, err error)
	IsWallTypeSoftDeleted(wallTypeUUID uuid.UUID) (isDeleted bool, err error)
	GetWallTypes(siteUUID uuid.UUID) (wts []*WallType, err error)
	SoftDeleteWallType(wallTypeUUID uuid.UUID) (err error)
	RestoreWallType(wallTypeUUID uuid.UUID) (err error)
	PatchUpdateWallType(wt *WallType) (err error)

	CreateWall(w *Wall) (id uuid.UUID, err error)
	GetWall(wallUUID uuid.UUID) (w *Wall, err error)
	IsWallSoftDeleted(wallUUID uuid.UUID) (isDeleted bool, err error)
	GetWalls(floorUUID uuid.UUID) (ws []*Wall, err error)
	GetWallsDetailed(floorUUID uuid.UUID) (walls []*WallDetailed, err error)
	SoftDeleteWall(wallUUID uuid.UUID) (err error)
	RestoreWall(wallUUID uuid.UUID) (err error)
	PatchUpdateWall(w *Wall) (err error)

	CreateAccessPointType(apt *AccessPointType) (id uuid.UUID, err error) // TODO: add color for apt
	GetAccessPointType(accessPointTypeUUID uuid.UUID) (apt *AccessPointType, err error)
	GetAccessPointTypeDetailed(accessPointTypeUUID uuid.UUID) (apt *AccessPointTypeDetailed, err error)
	IsAccessPointTypeSoftDeleted(accessPointTypeUUID uuid.UUID) (isDeleted bool, err error)
	GetAccessPointTypes(siteUUID uuid.UUID) (apts []*AccessPointType, err error)
	GetAccessPointTypesDetailed(siteUUID uuid.UUID) (aps []*AccessPointTypeDetailed, err error)
	SoftDeleteAccessPointType(accessPointTypeUUID uuid.UUID) (err error)
	RestoreAccessPointType(accessPointTypeUUID uuid.UUID) (err error)

	CreateRadioTemplate(r *RadioTemplate) (id uuid.UUID, err error)
	GetRadioTemplate(radioUUID uuid.UUID) (r RadioTemplate, err error)
	IsRadioTemplateSoftDeleted(radioUUID uuid.UUID) (isDeleted bool, err error)
	GetRadioTemplates(accessPointTypeID uuid.UUID) (rs []*RadioTemplate, err error)
	SoftDeleteRadioTemplate(radioUUID uuid.UUID) (err error)
	RestoreRadioTemplate(radioUUID uuid.UUID) (err error)
	PatchUpdateRadioTemplate(r *RadioTemplate) (err error)

	CreateRadio(r *Radio) (id uuid.UUID, err error)
	GetRadio(radioUUID uuid.UUID) (r Radio, err error)
	IsRadioSoftDeleted(radioUUID uuid.UUID) (isDeleted bool, err error)
	GetRadios(accessPointTypeUUID uuid.UUID) (rs []*Radio, err error)
	SoftDeleteRadio(radioUUID uuid.UUID) (err error)
	RestoreRadio(radioUUID uuid.UUID) (err error)
	PatchUpdateRadio(r *Radio) (err error)

	CreateAccessPoint(ap *AccessPoint) (id uuid.UUID, err error)
	GetAccessPoint(accessPointUUID uuid.UUID) (ap *AccessPoint, err error)
	GetAccessPointDetailed(accessPointUUID uuid.UUID) (ap *AccessPointDetailed, err error)
	IsAccessPointSoftDeleted(accessPointUUID uuid.UUID) (isDeleted bool, err error)
	GetAccessPoints(floorUUID uuid.UUID) (aps []*AccessPoint, err error)
	GetAccessPointsDetailed(floorUUID uuid.UUID) (aps []*AccessPointDetailed, err error)
	SoftDeleteAccessPoint(accessPointUUID uuid.UUID) (err error)
	RestoreAccessPoint(accessPointUUID uuid.UUID) (err error)
	PatchUpdateAccessPoint(ap *AccessPoint) (err error)

	//SetRadioState(rs *RadioState) (id uuid.UUID, err error)
	//GetRadioStates(accessPointID uuid.UUID) (radioStates []RadioState, err error)

	Health() map[string]string
}

type User struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Username  string     `json:"username" db:"username"`
	Password  string     `json:"password" db:"password"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type Role struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type UserRole struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"userId" db:"user_id"`
	RoleID    uuid.UUID  `json:"roleId" db:"role_id"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type RefreshToken struct {
	ID     uuid.UUID  `json:"id" db:"id"`
	Token  string     `json:"token" db:"token"`
	Expiry *time.Time `json:"expiry" db:"expiry"`
	UserID uuid.UUID  `json:"userId" db:"user_id"`
}

type Site struct {
	ID               uuid.UUID          `json:"id" db:"id"`
	Name             string             `json:"name" db:"name"`
	Description      string             `json:"description" db:"description"`
	CreatedAt        time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time          `json:"updatedAt" db:"updated_at"`
	DeletedAt        *time.Time         `json:"deletedAt" db:"deleted_at"`
	UserID           uuid.UUID          `json:"userId" db:"user_id"`
	Buildings        []*Building        `json:"buildings"`
	AccessPointTypes []*AccessPointType `json:"accessPointTypes"`
	WallTypes        []*WallType        `json:"wallTypes"`
}

type Building struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Country     string     `json:"country" db:"country"`
	City        string     `json:"city" db:"city"`
	Address     string     `json:"address" db:"address"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt   *time.Time `json:"deletedAt" db:"deleted_at"`
	SiteID      uuid.UUID  `json:"siteId" db:"site_id"`
	Floors      []*Floor   `json:"floors"`
}

type Floor struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	Name         *string                `json:"name" db:"name"`
	Number       *int                   `json:"number" db:"number"`
	Image        *string                `json:"image" db:"image"`
	Scale        *float64               `json:"scale" db:"scale"`
	CreatedAt    time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time              `json:"updatedAt" db:"updated_at"`
	DeletedAt    *time.Time             `json:"deletedAt" db:"deleted_at"`
	BuildingID   uuid.UUID              `json:"buildingId" db:"building_id"`
	AccessPoints []*AccessPointDetailed `json:"accessPoints"`
	Walls        []*WallDetailed        `json:"walls"`
	//Sensors      []*Sensor      `json:"sensors"`
}

type AccessPoint struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Name              string     `json:"name" db:"name"`
	X                 *int       `json:"x" db:"x"`
	Y                 *int       `json:"y" db:"y"`
	Z                 *float64   `json:"z" db:"z"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt         *time.Time `json:"deletedAt" db:"deleted_at"`
	FloorID           uuid.UUID  `json:"floorId" db:"floor_id"`
	AccessPointTypeID uuid.UUID  `json:"accessPointTypeId" db:"access_point_type_id"`
}

type AccessPointType struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Color     string     `json:"color" db:"color"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
	SiteID    uuid.UUID  `json:"siteId" db:"site_id"`
}

type RadioTemplate struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Number            *int       `json:"number" db:"number"`
	Channel           *int       `json:"channel" db:"channel"`
	WiFi              *string    `json:"wifi" db:"wifi"`
	Power             *int       `json:"power" db:"power"`
	Bandwidth         *string    `json:"bandwidth" db:"bandwidth"`
	GuardInterval     *int       `json:"guardInterval" db:"guard_interval"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt         *time.Time `json:"deletedAt" db:"deleted_at"`
	AccessPointTypeID uuid.UUID  `json:"accessPointTypeId" db:"access_point_type_id"`
}

type Radio struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Number        *int       `json:"number" db:"number"`
	Channel       *int       `json:"channel" db:"channel"`
	WiFi          *string    `json:"wifi" db:"wifi"`
	Power         *int       `json:"power" db:"power"`
	Bandwidth     *string    `json:"bandwidth" db:"bandwidth"`
	GuardInterval *int       `json:"guardInterval" db:"guard_interval"`
	IsActive      *bool      `json:"isActive" db:"is_active"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	AccessPointID uuid.UUID  `json:"accessPointId" db:"access_point_id"`
}

//type RadioState struct {
//	AccessPointID uuid.UUID `json:"accessPointId" db:"access_point_id"`
//	RadioID       uuid.UUID `json:"radioId" db:"radio_id"`
//	IsActive      bool      `json:"isActive" db:"is_active"`
//}

type AccessPointDetailed struct {
	AccessPoint
	AccessPointType *AccessPointType `json:"accessPointType"`
	Radios          []*Radio         `json:"radios"`
}

type AccessPointTypeDetailed struct {
	AccessPointType
	RadioTemplates []*RadioTemplate `json:"radioTemplates"`
}

type Wall struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	X1         *int       `json:"x1" db:"x1"`
	Y1         *int       `json:"y1" db:"y1"`
	X2         *int       `json:"x2" db:"x2"`
	Y2         *int       `json:"y2" db:"y2"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt  *time.Time `json:"deletedAt" db:"deleted_at"`
	FloorID    uuid.UUID  `json:"floorId" db:"floor_id"`
	WallTypeID uuid.UUID  `json:"wallTypeId" db:"wall_type_id"`
}

type WallType struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Color        string     `json:"color" db:"color"`
	Attenuation1 *float64   `json:"attenuation1" db:"attenuation1"`
	Attenuation2 *float64   `json:"attenuation2" db:"attenuation2"`
	Attenuation3 *float64   `json:"attenuation3" db:"attenuation3"`
	Thickness    *float64   `json:"thickness" db:"thickness"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
	SiteID       uuid.UUID  `json:"siteId" db:"site_id"`
}

type WallDetailed struct {
	Wall
	WallType *WallType `json:"wallType"`
}
