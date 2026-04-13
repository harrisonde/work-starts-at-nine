package middleware

import (
	"wsan/models"

	"github.com/cidekar/adele-framework"
)

type Middleware struct {
	App    *adele.Adele
	Models *models.Models
}
