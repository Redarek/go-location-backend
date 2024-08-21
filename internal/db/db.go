package db

import (
	"github.com/google/uuid"
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
	UpdateFloorHeatmap(floorUUID uuid.UUID, fileName string) (err error)

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

	CreateAccessPointType(apt *AccessPointType) (id uuid.UUID, err error)
	GetAccessPointType(accessPointTypeUUID uuid.UUID) (apt *AccessPointType, err error)
	GetAccessPointTypeDetailed(accessPointTypeUUID uuid.UUID) (apt *AccessPointTypeDetailed, err error)
	IsAccessPointTypeSoftDeleted(accessPointTypeUUID uuid.UUID) (isDeleted bool, err error)
	GetAccessPointTypes(siteUUID uuid.UUID) (apts []*AccessPointType, err error)
	GetAccessPointTypesDetailed(siteUUID uuid.UUID) (aps []*AccessPointTypeDetailed, err error)
	SoftDeleteAccessPointType(accessPointTypeUUID uuid.UUID) (err error)
	RestoreAccessPointType(accessPointTypeUUID uuid.UUID) (err error)
	PatchUpdateAccessPointType(apt *AccessPointType) (err error)

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

	CreateSensorType(s *SensorType) (id uuid.UUID, err error)
	GetSensorType(sensorTypeUUID uuid.UUID) (s *SensorType, err error)
	IsSensorTypeSoftDeleted(sensorTypeUUID uuid.UUID) (isDeleted bool, err error)
	GetSensorTypes(siteUUID uuid.UUID) (ss []*SensorType, err error)
	SoftDeleteSensorType(sensorTypeUUID uuid.UUID) (err error)
	RestoreSensorType(sensorTypeUUID uuid.UUID) (err error)
	PatchUpdateSensorType(s *SensorType) (err error)

	CreateSensor(s *Sensor) (id uuid.UUID, err error)
	GetSensor(sensorUUID uuid.UUID) (s *Sensor, err error)
	GetSensorDetailed(sensorUUID uuid.UUID) (s *SensorDetailed, err error)
	IsSensorSoftDeleted(sensorUUID uuid.UUID) (isDeleted bool, err error)
	GetSensors(floorUUID uuid.UUID) (ss []*Sensor, err error)
	GetSensorsDetailed(floorUUID uuid.UUID) (ss []*SensorDetailed, err error)
	SoftDeleteSensor(sensorUUID uuid.UUID) (err error)
	RestoreSensor(sensorUUID uuid.UUID) (err error)
	PatchUpdateSensor(s *Sensor) (err error)

	Health() map[string]string
}

// type RadiationDiagram struct {
// 	SensorID uuid.UUID       `json:"sensorId" db:"sensor_id"`

// }
