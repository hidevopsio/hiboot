package controller

import "github.com/hidevopsio/hiboot/pkg/starter/web"

// Health is the health check struct
type Health struct {
	Status string `json:"status"`
}

type healthController struct {
	web.Controller
}

func init() {
	web.Add(new(healthController))
}

// GET /health
func (c *healthController) Get() {

	health := Health{
		Status: "UP",
	}
	c.Ctx.JSON(health)
}