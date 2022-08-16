package repo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fomik2/ticket-system/internal/entities"
)

type Repo struct {
	ticketsPath, counterPath string
	counter                  uint32
}

func New(config map[string]string) *Repo {
	t := Repo{}
	t.counterPath = config["counter"]
	t.ticketsPath = config["tickets"]
	return &t
}

//SLAConfig устанавливает SLA заявки в зависимости от выбранного severity
func (t *Repo) SLAConfig(severity string) time.Time {
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

func (t *Repo) writeTicketsToFiles(arr []entities.Ticket) error {
	//open counter and write current counter
	f, err := os.Create(t.counterPath)
	if err != nil {
		log.Println("Не могу открыть файл для записи", err)
	}
	var s string = strconv.FormatUint(uint64(t.counter+1), 10)
	_, err = f.WriteString(s)
	if err != nil {
		log.Println("Не могу записать номера заявок в файл", err)
	}
	//open json and parse tickets
	file, err := json.MarshalIndent(arr, "", " ")
	if err != nil {
		log.Println("Ошибка при записи json в файл", err)
	}
	err = ioutil.WriteFile(t.ticketsPath, file, 0644)
	if err != nil {
		log.Println("Не могу записать тикеты в файл", err)
	}
	log.Println("Записываем данные в файлы")
	return err
}

func (t *Repo) readTicketsFromFiles() []entities.Ticket {
	ticketList := []entities.Ticket{}
	//read counter of tickets from file
	byteCounter, err := os.ReadFile(t.counterPath)
	if err != nil {
		log.Panicln("Не могу прочитать файл-счетчик", err)
	}
	strCounter := string(byteCounter)
	uint64number, err := strconv.ParseUint(strCounter, 10, 32)
	t.counter = uint32(uint64number)
	if err != nil {
		log.Panicln("Не могу прочитать счетчик", err)
	} else {
		log.Println("Считываем счетчик тикетов")
	}
	//read all tickets from json
	jsonFile, err := os.Open(t.ticketsPath)
	if err != nil {
		log.Panicln("Не могу открыть файл с заявками", err)
	} else {
		log.Println("Считываем заявки из базы данных")
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Panicln("Не могу прочитать файл с заявками", err)
	}
	err = json.Unmarshal(byteValue, &ticketList)
	if err != nil {
		log.Panicln("Не могу записать полученный json в структуру", err)
	}
	return ticketList
}

func (t *Repo) Get(id int) (entities.Ticket, error) {
	ticketList := t.readTicketsFromFiles()
	for _, ticket := range ticketList {
		if ticket.Number == uint32(id) {
			return ticket, nil
		}
	}
	return entities.Ticket{}, nil
}

func (t *Repo) List() ([]entities.Ticket, error) {
	ticketList := t.readTicketsFromFiles()
	return ticketList, nil
}

func (t *Repo) Create(ticket entities.Ticket) (entities.Ticket, error) {
	ticketList := t.readTicketsFromFiles()
	ticket.Number = t.counter + 1
	ticket.SLA = t.SLAConfig(ticket.Severity)
	ticketList = append(ticketList, ticket)
	err := t.writeTicketsToFiles(ticketList)
	if err != nil {
		log.Println("Программа не смогла записать данные, проверьте, существуют ли файлы")
	}
	return ticket, nil
}

func (t *Repo) Delete(id int) error {
	ticketList := t.readTicketsFromFiles()
	for i, ticket := range ticketList {
		if ticket.Number == uint32(id) {
			ticketList = removeIndex(ticketList, i)
		}
		err := t.writeTicketsToFiles(ticketList)
		if err != nil {
			log.Println("Программа не смогла записать данные, проверьте, существуют ли файлы")
		}
	}
	return nil
}

func (t *Repo) Update(ticket entities.Ticket) (entities.Ticket, error) {
	ticketList := t.readTicketsFromFiles()
	for i, elem := range ticketList {
		if elem.Number == uint32(ticket.Number) {
			ticketList[i] = elem
			err := t.writeTicketsToFiles(ticketList)
			if err != nil {
				log.Println("Программа не смогла записать данные, проверьте, существуют ли файлы")
			}
			return elem, nil
		}
	}
	return entities.Ticket{}, nil
}

//RemoveIndex удаляет из слайса элемент заявки по индексу
func removeIndex(s []entities.Ticket, index int) []entities.Ticket {
	if len(s) == 1 {
		return []entities.Ticket{}
	}
	return append(s[:index], s[index+1:]...)
}
