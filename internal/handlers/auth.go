package handlers

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//HashPassword bcrypt password hashing
func (h *Handlers) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return string(bytes), err
	}
	return string(bytes), err
}

//CheckPasswordHash bcrypt password checking
func (h *Handlers) CheckPasswordHash(hashedPassFromDB, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassFromDB), []byte(password))
	if err != nil {
		return false, err
	}
	return true, err
}

//Authentication middleware providing authentiction check for all handlers
func (h *Handlers) Authentication(f http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		session, err := h.sessionStore.Get(r, "session.id")
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server error"))
			return
		}
		authenticated := session.Values["authenticated"]
		if authenticated != nil && authenticated != false {
			f(writer, r)
		} else {
			if r.RequestURI != "/login" {
				http.Redirect(writer, r, "/login", http.StatusSeeOther)
			}
		}
	}
}

func (h *Handlers) Login(writer http.ResponseWriter, r *http.Request) {
	log.Println("Login handler in action....", r.Method)
	h.templs["auth"].Execute(writer, nil)
}

func (h *Handlers) LoginHandler(writer http.ResponseWriter, r *http.Request) {
	log.Println("LoginHandler in action....Authentication process...", r.Method)
	// ParseForm parses the raw query from the URL and updates r.Form
	err := r.ParseForm()
	if err != nil {
		http.Error(writer, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}

	user, err := h.repo.FindUser(r.Form["username"][0])

	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error while find user in DB"))
		return
	}

	if user.Name == "" {
		writer.Write([]byte("Unauthorized. (No user found)"))
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	checkPass, err := h.CheckPasswordHash(user.Password, r.Form["password"][0])
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error while compare password with hash"))
		return
	}

	if checkPass {
		// It returns a new session if the sessions doesn't exist
		session, err := h.sessionStore.Get(r, "session.id")
		if err != nil {
			log.Println(err)
			writer.Write([]byte("Internal server error. Session cound not been decoded."))
			return
		}
		session.Values["authenticated"] = true
		session.Values["email"] = user.Email
		// Saves all sessions used during the current request
		session.Save(r, writer)
		http.Redirect(writer, r, "/", http.StatusSeeOther)

	} else {
		http.Error(writer, "Invalid Credentials", http.StatusUnauthorized)
	}

}

func (h *Handlers) LogoutHandler(writer http.ResponseWriter, r *http.Request) {
	// Get registers and returns a session for the given name and session store.
	session, _ := h.sessionStore.Get(r, "session.id")
	// Set the authenticated value on the session to false
	session.Values["authenticated"] = false
	err := session.Save(r, writer)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error. Session cound not been saved."))
		return
	}
	http.Redirect(writer, r, "/login", http.StatusSeeOther)

}
