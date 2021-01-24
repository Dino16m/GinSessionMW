package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionMiddleware struct {
	PayloadFunc  PayloadFunc
	Unauthorized UnauthorizedFunc
	Options      sessions.Options
	Repo         UserRepo
	SessionFunc  SessionFunc
}

const MapKey = "AuthUser"
const SessionAlias = "session_id"

func New(
	// Function that returns a unique key when called with a user instance
	PayloadFunc PayloadFunc,
	/** Function called when an unauuthorized user tries to access a
		protected resource
		Function is called with *gin.Context type as argument
		Usually, you should use it to return an error page or a json response
	**/
	Unauthorized UnauthorizedFunc,
	Options sessions.Options,
	// Function that returns a user model when called with the unique key
	// from PayloadFunc
	// Should return nil on user not found
	Repo UserRepo,
	// Function that returns a session instance when called with context
	// Defaults to github.com/gin-contrib/sessions instance
	SessionFunc SessionFunc,
) *SessionMiddleware {
	sm := SessionMiddleware{
		PayloadFunc:  PayloadFunc,
		Unauthorized: Unauthorized,
		Options:      Options,
		Repo:         Repo,
		SessionFunc:  SessionFunc,
	}
	if sm.SessionFunc == nil {
		sm.SessionFunc = func(c *gin.Context) Session {
			session := sessions.Default(c)
			return session
		}
	}
	return &sm
}
func (sm *SessionMiddleware) GetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		sm.handleRequest(c)
	}
}

func (sm *SessionMiddleware) Login(c *gin.Context, user interface{}) {
	identityKey := sm.PayloadFunc(user)
	sm.login(c, identityKey)
}

func (sm *SessionMiddleware) login(c *gin.Context, identityKey interface{}) {
	session := sm.SessionFunc(c)
	if sm.Options != (sessions.Options{}) {
		session.Options(sm.Options)
	}
	session.Set(SessionAlias, identityKey)
	session.Save()
}

func (sm *SessionMiddleware) Logout(c *gin.Context) {
	session := sm.SessionFunc(c)
	session.Clear()
	session.Save()
}
func (sm *SessionMiddleware) GetAuthUser(c *gin.Context) interface{} {
	userMap := c.GetStringMap(MapKey)
	if userMap == nil {
		return nil
	}
	return userMap[MapKey]
}

func (sm *SessionMiddleware) handleRequest(c *gin.Context) {
	session := sm.SessionFunc(c)
	identityKey := session.Get(SessionAlias)

	if identityKey == nil {
		sm.Unauthorized(c)
		return
	}

	user := sm.Repo(identityKey)
	if user == nil {
		sm.Unauthorized(c)
	} else {
		userMap := map[string]interface{}{
			MapKey: user,
		}
		c.Set(MapKey, userMap)
		c.Next()
	}
}
