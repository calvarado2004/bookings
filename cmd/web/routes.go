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

	mux.Get("/bookings", handlers.Repo.Home)
	mux.Get("/bookings/", handlers.Repo.Home)
	mux.Get("/bookings/about", handlers.Repo.About)
	mux.Get("/bookings/generals-quarters", handlers.Repo.Generals)
	mux.Get("/bookings/majors-suite", handlers.Repo.Majors)
	mux.Get("/bookings/specials", handlers.Repo.Specials)

	mux.Get("/bookings/reservation", handlers.Repo.Reservation)
	mux.Post("/bookings/reservation", handlers.Repo.PostReservation)
	mux.Get("/bookings/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Get("/bookings/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/bookings/book-room", handlers.Repo.BookRoom)

	mux.Get("/bookings/contact", handlers.Repo.Contact)
	mux.Get("/bookings/availability", handlers.Repo.Availability)
	mux.Post("/bookings/availability", handlers.Repo.PostAvailability)
	mux.Post("/bookings/availability-json", handlers.Repo.AvailabilityJSON)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/bookings/static/*", http.StripPrefix("/bookings/static", fileServer))

	return mux
}
