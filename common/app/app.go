package app

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/thoas/go-funk"

	"github.com/blackdreamers/core/cache/redis"
	"github.com/blackdreamers/core/db"
	errs "github.com/blackdreamers/core/errors"
	"github.com/blackdreamers/core/utils"
)

const (
	appCacheKey = "apps"
)

type Apps []App

type App struct {
	db.Model
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	State bool   `json:"state"`
	Ex
}

type Ex struct {
	Srv    []string `json:"srv"`
	Config Config   `json:"config"`
}

type Config map[string]interface{}

func (c *Config) Parse(target interface{}, keys ...string) error {
	var err error
	var data []byte

	if len(keys) == 0 {
		data, err = json.Marshal(c)
	} else {
		conf := *c
		for _, key := range keys {
			if _, ok := conf[key]; !ok {
				return errors.New(errs.ConfigNotFound)
			}
			conf = conf[key].(map[string]interface{})
		}
		data, err = json.Marshal(conf)
	}

	if err != nil {
		return err
	}
	return json.Unmarshal(data, &target)
}

func GetApps(ctx context.Context) (apps Apps, err error) {
	var appsJson string
	if appsJson, err = redis.Client().Get(ctx, appCacheKey).Result(); err != nil && redis.IsNotNil(err) {
		return nil, err
	}
	if redis.IsNil(err) {
		return apps, nil
	}
	if err = json.Unmarshal([]byte(appsJson), &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func GetByName(ctx context.Context, name string) (*App, error) {
	apps, err := GetApps(ctx)
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		if app.Name == name {
			return &app, nil
		}
	}
	return nil, errors.New(errs.AppNotFound)
}

func GetBySrv(ctx context.Context, srv string) (*App, error) {
	apps, err := GetApps(ctx)
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		if funk.ContainsString(app.Srv, srv) {
			return &app, nil
		}
	}
	return nil, errors.New(errs.AppNotFound)
}

func GetConfigByName(ctx context.Context, name string) (*Config, error) {
	apps, err := GetApps(ctx)
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		if app.Name == name {
			return &app.Config, nil
		}
	}
	return nil, errors.New(errs.AppNotFound)
}

func GetConfigBySrv(ctx context.Context, srv string) (conf *Config, err error) {
	apps, err := GetApps(ctx)
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		if funk.ContainsString(app.Srv, srv) {
			return &app.Config, nil
		}
	}
	return nil, errors.New(errs.AppNotFound)
}

func GetByClient(ctx context.Context) (app *App, err error) {
	cliName, ok := utils.GetClientSrvName(ctx)
	if !ok {
		return nil, errors.New(errs.ParseClientNameErr)
	}

	if app, err = GetBySrv(ctx, cliName); err != nil {
		return nil, err
	}

	return app, nil
}

func GetConfigByClient(ctx context.Context) (conf *Config, err error) {
	cliName, ok := utils.GetClientSrvName(ctx)
	if !ok {
		return nil, errors.New(errs.ParseClientNameErr)
	}

	if conf, err = GetConfigBySrv(ctx, cliName); err != nil {
		return nil, err
	}

	return conf, nil
}

func SetApps(ctx context.Context, apps Apps) error {
	if appsJson, err := json.Marshal(&apps); err != nil {
		return err
	} else {
		return redis.Client().Set(ctx, appCacheKey, string(appsJson), -1).Err()
	}
}

func SetApp(ctx context.Context, app *App) error {
	apps, err := GetApps(ctx)
	if err != nil {
		return err
	}

	exist := false
	for i, a := range apps {
		if a.Name == app.Name {
			exist = true
			apps[i] = *app
			break
		}
	}

	if !exist {
		apps = append(apps, *app)
	}

	return SetApps(ctx, apps)
}

func RemoveApp(ctx context.Context, app *App) error {
	apps, err := GetApps(ctx)
	if err != nil {
		return err
	}

	for i, a := range apps {
		if a.Name == app.Name {
			apps = append(apps[:i], apps[i+1:]...)
			break
		}
	}

	return SetApps(ctx, apps)
}
