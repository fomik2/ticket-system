package repo

import (
	"github.com/fomik2/ticket-system/internal/entities"
)

type Repo struct {
}

func (t *Repo) Get(id int) (entities.Ticket, error) {
	for _, ticket := range entities.TicketList {
		if ticket.Number == uint32(id) {
			return ticket, nil
		}
	}
	return entities.Ticket{}, nil
}

func (t *Repo) List() ([]entities.Ticket, error) {
	return entities.TicketList, nil
}

func (t *Repo) Create(ticket entities.Ticket) (entities.Ticket, error) {
	entities.TicketList = append(entities.TicketList, ticket)
	return ticket, nil
}

func (t *Repo) Delete(id int) error {
	for i, ticket := range entities.TicketList {
		if ticket.Number == uint32(id) {
			entities.TicketList = removeIndex(entities.TicketList, i)
		}
	}
	return nil
}

func (t *Repo) Update(ticket entities.Ticket) (entities.Ticket, error) {
	for i, t := range entities.TicketList {
		if t.Number == uint32(ticket.Number) {
			entities.TicketList[i] = ticket
		}
	}
	return ticket, nil
}

//RemoveIndex удаляет из слайса элемент заявки по индексу
func removeIndex(s []entities.Ticket, index int) []entities.Ticket {
	if len(s) == 1 {
		return []entities.Ticket{}
	}
	return append(s[:index], s[index+1:]...)
}
