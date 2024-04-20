package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"location-backend/internal/db"
)

// CreateSite creates a site
func (s *Fiber) CreateSite(c *fiber.Ctx) (err error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userUUID, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse user uuid")
		return
	}
	siteInput := new(db.Site)
	err = c.BodyParser(siteInput)
	if err != nil {
		return err
	}

	siteID, err := s.db.CreateSite(userUUID, siteInput)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"id": siteID,
	})
}

// GetSite retrieves a site
func (s *Fiber) GetSite(c *fiber.Ctx) (err error) {
	siteUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse site uuid")
		return
	}
	site, err := s.db.GetSite(siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get site")
		return
	}
	return c.JSON(fiber.Map{
		"data": site,
	})
}

// GetSites retrieves sites
func (s *Fiber) GetSites(c *fiber.Ctx) (err error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userUUID, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse user uuid")
		return
	}
	sites, err := s.db.GetSites(userUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sites")
		return
	}
	return c.JSON(fiber.Map{
		"data": sites,
	})
}

// GetSitesDetailed retrieves sites detailed
func (s *Fiber) GetSitesDetailed(c *fiber.Ctx) (err error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userUUID, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse user uuid")
		return
	}
	sites, err := s.db.GetSites(userUUID)
	for _, site := range sites {
		log.Debug().Msgf("Site: %s", site)
		buildings, err := s.db.GetBuildings(site.ID)
		if err != nil {
			continue
		}

		for _, building := range buildings {
			log.Debug().Msgf("Building: %s", building)
			floors, err := s.db.GetFloors(building.ID)
			if err != nil {
				continue
			}

			for _, floor := range floors {
				log.Debug().Msgf("Floor: %s", floor)
				aps, err := s.db.GetAccessPointsDetailed(floor.ID)
				if err != nil {
					continue
				}
				walls, err := s.db.GetWallsDetailed(floor.ID)
				floor.AccessPoints = aps
				floor.Walls = walls
			}
			building.Floors = floors
		}
		site.Buildings = buildings

		wallTypes, err := s.db.GetWallTypes(site.ID)
		if err != nil {
			continue
		}
		site.WallTypes = wallTypes

		accessPointTypes, err := s.db.GetAccessPointTypes(site.ID)
		if err != nil {
			continue
		}
		site.AccessPointTypes = accessPointTypes

	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sites")
		return
	}
	return c.JSON(fiber.Map{
		"data": sites,
	})
}

// SoftDeleteSite soft delete a site
func (s *Fiber) SoftDeleteSite(c *fiber.Ctx) (err error) {
	siteUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse site uuid")
		return
	}
	isDeleted, err := s.db.IsSiteSoftDeleted(siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted site")
		return
	}
	if !isDeleted {
		err = s.db.SoftDeleteSite(siteUUID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to soft delete a site")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Site has already been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// RestoreSite restore a site
func (s *Fiber) RestoreSite(c *fiber.Ctx) (err error) {
	siteUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse site uuid")
		return
	}
	isDeleted, err := s.db.IsSiteSoftDeleted(siteUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get soft deleted site")
		return
	}
	if isDeleted {
		err = s.db.RestoreSite(siteUUID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to restore a site")
			return
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Site has not been soft deleted")
	}
	return c.SendStatus(fiber.StatusOK)
}

// PatchUpdateSite patch updates a site based on provided fields
func (s *Fiber) PatchUpdateSite(c *fiber.Ctx) error {
	var input db.Site
	if err := c.BodyParser(&input); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	if err := s.db.PatchUpdateSite(&input); err != nil {
		log.Error().Err(err).Msg("Failed to update site")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update site")
	}

	return c.SendStatus(fiber.StatusOK)
}
