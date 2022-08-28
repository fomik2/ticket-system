package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fomik2/ticket-system/internal/entities"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JWTCustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// APISignin create JWT token if user exist in DB, if not -- return unauthorized response code
func (h *Handlers) APISignin(c echo.Context) error {

	var claims JWTCustomClaims
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(c.Request().Body).Decode(&claims)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	user, err := h.repo.FindUser(claims.Username)

	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error while find user in DB"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	//if user not found in DB by login
	if user.Name == "" {
		c.Response().Write([]byte("Unauthorized. (No user found)"))
		c.Response().WriteHeader(http.StatusUnauthorized)
		return err
	}

	checkPass, err := h.CheckPasswordHash(user.Password, claims.Password)
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error while compare password with hash"))
		c.Response().WriteHeader(http.StatusUnauthorized)
		return err
	}

	//if password hashes doesn't match
	if !checkPass {
		c.Response().WriteHeader(http.StatusUnauthorized)
		return err
	}

	// Declare the expiration time of the token
	claims.ExpiresAt = time.Now().Add(5 * time.Minute).Unix()

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(h.jwtKey))
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "email"
	cookie.Value = user.Email
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
		"email": user.Email,
	})

}

func (h *Handlers) APICreateTicket(c echo.Context) error {
	ticket := entities.Ticket{}
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	json.Unmarshal(reqBody, &ticket)
	h.repo.CreateTicket(ticket)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	return c.JSONPretty(http.StatusCreated, ticket, "  ")
}

func (h *Handlers) APIGetTicket(c echo.Context) error {
	log.Println("APIGetTicket handler in action...", c.Request().Method)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Response().WriteHeader(http.StatusNotFound)
		return err
	}
	ticket, err := h.repo.GetTicket(id)
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server Error"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	if err != nil {
		fmt.Fprintf(c.Response(), "Kindly enter data with the event title and description only in order to get ticket")
		return err
	}
	return c.JSONPretty(http.StatusOK, ticket, "  ")
}

func (h *Handlers) APIDeleteTicket(c echo.Context) error {
	log.Println("APIDeleteTicket handler in action...", c.Request().Method)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Response().WriteHeader(http.StatusNotFound)
		return err
	}
	err = h.repo.DeleteTicket(id)
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server Error"))
		return err
	}

	if err != nil {
		fmt.Fprintf(c.Response(), "Kindly enter data with the event title and description only in order to get ticket")
		return err
	}
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func (h *Handlers) APIGetListTickets(c echo.Context) error {
	log.Println("APIGetListTickets handler in action...", c.Request().Method)
	ticketList, err := h.repo.ListTickets()
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server Error when retrive tickets list"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	return c.JSONPretty(http.StatusOK, ticketList, "  ")
}

func (h *Handlers) APIGetListTicketsByUser(c echo.Context) error {
	log.Println("APIGetListTicketsByUser handler in action...", c.Request().Method)
	cookie, err := c.Cookie("email")
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server Error when email from cookies"))
		c.Response().WriteHeader(http.StatusNotFound)
		return err
	}
	fmt.Println(cookie.Value)
	ticketList, err := h.repo.ListTicketsByUser(cookie.Value)
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server Error when retrive tickets list"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	return c.JSONPretty(http.StatusOK, ticketList, "  ")
}

func (h *Handlers) APIUpdateTicket(c echo.Context) error {
	log.Println("APIUpdateTicket handler in action...", c.Request().Method)
	ticket := entities.Ticket{}
	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Fprintf(c.Response(), "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &ticket)
	_, err = h.repo.UpdateTicket(ticket)
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server Error when updatetin particular ticket"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	return c.JSONPretty(http.StatusOK, ticket, "  ")
}
