package server

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"location-backend/internal/config"
)

func (s *Fiber) RegisterFiberRoutes() {

	s.App.Get("/", s.HelloWorldHandler)
	s.App.Static("/static", "./static")

	s.App.Get("/health", s.healthHandler)

	api := s.App.Group("/api")
	v1 := api.Group("/v1")

	user := v1.Group("/user")
	user.Post("/register", s.Register)
	user.Post("/login", s.Login)

	v1.Use(jwtware.New(jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte(config.App.JWTSecret)}})) // TODO: implement refresh token route

	site := v1.Group("/site")
	site.Post("/", s.CreateSite)
	site.Get("/", s.GetSite)
	site.Get("/all", s.GetSites)
	site.Patch("/", s.PatchUpdateSite)
	site.Patch("/sd", s.SoftDeleteSite)
	site.Patch("/restore", s.RestoreSite)

	building := v1.Group("/building")
	building.Post("/", s.CreateBuilding)
	building.Get("/", s.GetBuilding)
	building.Get("/all", s.GetBuildings)
	building.Patch("/", s.PatchUpdateBuilding)
	building.Patch("/sd", s.SoftDeleteBuilding)
	building.Patch("/restore", s.RestoreBuilding)

	floor := v1.Group("/floor")
	floor.Post("/", s.CreateFloor)
	floor.Get("/", s.GetFloor)
	floor.Get("/all", s.GetFloors)
	floor.Patch("/", s.PatchUpdateFloor)
	floor.Patch("/sd", s.SoftDeleteFloor)
	floor.Patch("/restore", s.RestoreFloor)

	wallType := v1.Group("/wallType")
	wallType.Post("/", s.CreateWallType)
	wallType.Get("/", s.GetWallType)
	wallType.Get("/all", s.GetWallTypes)
	wallType.Patch("/", s.PatchUpdateWallType)
	wallType.Patch("/sd", s.SoftDeleteWallType)
	wallType.Patch("/restore", s.RestoreWallType)

	wall := v1.Group("/wall")
	wall.Post("/", s.CreateWall)
	wall.Get("/", s.GetWall)
	wall.Get("/all", s.GetWalls)
	wall.Patch("/", s.PatchUpdateWall)
	wall.Patch("/sd", s.SoftDeleteWall)
	wall.Patch("/restore", s.RestoreWall)

	apt := v1.Group("/apt")
	apt.Post("/", s.CreateAccessPointType)
	apt.Get("/", s.GetAccessPointType)
	apt.Get("/all", s.GetAccessPointTypes)
	apt.Patch("/sd", s.SoftDeleteAccessPointType)
	apt.Patch("/restore", s.RestoreAccessPointType)

	radio := v1.Group("/radio")
	radio.Post("/", s.CreateRadio)
	radio.Get("/", s.GetRadio)
	radio.Get("/all", s.GetRadios)
	radio.Patch("/", s.PatchUpdateRadio)
	radio.Patch("/sd", s.SoftDeleteRadio)
	radio.Patch("/restore", s.RestoreRadio)

	ap := v1.Group("/ap")
	ap.Post("/", s.CreateAccessPoint)
	ap.Get("/", s.GetAccessPoint)
	ap.Get("/all", s.GetAccessPoints)
	ap.Patch("/", s.PatchUpdateAccessPoint)
	ap.Patch("/sd", s.SoftDeleteAccessPoint)
	ap.Patch("/restore", s.RestoreAccessPoint)

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
