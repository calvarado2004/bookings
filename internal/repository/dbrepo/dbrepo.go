package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/calvarado2004/bookings/internal/config"
	"github.com/calvarado2004/bookings/internal/models"
	"github.com/calvarado2004/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, app *config.AppConfig) repository.DatabaseRepo {

	return &postgresDBRepo{
		App: app,
		DB:  conn,
	}

}

//InsertReservation, inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
	         values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
