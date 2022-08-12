package entities

import (
	"time"
)

var (
	TicketList []Ticket //tickets сожержит все тикеты, которые есть в системе
)

/*Ticket описывает заявку*/
type Ticket struct {
	Title, Description, Status, CreatedAt, Severity string
	SLA                                             time.Time
	Number                                          uint32
}
