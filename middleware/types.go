package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type sessionErr struct {
	message string
	code    int
}

type UserRepo func(identityKey interface{}) interface{}
type PayloadFunc func(interface{}) interface{}
type UnauthorizedFunc func(c *gin.Context)
type SessionFunc func(c *gin.Context) Session

type Session interface {
	//this interface is implemented by github.com/gin-contrib/sessions
	Get(key interface{}) interface{}
	Set(key interface{}, val interface{})
	Delete(key interface{})
	Clear()
	Options(sessions.Options)
	Save() error
}
