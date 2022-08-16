package entities

import (
	"time"
)

/*Ticket описывает заявку*/
type Ticket struct {
	Title, Description, Status, Severity string
	SLA, CreatedAt                       time.Time
	Number                               uint32
}
