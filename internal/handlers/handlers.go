package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/fomik2/ticket-system/internal/entities"
)

type Tickets interface {
	Get(id int) (entities.Ticket, error)
	List() ([]entities.Ticket, error)
	Create(entities.Ticket) (entities.Ticket, error)
	Update(entities.Ticket) (entities.Ticket, error)
	Delete(id int) error
}

/*formData передеается в темплейт при вызове editHandler или welcomeHandler*/
type formData struct {
	entities.Ticket
	Errors     []string
	TicketList []entities.Ticket
}

type Handlers struct {
	repo            Tickets
	editorTemplPath string
	ticketsPath     string
	counterPath     string
	Templ           *template.Template
}

func New(index, editor, tickets, counter string, repo Tickets) (*Handlers, error) {
	var err error
	newHandler := Handlers{}
	newHandler.repo = repo
	newHandler.editorTemplPath = editor
	newHandler.ticketsPath = tickets
	newHandler.counterPath = counter
	newHandler.Templ, err = createTemplates(index, editor)
	if err != nil {
		return &Handlers{}, fmt.Errorf("error when try to parse templates %w", err)
	}
	return &newHandler, nil
}

//getTicketID берет реквест и возвращает ID тикета
func getTicketID(writer http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return id, fmt.Errorf("can't parse id.  %w", err)
	}
	return id, nil
}

func createTemplates(indexPath, editorPath string) (indexTempl *template.Template, err error) {
	indexTempl, err = template.ParseFiles(indexPath, editorPath)
	if err != nil {
		return nil, fmt.Errorf("can't parse id.  %w", err)
	}
	return indexTempl, nil
}

//GetTicketForEdit выбрать заявку для редактирования и показать её
func (h *Handlers) GetTicketForEdit(writer http.ResponseWriter, r *http.Request) {
	id, err := getTicketID(writer, r)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server Error"))
	}
	ticket, err := h.repo.Get(id)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server Error"))
	}
	h.Templ.Execute(writer, formData{
		Ticket: ticket, Errors: []string{},
	})
	if err != nil {
		log.Println("can't execute template", err)
	}
}

//EditHandler редактирование заявки
func (h *Handlers) EditHandler(writer http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Form["action"][0] {
	case "Редактировать":
		id, err := getTicketID(writer, r)
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server Error"))
		}
		currentTicket, err := h.repo.Get(id)
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server error"))
		}
		currentTicket.Description = r.Form["description"][0]
		currentTicket.Title = r.Form["title"][0]
		currentTicket.Severity = r.Form["severity"][0]
		_, err = h.repo.Update(currentTicket)
		if err != nil {
			log.Println("can't update ticket", err)
			writer.Write([]byte("Internal server error"))
		}
		http.Redirect(writer, r, "/", http.StatusSeeOther)
	case "Удалить":
		id, err := getTicketID(writer, r)
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server Error"))
		}
		err = h.repo.Delete(id)
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server Error"))
			return
		}
		http.Redirect(writer, r, "/", http.StatusSeeOther)
	}
}

//CreateTicket создание новой заявки
func (h *Handlers) CreateTicket(writer http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	responseData := entities.Ticket{
		Title:       r.Form["title"][0],
		Description: r.Form["description"][0],
		Severity:    r.Form["severity"][0],
		Status:      "Создана",
		CreatedAt:   time.Now().Local(),
	}
	errors := []string{}
	if responseData.Title == "" {
		errors = append(errors, "Введите название заявки")
	}
	if responseData.Description == "" {
		errors = append(errors, "Введите описание")
	}
	if len(errors) > 0 {
		tickets, err := h.repo.List()
		if err != nil {
			log.Println("can't load tickets, check data files", err)
			writer.Write([]byte("Internal server error"))
		}
		h.Templ.Execute(writer, formData{Ticket: responseData, Errors: errors, TicketList: tickets})
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server error, can't load template"))
		}
	} else {
		_, err := h.repo.Create(responseData)
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server error"))
		}
		http.Redirect(writer, r, "/", http.StatusSeeOther)
	}
}

//welcomeHandler отображение формы и списка всех заявок
func (h *Handlers) WelcomeHandler(writer http.ResponseWriter, r *http.Request) {
	tickets, err := h.repo.List()
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error"))
	} else {
		h.Templ.Execute(writer, formData{
			Ticket: entities.Ticket{}, Errors: []string{}, TicketList: tickets,
		})
	}
	if err != nil {
		log.Println("can't open template file", err)
		writer.Write([]byte("Internal server error"))
	}
}
