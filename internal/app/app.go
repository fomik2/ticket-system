package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fomik2/ticket-system/config"
	"github.com/fomik2/ticket-system/internal/entities"
)

/*formData передеается в темплейт при вызове editHandler или welcomeHandler*/
type formData struct {
	entities.Ticket
	Errors     []string
	TicketList []entities.Ticket
}

var ticketNumbers uint32 // счетчик заявок.

//ticketNumberPlus инкрементирует счетчик при создании заявки
func ticketNumberPlus() uint32 {
	ticketNumbers = ticketNumbers + 1
	return ticketNumbers
}

//editHandler редактирование заявки
func editHandler(writer http.ResponseWriter, r *http.Request) {
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
			writeTicketsToFiles(entities.TicketList)
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		case "Удалить":
			ticket.Delete(param1)
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		}
	}
}

//welcomeHandler создание новой заявки и вывод всех заявок
func welcomeHandler(writer http.ResponseWriter, r *http.Request) {

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
			Number:      ticketNumberPlus(),
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
			entities.TicketList = append(entities.TicketList, responseData)
			writeTicketsToFiles(entities.TicketList)
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		}
	}
}

func writeTicketsToFiles(arr []entities.Ticket) {
	//open counter and write current counter
	f, err := os.Create("counter")
	if err != nil {
		log.Println("Не могу открыть файл для записи")
		panic(err)
	}
	var s string = strconv.FormatUint(uint64(ticketNumbers), 10)
	_, err = f.WriteString(s)
	if err != nil {
		log.Println("Не могу записать номера заявок в файл")
		panic(err)
	}
	//open json and parse tickets
	file, err := json.MarshalIndent(arr, "", " ")
	if err != nil {
		log.Println("Не могу записать тикеты в файл")
		panic(err)
	}
	_ = ioutil.WriteFile("tickets.json", file, 0644)
	log.Println("Записываем данные в файлы")

}

func readTicketsFromFiles() {
	//read counter of tickets from file
	byteCounter, err := os.ReadFile("counter")
	if err != nil {
		fmt.Println("Не могу прочитать файл-счетчик")
		panic(err)
	}
	strCounter := string(byteCounter)
	uint64number, err := strconv.ParseUint(strCounter, 10, 32)
	ticketNumbers = uint32(uint64number)
	if err != nil {
		log.Panicln("Не могу прочитать счетчик", err)
	} else {
		log.Println("Считываем счетчик тикетов")
	}
	//read all tickets from json
	jsonFile, err := os.Open("tickets.json")
	if err != nil {
		log.Panicln("Не могу прочитать файл с заявками", err)
	} else {
		log.Println("Считываем заявки из базы данных")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &entities.TicketList)
	if err != nil {
		log.Panicln("Не могу записать полученный json в структуру", err)
	}

}

func Run(cfg *config.Config) {
	readTicketsFromFiles()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/tickets/", editHandler)
	fs := http.FileServer(http.Dir(cfg.CSS.Path))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	fmt.Println(cfg.HTTP.Port)
	err := http.ListenAndServe(cfg.HTTP.Port, nil)
	if err != nil {
		log.Fatal("Problem related to starting HTTP server", err)
	}
}
