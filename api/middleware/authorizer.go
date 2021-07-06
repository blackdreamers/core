package middleware

import (
	"encoding/gob"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	log "github.com/blackdreamers/go-micro/v3/logger"
)

const (
	UserAuthKey = "user_auth"
	// Anonymous no session
	Anonymous = "anonymous"
	// Unemployed expired session
	Unemployed = "unemployed"
)

type Authorizer struct{}

type UserAuth struct {
	State bool `json:"state"`
}

func (a *Authorizer) Init() ([]gin.HandlerFunc, error) {
	auth := &BasicAuthorizer{}
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			if strings.Contains(c.Request.URL.Path, "login") {
				return
			}

			role, state := auth.GetUser(c)
			if !state || !auth.CheckPermission(c, role) {
				switch role {
				case Anonymous:
					auth.RequireLogIn(c)
				case Unemployed:
					auth.RequireReLogIn(c)
				default:
					auth.RequirePermission(c)
				}
			}
		},
	}, nil
}

type BasicAuthorizer struct{}

func (a *BasicAuthorizer) GetUser(c *gin.Context) (role string, state bool) {
	if _, err := c.Cookie(SessionName); err != nil {
		return Anonymous, false
	}

	session := sessions.Default(c)
	user := session.Get(UserAuthKey)
	if user != nil {
		// rest ttl
		session.Set(SessionName, user)
		if err := session.Save(); err != nil {
			log.Error(err)
		}

		if !user.(*UserAuth).State {
			return "", false
		}

		return "", true
	}
	return Unemployed, false
}

func (a *BasicAuthorizer) CheckPermission(c *gin.Context, role string) bool {
	return true
}

func (a *BasicAuthorizer) RequireLogIn(c *gin.Context) {
	a.Abort(c, http.StatusUnauthorized, "need to login")
}

func (a *BasicAuthorizer) RequireReLogIn(c *gin.Context) {
	a.Abort(c, 499, "invalid login")
}

func (a *BasicAuthorizer) RequirePermission(c *gin.Context) {
	a.Abort(c, http.StatusForbidden, "no permission")
}

func (a *BasicAuthorizer) Abort(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{
		"code":      code,
		"msg":       msg,
		"data":      nil,
		"timestamp": time.Now().Unix(),
	})
}

func init() {
	gob.Register(&UserAuth{})
}
