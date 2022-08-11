package tickets

import (
	"github.com/fomik2/ticket-system/internal/entities"
)

type Tickets interface {
	Get(id int) (entities.Ticket, error)
	List() ([]entities.Ticket, error)
	Create(entities.Ticket) (entities.Ticket, error)
	Update(entities.Ticket) (entities.Ticket, error)
	Delete(id int) error
}
