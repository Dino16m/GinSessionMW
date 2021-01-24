package test

import (
	"net/http/httptest"
	"testing"

	"github.com/dino16m/GinSessionMW/middleware"
	"github.com/dino16m/GinSessionMW/test/mocks"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type dummyUser struct {
	Username string
	Password string
}

type SessionTestSuite struct {
	suite.Suite
	payloadFunc      mocks.PayloadFunc
	userRepo         mocks.UserRepo
	unauthorizedFunc mocks.UnauthorizedFunc
	sessionFunc      mocks.SessionFunc
	middleware       *middleware.SessionMiddleware
	ctx              *gin.Context
}

var user = dummyUser{
	Username: "dino",
	Password: "password",
}

func (suite *SessionTestSuite) SetupTest() {
	suite.payloadFunc = *new(mocks.PayloadFunc)
	suite.userRepo = *new(mocks.UserRepo)
	suite.unauthorizedFunc = *new(mocks.UnauthorizedFunc)
	suite.sessionFunc = *new(mocks.SessionFunc)
	suite.middleware = middleware.New(
		suite.payloadFunc.Execute,
		suite.unauthorizedFunc.Execute,
		sessions.Options{},
		suite.userRepo.Execute,
		suite.sessionFunc.Execute,
	)
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
}

func (suite *SessionTestSuite) TestHandlerCallsUnauthorizedOnNoSession() {
	mw := suite.middleware.GetHandler()
	session := new(mocks.Session)
	session.On("Get", mock.AnythingOfType("string")).Return(nil)
	suite.sessionFunc.On("Execute", suite.ctx).Return(session)
	suite.unauthorizedFunc.On("Execute", suite.ctx).Return(nil)
	mw(suite.ctx)
	suite.unauthorizedFunc.AssertExpectations(suite.T())

}

func (suite *SessionTestSuite) TestHandlerCallsUnauthorizedOnUserNotFound() {
	mw := suite.middleware.GetHandler()
	session := new(mocks.Session)
	session.On("Get", mock.AnythingOfType("string")).Return("Dino")
	suite.sessionFunc.On("Execute", suite.ctx).Return(session)
	suite.userRepo.On("Execute", "Dino").Return(nil)
	suite.unauthorizedFunc.On("Execute", suite.ctx).Return(nil)
	mw(suite.ctx)
	suite.unauthorizedFunc.AssertExpectations(suite.T())
}

func (suite *SessionTestSuite) TestUserReturnedWhenSesssionExists() {
	mw := suite.middleware.GetHandler()
	session := new(mocks.Session)
	session.On("Get", mock.AnythingOfType("string")).Return("Dino")
	suite.sessionFunc.On("Execute", suite.ctx).Return(session)
	suite.userRepo.On("Execute", "Dino").Return(user)
	mw(suite.ctx) //call middleware first
	testUser := suite.middleware.GetAuthUser(suite.ctx)
	suite.Equal(user, testUser)
}

func (suite *SessionTestSuite) TestLogin() {
	username := "Dino"
	suite.payloadFunc.On("Execute", user).Return(username)
	session := new(mocks.Session)
	session.On("Set", middleware.SessionAlias, username).Return(nil)
	session.On("Save").Return(nil)
	suite.sessionFunc.On("Execute", suite.ctx).Return(session)
	suite.middleware.Login(suite.ctx, user)
	session.AssertExpectations(suite.T())
}

func (suite *SessionTestSuite) TestLogout() {
	session := new(mocks.Session)
	session.On("Save").Return(nil)
	session.On("Clear").Return(nil)
	suite.sessionFunc.On("Execute", suite.ctx).Return(session)
	suite.middleware.Logout(suite.ctx)
	session.AssertExpectations(suite.T())
}

func TestSessionTestSuite(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}
