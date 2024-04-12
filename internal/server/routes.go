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
	v1 := api.Group("/v1") // TODO: implement refresh token route
	v1.Use(jwtware.New(jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte(config.App.JWTSecret)}}))

	user := v1.Group("/user")
	user.Post("/login", s.Login)
	user.Post("/register", s.Register)

	site := v1.Group("/site")
	site.Post("/", s.CreateSite) // create site
	site.Get("/", s.GetSite)     // get one site
	//site.Get("/full")     // get full one site
	site.Get("/all", s.GetSites) // get all sites
	//site.Get("/all/full") // get full all sites
	//site.Put("/")         // update one site
	site.Patch("/", s.PatchUpdateSite)
	site.Patch("/sd", s.SoftDeleteSite)   // soft delete one site
	site.Patch("/restore", s.RestoreSite) // restore one site
	//site.Delete("/")      // delete one site

	building := v1.Group("/building")
	building.Post("/", s.CreateBuilding)          // create building
	building.Get("/", s.GetBuilding)              // get one building
	building.Get("/all", s.GetBuildings)          // get all buildings
	building.Patch("/", s.PatchUpdateBuilding)    // patch update one floor
	building.Patch("/sd", s.SoftDeleteBuilding)   // soft delete one site
	building.Patch("/restore", s.RestoreBuilding) // restore one site
	//building.Put("/")    // update one building
	//building.Delete("/") // delete one building
	//
	floor := v1.Group("/floor")
	floor.Post("/", s.CreateFloor) // create floor
	floor.Get("/", s.GetFloor)     // get one floor
	floor.Get("/all", s.GetFloors) // get all floors
	//floor.Put("/")    // update one floor
	floor.Patch("/", s.PatchUpdateFloor)    // patch update one floor
	floor.Patch("/sd", s.SoftDeleteFloor)   // soft delete one site
	floor.Patch("/restore", s.RestoreFloor) // restore one site
	//floor.Delete("/") // delete one floor
	//
	//wallType := v1.Group("/wallType")
	//wallType.Post("/")   // create wallType
	//wallType.Get("/")    // get one wallType
	//wallType.Get("/all") // get all wallTypes
	//wallType.Put("/")    // update one wallType
	//wallType.Delete("/") // delete one wallType
	//
	//wall := v1.Group("/wall")
	//wall.Post("/")   // create wall
	//wall.Get("/")    // get one wall
	//wall.Get("/all") // get all walls
	//wall.Put("/")    // update one wall
	//wall.Patch("/")  // patch update one wall
	//wall.Delete("/") // delete one wall
	//
	//ap := v1.Group("/ap")
	//ap.Post("/")   // create access point
	//ap.Get("/")    // get one access point
	//ap.Get("/all") // get all access points
	//ap.Put("/")    // update one access point
	//ap.Patch("/")  // patch update one access point
	//ap.Delete("/") // delete one access point
	//
	//sensor := v1.Group("/sensor")
	//sensor.Post("/")   // create sensor
	//sensor.Get("/")    // get one sensor
	//sensor.Get("/all") // get all sensors
	//sensor.Put("/")    // update one sensor
	//sensor.Patch("/")  // patch update one sensor
	//sensor.Delete("/") // delete one sensor
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
