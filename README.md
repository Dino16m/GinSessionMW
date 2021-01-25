# GinSessionMiddleware

This is a session middleare for Golang's gin-gonic

This middleware:

* Guards user routes that require authentication
* Provides a login method to create a session for a user
* Provides a logout method to end a user's session 




## Installation

Download and install the package from PyPi:
````bash
go get github.com/dino16m/GinSessionMW
````

### Example use for middleware
````golang
package example

import (
	"github.com/dino16m/GinSessionMW/middleware"
    "github.com/gin-contrib/sessions"
    github.com/gin-gonic/gin
)

type struct User{
	Email string
    Password string
}

func main(){
	middlwr := middleware.New(
    	PayloadFunc:  func(user interface{}) interface{}{return user.(*User).Email},
        UnauthorizedFunc: func(c *gin.Context){c.JSON(403, gin.H{"msg": "Unauthorized"})},
        Options: sessions.Options{}, //default
        Repo: func(username interface{}) interface{}{return &User{"example.com", "dummypassword"}},
        sessionFunc: func(c *gin.Context) middleware.Session,
        )
        
        r := gin.Default()
        store, _ := sessions.redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
		r.Use(sessions.Sessions("mysession", store))
        
        auth := r.Group("/protected")
        auth.Use(middlwr.GetHandler()){
        	auth.GET("/get", func(c *gin.Context)){
            	c.JSON(200)
            }
        } 	
}
````

### Example use for Login method
````golang
package example

import (
	"github.com/dino16m/GinSessionMW/middleware"
    "github.com/gin-contrib/sessions"
    github.com/gin-gonic/gin
)

func main(){
	middlwr := middleware.New(
    	PayloadFunc:  func(user interface{}) interface{}{return user.(*User).Email},
        UnauthorizedFunc: func(c *gin.Context){c.JSON(403, gin.H{"msg": "Unauthorized"})},
        Options: sessions.Options{}, //default
        Repo: func(username interface{}) interface{}{return &User{"example.com", "dummypassword"}},
        sessionFunc: func(c *gin.Context) middleware.Session,
        )
        
	r := gin.Default()
    r.GET("/login", func(c *gin.Context)){
      user = User{}// Get User Somehow
      middlwr.Login(c *gin.Context)
      c.JSON(200)
    }
}
````
#### There are other important methods defined on ths struct