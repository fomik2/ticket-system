package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fomik2/ticket-system/internal/entities"
	"github.com/fomik2/ticket-system/internal/filerw"
)

/*formData передеается в темплейт при вызове editHandler или welcomeHandler*/
type formData struct {
	entities.Ticket
	Errors     []string
	TicketList []entities.Ticket
}

//editHandler редактирование заявки
func EditHandler(writer http.ResponseWriter, r *http.Request) {
	ticket := &entities.Ticket{}
	param1, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Something wrong with convertion string to int", err)
	}
	if r.Method == http.MethodGet {
		createTemplate, err := template.ParseFiles("./templates/editor.html")
		if err != nil {
			log.Panicln("Проблема с загрузкой темплейта", err)
		}
		ticket, err := ticket.Get(param1)
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
			ticket, err := ticket.Get(param1)
			if err != nil {
				log.Panicln("Пока что она всегда nil")
			}
			ticket.Description = r.Form["description"][0]
			ticket.Title = r.Form["title"][0]
			ticket.Severity = r.Form["severity"][0]
			_, err = ticket.Update(ticket)
			if err != nil {
				log.Panicln("Пока что она всегда nil")
			}
			filerw.WriteTicketsToFiles(entities.TicketList)
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		case "Удалить":
			ticket.Delete(param1)
			http.Redirect(writer, r, "/", http.StatusSeeOther)
			filerw.WriteTicketsToFiles(entities.TicketList)
		}
	}
}

//welcomeHandler создание новой заявки и вывод всех заявок
func WelcomeHandler(writer http.ResponseWriter, r *http.Request) {
	createTemplate, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Println("Проблема с загрузкой темплейта", err)
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
			CreatedAt:   time.Now().Format("02/01/2006 15:04"),
			SLA:         time.Now(),
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
			responseData.Create(responseData)
			filerw.WriteTicketsToFiles(entities.TicketList)
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		}
	}
}
