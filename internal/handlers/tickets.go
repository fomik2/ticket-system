package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"

	"github.com/fomik2/ticket-system/internal/entities"
)

type RepoInterface interface {
	Tickets
	Users
}

type Tickets interface {
	GetTicket(id int) (entities.Ticket, error)
	ListTickets() ([]entities.Ticket, error)
	CreateTicket(entities.Ticket) (entities.Ticket, error)
	UpdateTicket(entities.Ticket) (entities.Ticket, error)
	DeleteTicket(id int) error
}

type Users interface {
	GetUser(id int) (entities.Users, error)
	ListUsers() ([]entities.Users, error)
	ListTicketsByUser(email string) ([]entities.Ticket, error)
	CreateUser(entities.Users) (entities.Users, error)
	UpdateUser(entities.Users) (entities.Users, error)
	DeleteUser(id int) error
	FindUser(username string) (entities.Users, error)
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
	jwtKey          []byte
	sessionStore    *sessions.CookieStore
	templs          map[string]*template.Template
}

func New(index, layout, editor, auth, user_create, secret string, repo RepoInterface) (*Handlers, error) {
	var err error
	newHandler := Handlers{}
	newHandler.repo = repo
	newHandler.layoutTemplPath = layout
	newHandler.jwtKey = []byte(secret)

	// store the secret key in env variable in production
	newHandler.sessionStore = sessions.NewCookieStore([]byte(secret))
	newHandler.templs, err = newHandler.parseTemplates(index, editor, auth, user_create)
	if err != nil {
		return &Handlers{}, fmt.Errorf("error when try to parse templates %w", err)
	}
	return &newHandler, nil
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

// GetTicketForEdit выбрать заявку для редактирования и показать её
func (h *Handlers) GetTicketForEdit(c echo.Context) error {
	log.Println("GetTicketForEdit Handler in action....", c.Request().Method)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error while parsing ticket ID"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	ticket, err := h.repo.GetTicket(id)
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	h.templs["editor"].Execute(c.Response(), formData{
		Ticket: ticket, Errors: []string{},
	})
	if err != nil {
		log.Println("can't execute template", err)
		return err
	}
	return nil
}

// EditHandler редактирование заявки
func (h *Handlers) EditHandler(c echo.Context) error {
	log.Println("Edit Handler in action....", c.Request().Method)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error while parsing ticket ID"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	currentTicket, err := h.repo.GetTicket(id)
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	currentTicket.Description = c.FormValue("description")
	currentTicket.Title = c.FormValue("title")
	currentTicket.Severity = c.FormValue("severity")
	_, err = h.repo.UpdateTicket(currentTicket)
	if err != nil {
		log.Println("can't update ticket", err)
		c.Response().Write([]byte("Internal server error"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	http.Redirect(c.Response(), c.Request(), "/", http.StatusSeeOther)
	return nil
}

func (h *Handlers) DeleteHandler(c echo.Context) error {
	log.Println("DeleteHandler in action....", c.Request().Method)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server eror while parsing ticket ID"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	err = h.repo.DeleteTicket(id)
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server eror"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	http.Redirect(c.Response(), c.Request(), "/", http.StatusSeeOther)
	return nil
}

// CreateTicket создание новой заявки
func (h *Handlers) CreateTicket(c echo.Context) error {
	log.Println("CreateTicket handler in action....", c.Request().Method)
	//get session values
	session, err := h.sessionStore.Get(c.Request(), "session.id")
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	userEmail := session.Values["email"]
	responseData := entities.Ticket{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Severity:    c.FormValue("severity"),
		Status:      "Создана",
		CreatedAt:   time.Now().Local(),
		OwnerEmail:  userEmail,
	}
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
			c.Response().Write([]byte("Internal server error"))
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
		h.templs["index"].Execute(c.Response(), formData{Ticket: responseData, Errors: errors, TicketList: tickets})
		if err != nil {
			log.Println(err)
			c.Response().Write([]byte("Internal server error, can't load template"))
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
	} else {
		_, err := h.repo.CreateTicket(responseData)
		if err != nil {
			log.Println(err)
			c.Response().Write([]byte("Internal server error"))
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
		http.Redirect(c.Response(), c.Request(), "/", http.StatusSeeOther)
	}
	return nil
}

// welcomeHandler отображение формы и списка всех заявок
func (h *Handlers) WelcomeHandler(c echo.Context) error {
	log.Println("Welcome handler in action....", c.Request().Method)
	tickets, err := h.repo.ListTickets()
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	} else {
		h.templs["index"].Execute(c.Response(), formData{
			Ticket: entities.Ticket{}, Errors: []string{}, TicketList: tickets,
		})
	}
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}
