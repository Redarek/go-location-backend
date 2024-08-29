package model

type WallDetailed struct {
	Wall
	WallType *WallType `json:"wallType"`
}
