package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/calvarado2004/bookings/internal/config"
	"github.com/calvarado2004/bookings/internal/driver"
	"github.com/calvarado2004/bookings/internal/forms"
	"github.com/calvarado2004/bookings/internal/helpers"
	"github.com/calvarado2004/bookings/internal/models"
	"github.com/calvarado2004/bookings/internal/render"
	"github.com/calvarado2004/bookings/internal/repository"
	"github.com/calvarado2004/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

//	About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.html", &models.TemplateData{})

}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

func (m *Repository) Specials(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "specials.page.html", &models.TemplateData{})
}

//Reservation get the reservation
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})

	data["reservation"] = res

	render.Template(w, r, "reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

//Post reservation
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		Room:      room,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, r, "reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	htmlMessage := fmt.Sprintf(`
	<strong>Reservation Confirmation</strong>
	Dear %s, <br>
	This email is your reservation confirmation from Bookings in Go.
	Dates selected %s to %s
	
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	//send notifications
	msg := models.MailData{
		To:       reservation.Email,
		From:     "admin@bookings.com",
		Subject:  "Reservation Confirmation from Bookings",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	//sends the message to the channel
	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/bookings/reservation-summary", http.StatusSeeOther)

}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/bookings/availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

//AvailabilityJSON generates a json output of available rooms
// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	// need to parse request body
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		// got a database error, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "     ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

//ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/bookings/reservation-summary", http.StatusSeeOther)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	//Grab the fourth element [3] of the result of split the URL by its slashes
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[3])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/bookings", http.StatusSeeOther)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/bookings/reservation", http.StatusSeeOther)

}

// BookRoom takes URL parameters, builds a sessional variable, and takes user to make res screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get room from db!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/bookings/reservation", http.StatusSeeOther)
}

func (m *Repository) Showlogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

//PostShowlogin handles logging the user in
func (m *Repository) PostShowlogin(w http.ResponseWriter, r *http.Request) {

	//Renew the token
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/bookings/user/login", http.StatusSeeOther)
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "logged in successfully")
	http.Redirect(w, r, "/bookings/", http.StatusSeeOther)

}

//Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())

	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/bookings/user/login", http.StatusSeeOther)
}

//AdminDashboard Admin dasboard
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})
}

//AdminAllReservations gets all reservations
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})

	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.html", &models.TemplateData{
		Data: data,
	})
}

//AdminAllNewReservations gets all the new reservations
func (m *Repository) AdminAllNewReservations(w http.ResponseWriter, r *http.Request) {

	newreservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//log.Println(newreservations)

	data := make(map[string]interface{})

	data["newreservations"] = newreservations

	render.Template(w, r, "admin-new-reservations.page.html", &models.TemplateData{
		Data: data,
	})
}

//AdminCalendarReservations displays the reservation calendar
func (m *Repository) AdminCalendarReservations(w http.ResponseWriter, r *http.Request) {
	//assume that there is no month/year specified

	now := time.Now()

	if r.URL.Query().Get("y") != "" {

		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	}

	data := make(map[string]interface{})

	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)
	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")
	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)

	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	//get the first and last days of a month

	currentYear, currentMonth, _ := now.Date()

	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)

	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {

		//create maps

		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {

			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0

		}

		// get all the restrictions for the current room
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {

				//it is a reservation
				for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0, 0, 1) {

					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}

			} else {

				//it is a block
				//for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0, 0, 1) {

				//blockMap[d.Format("2006-01-2")] = y.ID
				//}
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID

			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)

	}

	render.Template(w, r, "admin-calendar-reservations.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

//AdminShowReservations show one reservation
func (m *Repository) AdminShowReservations(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[5])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//log.Println(id)

	src := exploded[4]

	stringMap := make(map[string]string)

	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["year"] = year
	stringMap["month"] = month

	//get reservation from the database

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})

	data["reservation"] = res

	render.Template(w, r, "admin-show-reservations.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})

}

func (m *Repository) AdminPostShowReservations(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[5])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[4]

	stringMap := make(map[string]string)

	stringMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservations(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/bookings/admin/reservations/%s", src), http.StatusSeeOther)

	} else {
		http.Redirect(w, r, fmt.Sprintf("/bookings/admin/reservations/calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

//AdminProcessReservations marks a reservation as processed
func (m *Repository) AdminProcessReservations(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := chi.URLParam(r, "src")

	err = m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/bookings/admin/reservations/%s", src), http.StatusSeeOther)

	} else {
		http.Redirect(w, r, fmt.Sprintf("/bookings/admin/reservations/calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

//AdminDeleteReservations marks a reservation as processed
func (m *Repository) AdminDeleteReservations(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := chi.URLParam(r, "src")

	err = m.DB.DeleteReservations(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted!")

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/bookings/admin/reservations/%s", src), http.StatusSeeOther)

	} else {
		http.Redirect(w, r, fmt.Sprintf("/bookings/admin/reservations/calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

func (m *Repository) AdminPostCalendar(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, err := strconv.Atoi(r.Form.Get("y"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	month, err := strconv.Atoi(r.Form.Get("m"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// process blocks

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	for _, x := range rooms {

		//Get the block map from the session. Loop through entire map, if we have an entry in the map
		//that does not exist in our posted data, and if the restriction id > 0, then it is a block we need to remove
		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)

		for name, value := range curMap {
			//ok will be false if the value is not in the map

			if val, ok := curMap[name]; ok {
				//only pay attention to values > 0, and that are not in the form post
				//the rest are placeholders for days without blocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {

						//delete the restriction by id
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}

					}
				}
			}

		}
	}

	//now handle new blocks

	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, err := strconv.Atoi(exploded[2])
			if err != nil {
				helpers.ServerError(w, err)
				return
			}

			t, err := time.Parse("2006-01-2", exploded[3])
			if err != nil {
				helpers.ServerError(w, err)
				return
			}

			//insert a new block
			err = m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				log.Println(err)
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved!")
	http.Redirect(w, r, fmt.Sprintf("/bookings/admin/reservations/calendar?y=%d&m=%d", year, month), http.StatusSeeOther)

}
