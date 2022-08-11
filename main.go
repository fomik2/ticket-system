package main

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
)

/*Ticket описывает заявку*/
type Ticket struct {
	Title, Description, Status, CreatedAt, Severity string
	SLA                                             time.Time
	Number                                          uint32
}

/*formData передеается в темплейт при вызове editHandler или welcomeHandler*/
type formData struct {
	Ticket
	Errors     []string
	TicketList []Ticket
}

var tickets []Ticket //tickets сожержит все тикеты, которые есть в системе

var ticketNumbers uint32 // счетчик заявок.

//RemoveIndex удаляет из слайса элемент заявки по индексу
func RemoveIndex(s []Ticket, index int) []Ticket {
	if len(s) == 1 {
		return []Ticket{}
	}
	return append(s[:index], s[index+1:]...)
}

//ticketNumberPlus инкрементирует счетчик при создании заявки
func ticketNumberPlus() uint32 {
	ticketNumbers = ticketNumbers + 1
	return ticketNumbers
}

//editHandler редактирование заявки
func editHandler(writer http.ResponseWriter, r *http.Request) {
	param1, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Something wrong with convertion string to int", err)
	}
	if r.Method == http.MethodGet {
		createTemplate, err := template.ParseFiles("./templates/editor.html")
		if err != nil {
			log.Println("Проблема с загрузкой темплейта", err)
		}
		for _, ticket := range tickets {
			if ticket.Number == uint32(param1) {
				createTemplate.Execute(writer, formData{
					Ticket: ticket, Errors: []string{},
				})
			}
		}
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		switch r.Form["action"][0] {
		case "Редактировать":
			for _, ticket := range tickets {
				if ticket.Number == uint32(param1) {
					ticket.Description = r.Form["description"][0]
					ticket.Title = r.Form["title"][0]
					ticket.Severity = r.Form["severity"][0]
					writeTicketsToFiles(tickets)
				}
			}
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		case "Удалить":
			for i, ticket := range tickets {
				if ticket.Number == uint32(param1) {
					tickets = RemoveIndex(tickets, i)
					writeTicketsToFiles(tickets)
				}
			}
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
			Ticket: Ticket{}, Errors: []string{}, TicketList: tickets,
		})
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		responseData := Ticket{
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
			createTemplate.Execute(writer, formData{Ticket: responseData, Errors: errors, TicketList: tickets})
		} else {
			tickets = append(tickets, responseData)
			writeTicketsToFiles(tickets)
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		}
	}
}

func writeTicketsToFiles(arr []Ticket) {
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
		log.Println("Не могу прочитать счетчик", err)
		panic(err)
	} else {
		log.Println("Считываем счетчик тикетов")
	}
	//read all tickets from json
	jsonFile, err := os.Open("tickets.json")
	if err != nil {
		log.Println("Не могу прочитать файл с заявками", err)
		panic(err)
	} else {
		log.Println("Считываем заявки из базы данных")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &tickets)
	if err != nil {
		log.Println("Не могу прочитать тикеты из файлы", err)
		panic(err)
	}

}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	readTicketsFromFiles()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/tickets/", editHandler)
	fs := http.FileServer(http.Dir(cfg.CSS.Path))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	err = http.ListenAndServe(cfg.HTTP.Port, nil)
	if err != nil {
		log.Fatal("Problem related to starting HTTP server", err)
	}

}
