package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// store the secret key in env variable in production
var store = sessions.NewCookieStore([]byte("my_secret_key"))

//Authentication middleware providing authentiction check for all handlers
func (h *Handlers) Authentication(f http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session.id")
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
	user, _ := h.repo.FindUser(r.Form["username"][0], r.Form["password"][0])
	if err != nil {
		log.Println(err)
		writer.Write([]byte("Internal server error"))
		return
	}
	if user.Name != "" {
		// It returns a new session if the sessions doesn't exist
		session, err := store.Get(r, "session.id")
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
	session, _ := store.Get(r, "session.id")
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
