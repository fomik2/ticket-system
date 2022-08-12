package filerw

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/fomik2/ticket-system/internal/entities"
)

var (
	ticketNumbers uint32 // счетчик заявок.
)

//ticketNumberPlus инкрементирует счетчик при создании заявки
func TicketNumberPlus() uint32 {
	ticketNumbers = ticketNumbers + 1
	return ticketNumbers
}

func WriteTicketsToFiles(arr []entities.Ticket, ticketsFile string) {
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
	_ = ioutil.WriteFile(ticketsFile, file, 0644)
	log.Println("Записываем данные в файлы")

}

func ReadTicketsFromFiles(ticketsFile, counterFile string) {
	//read counter of tickets from file
	byteCounter, err := os.ReadFile(counterFile)
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
	jsonFile, err := os.Open(ticketsFile)
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
