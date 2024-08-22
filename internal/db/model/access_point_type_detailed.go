package model

type AccessPointTypeDetailed struct {
	AccessPointType
	RadioTemplates []*RadioTemplate `json:"radioTemplates"`
}
