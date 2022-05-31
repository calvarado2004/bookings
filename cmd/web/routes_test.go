package main

import (
	"fmt"
	"testing"

	"github.com/calvarado2004/bookings/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch mux.(type) {
	case *chi.Mux:
		//nothing
	default:
		t.Logf(fmt.Sprintf("type is no the expected but it is %T", mux))
	}
}
