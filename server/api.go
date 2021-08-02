package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	coreapi "github.com/blackdreamers/core/api"
	"github.com/blackdreamers/core/api/middleware"
	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/constant"
	"github.com/blackdreamers/core/utils"
	"github.com/blackdreamers/go-micro/v3"
	log "github.com/blackdreamers/go-micro/v3/logger"
)

var (
	api = &apiEntry{}
)

type apiEntry struct {
	r       *gin.Engine
	s       *http.Server
	routers []Router
}

type Router interface {
	Router(r *gin.Engine)
}

func AddRouters(apiRouters ...Router) {
	api.routers = append(api.routers, apiRouters...)
}

func ApiEngine() *gin.Engine {
	return api.r
}

func (a *apiEntry) init(opts ...micro.Option) error {
	a.r = gin.New()
	a.r.ForwardedByClientIP = true
	gin.ForceConsoleColor()
	env := config.Env
	switch env {
	case constant.Dev:
		env = gin.DebugMode
	case constant.Test:
		env = gin.TestMode
	case constant.Prod,
		constant.Release:
		env = gin.ReleaseMode
	}
	gin.SetMode(env)

	capi := &coreapi.API{}
	a.r.NoRoute(capi.API404)
	capi.Validator()

	if config.IsDevEnv() {
		a.r.Use(gin.LoggerWithFormatter(logFormatter), gin.Recovery())
	} else {
		mLog := &middleware.Log{}
		ml, err := mLog.Init()
		if err != nil {
			return err
		}
		a.r.Use(ml...)
	}

	return nil
}

func (a *apiEntry) run() error {
	for _, m := range middleware.Ms {
		ms, err := m.Init()
		if err != nil {
			return err
		}
		a.r.Use(ms...)
	}

	pprof.Register(a.r)

	for _, r := range a.routers {
		r.Router(a.r)
	}

	a.s = &http.Server{
		Addr:    config.Service.Addr,
		Handler: a.r,
	}

	if log.V(log.InfoLevel, log.DefaultLogger) {
		log.Infof("HTTP API Listening on %s", config.Service.Addr)
	}

	if err := a.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func logFormatter(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}

	var bodyStr string
	if log.V(log.DebugLevel, log.DefaultLogger) && param.Request.Header.Get("Content-Type") == "application/json" {
		body, _ := ioutil.ReadAll(param.Request.Body)
		bodyJson := utils.JsonIndent(body)
		if len(bodyJson) > 0 {
			bodyStr = "\n" + bodyJson
		}
	}

	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v %v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		bodyStr,
		param.ErrorMessage,
	)
}
