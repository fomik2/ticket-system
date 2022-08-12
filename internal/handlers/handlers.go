package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/fomik2/ticket-system/internal/entities"
	"github.com/fomik2/ticket-system/internal/filerw"
	"github.com/fomik2/ticket-system/internal/repo"
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

//getTicketID берет реквест и возвращает ID тикета
func getTicketID(writer http.ResponseWriter, r *http.Request, config map[string]string) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Panicln("Не могу распарсить ID тикета", err)
	}
	return id
}

//editHandler редактирование заявки
func EditHandler(writer http.ResponseWriter, r *http.Request, config map[string]string) {
	ticket := &repo.Repo{}
	if r.Method == http.MethodGet {
		createTemplate, err := template.ParseFiles(config["editor"])
		if err != nil {
			log.Panicln("Проблема с загрузкой темплейта", err)
		}
		id := getTicketID(writer, r, config)
		ticket, err := ticket.Get(id)
		if err != nil {
			log.Panicln("Пока что она всегда nil")
		}
		createTemplate.Execute(writer, formData{
			Ticket: ticket, Errors: []string{},
		})

	} else if r.Method == http.MethodPost {
		r.ParseForm()
		switch r.Form["action"][0] {
		case "Редактировать":
			id := getTicketID(writer, r, config)
			ticket := &repo.Repo{}
			currentTicket, err := ticket.Get(id)
			if err != nil {
				log.Panicln("Пока что она всегда nil")
			}
			currentTicket.Description = r.Form["description"][0]
			currentTicket.Title = r.Form["title"][0]
			currentTicket.Severity = r.Form["severity"][0]
			_, err = ticket.Update(currentTicket)
			if err != nil {
				log.Panicln("Пока что она всегда nil")
			}
			filerw.WriteTicketsToFiles(entities.TicketList, config["tickets"])
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		case "Удалить":
			id := getTicketID(writer, r, config)
			ticket.Delete(id)
			http.Redirect(writer, r, "/", http.StatusSeeOther)
			filerw.WriteTicketsToFiles(entities.TicketList, config["tickets"])
		}
	}
}

func SLAConfig(severity string) time.Time {
	curTime := time.Now().Local()
	var SLATime time.Time
	switch {
	case severity == "5":
		SLATime = curTime.Add(3 * time.Hour)
	case severity == "4":
		SLATime = curTime.Add(4 * time.Hour)
	case severity == "3":
		SLATime = curTime.Add(5 * time.Hour)
	case severity == "2":
		SLATime = curTime.Add(6 * time.Hour)
	case severity == "1":
		SLATime = curTime.Add(7 * time.Hour)
	}
	return SLATime

}

//welcomeHandler создание новой заявки и вывод всех заявок
func WelcomeHandler(writer http.ResponseWriter, r *http.Request, config map[string]string) {

	ticket := &repo.Repo{}
	createTemplate, err := template.ParseFiles(config["index"])
	if err != nil {
		log.Panicln("Проблема с загрузкой темплейта", err)
	}

	if r.Method == http.MethodGet {
		createTemplate.Execute(writer, formData{
			Ticket: entities.Ticket{}, Errors: []string{}, TicketList: entities.TicketList,
		})
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		responseData := entities.Ticket{
			Title:       r.Form["title"][0],
			Description: r.Form["description"][0],
			Severity:    r.Form["severity"][0],
			Status:      "Создана",
			CreatedAt:   time.Now().Local(),
			SLA:         SLAConfig(r.Form["severity"][0]),
			Number:      filerw.TicketNumberPlus(),
		}
		errors := []string{}
		if responseData.Title == "" {
			errors = append(errors, "Введите название заявки")
		}
		if responseData.Description == "" {
			errors = append(errors, "Введите описание")
		}
		if len(errors) > 0 {
			createTemplate.Execute(writer, formData{Ticket: responseData, Errors: errors, TicketList: entities.TicketList})
		} else {
			ticket.Create(responseData)
			filerw.WriteTicketsToFiles(entities.TicketList, config["tickets"])
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		}
	}
}
