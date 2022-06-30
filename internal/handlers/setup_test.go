package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/calvarado2004/bookings/internal/config"
	"github.com/calvarado2004/bookings/internal/models"
	"github.com/calvarado2004/bookings/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

var functions = template.FuncMap{
	"humanDate":  render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate":    render.Iterate,
	"add":        render.Add,
}

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// change this to true when in production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)

	listenForMail()

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Route("/bookings", func(mux chi.Router) {

		mux.Get("/", Repo.Home)
		mux.Get("/about", Repo.About)
		mux.Get("/generals-quarters", Repo.Generals)
		mux.Get("/majors-suite", Repo.Majors)

		mux.Get("/availability", Repo.Availability)
		mux.Post("/availability", Repo.PostAvailability)
		mux.Post("/availability-json", Repo.AvailabilityJSON)

		mux.Get("/contact", Repo.Contact)

		mux.Get("/reservation", Repo.Reservation)
		mux.Post("/reservation", Repo.PostReservation)
		mux.Get("/reservation-summary", Repo.ReservationSummary)
		mux.Get("/specials", Repo.Specials)

		mux.Get("/user/login", Repo.Showlogin)
		mux.Post("/user/login", Repo.PostShowlogin)

		mux.Get("/user/logout", Repo.Logout)

		mux.Get("/admin/dashboard", Repo.AdminDashboard)

		mux.Get("/admin/reservations/new", Repo.AdminAllNewReservations)
		mux.Get("/admin/reservations/all", Repo.AdminAllReservations)
		mux.Get("/admin/reservations/calendar", Repo.AdminCalendarReservations)
		mux.Post("/admin/reservations/calendar", Repo.AdminPostCalendar)
		mux.Get("/admin/process-reservation/{src}/{id}/do", Repo.AdminProcessReservations)
		mux.Get("/admin/delete-reservation/{src}/{id}/do", Repo.AdminDeleteReservations)

		mux.Get("/admin/reservations/{src}/{id}/show", Repo.AdminShowReservations)
		mux.Post("/admin/reservations/{src}/{id}", Repo.AdminPostShowReservations)

		fileServer := http.FileServer(http.Dir("./../../static/"))
		mux.Handle("/static/*", http.StripPrefix("/bookings/static", fileServer))

	})

	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		log.Println(err)
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			log.Println(err)
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			log.Println(err)
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				log.Println(err)
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
