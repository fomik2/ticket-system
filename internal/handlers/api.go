package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fomik2/ticket-system/internal/entities"
	"github.com/golang-jwt/jwt/v4"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Create the Signin handler
func (h *Handlers) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
	}
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

//JWT middleware providing authentiction check for all handlers
func (h *Handlers) JWTAuthMiddleWare(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tknValid, err := h.CheckJWT(w, r)
		if tknValid == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil || !tknValid.Valid {
			if err == http.ErrNoCookie || err == jwt.ErrSignatureInvalid || !tknValid.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		f(w, r)
	}
}

//CheckJW provide JWT checking
func (h *Handlers) CheckJWT(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}
	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return tkn, nil

	// Finally, return the welcome message to the user, along with their
	// username given in the token
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

func (h *Handlers) APICreateTicket(w http.ResponseWriter, r *http.Request) {
	ticket := entities.Ticket{}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &ticket)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ticket)
}

func (h *Handlers) APIGetTicket(w http.ResponseWriter, r *http.Request) {
	id, err := getTicketID(w, r)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Internal server Error"))
		return
	}
	ticket, err := h.repo.GetTicket(id)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Internal server Error"))
		return
	}

	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to get ticket")
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ticket)
}
