package models

type WallDetailed struct {
	Wall
	WallType *WallType `json:"wallType"`
}
