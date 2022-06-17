package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/calvarado2004/bookings/internal/models"
)

//AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	//Creates the channel MailChan from the model MailData
	MailChan chan models.MailData
}
