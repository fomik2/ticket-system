package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword bcrypt password hashing
func (h *Handlers) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return string(bytes), err
	}
	return string(bytes), err
}

// CheckPasswordHash bcrypt password checking
func (h *Handlers) CheckPasswordHash(hashedPassFromDB, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassFromDB), []byte(password))
	if err != nil {
		return false, err
	}
	return true, err
}

// Authentication middleware providing authentiction check for all handlers
func (h *Handlers) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := h.sessionStore.Get(c.Request(), "session.id")
		if err != nil {
			log.Println(err)
			c.Response().Write([]byte("Internal server error"))
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
		authenticated := session.Values["authenticated"]
		if authenticated != nil && authenticated != false {
			next(c)
		} else {
			if c.Request().RequestURI != "/login" {
				http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
			}
		}
		return nil
	}
}

func (h *Handlers) Login(c echo.Context) error {
	log.Println("Login handler in action....", c.Request().Method)
	h.templs["auth"].Execute(c.Response(), nil)
	return nil
}

func (h *Handlers) LoginHandler(c echo.Context) error {
	log.Println("LoginHandler in action....Authentication process...", c.Request().Method)
	// ParseForm parses the raw query from the URL and updates r.Form

	user, err := h.repo.FindUser(c.FormValue("username"))

	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	if user.Name == "" {
		c.Response().Write([]byte("Unauthorized. (No user found)"))
		c.Response().WriteHeader(http.StatusUnauthorized)
		return errors.New("User unauthorized")
	}

	checkPass, err := h.CheckPasswordHash(user.Password, c.FormValue("password"))
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error while compare password with hash"))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return errors.New("User unauthorized")
	}

	if checkPass {
		// It returns a new session if the sessions doesn't exist
		session, err := h.sessionStore.Get(c.Request(), "session.id")
		if err != nil {
			log.Println(err)
			c.Response().Write([]byte("Internal server error. Session cound not been decoded."))
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}
		session.Values["authenticated"] = true
		session.Values["email"] = user.Email
		// Saves all sessions used during the current request
		session.Save(c.Request(), c.Response())
		http.Redirect(c.Response(), c.Request(), "/", http.StatusSeeOther)

	} else {
		http.Error(c.Response(), "Invalid Credentials", http.StatusUnauthorized)
	}
	return nil
}

func (h *Handlers) LogoutHandler(c echo.Context) error {
	// Get registers and returns a session for the given name and session store.
	session, _ := h.sessionStore.Get(c.Request(), "session.id")
	// Set the authenticated value on the session to false
	session.Values["authenticated"] = false
	err := session.Save(c.Request(), c.Response())
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("Internal server error. Session cound not been saved."))
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
	return nil
}
