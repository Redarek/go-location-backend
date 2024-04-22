package server

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"location-backend/internal/config"
)

func (s *Fiber) RegisterFiberRoutes() {

	// Helmet middleware helps secure your apps by setting various HTTP headers.
	s.App.Use(helmet.New())

	s.App.Get("/", s.HelloWorldHandler)
	s.App.Static("/static", "./static")

	s.App.Get("/health", s.healthHandler)

	// Initialize default config (Assign the middleware to /metrics)
	s.App.Get("/metrics", monitor.New())

	// Liveness and readiness probes middleware for Fiber that provides two endpoints
	// for checking the liveness and readiness state of HTTP applications.
	// default endpoints: /livez and /readyz
	s.App.Use(healthcheck.New())

	api := s.App.Group("/api")
	v1 := api.Group("/v1")

	u := v1.Group("/user")
	u.Post("/register", s.Register)
	u.Post("/login", s.Login)

	v1.Use(jwtware.New(jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte(config.App.JWTSecret)}})) // TODO: implement refresh token route

	site := v1.Group("/site")
	site.Post("/", s.CreateSite)
	site.Get("/", s.GetSite)
	site.Get("/all", s.GetSites)
	site.Get("/all/detailed", s.GetSitesDetailed)
	site.Patch("/", s.PatchUpdateSite)
	site.Patch("/sd", s.SoftDeleteSite)
	site.Patch("/restore", s.RestoreSite)

	b := v1.Group("/building")
	b.Post("/", s.CreateBuilding)
	b.Get("/", s.GetBuilding)
	b.Get("/all", s.GetBuildings)
	b.Patch("/", s.PatchUpdateBuilding)
	b.Patch("/sd", s.SoftDeleteBuilding)
	b.Patch("/restore", s.RestoreBuilding)

	f := v1.Group("/floor")
	f.Post("/", s.CreateFloor)
	f.Get("/", s.GetFloor)
	f.Get("/all", s.GetFloors)
	f.Patch("/", s.PatchUpdateFloor)
	f.Patch("/sd", s.SoftDeleteFloor)
	f.Patch("/restore", s.RestoreFloor)

	wt := v1.Group("/wallType")
	wt.Post("/", s.CreateWallType)
	wt.Get("/", s.GetWallType)
	wt.Get("/all", s.GetWallTypes)
	wt.Patch("/", s.PatchUpdateWallType)
	wt.Patch("/sd", s.SoftDeleteWallType)
	wt.Patch("/restore", s.RestoreWallType)

	w := v1.Group("/wall")
	w.Post("/", s.CreateWall)
	w.Get("/", s.GetWall)
	w.Get("/all", s.GetWalls)
	w.Patch("/", s.PatchUpdateWall)
	w.Patch("/sd", s.SoftDeleteWall)
	w.Patch("/restore", s.RestoreWall)

	apt := v1.Group("/apt")
	apt.Post("/", s.CreateAccessPointType)
	apt.Get("/", s.GetAccessPointType)
	apt.Get("/all", s.GetAccessPointTypes)
	apt.Patch("/sd", s.SoftDeleteAccessPointType)
	apt.Patch("/restore", s.RestoreAccessPointType)

	rt := v1.Group("/radioTemplate")
	rt.Post("/", s.CreateRadioTemplate)
	rt.Get("/", s.GetRadioTemplate)
	rt.Get("/all", s.GetRadioTemplates)
	rt.Patch("/", s.PatchUpdateRadioTemplate)
	rt.Patch("/sd", s.SoftDeleteRadioTemplate)
	rt.Patch("/restore", s.RestoreRadioTemplate)

	ap := v1.Group("/ap")
	ap.Post("/", s.CreateAccessPoint)
	ap.Get("/", s.GetAccessPoint)
	ap.Get("/detailed", s.GetAccessPointDetailed)
	ap.Get("/all", s.GetAccessPoints)
	ap.Get("/all/detailed", s.GetAccessPointsDetailed)
	ap.Patch("/", s.PatchUpdateAccessPoint)
	ap.Patch("/sd", s.SoftDeleteAccessPoint)
	ap.Patch("/restore", s.RestoreAccessPoint)

	r := v1.Group("/radio")
	r.Post("/", s.CreateRadio)
	r.Get("/", s.GetRadio)
	r.Get("/all", s.GetRadios)
	r.Patch("/", s.PatchUpdateRadio)
	r.Patch("/sd", s.SoftDeleteRadio)
	r.Patch("/restore", s.RestoreRadio)

	//sensor := v1.Group("/sensor")
	//sensor.Post("/", s.CreateSensor)
	//sensor.Get("/", s.GetSensor)
	//sensor.Get("/all", s.GetSensor)
	//sensor.Patch("/", s.PatchUpdateSensor)
	//sensor.Patch("/sd", s.SoftDeleteSensor)
	//sensor.Patch("/restore", s.RestoreSensor)

}

func (s *Fiber) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *Fiber) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
