package model

type AccessPointDetailed struct {
	AccessPoint
	AccessPointType *AccessPointType `json:"accessPointType"`
	Radios          []*Radio         `json:"radios"`
}
