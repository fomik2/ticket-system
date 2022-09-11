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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	repo    mocks.RepoInterface
	handler Handlers
}

// //It's executed before every test function

func (suite *UnitTestSuite) SetupTest() {
	suite.repo = mocks.RepoInterface{}
	suite.handler = Handlers{
		repo:         &suite.repo,
		sessionStore: sessions.NewCookieStore([]byte("my_secret_key")),
	}
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
func (suite *UnitTestSuite) TestApiGetTicket() {
	suite.repo.On("GetTicket", 10).Return(entities.Ticket{
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
	err := suite.handler.APIGetTicket(c)
	suite.repo.AssertExpectations(suite.T())
	suite.Nil(err)

}

func (suite *UnitTestSuite) TestCreateTicket() {

	ticket := entities.Ticket{
		Title:       "Test Title",
		Description: "Test Description",
		Status:      "Создана",
		Severity:    "5",
		CreatedAt:   time.Now().Local(),
		Number:      10,
		OwnerEmail:  "test@test.com",
	}
	//create an expectation
	suite.repo.On("CreateTicket", mock.Anything).Return(ticket, nil)

	e := echo.New()
	//create form values which are passed to CreateTicket method
	form := make(url.Values)
	form.Set("title", "Test Title")
	form.Set("description", "Test description")
	form.Set("severity", "5")
	//create new request
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Form = form

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	//create a session and save it
	session, err := suite.handler.sessionStore.Get(c.Request(), "session.id")
	if err != nil {
		suite.T().Errorf("can't create test session")
	}
	session.Values["email"] = "test@test.com"
	session.Save(c.Request(), c.Response())
	c.SetPath("/")
	err = suite.handler.CreateTicket(c)
	suite.repo.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err)
}
