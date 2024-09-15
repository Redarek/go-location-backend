package v1

import "location-backend/internal/router"

type Handler interface {
	Register(router *router.Router)
}
