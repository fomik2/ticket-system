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

func (t Ticket) Get(id int) (Ticket, error) {
	for _, ticket := range TicketList {
		if ticket.Number == uint32(id) {
			return ticket, nil
		}
	}
	return Ticket{}, nil
}

func (t Ticket) List() ([]Ticket, error) {
	return TicketList, nil
}

func (t Ticket) Create(ticket Ticket) (Ticket, error) {
	TicketList = append(TicketList, ticket)
	return ticket, nil
}

func (t Ticket) Delete(id int) error {
	for i, ticket := range TicketList {
		if ticket.Number == uint32(id) {
			TicketList = RemoveIndex(TicketList, i)
		}
	}
	return nil
}

func (t Ticket) Update(ticket Ticket) (Ticket, error) {
	for i, t := range TicketList {
		if t.Number == uint32(ticket.Number) {
			TicketList[i] = ticket
		}
	}
	return ticket, nil
}

//RemoveIndex удаляет из слайса элемент заявки по индексу
func RemoveIndex(s []Ticket, index int) []Ticket {
	if len(s) == 1 {
		return []Ticket{}
	}
	return append(s[:index], s[index+1:]...)
}
