package entities

import (
	"time"
)

/*Ticket описывает заявку*/
type Ticket struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Status      string `json:"Status"`
	Severity    string `json:"Severity"`
	SLA         time.Time
	CreatedAt   time.Time `json:"CreatedAt"`
	Number      uint32    `json:"ID"`
	OwnerEmail  interface{}
}

type Users struct {
	ID        uint32
	Name      string
	Password  string
	Email     string
	CreatedAt time.Time
}
