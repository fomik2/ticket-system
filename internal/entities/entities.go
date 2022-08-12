package entities

import (
	"time"
)

var (
	TicketList []Ticket //tickets сожержит все тикеты, которые есть в системе
)

/*Ticket описывает заявку*/
type Ticket struct {
	Title, Description, Status, Severity string
	SLA, CreatedAt                       time.Time
	Number                               uint32
}
