package main

import (
	"wsan/handlers"
	"wsan/middleware"
	"wsan/models"

	"github.com/cidekar/adele-framework"
	"github.com/cidekar/adele-framework/mailer"
)

type application struct {
	App        *adele.Adele
	Handlers   *handlers.Handlers
	Mail       *mailer.Mail
	Middleware *middleware.Middleware
	Models     *models.Models
}
