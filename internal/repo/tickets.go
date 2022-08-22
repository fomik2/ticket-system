package repo

import (
	"database/sql"
	"time"

	"github.com/fomik2/ticket-system/internal/entities"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	t := Repo{}
	t.db = db
	return &t
}

//SLAConfig устанавливает SLA заявки в зависимости от выбранного severity
func SLAConfig(severity string) time.Time {
	curTime := time.Now().Local()
	var SLATime time.Time
	switch {
	case severity == "5":
		SLATime = curTime.Add(3 * time.Hour)
	case severity == "4":
		SLATime = curTime.Add(4 * time.Hour)
	case severity == "3":
		SLATime = curTime.Add(5 * time.Hour)
	case severity == "2":
		SLATime = curTime.Add(6 * time.Hour)
	case severity == "1":
		SLATime = curTime.Add(7 * time.Hour)
	}
	return SLATime
}

func (t *Repo) CreateTicket(ticket entities.Ticket) (entities.Ticket, error) {
	ticket.SLA = SLAConfig(ticket.Severity)
	_, err := t.db.Exec("INSERT INTO tickets VALUES(NULL,?,?,?,?,?,?,?);", ticket.Title, ticket.Description, ticket.Status, ticket.Severity, ticket.SLA.Format("2006-01-02 15:04"), ticket.CreatedAt.Format("2006-01-02 15:04"), ticket.OwnerEmail)
	if err != nil {
		return entities.Ticket{}, err
	}
	return ticket, nil
}

func (t *Repo) ListTickets() ([]entities.Ticket, error) {
	rows, err := t.db.Query("SELECT * FROM tickets ORDER BY id;")
	data := []entities.Ticket{}
	if err != nil {
		return []entities.Ticket{}, err
	}
	defer rows.Close()
	for rows.Next() {
		i := entities.Ticket{}
		err = rows.Scan(&i.Number, &i.Title, &i.Description, &i.Status, &i.Severity, &i.SLA, &i.CreatedAt, &i.OwnerEmail)

		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}
	return data, nil
}

func (t *Repo) GetTicket(id int) (entities.Ticket, error) {
	row := t.db.QueryRow("SELECT * FROM tickets WHERE id=?;", id)
	ticket := entities.Ticket{}
	err := row.Scan(&ticket.Number, &ticket.Title, &ticket.Description, &ticket.Status, &ticket.Severity, &ticket.SLA, &ticket.CreatedAt, &ticket.OwnerEmail)
	if err != nil {
		return entities.Ticket{}, err
	}
	return ticket, nil
}

func (t *Repo) DeleteTicket(id int) error {
	_, err := t.db.Exec("DELETE FROM tickets WHERE id=?;", id)
	if err != nil {
		return err
	}
	return nil
}

func (t *Repo) UpdateTicket(ticket entities.Ticket) (entities.Ticket, error) {
	_, err := t.db.Exec("UPDATE tickets SET title=?, description=?, severity=? WHERE id=?;", ticket.Title, ticket.Description, ticket.Severity, ticket.Number)
	if err != nil {
		return entities.Ticket{}, err
	}
	return ticket, nil
}
