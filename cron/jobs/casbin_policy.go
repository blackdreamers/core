package jobs

import (
	"github.com/blackdreamers/core/api/auth"
	"github.com/blackdreamers/core/consts"
	"github.com/blackdreamers/core/cron"
	log "github.com/blackdreamers/core/logger"
)

type CasbinPolicy struct{}

func (c *CasbinPolicy) Name() string {
	return "casbin_policy"
}

func (c *CasbinPolicy) Spec() string {
	return "@every 5m"
}

func (c *CasbinPolicy) Options() []cron.Option {
	return []cron.Option{
		cron.Wrapper(&cron.DelayIfStillRunning),
	}
}

func (c *CasbinPolicy) Run() {
	c.execute()
}

func (c *CasbinPolicy) execute() {
	if err := auth.LoadPolicy(); err != nil {
		log.Field(consts.ErrKey, err).Log(log.ErrorLevel, "casbin load policy")
	}
}
