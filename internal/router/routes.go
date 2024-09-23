package router

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"location-backend/internal/composites"
	"location-backend/internal/config"
	"location-backend/internal/middleware"
)

func RegisterRoutes(router *Router, handlerComposite *composites.HandlerComposite) {
	// CORS
	router.App.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     config.App.ClientURL,
		AllowCredentials: true,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		}, ","),
	}))

	// Public route
	router.App.Static("/public", "/public")

	api := router.App.Group("/api")
	v1 := api.Group("/v1")

	handlerComposite.HealthHandler.Register(&v1)

	user := v1.Group("/user")
	handlerComposite.UserHandler.Register(&user)

	role := v1.Group("/role", middleware.Auth)
	handlerComposite.RoleHandler.Register(&role)

	site := v1.Group("/site", middleware.Auth)
	handlerComposite.SiteHandler.Register(&site)

	building := v1.Group("/building", middleware.Auth)
	handlerComposite.BuildingHandler.Register(&building)
}

// import (
// 	"location-backend/internal/config"

// 	jwtware "github.com/gofiber/contrib/jwt"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/healthcheck"
// 	"github.com/gofiber/fiber/v2/middleware/monitor"
// )

// func (s *Fiber) RegisterFiberRoutes() {
// 	s.App.Use(cors.New(cors.Config{
// 		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
// 		AllowOrigins:     config.App.ClientURL,
// 		AllowCredentials: true,
// 		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
// 	}))

// 	//s.App.Use(limiter.New(limiter.Config{
// 	//	Max:        100,
// 	//	Expiration: 1 * time.Minute,
// 	//	KeyGenerator: func(c *fiber.Ctx) string {
// 	//		return c.IP()
// 	//	},
// 	//	LimitReached: func(c *fiber.Ctx) error {
// 	//		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
// 	//			"error": "Too Many Requests",
// 	//		})
// 	//	},
// 	//}))

// 	//s.App.Use(helmet.New())

// 	s.App.Static("/static", "./static")

// 	s.App.Get("/health", s.healthHandler)

// 	s.App.Get("/metrics", monitor.New())

// 	s.App.Use(healthcheck.New())

// 	api := s.App.Group("/api")
// 	v1 := api.Group("/v1")

// 	user := v1.Group("/user")
// 	user.Post("/register", s.Register)
// 	user.Post("/login", s.Login)

// 	v1.Use(jwtware.New(jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte(config.App.JWTSecret)}}))

// 	site := v1.Group("/site")
// 	site.Post("/", s.CreateSite)
// 	site.Get("/", s.GetSite)
// 	site.Get("/all", s.GetSites)
// 	site.Get("/all/detailed", s.GetSitesDetailed)
// 	site.Patch("/", s.PatchUpdateSite)
// 	site.Patch("/sd", s.SoftDeleteSite)
// 	site.Patch("/restore", s.RestoreSite)

// 	building := v1.Group("/building")
// 	building.Post("/", s.CreateBuilding)
// 	building.Get("/", s.GetBuilding)
// 	building.Get("/all", s.GetBuildings)
// 	building.Patch("/", s.PatchUpdateBuilding)
// 	building.Patch("/sd", s.SoftDeleteBuilding)
// 	building.Patch("/restore", s.RestoreBuilding)

// 	floor := v1.Group("/floor")
// 	floor.Post("/", s.CreateFloor)
// 	floor.Get("/", s.GetFloor)
// 	floor.Get("/all", s.GetFloors)
// 	floor.Patch("/", s.PatchUpdateFloor)
// 	floor.Patch("/sd", s.SoftDeleteFloor)
// 	floor.Patch("/restore", s.RestoreFloor)

// 	wallType := v1.Group("/wallType")
// 	wallType.Post("/", s.CreateWallType)
// 	wallType.Get("/", s.GetWallType)
// 	wallType.Get("/all", s.GetWallTypes)
// 	wallType.Patch("/", s.PatchUpdateWallType)
// 	wallType.Patch("/sd", s.SoftDeleteWallType)
// 	wallType.Patch("/restore", s.RestoreWallType)

