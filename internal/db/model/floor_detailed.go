package model

type FloorDetailed struct {
	Floor
	AccessPoints []*AccessPointDetailed `json:"accessPoints"`
	Walls        []*WallDetailed        `json:"walls"`
	Sensors      []*Sensor              `json:"sensors"`
}
