package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/fomik2/ticket-system/internal/entities"
)

type formDataUsers struct {
	entities.Users
	Errors []string
}

func (h *Handlers) CreateUserGet(writer http.ResponseWriter, r *http.Request) {
	log.Println("Create user handler in action....", r.Method)
	h.templs["user_create"].Execute(writer, formDataUsers{
		Users: entities.Users{}, Errors: []string{},
	})

}

func (h *Handlers) CreateUser(writer http.ResponseWriter, r *http.Request) {
	log.Println("Create user insert to table handler in action....", r.Method)
	r.ParseForm()
	newUser := entities.Users{}
	newUser.Name = r.Form["name"][0]
	newUser.Email = r.Form["email"][0]
	newUser.Password = r.Form["password"][0]
	newUser.CreatedAt = time.Now().Local()
	_, err := h.repo.CreateUser(newUser)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(writer, r, "/", http.StatusSeeOther)

}
