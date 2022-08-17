package repo

import (
	"encoding/json"
	"fmt"
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

func New(tickets, counter string) *Repo {
	t := Repo{}
	t.counterPath = counter
	t.ticketsPath = tickets
	return &t
}

//SLAConfig устанавливает SLA заявки в зависимости от выбранного severity
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

func (t *Repo) writeTicketsToFiles(arr []entities.Ticket) error {
	//open counter and write current counter
	f, err := os.Create(t.counterPath)
	if err != nil {
		return fmt.Errorf("can't open file for writing %w", err)
	}
	var s string = strconv.FormatUint(uint64(t.counter+1), 10)
	_, err = f.WriteString(s)
	if err != nil {
		return fmt.Errorf("can't write data to file %w", err)
	}
	//open json and parse tickets
	file, err := json.MarshalIndent(arr, "", " ")
	if err != nil {
		return fmt.Errorf("something wrong with json marshal. %w", err)
	}
	err = ioutil.WriteFile(t.ticketsPath, file, 0644)
	if err != nil {
		return fmt.Errorf("can't write data to file %w", err)
	}
	log.Println("Writing data to file...")
	return err
}

func (t *Repo) readTicketsFromFiles() ([]entities.Ticket, error) {
	ticketList := []entities.Ticket{}
	//read counter of tickets from file
	byteCounter, err := os.ReadFile(t.counterPath)
	if err != nil {
		return []entities.Ticket{}, fmt.Errorf("can't read file %w", err)
	}
	strCounter := string(byteCounter)
	uint64number, err := strconv.ParseUint(strCounter, 10, 32)
	t.counter = uint32(uint64number)
	if err != nil {
		return []entities.Ticket{}, fmt.Errorf("can't read counter from file %w", err)
	} else {
		log.Println("Read counter...")
	}
	//read all tickets from json
	jsonFile, err := os.Open(t.ticketsPath)
	if err != nil {
		return []entities.Ticket{}, fmt.Errorf("can't open ticket file. %w", err)
	} else {
		log.Println("Read JSON file with tickets...")
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []entities.Ticket{}, fmt.Errorf("can't read json. %w", err)
	}
	err = json.Unmarshal(byteValue, &ticketList)
	if err != nil {
		return []entities.Ticket{}, fmt.Errorf("can't unmarshal json. %w", err)
	}
	return ticketList, err
}

func (t *Repo) Get(id int) (entities.Ticket, error) {
	ticketList, err := t.readTicketsFromFiles()
	if err != nil {
		return entities.Ticket{}, err
	}
	for _, ticket := range ticketList {
		if ticket.Number == uint32(id) {
			return ticket, nil
		}
	}
	return entities.Ticket{}, nil
}

func (t *Repo) List() ([]entities.Ticket, error) {
	ticketList, err := t.readTicketsFromFiles()
	if err != nil {
		return []entities.Ticket{}, err
	}
	return ticketList, nil
}

func (t *Repo) Create(ticket entities.Ticket) (entities.Ticket, error) {
	ticketList, err := t.readTicketsFromFiles()
	if err != nil {
		return entities.Ticket{}, err
	}
	ticket.Number = t.counter + 1
	ticket.SLA = SLAConfig(ticket.Severity)
	ticketList = append(ticketList, ticket)
	err = t.writeTicketsToFiles(ticketList)
	if err != nil {
		return ticket, err
	}
	return ticket, nil
}

func (t *Repo) Delete(id int) error {
	ticketList, err := t.readTicketsFromFiles()
	if err != nil {
		return err
	}
	for i, ticket := range ticketList {
		if ticket.Number == uint32(id) {
			ticketList = removeIndex(ticketList, i)
		}
	}
	err = t.writeTicketsToFiles(ticketList)
	if err != nil {
		return err
	}
	return err
}

func (t *Repo) Update(ticket entities.Ticket) (entities.Ticket, error) {
	ticketList, err := t.readTicketsFromFiles()
	if err != nil {
		return entities.Ticket{}, err
	}
	for i, elem := range ticketList {
		if elem.Number == uint32(ticket.Number) {
			ticketList[i] = ticket
			err := t.writeTicketsToFiles(ticketList)
			if err != nil {
				return ticket, err
			}
			return ticket, nil
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
