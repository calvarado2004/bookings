package main

import (
	"net/http"

	"github.com/calvarado2004/bookings/internal/config"
	"github.com/calvarado2004/bookings/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	//mux := pat.New()
	//mux.Get("/hello-world", http.HandlerFunc(handlers.Repo.Home))
	//mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Mount("/bookings", bookingsRoutes())

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/bookings/static/*", http.StripPrefix("/bookings/static", fileServer))

	return mux
}

func bookingsRoutes() http.Handler {

	r := chi.NewRouter()
	r.Get("/", handlers.Repo.Home)
	r.Get("/about", handlers.Repo.About)
	r.Get("/generals-quarters", handlers.Repo.Generals)
	r.Get("/majors-suite", handlers.Repo.Majors)
	r.Get("/specials", handlers.Repo.Specials)
	r.Get("/reservation", handlers.Repo.Reservation)
	r.Post("/reservation", handlers.Repo.PostReservation)
	r.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	r.Get("/contact", handlers.Repo.Contact)
	r.Get("/availability", handlers.Repo.Availability)
	r.Post("/availability", handlers.Repo.PostAvailability)
	r.Post("/availability-json", handlers.Repo.AvailabilityJSON)

	return r
}
