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

	mux.Route("/bookings", func(mux chi.Router) {

		mux.Get("/", handlers.Repo.Home)
		mux.Get("/about", handlers.Repo.About)
		mux.Get("/generals-quarters", handlers.Repo.Generals)
		mux.Get("/majors-suite", handlers.Repo.Majors)
		mux.Get("/specials", handlers.Repo.Specials)

		mux.Get("/reservation", handlers.Repo.Reservation)
		mux.Post("/reservation", handlers.Repo.PostReservation)
		mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
		mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
		mux.Get("/book-room", handlers.Repo.BookRoom)

		mux.Get("/contact", handlers.Repo.Contact)
		mux.Get("/availability", handlers.Repo.Availability)
		mux.Post("/availability", handlers.Repo.PostAvailability)
		mux.Post("/availability-json", handlers.Repo.AvailabilityJSON)
		mux.Get("/user/login", handlers.Repo.Showlogin)
		mux.Post("/user/login", handlers.Repo.PostShowlogin)

		mux.Get("/user/logout", handlers.Repo.Logout)

		fileServer := http.FileServer(http.Dir("./static/"))

		mux.Handle("/static/*", http.StripPrefix("/bookings/static", fileServer))

	})

	mux.Route("/bookings-admin", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/reservations", handlers.Repo.AdminDashboard)
		mux.Get("/reservations/new", handlers.Repo.AdminAllNewReservations)
		mux.Get("/reservations/all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservations)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservations)
		mux.Get("/reservations/calendar", handlers.Repo.AdminCalendarReservations)
		mux.Post("/reservations/calendar", handlers.Repo.AdminPostCalendar)
		mux.Get("/reservations/process-reservation/{src}/{id}", handlers.Repo.AdminProcessReservations)
		mux.Get("/reservations/delete-reservation/{src}/{id}", handlers.Repo.AdminDeleteReservations)

	})

	return mux
}
