package models

type AccessPointTypeDetailed struct {
	AccessPointType
	RadioTemplates []*RadioTemplate `json:"radioTemplates"`
}
