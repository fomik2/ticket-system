package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/fomik2/ticket-system/internal/entities"
	"github.com/fomik2/ticket-system/mocks"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	repo RepoInterface
}

// //It's executed before every test function
// func (uts *UnitTestSuite) SetupTestsFunctions() {
// 	repo := mocks.RepoInterface{}
// 	uts.repo = &repo
// }

func TestApiGetTicket(t *testing.T) {
	repo := mocks.RepoInterface{}
	handler := Handlers{
		repo: &repo,
	}
	repo.On("GetTicket", 10).Return(entities.Ticket{
		Title:       "Test Title",
		Description: "Test Description",
		Status:      "Created",
		Severity:    "High",
		SLA:         time.Now().Local(),
		CreatedAt:   time.Now().Local(),
		Number:      10,
		OwnerEmail:  "test@test.com",
	}, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/tickets/:id")
	c.SetParamNames("id")
	c.SetParamValues("10")
	err := handler.APIGetTicket(c)
	repo.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestCreateTicket(t *testing.T) {
	repo := mocks.RepoInterface{}
	handler := Handlers{
		repo:         &repo,
		sessionStore: sessions.NewCookieStore([]byte("my_secret_key")),
	}

	ticket := entities.Ticket{
		Title:       "Test Title",
		Description: "Test Description",
		Status:      "Created",
		Severity:    "5",
		SLA:         time.Now().Local(),
		CreatedAt:   time.Now().Local(),
		Number:      10,
		OwnerEmail:  "test@test.com",
	}
	var ticketList []entities.Ticket
	ticketList = append(ticketList, ticket)
	repo.On("CreateTicket", ticket).Return(ticket, nil)
	repo.On("ListTickets").Return(ticketList, nil)
	e := echo.New()
	f := make(url.Values)
	f.Set("title", "Test Title")
	f.Set("description", "Test description")
	f.Set("severity", "5")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Saves all sessions used during the current request

	c := e.NewContext(req, rec)
	session, err := handler.sessionStore.Get(c.Request(), "session.id")
	if err != nil {
		t.Errorf("can't create test session")
	}
	session.Values["email"] = "test@test.com"
	session.Save(c.Request(), c.Response())
	c.SetPath("/")

	err = handler.CreateTicket(c)
	repo.AssertExpectations(t)
	assert.Nil(t, err)
}
