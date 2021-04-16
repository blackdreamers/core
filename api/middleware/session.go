package middleware

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/constant"
)

type Session struct{}

func (s *Session) Init() ([]gin.HandlerFunc, error) {
	var sessionStore redis.Store

	switch config.Session.Store {
	case constant.CookieStore:
		sessionStore = cookie.NewStore([]byte(config.Session.Secret))
	case constant.MemoryStore:
		sessionStore = memstore.NewStore([]byte(config.Session.Secret))
	case constant.RedisStore:
		var err error
		sessionStore, err = redis.NewStoreWithDB(
			512,
			"tcp",
			config.Redis.Addrs[0],
			config.Redis.Password,
			strconv.Itoa(config.Session.DB),
			[]byte(config.Session.Secret),
		)
		if err != nil {
			return nil, err
		}
	default:
	}

	sessionStore.Options(sessions.Options{
		MaxAge:   config.Session.MaxAge,
		HttpOnly: config.Session.HttpOnly,
		Path:     "/",
	})

	return []gin.HandlerFunc{sessions.Sessions("session", sessionStore)}, nil
}

func init() {
	Middlewares(&Session{})
}
