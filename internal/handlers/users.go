package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/fomik2/ticket-system/internal/entities"
	"github.com/labstack/echo/v4"
)

type formDataUsers struct {
	entities.Users
	Errors []string
}

func (h *Handlers) CreateUserGet(c echo.Context) error {
	log.Println("Create user handler in action....", c.Request().Method)
	err := h.templs["user_create"].Execute(c.Response(), formDataUsers{
		Users: entities.Users{}, Errors: []string{},
	})
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Can't show template"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *Handlers) CreateUser(c echo.Context) error {
	log.Println("Create user insert to table handler in action....", c.Request())
	var err error
	newUser := entities.Users{}
	newUser.Name = c.FormValue("name")
	newUser.Email = c.FormValue("email")
	newUser.Password, err = h.HashPassword(c.FormValue("password"))
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	newUser.CreatedAt = time.Now().Local()

	_, err = h.repo.CreateUser(newUser)
	if err != nil {
		c.Response().Write([]byte(err.Error()))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	http.Redirect(c.Response(), c.Request(), "/", http.StatusSeeOther)
	return nil
}