// 	wall := v1.Group("/wall")
// 	wall.Post("/", s.CreateWall)
// 	wall.Get("/", s.GetWall)
// 	wall.Get("/all", s.GetWalls)
// 	wall.Patch("/", s.PatchUpdateWall)
// 	wall.Patch("/sd", s.SoftDeleteWall)
// 	wall.Patch("/restore", s.RestoreWall)

// 	accessPointType := v1.Group("/apt")
// 	accessPointType.Post("/", s.CreateAccessPointType)
// 	accessPointType.Get("/", s.GetAccessPointType)
// 	accessPointType.Get("/all", s.GetAccessPointTypes)
// 	accessPointType.Patch("/", s.PatchUpdateAccessPointType)
// 	accessPointType.Patch("/sd", s.SoftDeleteAccessPointType)
// 	accessPointType.Patch("/restore", s.RestoreAccessPointType)

// 	radioTemplate := v1.Group("/radioTemplate")
// 	radioTemplate.Post("/", s.CreateRadioTemplate)
// 	radioTemplate.Get("/", s.GetRadioTemplate)
// 	radioTemplate.Get("/all", s.GetRadioTemplates)
// 	radioTemplate.Patch("/", s.PatchUpdateRadioTemplate)
// 	radioTemplate.Patch("/sd", s.SoftDeleteRadioTemplate)
// 	radioTemplate.Patch("/restore", s.RestoreRadioTemplate)

// 	accessPoint := v1.Group("/ap")
// 	accessPoint.Post("/", s.CreateAccessPoint)
// 	accessPoint.Get("/", s.GetAccessPoint)
// 	accessPoint.Get("/detailed", s.GetAccessPointDetailed)
// 	accessPoint.Get("/all", s.GetAccessPoints)
// 	accessPoint.Get("/all/detailed", s.GetAccessPointsDetailed)
// 	accessPoint.Patch("/", s.PatchUpdateAccessPoint)
// 	accessPoint.Patch("/sd", s.SoftDeleteAccessPoint)
// 	accessPoint.Patch("/restore", s.RestoreAccessPoint)

// 	radio := v1.Group("/radio")
// 	radio.Post("/", s.CreateRadio)
// 	radio.Get("/", s.GetRadio)
// 	radio.Get("/all", s.GetRadios)
// 	radio.Patch("/", s.PatchUpdateRadio)
// 	radio.Patch("/sd", s.SoftDeleteRadio)
// 	radio.Patch("/restore", s.RestoreRadio)

// 	sensorType := v1.Group("/sensorType")
// 	sensorType.Post("/", s.CreateSensorType)
// 	sensorType.Get("/", s.GetSensorType)
// 	sensorType.Get("/all", s.GetSensorTypes)
// 	sensorType.Patch("/", s.PatchUpdateSensorType)
// 	sensorType.Patch("/sd", s.SoftDeleteSensorType)
// 	sensorType.Patch("/restore", s.RestoreSensorType)

// 	sensor := v1.Group("/sensor")
// 	sensor.Post("/", s.CreateSensor)
// 	sensor.Get("/", s.GetSensor)
// 	sensor.Get("/detailed", s.GetSensorDetailed)
// 	sensor.Get("/all", s.GetSensors)
// 	sensor.Get("/all/detailed", s.GetSensorsDetailed)
// 	sensor.Patch("/", s.PatchUpdateSensor)
// 	sensor.Patch("/sd", s.SoftDeleteSensor)
// 	sensor.Patch("/restore", s.RestoreSensor)

// 	matrix := v1.Group("/matrix")
// 	matrix.Post("/", s.CreateMatrix)

// }

// func (s *Fiber) HelloWorldHandler(c *fiber.Ctx) error {
// 	resp := fiber.Map{
// 		"message": "Hello World",
// 	}

// 	return c.JSON(resp)
// }

// func (s *Fiber) healthHandler(c *fiber.Ctx) error {
// 	return c.JSON(s.db.Health())
// }
