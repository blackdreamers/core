package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/longhaoteng/melody"

	"github.com/blackdreamers/core/api"
	"github.com/blackdreamers/core/api/auth"
	"github.com/blackdreamers/core/consts"
	log "github.com/blackdreamers/core/logger"
)

var (
	WS *melody.Melody
)

type Websocket struct {
	*api.API
}

func (w *Websocket) Router(r *gin.Engine) {
	WS = melody.New()
	r.GET("/ws", w.Upgrade)

	WS.HandleConnect(func(s *melody.Session) {})

	WS.HandleDisconnect(func(s *melody.Session) {})

	WS.HandlePong(func(s *melody.Session) {
		_ = s.Write([]byte("pong"))
	})
}

func (w *Websocket) Upgrade(c *gin.Context) {
	keys := make(map[string]interface{})
	keys[auth.TokenKey] = w.Session.GetToken(c)
	if err := WS.HandleRequestWithKeys(c.Writer, c.Request, keys); err != nil {
		log.Fields(consts.ErrKey, err).Logf(log.ErrorLevel, "WS request")
	}
}
