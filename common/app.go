package common

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
	AppEx
}

type AppEx struct {
	Srv    []string  `json:"srv"`
	Config AppConfig `json:"app"`
}

type AppConfig struct {
	WxMini map[string]struct {
		AppID     string `json:"app_id"`
		AppSecret string `json:"app_secret"`
	} `json:"wx_mini"`

	Storage map[string]struct {
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
	} `json:"storage"`
}

func GetApps() (apps Apps, err error) {
	var appsJson string
	if appsJson, err = redis.Client().Get(context.Background(), appCacheKey).Result(); redis.IsNotNil(err) {
		return nil, err
	}

	if err = json.Unmarshal([]byte(appsJson), &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func GetByName(name string) (*App, error) {
	apps, err := GetApps()
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

func GetBySrv(srv string) (*App, error) {
	apps, err := GetApps()
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

func GetConfigByName(name string) (*AppConfig, error) {
	apps, err := GetApps()
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

func GetConfigBySrv(srv string) (conf *AppConfig, err error) {
	apps, err := GetApps()
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

	if app, err = GetBySrv(cliName); err != nil {
		return nil, err
	}

	return app, nil
}

func GetConfigByClient(ctx context.Context) (conf *AppConfig, err error) {
	cliName, ok := utils.GetClientSrvName(ctx)
	if !ok {
		return nil, errors.New(errs.ParseClientNameErr)
	}

	if conf, err = GetConfigBySrv(cliName); err != nil {
		return nil, err
	}

	return conf, nil
}

func SetApps(apps Apps) error {
	if appsJson, err := json.Marshal(&apps); err != nil {
		return err
	} else {
		return redis.Client().Set(context.Background(), appCacheKey, string(appsJson), -1).Err()
	}
}

func SetApp(app *App) error {
	apps, err := GetApps()
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

	return SetApps(apps)
}

func RemoveApp(app *App) error {
	apps, err := GetApps()
	if err != nil {
		return err
	}

	for i, a := range apps {
		if a.Name == app.Name {
			apps = append(apps[:i], apps[i+1:]...)
			break
		}
	}

	return SetApps(apps)
}
