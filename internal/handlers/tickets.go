package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/fomik2/ticket-system/internal/entities"
)

type RepoInterface interface {
	GetTicket(id int) (entities.Ticket, error)
	ListTickets() ([]entities.Ticket, error)
	CreateTicket(entities.Ticket) (entities.Ticket, error)
	UpdateTicket(entities.Ticket) (entities.Ticket, error)
	DeleteTicket(id int) error
	GetUser(id int) (entities.Users, error)
	ListUsers() ([]entities.Users, error)
	CreateUser(entities.Users) (entities.Users, error)
	UpdateUser(entities.Users) (entities.Users, error)
	DeleteUser(id int) error
	IsUserExistInDB(username, password string) (int, error)
	FindUser(username, password string) (entities.Users, error)
}

/*formData передеается в темплейт при вызове editHandler или welcomeHandler*/
type formData struct {
	entities.Ticket
	Errors     []string
	TicketList []entities.Ticket
}

type Handlers struct {
	repo            RepoInterface
	layoutTemplPath string
	templs          map[string]*template.Template
}

func New(index, layout, editor, auth, user_create string, repo RepoInterface) (*Handlers, error) {
	var err error
	newHandler := Handlers{}
	newHandler.repo = repo
	newHandler.layoutTemplPath = layout
	newHandler.templs, err = newHandler.parseTemplates(index, editor, auth, user_create)
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

func (h *Handlers) parseTemplates(templPathes ...string) (map[string]*template.Template, error) {
	templs := make(map[string]*template.Template)
	var err error
	for _, templ := range templPathes {
		templ_name := strings.Split(templ, "/")[2]
		templs[templ_name], err = template.ParseFiles(h.layoutTemplPath, templ)
		if err != nil {
			return nil, fmt.Errorf("can't parse index template.  %w", err)
		}
	}
	return templs, nil
}

//GetTicketForEdit выбрать заявку для редактирования и показать её
func (h *Handlers) GetTicketForEdit(writer http.ResponseWriter, r *http.Request) {
	log.Println("GetTicketForEdit Handler in action....", r.Method)

	id, err := getTicketID(writer, r)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server Error"))
		return
	}
	ticket, err := h.repo.GetTicket(id)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server Error"))
		return
	}
	h.templs["editor"].Execute(writer, formData{
		Ticket: ticket, Errors: []string{},
	})
	if err != nil {
		log.Println("can't execute template", err)
		return
	}
}

//EditHandler редактирование заявки
func (h *Handlers) EditHandler(writer http.ResponseWriter, r *http.Request) {
	log.Println("Edit Handler in action....", r.Method)
	r.ParseForm()
	id, err := getTicketID(writer, r)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server Error"))
		return
	}
	currentTicket, err := h.repo.GetTicket(id)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error"))
		return
	}
	currentTicket.Description = r.Form["description"][0]
	currentTicket.Title = r.Form["title"][0]
	currentTicket.Severity = r.Form["severity"][0]
	_, err = h.repo.UpdateTicket(currentTicket)
	if err != nil {
		log.Println("can't update ticket", err)
		writer.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(writer, r, "/", http.StatusSeeOther)
}

func (h *Handlers) DeleteHandler(writer http.ResponseWriter, r *http.Request) {
	log.Println("DeleteHandler in action....", r.Method)
	id, err := getTicketID(writer, r)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server Error"))
		return
	}
	err = h.repo.DeleteTicket(id)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server Error"))
		return
	}
	http.Redirect(writer, r, "/", http.StatusSeeOther)
}

//CreateTicket создание новой заявки
func (h *Handlers) CreateTicket(writer http.ResponseWriter, r *http.Request) {
	log.Println("CreateTicket handler in action....", r.Method)
	//get session values
	session, err := store.Get(r, "session.id")
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error"))
		return
	}
	userEmail := session.Values["email"]
	r.ParseForm()
	responseData := entities.Ticket{
		Title:       r.Form["title"][0],
		Description: r.Form["description"][0],
		Severity:    r.Form["severity"][0],
		Status:      "Создана",
		CreatedAt:   time.Now().Local(),
		OwnerEmail:  userEmail,
	}
	fmt.Println(responseData)
	errors := []string{}
	if responseData.Title == "" {
		errors = append(errors, "Введите название заявки")
	}
	if responseData.Description == "" {
		errors = append(errors, "Введите описание")
	}
	if len(errors) > 0 {
		tickets, err := h.repo.ListTickets()
		if err != nil {
			log.Println("can't load tickets, check data files", err)
			writer.Write([]byte("Internal server error"))
			return
		}
		h.templs["index"].Execute(writer, formData{Ticket: responseData, Errors: errors, TicketList: tickets})
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server error, can't load template"))
			return
		}
	} else {
		_, err := h.repo.CreateTicket(responseData)
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server error"))
			return
		}
		http.Redirect(writer, r, "/", http.StatusSeeOther)
	}
}

//welcomeHandler отображение формы и списка всех заявок
func (h *Handlers) WelcomeHandler(writer http.ResponseWriter, r *http.Request) {
	log.Println("Welcome handler in action....", r.Method)
	tickets, err := h.repo.ListTickets()
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error"))
		return
	} else {
		h.templs["index"].Execute(writer, formData{
			Ticket: entities.Ticket{}, Errors: []string{}, TicketList: tickets,
		})
	}
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error"))
		return
	}
}
