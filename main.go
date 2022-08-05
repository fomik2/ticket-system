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
)

type Ticket struct {
	Title, Description, Status, CreatedAt, Severity string
	SLA                                             time.Time
	Number                                          uint32
}

type formData struct {
	*Ticket
	Errors     []string
	TicketList []*Ticket
}

var tickets []*Ticket

var ticketNumbers uint32

func RemoveIndex(s []*Ticket, index int) []*Ticket {
	if len(s) == 1 {
		return []*Ticket{}
	}
	return append(s[:index], s[index+1:]...)
}

func ticketNumberPlus() uint32 {
	ticketNumbers = ticketNumbers + 1
	return ticketNumbers
}

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

func welcomeHandler(writer http.ResponseWriter, r *http.Request) {
	createTemplate, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Println("Проблема с загрузкой темплейта", err)
	}

	if r.Method == http.MethodGet {
		createTemplate.Execute(writer, formData{
			Ticket: &Ticket{}, Errors: []string{}, TicketList: tickets,
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
			createTemplate.Execute(writer, formData{Ticket: &responseData, Errors: errors, TicketList: tickets})
		} else {
			tickets = append(tickets, &responseData)
			writeTicketsToFiles(tickets)
			http.Redirect(writer, r, "/", http.StatusSeeOther)
		}
	}
}

func writeTicketsToFiles(arr []*Ticket) {
	//open counter and write current counter
	f, err := os.Create("counter")
	if err != nil {
		log.Println("Can't open file for writing")
		panic(err)
	}
	var s string = strconv.FormatUint(uint64(ticketNumbers), 10)
	_, err = f.WriteString(s)
	if err != nil {
		log.Println("Can't write counter to file ")
		panic(err)
	}
	//open json and parse tickets
	file, _ := json.MarshalIndent(arr, "", " ")
	_ = ioutil.WriteFile("tickets.json", file, 0644)
}

func readTicketsFromFiles() {
	//read counter of tickets from file
	counter, err := os.Open("counter")
	byteCounter, _ := ioutil.ReadAll(counter)
	if err != nil {
		fmt.Println("Не могу прочитать файл-счетчик")
		panic(err)
	}
	strCounter := string(byteCounter)
	uint64number, err := strconv.ParseUint(strCounter, 10, 32)
	ticketNumbers = uint32(uint64number)
	if err != nil {
		fmt.Println("Can't parse counter")
		panic(err)
	}
	//read all tickets from json
	jsonFile, err := os.Open("tickets.json")
	if err != nil {
		fmt.Println("Не могу прочитать файл с заявками")
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &tickets)
	defer jsonFile.Close()
	defer counter.Close()

}

func main() {
	readTicketsFromFiles()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/tickets/", editHandler)
	fs := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	err := http.ListenAndServe(":5002", nil)
	if err != nil {
		log.Fatal("Problem related to starting HTTP server", err)
	}

}
