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
	repo                   mocks.RepoInterface
	handler                Handlers
	ticket                 entities.Ticket
	user, unauthorizeduser entities.Users
}

// //It's executed before every test function

func (suite *UnitTestSuite) SetupTest() {
	suite.repo = mocks.RepoInterface{}
	suite.handler = Handlers{
		repo:         &suite.repo,
		sessionStore: sessions.NewCookieStore([]byte("my_secret_key")),
	}
	suite.ticket = entities.Ticket{
		Title:       "Test Title",
		Description: "Test Description",
		Status:      "Создана",
		Severity:    "5",
		CreatedAt:   time.Now().Local(),
		Number:      10,
		OwnerEmail:  "test@test.com",
	}
	suite.user = entities.Users{
		ID:       5,
		Name:     "admin",
		Password: "$2a$14$KSP3T3/6V9mX.H41neLsyu7JQEEBMnJtIGF2Ley9CykRu7GU8AczS",
		Email:    "test@test.ru",
	}
	suite.unauthorizeduser = entities.Users{
		ID:       5,
		Name:     "admin",
		Password: "$2a$14$KSP3T3/ABCD.H41neLsyu7JQEEBMnJtIGF2Ley9CykRu7GU8AczS",
		Email:    "test@test.ru",
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

	//create an expectation
	suite.repo.On("CreateTicket", mock.Anything).Return(suite.ticket, nil)

	e := echo.New()
	//create form values which are passed to CreateTicket method
	form := make(url.Values)
	form.Set("title", "Test Title3")
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

func (suite *UnitTestSuite) TestAPICreateTicket() {
	json := `{
		"Title":"Test title",
		"Description":"Test description",
		"Severity":"5"
	}`

	suite.repo.Mock.On("CreateTicket", mock.Anything).Return(entities.Ticket{
		Title:       "Test title",
		Description: "Test description",
		Severity:    "5",
	}, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(json))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/tickets")

	err := suite.handler.APICreateTicket(c)
	assert.Nil(suite.T(), err)

}

func (suite *UnitTestSuite) TestAPIGetListTickets() {
	ticketList := []entities.Ticket{}
	ticketList = append(ticketList, suite.ticket)
	suite.repo.Mock.On("ListTickets", mock.Anything).Return(ticketList, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/tickets/")
	err := suite.handler.APIGetListTickets(c)
	suite.repo.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err)
}

func (suite *UnitTestSuite) TestAPIGetListTicketsByUser() {
	ticketList := []entities.Ticket{}
	ticketList = append(ticketList, suite.ticket)
	cookie := &http.Cookie{
		Name:  "email",
		Value: "test@test.com",
	}
	suite.repo.Mock.On("ListTicketsByUser", "test@test.com").Return(ticketList, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := suite.handler.APIGetListTicketsByUser(c)
	suite.repo.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err)
}

func (suite *UnitTestSuite) TestAPIUpdateTicket() {
	json := `{
		"Title":"Test title edited",
		"Description":"Test description edited",
		"Severity":"5",
		"Number":"10"
	}`

	suite.repo.On("UpdateTicket", mock.Anything).Return(entities.Ticket{
		Title:       "Test title edited",
		Description: "Test description edited",
		Severity:    "5",
	}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(json))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/tickets/:id")
	c.SetParamNames("id")
	c.SetParamValues("10")
	err := suite.handler.APIUpdateTicket(c)
	suite.repo.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err)
}

func (suite *UnitTestSuite) TestAPIDeleteTicket() {
	suite.repo.On("DeleteTicket", mock.Anything).Return(nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/tickets/:id")
	c.SetParamNames("id")
	c.SetParamValues("10")
	err := suite.handler.APIDeleteTicket(c)
	suite.repo.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err)
}

func (suite *UnitTestSuite) TestLoginHandlerSUCCESS() {
	suite.repo.On("FindUser", mock.Anything).Return(suite.user, nil)
	e := echo.New()
	//create form values which are passed to CreateTicket method
	form := make(url.Values)
	form.Set("username", "admin")
	form.Set("password", "admin")
	//create new request
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Form = form
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := suite.handler.LoginHandler(c)
	suite.repo.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err)
}

func (suite *UnitTestSuite) TestLoginHandlerUnauthorizeBCrypt() {
	//When passwrod hash is incorrect
	suite.repo.On("FindUser", mock.Anything).Return(suite.unauthorizeduser, nil)
	e := echo.New()
	//create form values which are passed to CreateTicket method
	form := make(url.Values)
	form.Set("username", "admin")
	form.Set("password", "admin")
	//create new request
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Form = form
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := suite.handler.LoginHandler(c)
	suite.repo.AssertExpectations(suite.T())
	expectedError := "User unauthorized"
	assert.EqualErrorf(suite.T(), err, expectedError, "Error shoud be %v, but got %v", expectedError, err)
}

func (suite *UnitTestSuite) TestLoginHandlerUnauthorize() {
	//When user not found
	suite.repo.On("FindUser", mock.Anything).Return(entities.Users{
		Name:     "",
		Password: "somePass",
	}, nil)
	e := echo.New()
	//create form values which are passed to CreateTicket method
	form := make(url.Values)
	form.Set("username", "admin")
	form.Set("password", "admin")
	//create new request
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Form = form
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := suite.handler.LoginHandler(c)
	suite.repo.AssertExpectations(suite.T())
	expectedError := "User unauthorized"
	assert.EqualErrorf(suite.T(), err, expectedError, "Error shoud be %v, but got %v", expectedError, err)
}
